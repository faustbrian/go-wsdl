# Documentation

## Package map

| Package | Use |
| --- | --- |
| [`wsdl`](https://pkg.go.dev/github.com/faustbrian/go-wsdl) | Parse, validate, serialize, and inspect WSDL 1.1 and WSDL 2.0 descriptions. |
| [`builder`](https://pkg.go.dev/github.com/faustbrian/go-wsdl/builder) | Build validated WSDL documents through a single-owner builder. |
| [`codegen`](https://pkg.go.dev/github.com/faustbrian/go-wsdl/codegen) | Derive a bounded, deterministic, language-neutral generation model from a compiled set. |
| [`compile`](https://pkg.go.dev/github.com/faustbrian/go-wsdl/compile) | Resolve imports and XML Schemas through injected resolvers and compile an immutable graph. |
| [`compose`](https://pkg.go.dev/github.com/faustbrian/go-wsdl/compose) | Merge compatible same-version, same-namespace WSDL documents with explicit conflict reporting. |
| [`diff`](https://pkg.go.dev/github.com/faustbrian/go-wsdl/diff) | Compare compiled sets conservatively for semantic compatibility. |
| [`resolve`](https://pkg.go.dev/github.com/faustbrian/go-wsdl/resolve) | Supply explicit deny-by-default or in-memory WSDL resource resolution policies. |

All packages are in-process libraries. They start no background work, own no
transport or service lifecycle, and require no shutdown. Use the root package
for byte-local description work and `compile` when an application needs an
immutable component graph; any external resolution must be explicit and
bounded. Use `wire` for SOAP envelopes, `xsd` for schema ownership, and an HTTP
client for transport execution.

## Guides and project resources

- [API reference](https://pkg.go.dev/github.com/faustbrian/go-wsdl)
- [Executable example](../example_test.go)
- [Versions and conformance](conformance.md)
- [Specification decisions](specification-decisions.md)
- [Models and source locations](models.md)
- [Parsing and deterministic generation](parsing-and-generation.md)
- [Validation](validation.md)
- [Resolution and compilation](resolution-and-compilation.md)
- [Builders and composition](builders-and-composition.md)
- [SOAP bindings](soap-bindings.md)
- [HTTP and MIME bindings](http-and-mime-bindings.md)
- [XML Schema integration](xsd-integration.md)
- [Extensions](extensions.md)
- [Code generation model](code-generation.md)
- [Semantic compatibility](compatibility.md)
- [Security and limits](security.md)
- [Interoperability](interoperability.md)
- [Performance](performance.md)
- [Migration](migration.md)
- [Module identity reviews](module-identity-reviews.md)
- [Cookbook](cookbook.md)
- [FAQ](faq.md)
- [Troubleshooting](faq.md#troubleshooting)
- [Compatibility policy](../COMPATIBILITY.md)
- [Testing and contribution guidance](../CONTRIBUTING.md#verification)
- [Changelog](../CHANGELOG.md)
- [License](../LICENSE)
- [Support](../SUPPORT.md)
- [Security reporting](../SECURITY.md#reporting)
