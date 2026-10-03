package wsdl

import (
	"context"
	"encoding/xml"
	"errors"
	"strings"
	"testing"
	"time"
)

type cancelingDocumentReader struct {
	reader *strings.Reader
	cancel context.CancelFunc
}

func TestParseCancellationPreservesPrivacyAndOrdinaryModels(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	expired, stop := context.WithDeadline(context.Background(), time.Unix(0, 0))
	defer stop()
	for _, source := range []string{
		`<definitions xmlns="http://schemas.xmlsoap.org/wsdl/" name="Ordinary" targetNamespace="urn:test"/>`,
		`<description xmlns="http://www.w3.org/ns/wsdl" targetNamespace="urn:test"/>`,
	} {
		for _, test := range []struct {
			ctx   context.Context
			cause error
		}{{ctx, context.Canceled}, {expired, context.DeadlineExceeded}} {
			document, err := Parse(test.ctx, []byte(source), ParseOptions{})
			if document != nil || !errors.Is(err, test.cause) || err.Error() != "wsdl: parse failed" {
				t.Fatalf("canceled Parse = %v, %v", document, err)
			}
		}
		document, err := Parse(context.Background(), []byte(source), ParseOptions{})
		if document == nil || err != nil {
			t.Fatalf("ordinary Parse = %v, %v", document, err)
		}
		if document.Version() == Version11 {
			model, ok := document.Definitions11()
			if !ok || model.Name != "Ordinary" || model.TargetNamespace != "urn:test" {
				t.Fatal("ordinary WSDL11 model changed")
			}
		} else {
			model, ok := document.Description20()
			if !ok || model.TargetNamespace != "urn:test" {
				t.Fatal("ordinary WSDL20 model changed")
			}
		}
	}
	if _, err := Parse(ctx, []byte(`<description xmlns="http://www.w3.org/ns/wsdl"/>`), ParseOptions{MaxDepth: -1}); err == nil || errors.Is(err, context.Canceled) {
		t.Fatal("early invalid-option precedence changed")
	}
}

func TestTreeHelpersHonorOwnerAndRetainContextFreeBehavior(t *testing.T) {
	node := &xmlNode{name: xml.Name{Space: NamespaceWSDL20, Local: "description"}, namespaces: map[string]string{}, children: []*xmlNode{
		{name: xml.Name{Space: NamespaceWSDL20, Local: "interface"}, attributes: []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "Ordinary"}}, namespaces: map[string]string{}},
	}}
	ctx, cancel := context.WithCancel(context.Background())
	state := &parseState{ctx: ctx}
	options := ParseOptions{MaxImports: 1, MaxOperations: 1, MaxBindings: 1, MaxEndpoints: 1, MaxExtensions: 1}
	for _, run := range []func() error{
		func() error { return validateCoreNCNames(node, NamespaceWSDL20, state) },
		func() error { return enforceComponentLimits(node, NamespaceWSDL20, options, state) },
		func() error { return assignBaseURIs(node, "urn:test", state) },
	} {
		if err := run(); err != nil {
			t.Fatal(err)
		}
	}
	cancel()
	for _, run := range []func() error{
		func() error { return validateCoreNCNames(node, NamespaceWSDL20, state) },
		func() error { return enforceComponentLimits(node, NamespaceWSDL20, options, state) },
		func() error { return assignBaseURIs(node, "urn:test", state) },
	} {
		if err := run(); !errors.Is(err, context.Canceled) {
			t.Fatalf("canceled tree owner = %v", err)
		}
	}
	if err := validateCoreNCNames(node, NamespaceWSDL20); err != nil {
		t.Fatal(err)
	}
	if err := enforceComponentLimits(node, NamespaceWSDL20, options); err != nil {
		t.Fatal(err)
	}
	if err := assignBaseURIs(node, "urn:ordinary"); err != nil {
		t.Fatal(err)
	}
}

func (r cancelingDocumentReader) Read(buffer []byte) (int, error) {
	n, err := r.reader.Read(buffer)
	if n > 0 && r.cancel != nil {
		r.cancel()
	}
	return n, err
}

func TestVersionParserHonorsCancellationBeforeBufferedTreeConversion(t *testing.T) {
	for _, version := range []struct {
		name   string
		source string
		parse  func(*xml.Decoder, xml.StartElement, *parseState) (*Document, error)
	}{
		{"WSDL11", `<definitions xmlns="http://schemas.xmlsoap.org/wsdl/" name="Ordinary" targetNamespace="urn:test"><message name="Input"/></definitions>`, parseDefinitions11},
		{"WSDL20", `<description xmlns="http://www.w3.org/ns/wsdl" targetNamespace="urn:test"><interface name="Ordinary"/></description>`, parseDescription20},
	} {
		for _, canceled := range []bool{false, true} {
			t.Run(version.name+map[bool]string{false: " ordinary", true: " canceled"}[canceled], func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				reader := cancelingDocumentReader{reader: strings.NewReader(version.source)}
				if canceled {
					reader.cancel = cancel
				}
				decoder := xml.NewDecoder(&contextReader{ctx: ctx, reader: reader})
				token, err := decoder.Token()
				if err != nil {
					t.Fatal(err)
				}
				start, ok := token.(xml.StartElement)
				if !ok {
					t.Fatalf("ordinary root token = %T", token)
				}
				state := &parseState{ctx: ctx, options: ParseOptions{
					MaxDepth: 4, MaxElements: 4, MaxAttributes: 10, MaxTextBytes: 100,
					MaxImports: 4, MaxOperations: 4, MaxBindings: 4, MaxEndpoints: 4, MaxExtensions: 4,
				}}
				document, err := version.parse(decoder, start, state)
				if canceled {
					if document != nil || !errors.Is(err, context.Canceled) {
						t.Fatalf("buffered canceled conversion = %v, %v; want nil Document, context.Canceled", document, err)
					}
				} else if document == nil || err != nil {
					t.Fatalf("ordinary conversion = %v, %v", document, err)
				}
			})
		}
	}
}
