// Package resolve defines explicit resource resolution for WSDL graphs.
// Parsing never invokes a Resolver. Start with Deny, NewMemory or an injected
// policy resolver. URI identity does not authorize schemes, hosts or redirects.
// Built-ins require non-nil contexts and own no files, connections or workers.
// Memory copies input/returned bytes; Chain borrows child implementations and
// cannot preempt callbacks or make them concurrency-safe.
package resolve

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/faustbrian/go-wsdl/v3/internal/errorprivacy"
)

var (
	// ErrAccessDenied identifies a resource class rejected by policy.
	ErrAccessDenied = errors.New("wsdl resolve: access denied")
	// ErrNotFound means a resolver does not contain the requested resource.
	ErrNotFound = errors.New("wsdl resolve: resource not found")
)

// Kind identifies the WSDL composition operation requesting a resource.
type Kind string

const (
	// KindImport identifies a WSDL import.
	KindImport Kind = "import"
	// KindInclude identifies a WSDL 2.0 include.
	KindInclude Kind = "include"
)

// Request identifies a resource after base-URI resolution.
type Request struct {
	// URI is the resource identity, not a filesystem path or network authorization decision.
	URI string
	// Namespace is the namespace URI used in expanded identity, not a lexical prefix.
	Namespace string
	// Kind identifies the resource edge or component category, not a protocol action.
	Kind Kind
	// Version identifies the applicable language or protocol version of this component.
	Version string
}

// Resource is an owned WSDL resource. Resolver implementations must not
// reuse Content storage across calls.
type Resource struct {
	// URI is the resource identity, not a filesystem path or network authorization decision.
	URI string
	// Content contains caller-supplied resource bytes; keep them stable while consumed.
	Content []byte
}

// Resolver obtains caller-authorized WSDL resources.
type Resolver interface {
	// Resolve receives a non-nil context and exact resolved identity. Implementations
	// own scheme/host/redirect/credential policy, finite response admission and context
	// cooperation. Return owned Content that remains stable; compilation requires
	// Resource.URI to equal Request.URI exactly, not just identify an equivalent URL.
	Resolve(context.Context, Request) (Resource, error)
}

// Deny returns the secure default resolver, which rejects every request.
func Deny() Resolver { return denyResolver{} }

type denyResolver struct{}

func (denyResolver) Resolve(ctx context.Context, request Request) (Resource, error) {
	if err := ctx.Err(); err != nil {
		return Resource{}, err
	}
	return Resource{}, errorprivacy.Wrap("wsdl resolve: access denied", fmt.Errorf("%w: %s", ErrAccessDenied, request.URI))
}

// Memory is an immutable collection of caller-supplied WSDL resources.
type Memory struct {
	resources map[string][]byte
}

// NewMemory copies resources into an immutable resolver. Keys must be valid,
// absolute resource URIs without fragments. Keys are stored as url.URL.String
// output and Resolve looks up the requested string exactly. The complete input
// is copied without a separate count/byte cap; callers bound admission and keep
// maps/bytes stable until return. Nil input constructs an empty resolver.
func NewMemory(resources map[string][]byte) (*Memory, error) {
	owned := make(map[string][]byte, len(resources))
	for identity, content := range resources {
		uri, err := url.Parse(identity)
		if err != nil || !uri.IsAbs() || uri.Fragment != "" {
			return nil, errorprivacy.Wrap("wsdl resolve: invalid resource URI", fmt.Errorf("wsdl resolve: invalid resource URI %q", identity))
		}
		owned[uri.String()] = append([]byte(nil), content...)
	}
	return &Memory{resources: owned}, nil
}

// Resolve returns a fresh copy of the configured resource bytes.
func (r *Memory) Resolve(ctx context.Context, request Request) (Resource, error) {
	if err := ctx.Err(); err != nil {
		return Resource{}, err
	}
	if r == nil {
		return Resource{}, errorprivacy.Wrap("wsdl resolve: access denied", fmt.Errorf("%w: %s", ErrAccessDenied, request.URI))
	}
	content, ok := r.resources[request.URI]
	if !ok {
		return Resource{}, errorprivacy.Wrap("wsdl resolve: resource not found", fmt.Errorf("%w: %s", ErrNotFound, request.URI))
	}
	return Resource{
		URI: request.URI, Content: append([]byte(nil), content...),
	}, nil
}

// Chain returns a resolver that tries children in order. Only ErrNotFound
// advances to the next resolver while the context remains active; policy and
// operational errors stop. The child list is copied; implementations are borrowed.
// Literal nil children are skipped, but typed-nil collaborators are not detected.
func Chain(resolvers ...Resolver) Resolver {
	return chain{resolvers: append([]Resolver(nil), resolvers...)}
}

type chain struct {
	resolvers []Resolver
}

func (c chain) Resolve(ctx context.Context, request Request) (Resource, error) {
	if err := ctx.Err(); err != nil {
		return Resource{}, err
	}
	for _, resolver := range c.resolvers {
		if resolver == nil {
			continue
		}
		resource, err := resolver.Resolve(ctx, request)
		if err == nil {
			return resource, nil
		}
		if !errors.Is(err, ErrNotFound) {
			return Resource{}, err
		}
		if err := ctx.Err(); err != nil {
			return Resource{}, err
		}
	}
	return Resource{}, errorprivacy.Wrap("wsdl resolve: resource not found", fmt.Errorf("%w: %s", ErrNotFound, request.URI))
}
