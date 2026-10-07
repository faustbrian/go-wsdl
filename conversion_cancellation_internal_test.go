package wsdl

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"io"
	"reflect"
	"testing"
)

func ordinaryConversionRoot(namespace, rootName, childName string) *xmlNode {
	return &xmlNode{
		name: xml.Name{Space: namespace, Local: rootName},
		attributes: []xml.Attr{
			{Name: xml.Name{Local: "name"}, Value: "Ordinary"},
			{Name: xml.Name{Local: "targetNamespace"}, Value: "urn:test"},
		},
		namespaces: map[string]string{"": namespace},
		children: []*xmlNode{{
			name:       xml.Name{Space: namespace, Local: childName},
			attributes: []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "Input"}},
			namespaces: map[string]string{"": namespace},
		}},
	}
}

func TestNestedConversionStopsAfterOrdinaryExtensionSerialization(t *testing.T) {
	original := marshalNode
	t.Cleanup(func() { marshalNode = original })
	for _, version := range []Version{Version11, Version20} {
		for _, canceled := range []bool{false, true} {
			ctx, cancel := context.WithCancel(context.Background())
			marshalNode = func(node *xmlNode) ([]byte, error) {
				payload, err := original(node)
				if canceled && err == nil {
					cancel()
				}
				return payload, err
			}
			namespace, rootName, childName := NamespaceWSDL11, "definitions", "message"
			if version == Version20 {
				namespace, rootName, childName = NamespaceWSDL20, "description", "interface"
			}
			root := ordinaryConversionRoot(namespace, rootName, childName)
			root.children[0].children = []*xmlNode{{
				name:       xml.Name{Space: "urn:extension", Local: "extra"},
				namespaces: map[string]string{"e": "urn:extension"},
				content:    []xmlContent{{text: []byte("ordinary")}},
			}}
			var extensions []Extension
			var err error
			if version == Version11 {
				value, conversionErr := decodeDefinitions11(root, &parseState{ctx: ctx})
				err = conversionErr
				if canceled && !reflect.DeepEqual(value, Definitions11{}) {
					t.Fatal("canceled nested WSDL11 conversion published a partial model")
				}
				if !canceled {
					extensions = value.Messages[0].Extensions
				}
			} else {
				value, conversionErr := decodeDescription20(ctx, root, ParseOptions{})
				err = conversionErr
				if canceled && !reflect.DeepEqual(value, Description20{}) {
					t.Fatal("canceled nested WSDL20 conversion published a partial model")
				}
				if !canceled {
					extensions = value.Interfaces[0].Extensions
				}
			}
			cancel()
			if canceled {
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("nested conversion error = %v", err)
				}
			} else if err != nil || len(extensions) != 1 || string(extensions[0].XML) != `<e:extra xmlns:e="urn:extension">ordinary</e:extra>` {
				t.Fatalf("ordinary nested conversion changed: error %v, extensions %#v", err, extensions)
			}
		}
	}
}

func TestOwnedXMLSerializationStopsAfterNestedEscaping(t *testing.T) {
	original := escapeXMLText
	t.Cleanup(func() { escapeXMLText = original })
	for _, canceled := range []bool{false, true} {
		ctx, cancel := context.WithCancel(context.Background())
		escapeXMLText = func(output io.Writer, text []byte) error {
			err := original(output, text)
			if canceled && err == nil && bytes.Equal(text, []byte("ordinary")) {
				cancel()
			}
			return err
		}
		child := &xmlNode{name: xml.Name{Local: "child"}, content: []xmlContent{{text: []byte("ordinary")}}}
		root := &xmlNode{name: xml.Name{Local: "root"}, children: []*xmlNode{child}, content: []xmlContent{{child: child}}}
		if err := root.useOwner(&parseState{ctx: ctx}); err != nil {
			t.Fatal(err)
		}
		payload, err := marshalXMLNode(root)
		cancel()
		if canceled {
			if !errors.Is(err, context.Canceled) || payload != nil {
				t.Fatalf("canceled serializer published bytes: %q, error %v", payload, err)
			}
		} else if err != nil || string(payload) != "<root><child>ordinary</child></root>" {
			t.Fatalf("ordinary serialization changed: %q, error %v", payload, err)
		}
	}
}

