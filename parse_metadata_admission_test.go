package wsdl_test

import (
	"context"
	"errors"
	"testing"

	wsdl "github.com/faustbrian/go-wsdl/v3"
)

func TestParseCanceledContextPrecedesMalformedXMLAdmission(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	source := []byte(`<description`)
	document, err := wsdl.Parse(ctx, source, wsdl.ParseOptions{})
	if document != nil || !errors.Is(err, context.Canceled) || err.Error() != "wsdl: parse failed" {
		t.Fatalf("canceled malformed XML = %v, %v; want nil and private cancellation error", document, err)
	}
	if document, err := wsdl.Parse(context.Background(), source, wsdl.ParseOptions{}); document != nil || err == nil || errors.Is(err, context.Canceled) {
		t.Fatalf("ordinary malformed XML = %v, %v; want non-cancellation failure", document, err)
	}
	if document, err := wsdl.Parse(ctx, source, wsdl.ParseOptions{MaxDepth: -1}); document != nil || err == nil || errors.Is(err, context.Canceled) {
		t.Fatalf("invalid options = %v, %v; want option failure before cancellation", document, err)
	}
}

func TestParsePreservesDirectMIMEXMLPart(t *testing.T) {
	source := []byte(`<definitions xmlns="http://schemas.xmlsoap.org/wsdl/" xmlns:tns="urn:test" xmlns:mime="http://schemas.xmlsoap.org/wsdl/mime/" targetNamespace="urn:test"><binding name="Binding" type="tns:API"><operation name="Call"><input><mime:mimeXml part="body"/></input></operation></binding></definitions>`)
	document, err := wsdl.Parse(context.Background(), source, wsdl.ParseOptions{})
	if err != nil {
		t.Fatal(err)
	}
	model, ok := document.Definitions11()
	if !ok || len(model.Bindings) != 1 || len(model.Bindings[0].Operations) != 1 || model.Bindings[0].Operations[0].Input == nil {
		t.Fatalf("missing binding input: %#v", model)
	}
	mime := model.Bindings[0].Operations[0].Input.MIME
	if mime == nil || len(mime.XML) != 1 || mime.XML[0].Part != "body" {
		t.Fatalf("MIME XML = %#v; want one body part", mime)
	}
}

func TestParsePreservesHTTPFaultReferenceTransferCoding(t *testing.T) {
	source := []byte(`<description xmlns="http://www.w3.org/ns/wsdl" xmlns:tns="urn:test" xmlns:whttp="http://www.w3.org/ns/wsdl/http" targetNamespace="urn:test"><binding name="HTTP" interface="tns:API" type="http://www.w3.org/ns/wsdl/http"><operation ref="tns:Call"><outfault ref="tns:Failure" whttp:transferCoding="chunked"/></operation></binding></description>`)
	document, err := wsdl.Parse(context.Background(), source, wsdl.ParseOptions{})
	if err != nil {
		t.Fatal(err)
	}
	model, ok := document.Description20()
	if !ok || len(model.Bindings) != 1 || len(model.Bindings[0].Operations) != 1 || len(model.Bindings[0].Operations[0].OutFaults) != 1 {
		t.Fatalf("missing binding fault reference: %#v", model)
	}
	http := model.Bindings[0].Operations[0].OutFaults[0].HTTP
	if http == nil || !http.TransferCodingSet || http.TransferCoding != "chunked" {
		t.Fatalf("HTTP fault reference = %#v; want explicit chunked transfer coding", http)
	}
}

func TestParsePreservesExplicitSOAPFaultSubcode(t *testing.T) {
	source := []byte(`<description xmlns="http://www.w3.org/ns/wsdl" xmlns:tns="urn:test" xmlns:wsoap="http://www.w3.org/ns/wsdl/soap" targetNamespace="urn:test"><binding name="SOAP" interface="tns:API" type="http://www.w3.org/ns/wsdl/soap"><fault ref="tns:Failure" wsoap:subcodes="tns:Specific"/></binding></description>`)
	document, err := wsdl.Parse(context.Background(), source, wsdl.ParseOptions{})
	if err != nil {
		t.Fatal(err)
	}
	model, ok := document.Description20()
	if !ok || len(model.Bindings) != 1 || len(model.Bindings[0].Faults) != 1 {
		t.Fatalf("missing binding fault: %#v", model)
	}
	soap := model.Bindings[0].Faults[0].SOAP
	if soap == nil || !soap.SubcodesSet || soap.SubcodesAny || len(soap.Subcodes) != 1 || soap.Subcodes[0] != (wsdl.QName{Namespace: "urn:test", Local: "Specific"}) {
		t.Fatalf("SOAP fault = %#v; want explicit {urn:test}Specific subcode", soap)
	}
}
