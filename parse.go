package wsdl

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/faustbrian/go-wire/v3/xmlwire"
	"github.com/faustbrian/go-wsdl/v3/internal/errorprivacy"
)

const (
	namespace11            = "http://schemas.xmlsoap.org/wsdl/"
	namespace20            = "http://www.w3.org/ns/wsdl"
	defaultMaxDocumentSize = 8 << 20
	defaultMaxDepth        = 256
	defaultMaxElements     = 100000
	defaultMaxAttributes   = 1000000
	defaultMaxTextBytes    = 8 << 20
	defaultMaxSchemas      = 64
	defaultMaxImports      = 4096
	defaultMaxOperations   = 10000
	defaultMaxBindings     = 10000
	defaultMaxEndpoints    = 10000
	defaultMaxExtensions   = 10000
)

var decoderToken = func(decoder *xml.Decoder) (xml.Token, error) {
	return decoder.Token()
}

// ParseOptions controls bounded parsing of one caller-supplied document.
type ParseOptions struct {
	// MaxDocumentBytes bounds supplied bytes and each charset-conversion pass; zero selects 8
	// MiB.
	MaxDocumentBytes int64
	// MaxDepth bounds XML element depth; zero selects 256.
	MaxDepth int
	// MaxElements bounds XML element count; zero selects 100,000.
	MaxElements int
	// MaxAttributes bounds cumulative XML attributes; zero selects 1,000,000.
	MaxAttributes int
	// MaxTextBytes bounds decoded character-data bytes; zero selects 8 MiB.
	MaxTextBytes int64
	// MaxSchemas bounds inline schemas per document; zero selects 64.
	MaxSchemas int
	// MaxImports bounds import/include elements per document; zero selects 4,096.
	MaxImports int
	// MaxOperations bounds operation elements per document; zero selects 10,000.
	MaxOperations int
	// MaxBindings bounds binding elements per document; zero selects 10,000.
	MaxBindings int
	// MaxEndpoints bounds port/endpoint elements per document; zero selects 10,000.
	MaxEndpoints int
	// MaxExtensions bounds foreign elements/attributes counted at core boundaries; zero
	// selects 10,000.
	MaxExtensions int
	// MaxSchemaNamespaceEntries and MaxSchemaModelBytes independently bound
	// each inline schema's owned namespace and model-string/copy work during
	// its first parse. Zero selects XSD's finite defaults (1,000,000 entries
	// and 64 MiB). Negative values are invalid. They are not heap limits and
	// do not replace the enclosing WSDL document, tree, or schema-count limits.
	MaxSchemaNamespaceEntries int
	// MaxSchemaModelBytes bounds each inline schema's model-string/copy work; zero selects
	// XSD's 64 MiB default.
	MaxSchemaModelBytes int64
	// SystemID is caller-supplied source identity used for URI bases/diagnostics, not
	// authorization.
	SystemID string
}

// Parse borrows source until return, decodes one owned mutable WSDL model and
// resolves only lexical URI references. It does not perform full semantic
// validation or resource loading; use Validate/compile for those phases.
// Context must be non-nil: valid XML can panic in the second-pass reader when
// nil is supplied. Cancellation is cooperative, including a completion check,
// not a deadline on every upstream XML/charset operation. Errors return no model
// and use a fixed category while retaining inspectable typed causes.
func Parse(ctx context.Context, source []byte, options ParseOptions) (*Document, error) {
	document, err := parse(ctx, source, options)
	if err == nil && ctx != nil {
		if stopped := ctx.Err(); stopped != nil {
			document, err = nil, stopped
		}
	}
	return document, errorprivacy.Wrap("wsdl: parse failed", err)
}

