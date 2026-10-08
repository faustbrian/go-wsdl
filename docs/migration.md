# Migration

Main prepares the v2 API; its public release remains pending. Pin a released
version, review the
[changelog](../CHANGELOG.md), and evaluate semantic diffs when upgrading.
Construct documents through `NewDocument11` or `NewDocument20`; do not depend
on lexical namespace prefixes or source attribute order.

Older parsers that loaded imports implicitly must move loading behind an
injected resolver. A nil resolver now means deny. Callers that treated all
WSDL versions as one model must branch on `Document.Version` and use the
version-specific accessor.

Generated-client users should compile a set and consume package `codegen`
rather than coupling generators to parser internals.

## V2 module and default-error contract

The unreleased implementation returns categorical library-produced default
error text rather than embedding input URIs, component names, lexical values
or collaborator causes. Replace text matching with `errors.Is`/`errors.As`.
Inspect `Diagnostics` and `ConflictError.Conflicts` explicitly when detailed
diagnostics are required. Wrapped causes remain available through standard
unwrapping; redact or authorize that opt-in detail before recording it.

This is not a v1 patch contract. Main now uses the official
`github.com/faustbrian/go-wsdl/v2` module and import suffix, without
version-specific source directories or branches. Update root and subpackage
imports together. Public v2 publication remains pending; published v1.0.0
retains the previous default text.

The XML supplier is `github.com/faustbrian/go-wire/v3 v3.0.0`. Applications
inspecting unwrapped Wire diagnostics must use its `/v3` types and sentinels.
The XML Schema supplier remains `github.com/faustbrian/go-xsd v1.0.0`;
this migration does not change that public type identity.
