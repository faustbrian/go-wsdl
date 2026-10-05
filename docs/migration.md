# Migration

The module has a stable v1 API. Pin a released v1 version, review the
[changelog](../CHANGELOG.md), and evaluate semantic diffs when upgrading.
Construct documents through `NewDocument11` or `NewDocument20`; do not depend
on lexical namespace prefixes or source attribute order.

Older parsers that loaded imports implicitly must move loading behind an
injected resolver. A nil resolver now means deny. Callers that treated all
WSDL versions as one model must branch on `Document.Version` and use the
version-specific accessor.

Generated-client users should compile a set and consume package `codegen`
rather than coupling generators to parser internals.

## Pending next-major default-error contract

The unreleased implementation returns categorical library-produced default
error text rather than embedding input URIs, component names, lexical values
or collaborator causes. Replace text matching with `errors.Is`/`errors.As`.
Inspect `Diagnostics` and `ConflictError.Conflicts` explicitly when detailed
diagnostics are required. Wrapped causes remain available through standard
unwrapping; redact or authorize that opt-in detail before recording it.

This is not a v1 patch contract. A future major release must adopt the official
`github.com/faustbrian/go-wsdl/v2` module and import suffix from main, without
version-specific source directories or branches. That migration and publication
have not occurred; published v1.0.0 retains the previous default text.
