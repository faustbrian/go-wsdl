package compile_test

import (
	"context"
	"errors"
	"testing"

	wsdlcompile "github.com/faustbrian/go-wsdl/v2/compile"
)

func TestCompilerRemainsReusableAfterParserCancellation(t *testing.T) {
	compiler, err := wsdlcompile.New(wsdlcompile.Options{})
	if err != nil {
		t.Fatal(err)
	}
	source := wsdlcompile.Source{URI: "urn:ordinary", Content: []byte(`<description xmlns="http://www.w3.org/ns/wsdl" targetNamespace="urn:test"><interface name="Ordinary"/></description>`)}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if set, err := compiler.Compile(ctx, source); set != nil || !errors.Is(err, context.Canceled) || err.Error() != "wsdl compile: failed" {
		t.Fatalf("canceled compile = %v, %v", set, err)
	}
	set, err := compiler.Compile(context.Background(), source)
	if set == nil || err != nil || len(set.Interfaces()) != 1 || set.Interfaces()[0].Name.Local != "Ordinary" {
		t.Fatalf("ordinary compile after cancellation = %v, %v", set, err)
	}
}