func TestPublishedModelsDoNotRetainConversionContext(t *testing.T) {
	for _, source := range []string{
		`<definitions xmlns="http://schemas.xmlsoap.org/wsdl/" name="Ordinary" targetNamespace="urn:test"><message name="Input"/></definitions>`,
		`<description xmlns="http://www.w3.org/ns/wsdl" targetNamespace="urn:test"><interface name="Input"/></description>`,
	} {
		ctx, cancel := context.WithCancel(context.Background())
		document, err := Parse(ctx, []byte(source), ParseOptions{})
		cancel()
		if err != nil || document == nil {
			t.Fatalf("ordinary Parse = %v, %v", document, err)
		}
		if document.Version() == Version11 {
			value, ok := document.Definitions11()
			if !ok || len(value.Messages) != 1 || value.Messages[0].Name != "Input" {
				t.Fatal("published WSDL11 access changed after cancellation")
			}
		} else {
			value, ok := document.Description20()
			if !ok || len(value.Interfaces) != 1 || value.Interfaces[0].Name != "Input" {
				t.Fatal("published WSDL20 access changed after cancellation")
			}
		}
		payload, err := Marshal(document, MarshalOptions{})
		if err != nil || len(payload) == 0 {
			t.Fatalf("published model Marshal = %q, %v", payload, err)
		}
		if _, err := Parse(context.Background(), payload, ParseOptions{}); err != nil {
			t.Fatalf("published model round trip = %v", err)
		}
	}
}

func TestConversionSerializationRetainsTerminalErrorPrecedence(t *testing.T) {
	original := marshalNode
	t.Cleanup(func() { marshalNode = original })
	cause := errors.New("ordinary serialization collaborator unavailable")
	for _, source := range []string{
		`<definitions xmlns="http://schemas.xmlsoap.org/wsdl/" xmlns:e="urn:extension"><message name="Input"><e:extra/></message></definitions>`,
		`<description xmlns="http://www.w3.org/ns/wsdl" xmlns:e="urn:extension"><interface name="Input"><e:extra/></interface></description>`,
	} {
		ctx, cancel := context.WithCancel(context.Background())
		marshalNode = func(*xmlNode) ([]byte, error) {
			cancel()
			return nil, cause
		}
		document, err := Parse(ctx, []byte(source), ParseOptions{})
		cancel()
		if document != nil || !errors.Is(err, cause) || errors.Is(err, context.Canceled) || err.Error() != "wsdl: parse failed" {
			t.Fatalf("terminal serialization precedence/privacy = %v, %v", document, err)
		}
	}
}

func TestPublicParseStopsAfterNestedConversionSerialization(t *testing.T) {
	original := marshalNode
	t.Cleanup(func() { marshalNode = original })
	for _, source := range []string{
		`<definitions xmlns="http://schemas.xmlsoap.org/wsdl/" xmlns:e="urn:extension"><message name="Input"><e:extra/></message></definitions>`,
		`<description xmlns="http://www.w3.org/ns/wsdl" xmlns:e="urn:extension"><interface name="Input"><e:extra/></interface></description>`,
	} {
		ctx, cancel := context.WithCancel(context.Background())
		marshalNode = func(node *xmlNode) ([]byte, error) {
			payload, err := original(node)
			if err == nil {
				cancel()
			}
			return payload, err
		}
		document, err := Parse(ctx, []byte(source), ParseOptions{})
		cancel()
		if document != nil || !errors.Is(err, context.Canceled) || err.Error() != "wsdl: parse failed" {
			t.Fatalf("public nested cancellation/privacy = %v, %v", document, err)
		}
	}
}

func TestSharedConversionHelpersHonorOwnedCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	node := &xmlNode{
		owner:      &parseState{ctx: ctx},
		name:       xml.Name{Space: NamespaceWSDL20, Local: "documentation"},
		attributes: []xml.Attr{{Name: xml.Name{Local: "names"}, Value: "One Two"}},
	}
	names, err := node.qnamesAttribute("names")
	if err != nil || len(names) != 2 || names[0].Local != "One" || names[1].Local != "Two" || node.documentation() == nil {
		t.Fatalf("ordinary shared QName/documentation conversion = %#v, %v", names, err)
	}
	cancel()
	if names, err := node.qnamesAttribute("names"); names != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("owned QName cancellation = %#v, %v", names, err)
	}
	if node.documentation() != nil || node.attribute("names") != "" {
		t.Fatal("canceled documentation/attribute conversion published values")
	}
	if value, err := decodeHTTPHeaders20(node); value != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("owned HTTP header cancellation = %#v, %v", value, err)
	}
	if value, err := decodeSOAPModules20(node); value != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("owned SOAP module cancellation = %#v, %v", value, err)
	}
	if value := decodeMIMEMultipart11(node); !reflect.DeepEqual(value, MIMEMultipart11{}) {
		t.Fatalf("owned MIME conversion published values: %#v", value)
	}
	// Constructed nodes without an invocation owner keep their historical use.
	node.owner = nil
	if names, err := node.qnamesAttribute("names"); err != nil || len(names) != 2 || node.documentation() == nil {
		t.Fatalf("context-free shared conversion = %#v, %v", names, err)
	}
}

