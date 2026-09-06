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
