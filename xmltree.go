package wsdl

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/faustbrian/go-xsd/datatype"
)

type xmlNode struct {
	// owner belongs to the temporary parse tree, never to a published model.
	owner      *parseState
	name       xml.Name
	attributes []xml.Attr
	children   []*xmlNode
	text       strings.Builder
	content    []xmlContent
	namespaces map[string]string
	location   Location
	baseURI    string
}

func (n *xmlNode) contextError() error {
	return n.owner.contextError()
}

// A terminal conversion error keeps its established precedence. A successful
// conversion must not publish its partially built value after cancellation.
func finishConversion[T any](node *xmlNode, value *T, err *error) {
	if *err == nil {
		*err = node.contextError()
		if *err != nil {
			var zero T
			*value = zero
		}
	}
}

func finishConversionValue[T any](node *xmlNode, value *T) {
	if node.contextError() != nil {
		var zero T
		*value = zero
	}
}

func finishConversionError(node *xmlNode, err *error) {
	if *err == nil {
		*err = node.contextError()
	}
}

// Parsed children already share their invocation owner. Constructed internal
// trees acquire it only when explicitly admitted by a context-bearing decoder.
func (n *xmlNode) useOwner(owner *parseState) error {
	if err := owner.contextError(); err != nil {
		return err
	}
	if n.owner == owner {
		return nil
	}
	n.owner = owner
	for _, child := range n.children {
		if err := child.useOwner(owner); err != nil {
			return err
		}
	}
	return owner.contextError()
}

var marshalNode = marshalXMLNode

var (
	writeNode     = writeXMLNode
	escapeXMLText = xml.EscapeText
)

func assignBaseURIs(node *xmlNode, inherited string, owner ...*parseState) error {
	if err := parseOwnerError(owner); err != nil {
		return err
	}
	base := inherited
	if reference, exists := xmlBaseReference(node.attributes, node.owner); exists {
		resolved, err := resolveURI(inherited, reference)
		if err != nil {
			return fmt.Errorf("wsdl: resolve xml:base %q: %w", reference, err)
		}
		base = resolved
	}
	node.baseURI = base
	for _, child := range node.children {
		if err := assignBaseURIs(child, base, owner...); err != nil {
			return err
		}
	}
	return nil
}

func xmlBaseReference(attributes []xml.Attr, owner ...*parseState) (string, bool) {
	if parseOwnerError(owner) != nil {
		return "", false
	}
	for _, attribute := range attributes {
		if parseOwnerError(owner) != nil {
			return "", false
		}
		if attribute.Name == (xml.Name{
			Space: "http://www.w3.org/XML/1998/namespace",
			Local: "base",
		}) {
			return attribute.Value, true
		}
	}
	return "", false
}

