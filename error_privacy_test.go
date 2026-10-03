package wsdl_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	wsdl "github.com/faustbrian/go-wsdl"
	"github.com/faustbrian/go-wsdl/builder"
	wsdlcompile "github.com/faustbrian/go-wsdl/compile"
	"github.com/faustbrian/go-wsdl/compose"
	"github.com/faustbrian/go-wsdl/resolve"
)

func privacyCategory(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatal("default error is nil")
	}
	if err.Error() != want {
		t.Errorf("default error = %v, want category %q", err, want)
	}
}

func ordinaryDiagnosticCount(t *testing.T, err error, code string) int {
	t.Helper()
	var diagnostics wsdl.Diagnostics
	if !errors.As(err, &diagnostics) {
		t.Fatalf("error has no typed diagnostics: %v", err)
	}
	count := 0
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code {
			count++
		}
	}
	return count
}

func TestOrdinaryDefaultErrorPrivacy(t *testing.T) {
	const uri = "https://example.test/catalog.wsdl"
	t.Run("built-in resolution", func(t *testing.T) {
		memory, err := resolve.NewMemory(nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, test := range []struct {
			resolver resolve.Resolver
			sentinel error
		}{
			{resolve.Deny(), resolve.ErrAccessDenied},
			{memory, resolve.ErrNotFound},
			{resolve.Chain(), resolve.ErrNotFound},
		} {
			resource, err := test.resolver.Resolve(context.Background(), resolve.Request{URI: uri})
			privacyCategory(t, err, test.sentinel.Error())
			if !errors.Is(err, test.sentinel) || resource.URI != "" || resource.Content != nil {
				t.Fatal("resolution category/result changed")
			}
		}
	})
	t.Run("semantic diagnostics", func(t *testing.T) {
		const source = `<definitions xmlns="http://schemas.xmlsoap.org/wsdl/" xmlns:tns="urn:catalog" targetNamespace="urn:catalog"><service name="Catalog"><port name="Public" binding="tns:Missing"/></service></definitions>`
		document, err := wsdl.Parse(context.Background(), []byte(source), wsdl.ParseOptions{SystemID: uri})
		if err != nil {
			t.Fatal(err)
		}
		diagnostics := wsdl.Validate(document, wsdl.ValidationOptions{})
		privacyCategory(t, diagnostics.Err(), "wsdl: validation failed")
		var inspected wsdl.Diagnostics
		if !errors.As(diagnostics.Err(), &inspected) || len(inspected) != 1 || inspected[0].Code != "WSDL11_BINDING_REFERENCE" || !strings.Contains(inspected[0].Message, "Missing") || !strings.Contains(inspected[0].Path, "Catalog") || inspected[0].Location.SystemID != uri {
			t.Fatalf("explicit diagnostics lost: %#v", inspected)
		}
		model, _ := document.Definitions11()
		_, err = wsdl.NewDocument11(model, wsdl.ValidationOptions{})
		privacyCategory(t, err, "wsdl: validate model")
		if !errors.As(err, &inspected) || inspected[0].Message != diagnostics[0].Message {
			t.Fatal("model diagnostic cause lost")
		}
	})
	t.Run("compiler identity and cause", func(t *testing.T) {
		const source = `<description xmlns="http://www.w3.org/ns/wsdl" targetNamespace="urn:catalog"><include location="child.wsdl"/></description>`
		cause := &ordinaryResolverError{message: "catalog lookup unavailable"}
		for _, test := range []struct {
			resource resolve.Resource
			cause    error
			sentinel error
		}{
			{resolve.Resource{URI: "https://example.test/other.wsdl"}, nil, wsdlcompile.ErrResourceIdentity},
			{resolve.Resource{}, cause, cause},
		} {
			compiler, err := wsdlcompile.New(wsdlcompile.Options{Resolver: ordinaryPrivacyResolver{resource: test.resource, err: test.cause}})
			if err != nil {
				t.Fatal(err)
			}
			set, err := compiler.Compile(context.Background(), wsdlcompile.Source{URI: uri, Content: []byte(source)})
			privacyCategory(t, err, "wsdl compile: failed")
			if set != nil || !errors.Is(err, test.sentinel) {
				t.Fatal("compiler result/cause changed")
			}
			if test.cause != nil {
				var inspected *ordinaryResolverError
				if !errors.As(err, &inspected) || inspected != cause || errors.Unwrap(err) == nil {
					t.Fatal("WSDL boundary or typed collaborator cause changed")
				}
			}
		}
	})
	t.Run("composition and builder", func(t *testing.T) {
		document, err := wsdl.NewDocument20(wsdl.Description20{TargetNamespace: "urn:catalog", Interfaces: []wsdl.Interface20{{Name: "Catalog"}}}, wsdl.ValidationOptions{})
		if err != nil {
			t.Fatal(err)
		}
		merged, err := compose.Merge(document, document)
		privacyCategory(t, err, compose.ErrConflict.Error())
		var conflicts *compose.ConflictError
		if merged != nil || !errors.Is(err, compose.ErrConflict) || !errors.As(err, &conflicts) || len(conflicts.Conflicts) != 1 || conflicts.Conflicts[0].Name != "Catalog" {
			t.Fatal("explicit conflicts lost")
		}
		value := builder.New20("urn:catalog")
		if err := value.AddInterface(wsdl.Interface20{Name: "Catalog"}); err != nil {
			t.Fatal(err)
		}
		err = value.AddInterface(wsdl.Interface20{Name: "Catalog"})
		privacyCategory(t, err, builder.ErrDuplicateComponent.Error())
		if !errors.Is(err, builder.ErrDuplicateComponent) {
			t.Fatal("builder category lost")
		}
	})
	t.Run("parser", func(t *testing.T) {
		document, err := wsdl.Parse(context.Background(), []byte(`<catalog xmlns="urn:catalog"/>`), wsdl.ParseOptions{SystemID: uri})
		privacyCategory(t, err, "wsdl: parse failed")
		if document != nil || errors.Unwrap(err) == nil {
			t.Fatal("parser result/cause changed")
		}
	})
}

type ordinaryResolverError struct{ message string }

func (e *ordinaryResolverError) Error() string { return e.message }

type ordinaryPrivacyResolver struct {
	resource resolve.Resource
	err      error
}

func (r ordinaryPrivacyResolver) Resolve(context.Context, resolve.Request) (resolve.Resource, error) {
	return r.resource, r.err
}
