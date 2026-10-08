package wsdl_test

import (
	"context"
	"errors"
	"testing"

	wsdl "github.com/faustbrian/go-wsdl/v3"
	wsdlcompile "github.com/faustbrian/go-wsdl/v3/compile"
	xsd "github.com/faustbrian/go-xsd/v2"
	xsdcompile "github.com/faustbrian/go-xsd/v2/compile"
	xsdresolve "github.com/faustbrian/go-xsd/v2/resolve"
)

func TestPublicXSDComposition(t *testing.T) {
	for _, version := range []string{"1.1", "2.0"} {
		t.Run(version, func(t *testing.T) {
			root, closeRoot := "definitions", "definitions"
			namespace := "http://schemas.xmlsoap.org/wsdl/"
			if version == "2.0" {
				root, closeRoot = "description", "description"
				namespace = "http://www.w3.org/ns/wsdl"
			}
			source := []byte(`<` + root + ` xmlns="` + namespace + `" targetNamespace="urn:api">` +
				`<types><xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema" targetNamespace="urn:api">` +
				`<xs:element name="Request" type="xs:string"/></xs:schema></types></` + closeRoot + `>`)
			document, err := wsdl.Parse(context.Background(), source, wsdl.ParseOptions{})
			if err != nil {
				t.Fatal(err)
			}
			var schemas []*xsd.Document
			if value, ok := document.Definitions11(); ok {
				schemas = value.Types.Schemas
			} else if value, ok := document.Description20(); ok {
				schemas = value.Types.Schemas
				var references []xsd.SchemaReference = value.Types.Imports
				if len(references) != 0 {
					t.Fatal("unexpected schema import")
				}
			}
			if len(schemas) != 1 || schemas[0].Elements[0].Name != "Request" {
				t.Fatal("inline schema composition lost the ordinary element")
			}
			resolver, err := xsdresolve.NewMemory(nil)
			if err != nil {
				t.Fatal(err)
			}
			compiler, err := wsdlcompile.New(wsdlcompile.Options{
				SchemaResolver: resolver, SchemaLimits: xsdcompile.Limits{},
			})
			if err != nil {
				t.Fatal(err)
			}
			set, err := compiler.Compile(context.Background(), wsdlcompile.Source{
				URI: "https://example.test/root.wsdl", Content: source,
			})
			if err != nil {
				t.Fatal(err)
			}
			var compiled *xsdcompile.Set = set.Schemas()
			if compiled == nil {
				t.Fatal("compiled inline schemas missing")
			}
			encoded, err := wsdl.Marshal(document, wsdl.MarshalOptions{})
			if err != nil {
				t.Fatal(err)
			}
			if _, err := wsdl.Parse(context.Background(), encoded, wsdl.ParseOptions{}); err != nil {
				t.Fatal("schema roundtrip failed:", err)
			}
		})
	}
}

func TestInlineSchemaParseAdmission(t *testing.T) {
	for _, root := range []struct{ name, namespace string }{
		{"definitions", "http://schemas.xmlsoap.org/wsdl/"},
		{"description", "http://www.w3.org/ns/wsdl"},
	} {
		t.Run(root.name, func(t *testing.T) {
			source := []byte(`<` + root.name + ` xmlns="` + root.namespace + `" targetNamespace="urn:api">` +
				`<types><xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema" targetNamespace="urn:api">` +
				`<xs:element name="Request" type="xs:string"/></xs:schema></types></` + root.name + `>`)
			for name, options := range map[string]wsdl.ParseOptions{
				"namespace": {MaxSchemaNamespaceEntries: 1},
				"model":     {MaxSchemaModelBytes: 1},
			} {
				document, err := wsdl.Parse(context.Background(), source, options)
				if document != nil || !errors.Is(err, xsd.ErrLimitExceeded) {
					t.Fatalf("%s allowance: document present=%v, error=%v", name, document != nil, err)
				}
			}
			for _, options := range []wsdl.ParseOptions{
				{MaxSchemaNamespaceEntries: -1}, {MaxSchemaModelBytes: -1},
			} {
				if document, err := wsdl.Parse(context.Background(), source, options); document != nil || err == nil {
					t.Fatal("negative schema allowance accepted")
				}
			}
			for _, options := range []wsdl.ParseOptions{
				{}, {MaxSchemaNamespaceEntries: 2, MaxSchemaModelBytes: 1 << 16},
			} {
				if document, err := wsdl.Parse(context.Background(), source, options); document == nil || err != nil {
					t.Fatal("admitted ordinary schema failed:", err)
				}
			}
		})
	}
}
