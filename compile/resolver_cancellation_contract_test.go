package compile_test

import (
	"context"
	"errors"
	"testing"

	wsdl "github.com/faustbrian/go-wsdl/v3"
	wsdlcompile "github.com/faustbrian/go-wsdl/v3/compile"
	"github.com/faustbrian/go-wsdl/v3/resolve"
)

// This uses the public resolver handoff, not a context with synthetic checks.
// Cancellation belongs to one operation; it must not poison the compiler.
func TestCompilerResolverCancellationContract(t *testing.T) {
	for _, fixture := range []struct {
		name        string
		version     wsdl.Version
		root, child string
	}{
		{
			name: "WSDL11", version: wsdl.Version11,
			root:  `<definitions xmlns="http://schemas.xmlsoap.org/wsdl/" targetNamespace="urn:root"><import namespace="urn:child" location="child.wsdl"/></definitions>`,
			child: `<definitions xmlns="http://schemas.xmlsoap.org/wsdl/" targetNamespace="urn:child"><portType name="Child"/></definitions>`,
		},
		{
			name: "WSDL20", version: wsdl.Version20,
			root:  `<description xmlns="http://www.w3.org/ns/wsdl" targetNamespace="urn:root"><import namespace="urn:child" location="child.wsdl"/></description>`,
			child: `<description xmlns="http://www.w3.org/ns/wsdl" targetNamespace="urn:child"><interface name="Child"/></description>`,
		},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			const rootURI = "https://example.test/root.wsdl"
			const childURI = "https://example.test/child.wsdl"
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			resolver := resolverCancellationContractFunc(func(requestCtx context.Context, request resolve.Request) (resolve.Resource, error) {
				if request.URI != childURI || request.Namespace != "urn:child" || request.Kind != resolve.KindImport || request.Version != string(fixture.version) {
					t.Fatalf("resolver request = %#v", request)
				}
				if requestCtx == ctx {
					cancel()
				}
				return resolve.Resource{URI: childURI, Content: []byte(fixture.child)}, nil
			})
			compiler, err := wsdlcompile.New(wsdlcompile.Options{Resolver: resolver})
			if err != nil {
				t.Fatal(err)
			}
			source := wsdlcompile.Source{URI: rootURI, Content: []byte(fixture.root)}
			set, err := compiler.Compile(ctx, source)
			if set != nil || !errors.Is(err, context.Canceled) || err.Error() != "wsdl compile: failed" {
				t.Fatalf("canceled resolver handoff = %v, %v; want nil Set and context.Canceled with fixed public text", set, err)
			}
			set, err = compiler.Compile(context.Background(), source)
			if err != nil || set == nil {
				t.Fatalf("compiler reuse = %v, %v", set, err)
			}
			documents := set.Documents()
			if len(documents) != 2 || documents[0].URI != childURI || documents[1].URI != rootURI ||
				documents[0].Namespace != "urn:child" || documents[1].Namespace != "urn:root" ||
				documents[0].Version != fixture.version || documents[1].Version != fixture.version {
				t.Fatalf("reused compiler document identities = %#v", documents)
			}
		})
	}
}

type resolverCancellationContractFunc func(context.Context, resolve.Request) (resolve.Resource, error)

func (f resolverCancellationContractFunc) Resolve(ctx context.Context, request resolve.Request) (resolve.Resource, error) {
	return f(ctx, request)
}
