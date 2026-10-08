# Migration

Public v3.0.0 is available from main. Pin a released version, review the
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

Published v2.0.0 returns categorical library-produced default
error text rather than embedding input URIs, component names, lexical values
or collaborator causes. Replace text matching with `errors.Is`/`errors.As`.
Inspect `Diagnostics` and `ConflictError.Conflicts` explicitly when detailed
diagnostics are required. Wrapped causes remain available through standard
unwrapping; redact or authorize that opt-in detail before recording it.

This is not a v1 patch contract. Public v2.0.0 uses the official
`github.com/faustbrian/go-wsdl/v2` module and import suffix, without
version-specific source directories or branches. Update root and subpackage
imports together. Public v2.0.0 is available; published v1.0.0
retains the previous default text.

The XML supplier is `github.com/faustbrian/go-wire/v3 v3.0.0`. Applications
inspecting unwrapped Wire diagnostics must use its `/v3` types and sentinels.
Public WSDL v2.0.0 retains `github.com/faustbrian/go-xsd v1.0.0` types.

## V3 XML Schema adoption

Published v3.0.0 uses `github.com/faustbrian/go-wsdl/v3` with
`github.com/faustbrian/go-xsd/v2 v2.0.0`, without version-specific source
directories or branches. Both releases are publicly available; clean public
consumer checks exercise their composition. Update every WSDL and XSD
root/subpackage import together.
`Types11.Schemas`, `Types20.Schemas`, `Types20.Imports`, `SchemaResolver`,
`SchemaLimits` and `Set.Schemas()` expose XSD v2 named types. Old-major values
are not assignable; reparse or explicitly migrate application-owned models.
Wire remains `/v3` and categorical WSDL default errors remain unchanged.

`ParseOptions.MaxSchemaNamespaceEntries` and `MaxSchemaModelBytes` bound each
inline schema's owned namespace and model-string/copy work. Zero selects XSD's
finite defaults of 1,000,000 entries and 64 MiB; negative limits are invalid.
Compiler `SchemaLimits.MaxParseNamespaceEntries` and `MaxParseModelBytes`
apply during the first root/import WSDL parse as well as later XSD compilation.
These per-schema allowances do not replace WSDL byte, tree or schema-count
limits and do not measure exact heap use or bound prior caller allocations.

The inline schema resolver retains independent XSD v2 constructor defaults:
256 resources, 64 MiB of identity plus content bytes, and 16 MiB per resource,
inclusive. Increasing graph limits does not raise these constructor limits.
Limit refusal returns no compiled set and retains the relevant XSD parser,
compiler or resolver `ErrLimitExceeded` through `errors.Is`. Review larger
workloads and split them into appropriately bounded compilation units rather
than assuming XSD v1 admission behavior remains supported. Explicit XSD
diagnostics and unwrapped causes remain trusted inspection surfaces.
