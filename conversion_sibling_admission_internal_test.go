package wsdl

import (
	"context"
	"encoding/xml"
	"testing"
)

// Constructed internal trees acquire cancellation ownership through useOwner.
// This is not a claim about a public Parse path with an already-owned XML tree.
func TestConstructedSiblingAdmissionPropagatesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	root := &xmlNode{children: []*xmlNode{
		{attributes: []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "First"}}},
		{attributes: []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "Second"}}},
	}}
	if err := root.useOwner(&parseState{ctx: ctx}); err != nil {
		t.Fatal(err)
	}
	if got := root.children[1].attribute("name"); got != "Second" {
		t.Fatalf("ordinary second sibling name = %q; want Second", got)
	}
	cancel()
	if got := root.children[1].attribute("name"); got != "" {
		t.Fatalf("canceled second sibling published name %q", got)
	}
}
