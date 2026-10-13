package wsdl_test

import (
	"context"
	"errors"
	"testing"
	"time"

	wsdl "github.com/faustbrian/go-wsdl/v3"
	wsdlcompile "github.com/faustbrian/go-wsdl/v3/compile"
	"github.com/faustbrian/go-wsdl/v3/resolve"
)

func TestNilContextAdmissionReturnsNoWSDLResult(t *testing.T) {
	compiler, err := wsdlcompile.New(wsdlcompile.Options{})
	if err != nil {
		t.Fatal(err)
	}
	for name, source := range map[string]string{
		"WSDL11":    `<definitions xmlns="http://schemas.xmlsoap.org/wsdl/" targetNamespace="urn:test"><portType name="Ordinary"/></definitions>`,
		"WSDL20":    `<description xmlns="http://www.w3.org/ns/wsdl" targetNamespace="urn:test"><interface name="Ordinary"/></description>`,
		"malformed": `<description`,
		"empty":     "",
	} {
		t.Run(name, func(t *testing.T) {
			for _, operation := range []struct {
				name, category string
				run            func(context.Context) (bool, error)
			}{
				{"Parse", "wsdl: parse failed", func(ctx context.Context) (bool, error) {
					document, err := wsdl.Parse(ctx, []byte(source), wsdl.ParseOptions{})
					return document != nil, err
				}},
				{"Compile", "wsdl compile: failed", func(ctx context.Context) (bool, error) {
					set, err := compiler.Compile(ctx, wsdlcompile.Source{URI: "urn:ordinary", Content: []byte(source)})
					return set != nil, err
				}},
			} {
				t.Run(operation.name, func(t *testing.T) {
					defer func() {
						if recovered := recover(); recovered != nil {
							t.Errorf("nil-context admission panicked instead of returning a failure: %T", recovered)
						}
					}()
					published, err := operation.run(nil)
					if published || err == nil || err.Error() != operation.category {
						t.Fatalf("nil-context admission published=%v, failure=%v", published, err)
					}
					if errors.Unwrap(err) == nil {
						t.Fatal("failure lost its inspectable cause")
					}
				})
			}
		})
	}
}

func TestContextAdmissionRetainsOrdinaryWSDLAndCancellation(t *testing.T) {
	compiler, err := wsdlcompile.New(wsdlcompile.Options{})
	if err != nil {
		t.Fatal(err)
	}
	for _, source := range []string{
		`<definitions xmlns="http://schemas.xmlsoap.org/wsdl/" targetNamespace="urn:test"><portType name="Ordinary"/></definitions>`,
		`<description xmlns="http://www.w3.org/ns/wsdl" targetNamespace="urn:test"><interface name="Ordinary"/></description>`,
	} {
		if document, err := wsdl.Parse(context.Background(), []byte(source), wsdl.ParseOptions{}); document == nil || err != nil {
			t.Fatalf("ordinary Parse = %v, %v", document, err)
		}
		if set, err := compiler.Compile(context.Background(), wsdlcompile.Source{URI: "urn:ordinary", Content: []byte(source)}); set == nil || err != nil || len(set.Interfaces()) != 1 || set.Interfaces()[0].Name.Local != "Ordinary" {
			t.Fatalf("ordinary Compile = %v, %v", set, err)
		}
		canceled, cancel := context.WithCancel(t.Context())
		cancel()
		expired, stop := context.WithDeadline(t.Context(), time.Unix(0, 0))
		t.Cleanup(stop)
		for _, test := range []struct {
			ctx   context.Context
			cause error
		}{
			{canceled, context.Canceled}, {expired, context.DeadlineExceeded},
		} {
			if document, err := wsdl.Parse(test.ctx, []byte(source), wsdl.ParseOptions{}); document != nil || !errors.Is(err, test.cause) || err.Error() != "wsdl: parse failed" {
				t.Fatalf("stopped Parse = %v, %v", document, err)
			}
			if set, err := compiler.Compile(test.ctx, wsdlcompile.Source{URI: "urn:ordinary", Content: []byte(source)}); set != nil || !errors.Is(err, test.cause) || err.Error() != "wsdl compile: failed" {
				t.Fatalf("stopped Compile = %v, %v", set, err)
			}
		}
	}
}