func resolveURI(base, reference string) (string, error) {
	referenceURL, err := url.Parse(reference)
	if err != nil {
		return "", err
	}
	if base == "" {
		return referenceURL.String(), nil
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	return baseURL.ResolveReference(referenceURL).String(), nil
}

type xmlContent struct {
	text  []byte
	child *xmlNode
}

type parseState struct {
	ctx        context.Context
	options    ParseOptions
	elements   int
	attributes int
	textBytes  int64
}

func (s *parseState) contextError() error {
	if s == nil || s.ctx == nil {
		return nil
	}
	return s.ctx.Err()
}

// Tree helpers also serve intentionally context-free internal model fixtures.
func parseOwnerError(owner []*parseState) error {
	if len(owner) == 0 {
		return nil
	}
	return owner[0].contextError()
}

type componentCounts struct {
	imports    int
	operations int
	bindings   int
	endpoints  int
	extensions int
}

func validateCoreNCNames(node *xmlNode, coreNamespace string, owner ...*parseState) error {
	if err := parseOwnerError(owner); err != nil {
		return err
	}
	if node.name.Space == coreNamespace {
		for _, local := range []string{"name", "messageLabel"} {
			if err := parseOwnerError(owner); err != nil {
				return err
			}
			value, exists := node.namespacedAttribute("", local)
			if !exists {
				continue
			}
			if datatype.ValidateBuiltInLexical("NCName", value) != nil {
				return fmt.Errorf(
					"wsdl: {%s}%s attribute %s value %q is not an NCName",
					node.name.Space,
					node.name.Local,
					local,
					value,
				)
			}
		}
	}
	for _, child := range node.children {
		if err := validateCoreNCNames(child, coreNamespace, owner...); err != nil {
			return err
		}
	}
	return nil
}

func enforceComponentLimits(
	root *xmlNode,
	coreNamespace string,
	options ParseOptions,
	owner ...*parseState,
) error {
	if err := parseOwnerError(owner); err != nil {
		return err
	}
	counts := componentCounts{}
	countComponents(root, coreNamespace, false, &counts, owner...)
	if err := parseOwnerError(owner); err != nil {
		return err
	}
	limits := []struct {
		name  string
		count int
		max   int
	}{
		{name: "imports", count: counts.imports, max: options.MaxImports},
		{name: "operations", count: counts.operations, max: options.MaxOperations},
		{name: "bindings", count: counts.bindings, max: options.MaxBindings},
		{name: "endpoints", count: counts.endpoints, max: options.MaxEndpoints},
		{name: "extensions", count: counts.extensions, max: options.MaxExtensions},
	}
	for _, limit := range limits {
		if limit.count > limit.max {
			return fmt.Errorf(
				"%w: %s exceed %d",
				ErrLimitExceeded,
				limit.name,
				limit.max,
			)
		}
	}
	return nil
}

func countComponents(
	node *xmlNode,
	coreNamespace string,
	parentCore bool,
	counts *componentCounts,
	owner ...*parseState,
) {
	if parseOwnerError(owner) != nil {
		return
	}
	core := node.name.Space == coreNamespace
	if core {
		switch node.name.Local {
		case "import", "include":
			counts.imports++
		case "operation":
			counts.operations++
		case "binding":
			counts.bindings++
		case "port", "endpoint":
			counts.endpoints++
		}
		for _, attribute := range node.attributes {
			if parseOwnerError(owner) != nil {
				return
			}
			if attribute.Name.Space != "" && attribute.Name.Space != "xmlns" &&
				attribute.Name.Space != coreNamespace &&
				attribute.Name.Space != "http://www.w3.org/XML/1998/namespace" {
				counts.extensions++
			}
		}
	} else if parentCore &&
		(node.name.Space != NamespaceXMLSchema || node.name.Local != "schema") {
		counts.extensions++
	}
	for _, child := range node.children {
		if parseOwnerError(owner) != nil {
			return
		}
		countComponents(child, coreNamespace, core, counts, owner...)
	}
}

func marshalXMLNode(node *xmlNode) (converted []byte, conversionErr error) {
	if err := node.contextError(); err != nil {
		return nil, err
	}
	defer finishConversion(node, &converted, &conversionErr)
	var output bytes.Buffer
	if err := writeNode(&output, node, true); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func writeXMLNode(output *bytes.Buffer, node *xmlNode, root bool) (conversionErr error) {
	if err := node.contextError(); err != nil {
		return err
	}
	defer finishConversionError(node, &conversionErr)
	return writeXMLNodeScoped(output, node, nil, root)
}

func writeXMLNodeScoped(
	output *bytes.Buffer,
	node *xmlNode,
	inherited map[string]string,
	emitNamespaces bool,
) (conversionErr error) {
	if err := node.contextError(); err != nil {
		return err
	}
	defer finishConversionError(node, &conversionErr)
	name, err := lexicalXMLName(node.name, node.namespaces, false, node.owner)
	if err != nil {
		return err
	}
	output.WriteByte('<')
	output.WriteString(name)
	if emitNamespaces {
		prefixes := make([]string, 0, len(node.namespaces))
		for prefix, namespace := range node.namespaces {
			if err := node.contextError(); err != nil {
				return err
			}
			if inheritedNamespace, ok := inherited[prefix]; ok && inheritedNamespace == namespace {
				continue
			}
			prefixes = append(prefixes, prefix)
		}
		sort.Strings(prefixes)
		for _, prefix := range prefixes {
			if err := node.contextError(); err != nil {
				return err
			}
			output.WriteString(" xmlns")
			if prefix != "" {
				output.WriteByte(':')
				output.WriteString(prefix)
			}
			output.WriteString(`="`)
			if err := escapeXMLText(output, []byte(node.namespaces[prefix])); err != nil {
				return err
			}
			if err := node.contextError(); err != nil {
				return err
			}
			output.WriteByte('"')
		}
	}
	for _, attribute := range node.attributes {
		if err := node.contextError(); err != nil {
			return err
		}
		if attribute.Name.Space == "xmlns" ||
			(attribute.Name.Space == "" && attribute.Name.Local == "xmlns") {
			continue
		}
		attributeName, nameErr := lexicalXMLName(attribute.Name, node.namespaces, true, node.owner)
		if nameErr != nil {
			return nameErr
		}
		output.WriteByte(' ')
		output.WriteString(attributeName)
		output.WriteString(`="`)
		if err := escapeXMLText(output, []byte(attribute.Value)); err != nil {
			return err
		}
		if err := node.contextError(); err != nil {
			return err
		}
		output.WriteByte('"')
	}
	output.WriteByte('>')
	for _, content := range node.content {
		if err := node.contextError(); err != nil {
			return err
		}
		if content.child != nil {
			if err := writeXMLNodeScoped(
				output,
				content.child,
				node.namespaces,
				true,
			); err != nil {
				return err
			}
			continue
		}
		if err := escapeXMLText(output, content.text); err != nil {
			return err
		}
		if err := node.contextError(); err != nil {
			return err
		}
	}
	output.WriteString("</")
	output.WriteString(name)
	output.WriteByte('>')
	return nil
}

func lexicalXMLName(name xml.Name, namespaces map[string]string, attribute bool, owner ...*parseState) (string, error) {
	if err := parseOwnerError(owner); err != nil {
		return "", err
	}
	if name.Space == "" {
		return name.Local, nil
	}
	if name.Space == "http://www.w3.org/XML/1998/namespace" {
		return "xml:" + name.Local, nil
	}
	prefixes := make([]string, 0)
	for prefix, namespace := range namespaces {
		if err := parseOwnerError(owner); err != nil {
			return "", err
		}
		if namespace == name.Space && (!attribute || prefix != "") {
			prefixes = append(prefixes, prefix)
		}
	}
	if len(prefixes) == 0 {
		return "", fmt.Errorf("wsdl: namespace %q has no in-scope prefix", name.Space)
	}
	sort.Strings(prefixes)
	if prefixes[0] == "" {
		return name.Local, nil
	}
	return prefixes[0] + ":" + name.Local, nil
}

func readXMLNode(
	decoder *xml.Decoder,
	start xml.StartElement,
	state *parseState,
	depth int,
) (*xmlNode, error) {
	if err := state.contextError(); err != nil {
		return nil, err
	}
	state.elements++
	state.attributes += len(start.Attr)
	if depth > state.options.MaxDepth {
		return nil, fmt.Errorf("%w: element depth exceeds %d", ErrLimitExceeded, state.options.MaxDepth)
	}
	if state.elements > state.options.MaxElements {
		return nil, fmt.Errorf("%w: element count exceeds %d", ErrLimitExceeded, state.options.MaxElements)
	}
	if state.attributes > state.options.MaxAttributes {
		return nil, fmt.Errorf("%w: attribute count exceeds %d", ErrLimitExceeded, state.options.MaxAttributes)
	}
	line, column := decoder.InputPos()
	node := &xmlNode{
		owner:      state,
		name:       start.Name,
		attributes: append([]xml.Attr(nil), start.Attr...),
		namespaces: make(map[string]string),
		location: Location{
			SystemID: state.options.SystemID,
			Line:     line,
			Column:   column,
			Offset:   decoder.InputOffset(),
		},
	}
	for _, attribute := range start.Attr {
		if err := state.contextError(); err != nil {
			return nil, err
		}
		if attribute.Name.Space == "xmlns" {
			node.namespaces[attribute.Name.Local] = attribute.Value
			continue
		}
		if attribute.Name.Space == "" && attribute.Name.Local == "xmlns" {
			node.namespaces[""] = attribute.Value
		}
	}

	for {
		if err := state.contextError(); err != nil {
			return nil, err
		}
		token, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("wsdl: parse {%s}%s: %w", start.Name.Space, start.Name.Local, err)
		}
		if err := state.contextError(); err != nil {
			return nil, err
		}
		switch value := token.(type) {
		case xml.Directive:
			return nil, ErrDTDForbidden
		case xml.StartElement:
			child, err := readXMLNode(decoder, value, state, depth+1)
			if err != nil {
				return nil, err
			}
			inheritNamespaces(child, node.namespaces, state)
			if err := state.contextError(); err != nil {
				return nil, err
			}
			node.children = append(node.children, child)
			node.content = append(node.content, xmlContent{child: child})
		case xml.CharData:
			state.textBytes += int64(len(value))
			if state.textBytes > state.options.MaxTextBytes {
				return nil, fmt.Errorf("%w: text bytes exceed %d", ErrLimitExceeded, state.options.MaxTextBytes)
			}
			text := append([]byte(nil), value...)
			node.text.Write(text)
			node.content = append(node.content, xmlContent{text: text})
		case xml.EndElement:
			if value.Name == start.Name {
				return node, nil
			}
		}
	}
}

func inheritNamespaces(node *xmlNode, inherited map[string]string, owner ...*parseState) {
	if parseOwnerError(owner) != nil {
		return
	}
	for prefix, namespace := range inherited {
		if parseOwnerError(owner) != nil {
			return
		}
		if _, exists := node.namespaces[prefix]; !exists {
			node.namespaces[prefix] = namespace
		}
	}
	for _, child := range node.children {
		if parseOwnerError(owner) != nil {
			return
		}
		inheritNamespaces(child, node.namespaces, owner...)
	}
}

func (n *xmlNode) attribute(local string) string {
	if n.contextError() != nil {
		return ""
	}
	for _, attribute := range n.attributes {
		if n.contextError() != nil {
			return ""
		}
		if attribute.Name.Space == "" && attribute.Name.Local == local {
			return attribute.Value
		}
	}
	return ""
}

// Sentinel-only accessors cannot distinguish a missing attribute from stopped
// traversal. Validators must admit the captured value through its actual owner
// before classifying it or inserting it into a symbol table.
func (n *xmlNode) checkedAttribute(local string) (string, error) {
	value := n.attribute(local)
	if err := n.contextError(); err != nil {
		return "", err
	}
	return value, nil
}

func (n *xmlNode) hasAttribute(local string) bool {
	if n.contextError() != nil {
		return false
	}
	for _, attribute := range n.attributes {
		if n.contextError() != nil {
			return false
		}
		if attribute.Name.Space == "" && attribute.Name.Local == local {
			return true
		}
	}
	return false
}

func (n *xmlNode) namespacedAttribute(namespace, local string) (string, bool) {
	if n.contextError() != nil {
		return "", false
	}
	for _, attribute := range n.attributes {
		if n.contextError() != nil {
			return "", false
		}
		if attribute.Name == (xml.Name{Space: namespace, Local: local}) {
			return attribute.Value, true
		}
	}
	return "", false
}

func (n *xmlNode) qnameAttribute(local string) (converted QName, conversionErr error) {
	if err := n.contextError(); err != nil {
		return QName{}, err
	}
	defer finishConversion(n, &converted, &conversionErr)
	if !n.hasAttribute(local) {
		return QName{}, nil
	}
	return n.parseQName(n.attribute(local))
}

func (n *xmlNode) parseQName(value string) (converted QName, conversionErr error) {
	if err := n.contextError(); err != nil {
		return QName{}, err
	}
	defer finishConversion(n, &converted, &conversionErr)
	lexical := value
	if lexical == "" {
		return QName{}, fmt.Errorf("wsdl: invalid QName %q", lexical)
	}
	if strings.TrimSpace(lexical) != lexical {
		return QName{}, fmt.Errorf("wsdl: invalid QName %q", lexical)
	}
	parts := strings.Split(lexical, ":")
	if len(parts) > 2 || parts[0] == "" || (len(parts) == 2 && parts[1] == "") {
		return QName{}, fmt.Errorf("wsdl: invalid QName %q", lexical)
	}
	prefix := ""
	name := parts[0]
	if len(parts) == 2 {
		prefix, name = parts[0], parts[1]
	}
	if datatype.ValidateBuiltInLexical("NCName", name) != nil ||
		(prefix != "" && datatype.ValidateBuiltInLexical("NCName", prefix) != nil) {
		return QName{}, fmt.Errorf("wsdl: invalid QName %q", lexical)
	}
	namespace, exists := n.namespaces[prefix]
	if prefix != "" && !exists {
		return QName{}, fmt.Errorf("wsdl: QName %q uses undeclared prefix %q", lexical, prefix)
	}
	return QName{Namespace: namespace, Local: name}, nil
}

func (n *xmlNode) qnamesAttribute(local string) (converted []QName, conversionErr error) {
	if err := n.contextError(); err != nil {
		return nil, err
	}
	defer finishConversion(n, &converted, &conversionErr)
	values := splitSpaceSeparated(n.attribute(local))
	if err := n.contextError(); err != nil {
		return nil, err
	}
	result := make([]QName, 0, len(values))
	for _, value := range values {
		if err := n.contextError(); err != nil {
			return nil, err
		}
		name, err := n.parseQName(value)
		if err != nil {
			return nil, err
		}
		if err := n.contextError(); err != nil {
			return nil, err
		}
		result = append(result, name)
	}
	return result, nil
}

func (n *xmlNode) documentation() (converted *Documentation) {
	if n.contextError() != nil {
		return nil
	}
	defer finishConversionValue(n, &converted)
	if n.name.Local != "documentation" ||
		(n.name.Space != NamespaceWSDL11 && n.name.Space != NamespaceWSDL20) {
		return nil
	}
	language := ""
	for _, attribute := range n.attributes {
		if n.contextError() != nil {
			return nil
		}
		if attribute.Name == (xml.Name{Space: "http://www.w3.org/XML/1998/namespace", Local: "lang"}) {
			language = attribute.Value
		}
	}
	return &Documentation{
		Language: language,
		Content:  strings.TrimSpace(n.text.String()),
		Location: n.location,
	}
}

func splitSpaceSeparated(value string) []string {
	return strings.Fields(value)
}

func soapVersion(namespace string) Version {
	if namespace == NamespaceSOAP12Binding {
		return Version12
	}
	return Version11
}

// Version12 identifies SOAP 1.2 where a model carries a SOAP version.
const Version12 Version = "1.2"
