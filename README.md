# wsdl

[![CI](https://github.com/faustbrian/go-wsdl/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-wsdl/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-wsdl/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-wsdl.svg)](https://pkg.go.dev/github.com/faustbrian/go-wsdl)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-wsdl?sort=semver)](https://github.com/faustbrian/go-wsdl/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`wsdl` is a bounded, deterministic WSDL 1.1 and WSDL 2.0 description
toolkit for Go. It parses caller-supplied XML, preserves extension data and
presence semantics, validates component references and bindings, resolves
imports only through injected resolvers, compiles immutable graphs, composes
documents, builds code-generation models, and reports semantic differences.

It is deliberately not a SOAP client. `wire` owns SOAP envelope primitives,
`xsd` owns schema compilation, and transport belongs in `http-client` or
another consumer.

```go
compiler, err := compile.New(compile.Options{}) // resolution denied by default
if err != nil {
    return err
}
set, err := compiler.Compile(ctx, compile.Source{
    URI:     "https://example.test/service.wsdl",
    Content: source,
})
if err != nil {
    return err
}
service, ok := set.Service(wsdl.QName{Namespace: "urn:example", Local: "API"})
```

The [documentation](docs/README.md) covers the model, security boundaries,
version-specific conformance, builders, composition, code generation,
interoperability, and release evidence. Observable specification choices are
recorded in the [decision register](docs/specification-decisions.md).
`make check` runs the normal local gate; `make check-all` also runs coverage,
fuzzing, benchmarks, and mutation.

For ecosystem-wide selection and ownership guidance, see the versioned
[Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/README.md)
and its
[Protocols and descriptions family](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/design-language.md#package-families-and-selection).

## Stability

The API is pre-1.0. Supported behavior is recorded independently in the
[WSDL 1.1 matrix](specification/requirements/wsdl-1.1.tsv) and
[WSDL 2.0 matrix](specification/requirements/wsdl-2.0.tsv). Matrix rows marked
`partial` or `missing` are not conformance claims.
