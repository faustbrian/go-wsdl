package wsdl

import (
	"context"
	"encoding/xml"
	"errors"
	"reflect"
	"testing"
)

func TestConversionOwnerAdmissionImport11(t *testing.T) {
	for _, canceled := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "canceled"}[canceled], func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			node := &xmlNode{
				owner: &parseState{ctx: ctx},
				name:  xml.Name{Space: NamespaceWSDL11, Local: "import"},
				attributes: []xml.Attr{
					{Name: xml.Name{Local: "namespace"}, Value: "urn:test"},
					{Name: xml.Name{Local: "location"}, Value: "ordinary.wsdl"},
				},
				baseURI: "https://example.test/wsdl/root.wsdl",
			}
			if canceled {
				cancel()
			}
			got, err := decodeImport11(node)
			if canceled {
				if !errors.Is(err, context.Canceled) || !reflect.DeepEqual(got, Import11{}) {
					t.Fatalf("canceled import = %#v, %v; want zero and context.Canceled", got, err)
				}
				return
			}
			want := Import11{Extensibility: Extensibility{ExtensionAttributes: []ExtensionAttribute{}}, Namespace: "urn:test", Location: "ordinary.wsdl", URI: "https://example.test/wsdl/ordinary.wsdl"}
			if err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("ordinary import = %#v, %v; want %#v", got, err, want)
			}
		})
	}
}

func TestConversionOwnerAdmissionMessage11(t *testing.T) {
	for _, canceled := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "canceled"}[canceled], func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			node := &xmlNode{
				name:       xml.Name{Space: NamespaceWSDL11, Local: "message"},
				attributes: []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "Input"}},
				children: []*xmlNode{{
					name: xml.Name{Space: NamespaceWSDL11, Local: "part"},
					attributes: []xml.Attr{
						{Name: xml.Name{Local: "name"}, Value: "value"},
						{Name: xml.Name{Local: "type"}, Value: "xs:string"},
					},
					namespaces: map[string]string{"xs": NamespaceXMLSchema},
				}},
			}
			if err := node.useOwner(&parseState{ctx: ctx}); err != nil {
				t.Fatal(err)
			}
			if canceled {
				cancel()
			}
			got, err := decodeMessage11(node)
			if canceled {
				if !errors.Is(err, context.Canceled) || !reflect.DeepEqual(got, Message11{}) {
					t.Fatalf("canceled message = %#v, %v; want zero and context.Canceled", got, err)
				}
				return
			}
			want := Message11{
				Extensibility: Extensibility{ExtensionAttributes: []ExtensionAttribute{}},
				Name:          "Input",
				Parts: []Part11{{
					Extensibility: Extensibility{ExtensionAttributes: []ExtensionAttribute{}},
					Name:          "value", Type: QName{Namespace: NamespaceXMLSchema, Local: "string"},
				}},
			}
			if err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("ordinary message = %#v, %v; want %#v", got, err, want)
			}
		})
	}
}

func TestConversionOwnerAdmissionInterface20(t *testing.T) {
	for _, canceled := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "canceled"}[canceled], func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			node := &xmlNode{
				owner: &parseState{ctx: ctx},
				name:  xml.Name{Space: NamespaceWSDL20, Local: "interface"},
				attributes: []xml.Attr{
					{Name: xml.Name{Local: "name"}, Value: "Input"},
					{Name: xml.Name{Local: "extends"}, Value: "tns:Parent"},
				},
				namespaces: map[string]string{"tns": "urn:test"},
			}
			if canceled {
				cancel()
			}
			got, err := decodeInterface20(node)
			if canceled {
				if !errors.Is(err, context.Canceled) || !reflect.DeepEqual(got, Interface20{}) {
					t.Fatalf("canceled interface = %#v, %v; want zero and context.Canceled", got, err)
				}
				return
			}
			want := Interface20{
				Extensibility: Extensibility{ExtensionAttributes: []ExtensionAttribute{}},
				Name:          "Input", Extends: []QName{{Namespace: "urn:test", Local: "Parent"}},
				StyleDefault: []string{},
			}
			if err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("ordinary interface = %#v, %v; want %#v", got, err, want)
			}
		})
	}
}

func TestFinishConversionValuePublication(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	node := &xmlNode{owner: &parseState{ctx: ctx}}
	want := Documentation{Language: "en", Content: "ordinary"}
	got := want
	finishConversionValue(node, &got)
	if got != want {
		t.Fatalf("ordinary publication = %#v; want %#v", got, want)
	}
	cancel()
	finishConversionValue(node, &got)
	if got != (Documentation{}) || !errors.Is(node.contextError(), context.Canceled) {
		t.Fatalf("canceled publication = %#v, %v; want zero and context.Canceled", got, node.contextError())
	}
}

func TestFinishConversionEstablishedErrorPrecedence(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	node := &xmlNode{owner: &parseState{ctx: ctx}}
	want := Documentation{Content: "ordinary"}
	got := want
	terminal := errors.New("established terminal cause")
	err := terminal
	cancel()
	finishConversion(node, &got, &err)
	if err != terminal || got != want {
		t.Fatalf("terminal result = %#v, %v; want unchanged value and exact terminal cause", got, err)
	}
	err = nil
	finishConversion(node, &got, &err)
	if !errors.Is(err, context.Canceled) || got != (Documentation{}) {
		t.Fatalf("successful canceled result = %#v, %v; want zero and context.Canceled", got, err)
	}
}
