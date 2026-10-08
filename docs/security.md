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

## Default errors and explicit diagnostics (v2.0.0)

Library-produced default error strings contain a fixed failure category, not
resource URIs, component names, lexical values or collaborator error text.
`errors.Is` and `errors.As` retain error classification and typed causes.
`Diagnostic.Message`, `Path`, `Location`, `ConflictError.Conflicts` and unwrapped
causes remain intentionally inspectable and can contain input-derived details.
Applications own authorization and redaction when explicitly recording those
details. Custom resolvers still own their directly returned error text; WSDL
compiler wrappers do not print that text by default.

This behavior is included in public v2.0.0, not a guarantee of the
published v1.0.0 release. This source selects public Wire v3 for XML parsing
and serialization, with WSDL's document budget explicitly applied to charset
conversion in both XML passes. Inspected Wire diagnostic types and sentinels
use `/v3`; they are not type-identical to the v1 diagnostics. Main uses the
`/v3` WSDL identity for the pending XSD v2 adoption, not a published v3 claim.
Direct inline parsing independently forwards namespace and model-copy work
allowances; compiler schema parser allowances apply during the first root and
imported WSDL parse, before following further WSDL references. XSD's inline
resolver construction retains finite count, cumulative identity/content and
per-resource byte limits. See [migration](migration.md) for defaults and the
distinct nominal types and error classifications.
The [versioned threat model](security-threat-model.md) identifies owned
boundaries and conditional collaborator obligations; it is not release proof.

## Coverage and qualification

The root module collects coverage evidence under the approved risk-based
Golib assurance policy. All coverage-required packages remain instrumented,
the complete module test command must succeed, and malformed or missing
coverage profiles fail. A statement percentage does not certify adequacy:
public regression tests, caller coverage and independent review establish
the affected contracts. Defensive cooperative-cancellation checkpoints are
retained rather than removed merely to improve a coverage percentage.

CI builds the pinned development tooling source to support this explicit
policy; it does not claim that the public v1.8.5 executable implements it.
Native mutation, security scanners and the required workflow result still
gate qualification. Historical failures remain failures; source review and
focused tests do not establish release readiness or public consumption.
