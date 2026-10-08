package resolve_test

import (
	"context"
	"errors"
	"testing"

	"github.com/faustbrian/go-wsdl/v2/resolve"
)

func TestMemoryReturnsOwnedResourceCopies(t *testing.T) {
	t.Parallel()

	resolver, err := resolve.NewMemory(map[string][]byte{
		"https://example.test/service.wsdl": []byte("original"),
	})
	if err != nil {
		t.Fatalf("NewMemory() error = %v", err)
	}
	request := resolve.Request{URI: "https://example.test/service.wsdl"}
	resource, err := resolver.Resolve(context.Background(), request)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	resource.Content[0] = 'X'
	again, err := resolver.Resolve(context.Background(), request)
	if err != nil || string(again.Content) != "original" {
		t.Fatalf("second Resolve() = %#v, %v", again, err)
	}
}

func TestResolversRejectInvalidMissingAndDeniedResources(t *testing.T) {
	t.Parallel()

	invalid := []string{"relative.wsdl", "https://example.test/a.wsdl#part", "://bad"}
	for _, identity := range invalid {
		if _, err := resolve.NewMemory(map[string][]byte{identity: nil}); err == nil {
			t.Fatalf("NewMemory(%q) error = nil", identity)
		}
	}
	resolver, err := resolve.NewMemory(nil)
	if err != nil {
		t.Fatalf("NewMemory() error = %v", err)
	}
	request := resolve.Request{URI: "https://example.test/missing.wsdl"}
	if _, err := resolver.Resolve(context.Background(), request); !errors.Is(err, resolve.ErrNotFound) {
		t.Fatalf("Resolve(missing) error = %v", err)
	}
	if _, err := resolve.Deny().Resolve(context.Background(), request); !errors.Is(err, resolve.ErrAccessDenied) {
		t.Fatalf("Deny.Resolve() error = %v", err)
	}
}

func TestResolverChainHonorsCancellationAndErrorSemantics(t *testing.T) {
	t.Parallel()

	memory, err := resolve.NewMemory(map[string][]byte{
		"https://example.test/service.wsdl": []byte("service"),
	})
	if err != nil {
		t.Fatalf("NewMemory() error = %v", err)
	}
	request := resolve.Request{URI: "https://example.test/service.wsdl"}
	resource, err := resolve.Chain(nil, memory).Resolve(context.Background(), request)
	if err != nil || string(resource.Content) != "service" {
		t.Fatalf("Chain.Resolve() = %#v, %v", resource, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := memory.Resolve(ctx, request); !errors.Is(err, context.Canceled) {
		t.Fatalf("Resolve(canceled) error = %v", err)
	}
	if _, err := resolve.Chain(resolve.Deny(), memory).Resolve(
		context.Background(),
		request,
	); !errors.Is(err, resolve.ErrAccessDenied) {
		t.Fatalf("Chain(deny).Resolve() error = %v", err)
	}
}

func TestResolversCoverNilCancellationAndExhaustion(t *testing.T) {
	t.Parallel()

	request := resolve.Request{URI: "https://example.test/service.wsdl"}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := resolve.Deny().Resolve(ctx, request); !errors.Is(err, context.Canceled) {
		t.Fatalf("Deny.Resolve(canceled) error = %v", err)
	}
	var memory *resolve.Memory
	if _, err := memory.Resolve(context.Background(), request); !errors.Is(err, resolve.ErrAccessDenied) {
		t.Fatalf("nil Memory.Resolve() error = %v", err)
	}
	if _, err := resolve.Chain().Resolve(context.Background(), request); !errors.Is(err, resolve.ErrNotFound) {
		t.Fatalf("empty Chain.Resolve() error = %v", err)
	}
	if _, err := resolve.Chain().Resolve(ctx, request); !errors.Is(err, context.Canceled) {
		t.Fatalf("Chain.Resolve(canceled) error = %v", err)
	}
}

func TestResolverChainStopsBeforeNextLookupAfterCancellation(t *testing.T) {
	t.Parallel()

	request := resolve.Request{URI: "https://example.test/service.wsdl"}
	memory, err := resolve.NewMemory(map[string][]byte{request.URI: []byte("service")})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	first := resolverFunc(func(context.Context, resolve.Request) (resolve.Resource, error) {
		cancel()
		return resolve.Resource{}, resolve.ErrNotFound
	})
	nextLookups := 0
	next := resolverFunc(func(ctx context.Context, request resolve.Request) (resolve.Resource, error) {
		nextLookups++
		return memory.Resolve(ctx, request)
	})
	resource, err := resolve.Chain(first, next).Resolve(ctx, request)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Resolve() error = %v, want context.Canceled", err)
	}
	if resource.URI != "" || resource.Content != nil {
		t.Errorf("Resolve() returned resource after cancellation: %#v", resource)
	}
	if nextLookups != 0 {
		t.Errorf("next lookup invoked %d times after cancellation, want 0", nextLookups)
	}
}

type resolverFunc func(context.Context, resolve.Request) (resolve.Resource, error)

func (f resolverFunc) Resolve(ctx context.Context, request resolve.Request) (resolve.Resource, error) {
	return f(ctx, request)
}
