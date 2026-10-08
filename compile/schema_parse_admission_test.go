package compile_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	wsdlcompile "github.com/faustbrian/go-wsdl/v3/compile"
	"github.com/faustbrian/go-wsdl/v3/resolve"
	xsd "github.com/faustbrian/go-xsd/v2"
	xsdcompile "github.com/faustbrian/go-xsd/v2/compile"
	xsdresolve "github.com/faustbrian/go-xsd/v2/resolve"
)

func schemaAdmissionWSDL(version, namespace, reference string, inline bool) []byte {
	root, language := "definitions", "http://schemas.xmlsoap.org/wsdl/"
	if version == "2.0" {
		root, language = "description", "http://www.w3.org/ns/wsdl"
	}
	source := `<` + root + ` xmlns="` + language + `" targetNamespace="` + namespace + `">`
	if reference != "" {
		source += `<import namespace="urn:child" location="` + reference + `"/>`
	}
	if inline {
		source += `<types><xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema"` +
			` targetNamespace="` + namespace + `"><xs:element name="Request" type="xs:string"/>` +
			`</xs:schema></types>`
	}
	return []byte(source + `</` + root + `>`)
}

func TestCompilerRetainsFiniteInlineResolverAdmission(t *testing.T) {
	for _, count := range []int{256, 257} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			resources := make(map[string][]byte)
			var root []byte
			for start, document := 0, 0; start < count; document++ {
				var source strings.Builder
				fmt.Fprintf(&source, `<description xmlns="http://www.w3.org/ns/wsdl" targetNamespace="urn:wsdl:%d">`, document)
				if document == 0 {
					for child := 1; child*64 < count; child++ {
						fmt.Fprintf(&source, `<import namespace="urn:wsdl:%d" location="child-%d.wsdl"/>`, child, child)
					}
				}
				source.WriteString(`<types>`)
				for end := min(start+64, count); start < end; start++ {
					fmt.Fprintf(&source, `<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema" targetNamespace="urn:schema:%d"/>`, start)
				}
				source.WriteString(`</types></description>`)
				if document == 0 {
					root = []byte(source.String())
				} else {
					resources[fmt.Sprintf("https://example.test/child-%d.wsdl", document)] = []byte(source.String())
				}
			}
			resolver, err := resolve.NewMemory(resources)
			if err != nil {
				t.Fatal(err)
			}
			compiler, err := wsdlcompile.New(wsdlcompile.Options{
				Resolver: resolver, SchemaLimits: xsdcompile.Limits{MaxSchemas: 300},
			})
			if err != nil {
				t.Fatal(err)
			}
			set, err := compiler.Compile(context.Background(), wsdlcompile.Source{
				URI: "https://example.test/root.wsdl", Content: root,
			})
			if count == 257 {
				if set != nil || !errors.Is(err, xsdresolve.ErrLimitExceeded) {
					t.Fatalf("inline constructor bound: set present=%v, error=%v", set != nil, err)
				}
			} else if set == nil || err != nil || set.Schemas() == nil {
				t.Fatal("exact inline constructor allowance failed:", err)
			}
		})
	}
}

func TestCompilerAppliesSchemaAdmissionBeforeFollowingWSDLImports(t *testing.T) {
	for _, version := range []string{"1.1", "2.0"} {
		for _, imported := range []bool{false, true} {
			for name, limits := range map[string]xsdcompile.Limits{
				"namespace": {MaxParseNamespaceEntries: 1},
				"model":     {MaxParseModelBytes: 1},
			} {
				t.Run(version+"/"+name+map[bool]string{false: "/root", true: "/import"}[imported], func(t *testing.T) {
					root := schemaAdmissionWSDL(version, "urn:api", "missing.wsdl", true)
					resources := map[string][]byte{}
					if imported {
						root = schemaAdmissionWSDL(version, "urn:api", "child.wsdl", false)
						resources["https://example.test/child.wsdl"] =
							schemaAdmissionWSDL(version, "urn:child", "missing.wsdl", true)
					}
					resolver, err := resolve.NewMemory(resources)
					if err != nil {
						t.Fatal(err)
					}
					compiler, err := wsdlcompile.New(wsdlcompile.Options{Resolver: resolver, SchemaLimits: limits})
					if err != nil {
						t.Fatal(err)
					}
					set, err := compiler.Compile(context.Background(), wsdlcompile.Source{
						URI: "https://example.test/root.wsdl", Content: root,
					})
					// A late-only check instead reaches the missing WSDL import. The
					// typed XSD refusal proves admission during the first inline parse.
					if set != nil || !errors.Is(err, xsd.ErrLimitExceeded) {
						t.Fatalf("first schema parse: set present=%v, missing WSDL=%v, error=%v",
							set != nil, errors.Is(err, resolve.ErrNotFound), err)
					}
					set, err = compiler.Compile(context.Background(), wsdlcompile.Source{
						URI:     "https://example.test/healthy.wsdl",
						Content: schemaAdmissionWSDL(version, "urn:api", "", false),
					})
					if set == nil || err != nil {
						t.Fatal("schema refusal poisoned reusable compiler:", err)
					}
				})
			}
		}
	}
}
