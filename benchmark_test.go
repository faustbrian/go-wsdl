package wsdl_test

import (
	"context"
	"testing"

	wsdl "github.com/faustbrian/go-wsdl/v3"
	wsdlcompile "github.com/faustbrian/go-wsdl/v3/compile"
)

var benchmarkWSDL = []byte(`<description xmlns="http://www.w3.org/ns/wsdl"` +
	` xmlns:tns="urn:benchmark" targetNamespace="urn:benchmark">` +
	`<interface name="API"><operation name="Call"` +
	` pattern="http://www.w3.org/ns/wsdl/in-out"><input element="#none"/>` +
	`<output element="#none"/></operation></interface>` +
	`<binding name="Binding" interface="tns:API" type="urn:binding">` +
	`<operation ref="tns:Call"/></binding><service name="Service" interface="tns:API">` +
	`<endpoint name="Endpoint" binding="tns:Binding" address="https://example.test"/>` +
	`</service></description>`)

// BenchmarkParseWSDL measures both XML passes and model conversion for one fixed
// small WSDL 2.0 payload, not hostile-input allocation bounds or fleet throughput.
func BenchmarkParseWSDL(b *testing.B) {
	for b.Loop() {
		if _, err := wsdl.Parse(context.Background(), benchmarkWSDL, wsdl.ParseOptions{}); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMarshalWSDL excludes initial Parse and measures default-budget deterministic
// encoding of the fixed in-memory model, without network or file IO.
func BenchmarkMarshalWSDL(b *testing.B) {
	document, err := wsdl.Parse(context.Background(), benchmarkWSDL, wsdl.ParseOptions{})
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for b.Loop() {
		if _, err := wsdl.Marshal(document, wsdl.MarshalOptions{}); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCompileWSDL excludes New but includes root parsing and complete compilation
// of a small no-import graph; it does not measure resolver/transport latency.
func BenchmarkCompileWSDL(b *testing.B) {
	compiler, err := wsdlcompile.New(wsdlcompile.Options{})
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for b.Loop() {
		if _, err := compiler.Compile(context.Background(), wsdlcompile.Source{
			URI: "https://example.test/service.wsdl", Content: benchmarkWSDL,
		}); err != nil {
			b.Fatal(err)
		}
	}
}