func TestVersionConversionOwnersHonorRealCanceledContext(t *testing.T) {
	for _, canceled := range []bool{false, true} {
		ctx, cancel := context.WithCancel(context.Background())
		if canceled {
			cancel()
		}
		state := &parseState{ctx: ctx, options: ParseOptions{MaxSchemas: 4}}
		definitions, err11 := decodeDefinitions11(ordinaryConversionRoot(NamespaceWSDL11, "definitions", "message"), state)
		description, err20 := decodeDescription20(ctx, ordinaryConversionRoot(NamespaceWSDL20, "description", "interface"), state.options)
		cancel()
		if canceled {
			if !errors.Is(err11, context.Canceled) || !reflect.DeepEqual(definitions, Definitions11{}) {
				t.Errorf("canceled WSDL11 conversion published model: error%v", err11)
			}
			if !errors.Is(err20, context.Canceled) || !reflect.DeepEqual(description, Description20{}) {
				t.Errorf("canceled WSDL20 conversion published model: error%v", err20)
			}
		} else {
			if err11 != nil || definitions.Name != "Ordinary" || definitions.TargetNamespace != "urn:test" || len(definitions.Messages) != 1 || definitions.Messages[0].Name != "Input" {
				t.Errorf("ordinary WSDL11 conversion changed: %v", err11)
			}
			if err20 != nil || description.TargetNamespace != "urn:test" || len(description.Interfaces) != 1 || description.Interfaces[0].Name != "Input" {
				t.Errorf("ordinary WSDL20 conversion changed: %v", err20)
			}
		}
	}
}

func TestSOAPActionRequiredHandoffHonorsRealCancellation(t *testing.T) {
	for _, lexical := range []string{"true", "false", "TRUE"} {
		for _, canceled := range []bool{false, true} {
			ctx, cancel := context.WithCancel(context.Background())
			node := &xmlNode{
				owner:      &parseState{ctx: ctx},
				attributes: []xml.Attr{{Name: xml.Name{Local: "soapActionRequired"}, Value: lexical}},
			}
			if !node.hasAttribute("soapActionRequired") {
				t.Fatal("ordinary attribute presence was not established")
			}
			// This is the actual private owner's hasAttribute-to-lexical
			// handoff, not a scheduler-controlled public Parse interruption.
			if canceled {
				cancel()
			}
			required, err := decodeSOAPActionRequired11(node)
			cancel()
			if canceled {
				if required || !errors.Is(err, context.Canceled) {
					t.Errorf("canceled %q handoff = %v, %v; want false, context.Canceled", lexical, required, err)
				}
			} else if lexical == "TRUE" {
				if required || err == nil || errors.Is(err, context.Canceled) || err.Error() != `wsdl: invalid SOAP action required value "TRUE"` {
					t.Errorf("established invalid lexical behavior = %v, %v", required, err)
				}
			} else if err != nil || required != (lexical == "true") {
				t.Errorf("ordinary %q handoff = %v, %v", lexical, required, err)
			}
		}
	}
}

func TestCheckedAttributeRetainsOwnerCauseBeforeSymbolClassification(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	node := &xmlNode{
		owner:      &parseState{ctx: ctx},
		attributes: []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "Ordinary"}},
	}
	name, err := node.checkedAttribute("name")
	if err != nil || name != "Ordinary" {
		t.Fatalf("ordinary symbol-name admission = %q, %v", name, err)
	}
	symbols := map[string]struct{}{name: {}}
	if err := registerSymbol(symbols, "ordinary component", name); !errors.Is(err, ErrDuplicateSymbol) {
		t.Fatalf("genuine duplicate classification = %v", err)
	}
	cancel()
	if name, err := node.checkedAttribute("name"); name != "" || !errors.Is(err, context.Canceled) {
		t.Fatalf("stopped symbol-name admission = %q, %v", name, err)
	}
	if len(symbols) != 1 {
		t.Fatal("stopped attribute admission changed the existing symbol table")
	}
	node.owner = nil
	if name, err := node.checkedAttribute("name"); err != nil || name != "Ordinary" {
		t.Fatalf("context-free symbol-name admission = %q, %v", name, err)
	}
}