// A real injected resolver is reachable for an admitted imported document, but
// must not receive an operation that failed nil-context admission.
type auditAdmissionResolver struct {
	calls   int
	failure error
}

func (r *auditAdmissionResolver) Resolve(context.Context, resolve.Request) (resolve.Resource, error) {
	r.calls++
	return resolve.Resource{}, r.failure
}

func TestNilContextCompileDoesNotAdmitResolverWork(t *testing.T) {
	failure := errors.New("resolver unavailable")
	resolver := &auditAdmissionResolver{failure: failure}
	compiler, err := wsdlcompile.New(wsdlcompile.Options{Resolver: resolver})
	if err != nil {
		t.Fatal(err)
	}
	source := wsdlcompile.Source{URI: "https://example.test/root.wsdl", Content: []byte(`<description xmlns="http://www.w3.org/ns/wsdl" targetNamespace="urn:root"><import namespace="urn:child" location="child.wsdl"/></description>`)}
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Errorf("nil-context admission panicked: %T", recovered)
		}
	}()
	//lint:ignore SA1012 Nil is the rejected input under test; substituting TODO would skip admission.
	//nolint:staticcheck // SA1012: this single call exercises nil-context rejection, not production context use.
	if set, err := compiler.Compile(nil, source); set != nil || err == nil || resolver.calls != 0 {
		t.Fatalf("unadmitted Compile = %v, %v; resolver calls=%d", set, err, resolver.calls)
	}
	if set, err := compiler.Compile(context.Background(), source); set != nil || !errors.Is(err, failure) || resolver.calls != 1 {
		t.Fatalf("admitted resolver failure = %v, %v; resolver calls=%d", set, err, resolver.calls)
	}
}

func TestContextAdmissionPreservesPreflightPrecedence(t *testing.T) {
	stopped, cancel := context.WithCancel(t.Context())
	cancel()
	source := []byte(`<description xmlns="http://www.w3.org/ns/wsdl" targetNamespace="urn:test"/>`)
	invalid := wsdl.ParseOptions{MaxDepth: -1}
	_, ordinaryError := wsdl.Parse(context.Background(), source, invalid)
	ordinaryCause := errors.Unwrap(ordinaryError)
	if ordinaryCause == nil {
		t.Fatal("invalid options must retain their inspectable cause")
	}
	compiler, err := wsdlcompile.New(wsdlcompile.Options{Limits: wsdlcompile.Limits{MaxBytes: 1}})
	if err != nil {
		t.Fatal(err)
	}
	for _, ctx := range []context.Context{nil, stopped} {
		if document, err := wsdl.Parse(ctx, source, invalid); document != nil || err == nil || errors.Unwrap(err) == nil || errors.Unwrap(err).Error() != ordinaryCause.Error() {
			t.Fatalf("invalid options lost precedence: document=%v error=%v", document, err)
		}
		if document, err := wsdl.Parse(ctx, source, wsdl.ParseOptions{MaxDocumentBytes: 1}); document != nil || !errors.Is(err, wsdl.ErrLimitExceeded) {
			t.Fatalf("raw-byte admission lost precedence: document=%v error=%v", document, err)
		}
		if set, err := compiler.Compile(ctx, wsdlcompile.Source{URI: "relative.wsdl", Content: source}); set != nil || !errors.Is(err, wsdlcompile.ErrResourceIdentity) {
			t.Fatalf("root identity admission lost precedence: set=%v error=%v", set, err)
		}
		if set, err := compiler.Compile(ctx, wsdlcompile.Source{URI: "urn:ordinary", Content: source}); set != nil || !errors.Is(err, wsdlcompile.ErrLimitExceeded) {
			t.Fatalf("compile byte admission lost precedence: set=%v error=%v", set, err)
		}
	}
}
