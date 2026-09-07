# FAQ

## Does parsing download imports?

No. Parsing is byte-local. Compilation resolves only through explicitly
injected WSDL and schema resolvers; both defaults deny.

## Does this package send SOAP or HTTP requests?

No. It describes and validates bindings. Envelope primitives belong to
`wire`; transport belongs to a client package.

## Why are WSDL 1.1 and 2.0 separate models?

Their component models, operation semantics, and binding adjuncts differ.
Combining them would erase meaningful presence and direction information.

## Is a successful parse a conformance result?

No. Run validation and, for imports or schemas, compilation. Consult the
version matrix for the exact supported claim.

## Can unknown extensions round-trip?

Yes, within limits. A required extension still fails validation unless the
caller explicitly declares that QName understood.

## Troubleshooting

### Why does compilation reject an import that parsing accepted?

Parsing records references but performs no I/O. Compilation denies resolution
unless the caller injects an explicit WSDL resolver and, when schemas are
loaded, a separate XML Schema resolver. See
[resolution and compilation](resolution-and-compilation.md).

### How do I inspect validation failures?

Validation returns ordered diagnostics with stable codes and source locations
when available. Use `Diagnostics.Err()` when a workflow needs one combined
error, and retain the diagnostics when callers need each failure independently.
See [validation](validation.md).

### Where should I report a problem?

Use the repository's [support routes](../SUPPORT.md) for reproducible defects
and questions. Report suspected vulnerabilities only through the
[private security process](../SECURITY.md#reporting).
