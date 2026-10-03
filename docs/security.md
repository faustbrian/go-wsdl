# Security and limits

The XML reader rejects DTDs and custom entity declarations. Parsing and
serialization bound bytes, depth, elements, attributes, text, schemas,
imports, operations, bindings, endpoints, extensions, diagnostics, and output.
Compilation adds graph depth, documents, references, cumulative bytes, schema
limits, and components. Code-generation models have independent limits.

No package opens a file, follows a redirect, contacts a host, reads an
environment variable, consults a global registry, or starts background work.
Injected resolvers are the only external I/O seam. Production resolvers should
allowlist schemes and hosts, reject private and link-local address changes,
bound redirects and response bytes, and avoid forwarding credentials.

Cancellation stops parsing and graph resolution. Returned models and compiled
sets own their mutable slices so callers cannot mutate shared compiler state.

## Default errors and explicit diagnostics (unreleased)

Library-produced default error strings contain a fixed failure category, not
resource URIs, component names, lexical values or collaborator error text.
`errors.Is` and `errors.As` retain error classification and typed causes.
`Diagnostic.Message`, `Path`, `Location`, `ConflictError.Conflicts` and unwrapped
causes remain intentionally inspectable and can contain input-derived details.
Applications own authorization and redaction when explicitly recording those
details. Custom resolvers still own their directly returned error text; WSDL
compiler wrappers do not print that text by default.

This behavior is a pending next-major change on main, not a guarantee of the
published v1.0.0 release. The `/v2` migration and release, parser-preflight
cancellation, Wire/XSD major adoption and versioned threat model remain pending.
