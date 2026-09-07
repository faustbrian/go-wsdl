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

## Status

The module has a stable v1 API and supports Go 1.26.6 or later.

## Installation

```sh
go get github.com/faustbrian/go-wsdl@v1
```

Import only the packages your application uses. The root `wsdl` package owns
document parsing, validation, and shared description types; `compile` owns
explicit resource resolution and immutable compiled graphs.

## Quick start

```go
package main

import (
    "context"
    "fmt"

    wsdl "github.com/faustbrian/go-wsdl"
)

func main() {
    source := []byte(`<definitions xmlns="http://schemas.xmlsoap.org/wsdl/"` +
        ` name="Inventory" targetNamespace="urn:inventory"/>`)

    document, err := wsdl.Parse(context.Background(), source, wsdl.ParseOptions{})
    if err != nil {
        fmt.Println(err)
        return
    }
    definitions, ok := document.Definitions11()
    if !ok {
        fmt.Println("not a WSDL 1.1 document")
        return
    }

    fmt.Println(document.Version())
    fmt.Println(definitions.Name)
    fmt.Println(definitions.TargetNamespace)
}
```

See the [executable example](example_test.go) for the checked output of this
standalone parse and inspection flow. The [package map](docs/README.md) routes advanced
compilation, composition, code generation, and compatibility use cases.

The module has no runtime service lifecycle: it starts no background goroutines,
owns no external resources, and requires no shutdown. Parsing and compilation
observe caller cancellation between bounded internal steps. Injected resolvers
receive that context and must honor it and bound their own I/O. Builders are
single-owner while compiled sets support concurrent lookup; see
[security and limits](docs/security.md),
[validation](docs/validation.md), and
[resolution and compilation](docs/resolution-and-compilation.md) for errors,
defaults, cancellation, ownership, and concurrency details.

The [documentation](docs/README.md) covers the model, security boundaries,
version-specific conformance, builders, composition, code generation,
interoperability, and release evidence. Observable specification choices are
recorded in the [decision register](docs/specification-decisions.md).

Project resources include the [API reference](https://pkg.go.dev/github.com/faustbrian/go-wsdl),
[executable example](example_test.go), [testing guidance](CONTRIBUTING.md#verification),
[FAQ](docs/faq.md), [troubleshooting](docs/faq.md#troubleshooting),
[changelog](CHANGELOG.md), [license](LICENSE), [support](SUPPORT.md), and
[private security reporting](SECURITY.md#reporting).

`make check` runs the normal local gate; `make ci` runs the complete repository
contract, including coverage, fuzzing, benchmarks, mutation, conformance, and
interoperability.

For ecosystem-wide selection and ownership guidance, see the versioned
[Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/README.md)
and its
[Protocols and descriptions family](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/design-language.md#package-families-and-selection).

## Stability

The stable v1 API follows semantic versioning. Supported behavior is recorded
independently in the
[WSDL 1.1 matrix](specification/requirements/wsdl-1.1.tsv) and
[WSDL 2.0 matrix](specification/requirements/wsdl-2.0.tsv). Matrix rows marked
`partial` or `missing` are not conformance claims.
