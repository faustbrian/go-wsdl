package wsdl

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"io"
	"testing"
)

func TestOwnedXMLSerializationStopsBeforeContinuingAfterEscape(t *testing.T) {
	original := escapeXMLText
	t.Cleanup(func() { escapeXMLText = original })
	for _, test := range []struct {
		name    string
		node    *xmlNode
		escaped string
		stopped string
		whole   string
	}{
		{
			name: "namespace",
			node: &xmlNode{
				name:       xml.Name{Space: "urn:ordinary", Local: "root"},
				namespaces: map[string]string{"": "urn:ordinary"},
			},
			escaped: "urn:ordinary",
			stopped: `<root xmlns="urn:ordinary`,
			whole:   `<root xmlns="urn:ordinary"></root>`,
		},
		{
			name: "attribute",
			node: &xmlNode{
				name:       xml.Name{Local: "root"},
				attributes: []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "ordinary"}},
			},
			escaped: "ordinary",
			stopped: `<root name="ordinary`,
			whole:   `<root name="ordinary"></root>`,
		},
		{
			name: "text",
			node: &xmlNode{
				name:    xml.Name{Local: "root"},
				content: []xmlContent{{text: []byte("ordinary")}},
			},
			escaped: "ordinary",
			stopped: `<root>ordinary`,
			whole:   `<root>ordinary</root>`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, mode := range []string{"ordinary", "canceled", "failure"} {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				test.node.owner = &parseState{ctx: ctx}
				escapeXMLText = func(output io.Writer, text []byte) error {
					if mode == "failure" {
						err := original(&boundedFailureWriter{}, text)
						cancel()
						return err
					}
					err := original(output, text)
					if mode == "canceled" && err == nil && string(text) == test.escaped {
						cancel()
					}
					return err
				}
				var output bytes.Buffer
				err := writeXMLNode(&output, test.node, true)
				if mode == "failure" {
					if err != errInjectedWrite || errors.Is(err, context.Canceled) {
						t.Fatalf("escaping error lost precedence: %v; want exact writer error", err)
					}
				} else if mode == "canceled" {
					if !errors.Is(err, context.Canceled) || output.String() != test.stopped {
						t.Fatalf("canceled output = %q, %v; want %q and context.Canceled", output.String(), err, test.stopped)
					}
				} else if err != nil || output.String() != test.whole {
					t.Fatalf("ordinary output = %q, %v; want %q", output.String(), err, test.whole)
				}
			}
		})
	}
}
