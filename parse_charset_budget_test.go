package wsdl_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	wsdl "github.com/faustbrian/go-wsdl/v2"
)

func TestParseVendorEncodingUsesDocumentBudget(t *testing.T) {
	for _, version := range []struct {
		name      string
		root      string
		namespace string
	}{
		{"1.1", "definitions", "http://schemas.xmlsoap.org/wsdl/"},
		{"2.0", "description", "http://www.w3.org/ns/wsdl"},
	} {
		for _, size := range []struct {
			name string
			text int
		}{
			{"small", 16},
			{"above-wire-default", (1 << 20) + 1},
		} {
			t.Run(version.name+"/"+size.name, func(t *testing.T) {
				prefix := `<?xml version="1.0" encoding="windows-1252"?><` + version.root +
					` xmlns="` + version.namespace + `" targetNamespace="urn:vendor"><documentation>`
				source := []byte(prefix + strings.Repeat("a", size.text) + "\xe9" +
					`</documentation></` + version.root + `>`)
				want := strings.Repeat("a", size.text) + "é"
				for _, limit := range []struct {
					name string
					max  int64
				}{
					{"default", 0},
					{"explicit-exact", int64(len(source))},
				} {
					t.Run(limit.name, func(t *testing.T) {
						document, err := wsdl.Parse(context.Background(), source, wsdl.ParseOptions{MaxDocumentBytes: limit.max})
						if err != nil || document == nil {
							t.Fatalf("admitted vendor document = (nonNil=%t, err=%v)", document != nil, err)
						}
						var documentation *wsdl.Documentation
						if definitions, ok := document.Definitions11(); ok {
							documentation = definitions.Documentation
						} else if description, ok := document.Description20(); ok {
							documentation = description.Documentation
						}
						if documentation == nil || documentation.Content != want {
							t.Fatal("vendor documentation was not decoded exactly")
						}
					})
				}
				document, err := wsdl.Parse(context.Background(), source, wsdl.ParseOptions{MaxDocumentBytes: int64(len(source) - 1)})
				if document != nil || !errors.Is(err, wsdl.ErrLimitExceeded) {
					t.Fatalf("one-byte-short document budget = (nonNil=%t, err=%v), want nil/limit", document != nil, err)
				}
			})
		}
	}
}