func parse(ctx context.Context, source []byte, options ParseOptions) (*Document, error) {
	limit := options.MaxDocumentBytes
	if limit < 0 || options.MaxDepth < 0 || options.MaxElements < 0 ||
		options.MaxAttributes < 0 || options.MaxTextBytes < 0 ||
		options.MaxSchemas < 0 || options.MaxImports < 0 ||
		options.MaxOperations < 0 || options.MaxBindings < 0 ||
		options.MaxEndpoints < 0 || options.MaxExtensions < 0 ||
		options.MaxSchemaNamespaceEntries < 0 || options.MaxSchemaModelBytes < 0 {
		return nil, errors.New("wsdl: parse limits must not be negative")
	}
	if limit == 0 {
		limit = defaultMaxDocumentSize
	}
	if options.MaxDepth == 0 {
		options.MaxDepth = defaultMaxDepth
	}
	if options.MaxElements == 0 {
		options.MaxElements = defaultMaxElements
	}
	if options.MaxAttributes == 0 {
		options.MaxAttributes = defaultMaxAttributes
	}
	if options.MaxTextBytes == 0 {
		options.MaxTextBytes = defaultMaxTextBytes
	}
	if options.MaxSchemas == 0 {
		options.MaxSchemas = defaultMaxSchemas
	}
	if options.MaxImports == 0 {
		options.MaxImports = defaultMaxImports
	}
	if options.MaxOperations == 0 {
		options.MaxOperations = defaultMaxOperations
	}
	if options.MaxBindings == 0 {
		options.MaxBindings = defaultMaxBindings
	}
	if options.MaxEndpoints == 0 {
		options.MaxEndpoints = defaultMaxEndpoints
	}
	if options.MaxExtensions == 0 {
		options.MaxExtensions = defaultMaxExtensions
	}
	if int64(len(source)) > limit {
		return nil, fmt.Errorf("%w: document bytes exceed %d", ErrLimitExceeded, limit)
	}

	if ctx != nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
	}
	charsetReader := func(label string, input io.Reader) (io.Reader, error) {
		return xmlwire.CharsetReaderWithLimit(label, input, limit)
	}
	// Both XML passes share the document-sized charset allowance; a declared
	// vendor encoding must not evade the raw-byte boundary through expansion.
	if _, err := xmlwire.Root(source, xmlwire.DecodeOptions{
		MaxBytes: limit, MaxDepth: options.MaxDepth, CharsetReader: charsetReader,
	}); err != nil {
		if errors.Is(err, xmlwire.ErrPayloadTooLarge) ||
			errors.Is(err, xmlwire.ErrNestingTooDeep) {
			return nil, fmt.Errorf("%w: %v", ErrLimitExceeded, err)
		}
		return nil, fmt.Errorf("wsdl: validate XML: %w", err)
	}

	state := parseState{ctx: ctx, options: options}
	decoder := xml.NewDecoder(&contextReader{ctx: ctx, reader: bytes.NewReader(source)})
	decoder.Strict = true
	decoder.Entity = map[string]string{}
	decoder.CharsetReader = charsetReader
	for {
		if err := state.contextError(); err != nil {
			return nil, err
		}
		token, err := decoderToken(decoder)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, errors.New("wsdl: document has no root element")
			}
			return nil, fmt.Errorf("wsdl: parse XML: %w", err)
		}
		if err := state.contextError(); err != nil {
			return nil, err
		}
		switch value := token.(type) {
		case nil:
			return nil, errors.New("wsdl: parser returned a nil token")
		case xml.Directive:
			return nil, ErrDTDForbidden
		case xml.StartElement:
			switch value.Name {
			case xml.Name{Space: namespace11, Local: "definitions"}:
				return parseDefinitions11(decoder, value, &state)
			case xml.Name{Space: namespace20, Local: "description"}:
				return parseDescription20(decoder, value, &state)
			default:
				return nil, fmt.Errorf(
					"wsdl: unsupported root element {%s}%s",
					value.Name.Space,
					value.Name.Local,
				)
			}
		}
	}
}

func parseDefinitions11(
	decoder *xml.Decoder,
	start xml.StartElement,
	state *parseState,
) (*Document, error) {
	root, err := readXMLNode(decoder, start, state, 1)
	if err != nil {
		return nil, err
	}
	if err := validateCoreNCNames(root, NamespaceWSDL11, state); err != nil {
		return nil, err
	}
	if err := enforceComponentLimits(root, NamespaceWSDL11, state.options, state); err != nil {
		return nil, err
	}
	if err := assignBaseURIs(root, state.options.SystemID, state); err != nil {
		return nil, err
	}
	if err := state.contextError(); err != nil {
		return nil, err
	}
	definitions, err := decodeDefinitions11(root, state)
	if err != nil {
		return nil, err
	}
	if err := state.contextError(); err != nil {
		return nil, err
	}
	return &Document{
		version: Version11, definitions11: &definitions,
	}, nil
}

func parseDescription20(
	decoder *xml.Decoder,
	start xml.StartElement,
	state *parseState,
) (*Document, error) {
	root, err := readXMLNode(decoder, start, state, 1)
	if err != nil {
		return nil, err
	}
	if err := validateCoreNCNames(root, NamespaceWSDL20, state); err != nil {
		return nil, err
	}
	if err := enforceComponentLimits(root, NamespaceWSDL20, state.options, state); err != nil {
		return nil, err
	}
	if err := assignBaseURIs(root, state.options.SystemID, state); err != nil {
		return nil, err
	}
	if err := state.contextError(); err != nil {
		return nil, err
	}
	description, err := decodeDescription20(state.ctx, root, state.options)
	if err != nil {
		return nil, err
	}
	if err := state.contextError(); err != nil {
		return nil, err
	}
	return &Document{
		version:       Version20,
		description20: &description,
	}, nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r *contextReader) Read(buffer []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(buffer)
}
