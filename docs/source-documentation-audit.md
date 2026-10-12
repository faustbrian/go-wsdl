# Source documentation audit

## Scope and reproduction

Tier A source documentation pass from `fb0ede53d9fc40efdce2174437f9c371b0d84df5`. Reproduce the complete
inventory with `git ls-files "*.go"`, module boundaries with
`git ls-files "*go.mod"`, and assurance entry points with
`rg "^func (Test|Fuzz|Benchmark|Example)" --glob "*_test.go"`.

All 71 tracked Go files in 1 module are in scope, including private
parser/conversion/compiler state, tests, helpers, examples, fuzz and benchmarks.
The root and builder, codegen, compile, compose, diff, resolve and
internal/errorprivacy packages are audited independently. No runtime, dependency,
fixture, pinned specification/provenance, API baseline or workflow change is planned.

## Initial gaps and bounded plan

The native source checker reports 565 missing exported comments, primarily
model fields and presence flags. Missing prose is distinct from unresolved
technical interpretation. WSDL 1.1 (2001-03-15) and WSDL 2.0 (2007-06-26)
authority, decision records and versioned fixtures are retained, not rewritten.

Investigate lexical versus resolved URI identity and resolver policy; XML/charset
budgets and allocation admission; foreign extension replay; shallow model access
versus owned compiler/generation snapshots; diagnostic privacy; context checkpoints
versus collaborator blocking; deterministic serialization versus XML canonicalization;
and builder/composition/diff failure and ownership contracts before writing them.

Selected proof is native pinned formatting, exact executable tokens/AST/directives,
package rendering and existing examples where useful, source/local-doc/specification
structure, links, changed-content secrets and ordinary complete final-diff review.
No broad unchanged runtime, conformance, fuzz, benchmark or interoperability campaign.

API mode is apidiff (`api/v3-baseline.txt`); no full-prose baseline is configured.
Native verification/package.mk defines explicit benchmark/conformance/interoperability
operations, not a docs generator. Imported fixture/tool pins do not authorize execution.

## Complete tracked review inventory

| Go file | Initial documentation gaps | Source or assurance boundary |
| --- | --- | --- |
| `accessors_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `benchmark_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `binding20_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `builder/builder.go` | 1 | borrowed mutable construction and canonical Build |
| `builder/builder_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `builder/fuzz_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `codegen/model.go` | 54 | owned generation projection, post-copy count limits and no emitter |
| `codegen/model_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `codegen/model_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `compile/compile.go` | 60 | exact injected identities, invocation-local walks, immutable graph copies and schema split |
| `compile/compile_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `compile/compile_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `compile/completion_cancellation_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `compile/parser_cancellation_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `compile/resolver_cancellation_contract_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `compile/schema_parse_admission_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `compose/compose.go` | 9 | duplicate identity conflicts, lexical dedup and canonical output |
| `compose/compose_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `compose/compose_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `conversion_cancellation_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `conversion_owner_admission_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `conversion_sibling_admission_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `diagnostic.go` | 7 | retained severity classification versus private default errors |
| `diff/diff.go` | 13 | nil-as-empty structural comparison and wire/application caveats |
| `diff/diff_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `diff/diff_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `error_privacy_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `error_privacy_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `errors.go` | 0 | typed resource/directive/duplicate classifications |
| `example_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `extensions_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `fuzz_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `http20_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `internal/errorprivacy/error.go` | 0 | fixed categories versus borrowed inspectable causes |
| `internal/errorprivacy/error_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `interop_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `model.go` | 26 | canonical constructors and shallow mutable accessor lifetime |
| `model11.go` | 172 | WSDL 1.1 core/SOAP/HTTP/MIME vocabulary and presence |
| `model20.go` | 198 | WSDL 2.0 core/adjunct defaults, legacy first-message views and presence |
| `model_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `parse.go` | 13 | two XML passes, shared charset allowance, nil context and publication checks |
| `parse11.go` | 0 | core/adjunct conversion, overload labels and owned foreign data |
| `parse11_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `parse20.go` | 0 | core/adjunct conversion, extension consumption and per-schema admission |
| `parse20_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `parse_cancellation_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `parse_charset_budget_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `parse_completion_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `parse_metadata_admission_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `parse_strict_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `parse_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `parse_versions_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `resolve/resolve.go` | 7 | borrowed policy collaborators, deny defaults and owned memory copies |
| `resolve/resolve_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `rpc20_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `serialization_cooperation_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `serialize.go` | 3 | deterministic prefix prework versus encoded output limits |
| `serialize11.go` | 0 | style-dependent child order and raw extension/token errors |
| `serialize20.go` | 0 | typed adjunct presence and complete direction/collection emission |
| `serialize_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `serialize_limits_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `serialize_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `soap20_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `soap_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `validate.go` | 2 | local references/styles, retained diagnostic cap and no context budget |
| `validate20_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `validate_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `validate_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `xmltree.go` | 0 | temporary context owner, ephemeral token copying, lexical URI and namespace fanout |
| `xmltree_internal_test.go` | 0 | assurance; cases and helper ownership reviewed |
| `xsd_composition_test.go` | 0 | assurance; cases and helper ownership reviewed |

## Investigated rationale and unresolved questions

- Parse performs Wire root admission and a strict encoding/xml tree pass, sharing
  the document charset allowance; DTDs/directives are rejected, XML predefined
  entities remain distinct from custom entities, and imports only resolve lexical
  xml:base identities. Parse is not complete semantic validation or resource IO.
- Nil Parse context is not generally accepted: valid XML reaches contextReader.Read,
  which calls ctx.Err without a nil guard. A minimal definitions document with nil
  ctx is the concrete boundary; use Background or another real context. Compile
  inherits this boundary. No panic repair or invented nil-classification is included.
- Raw Document accessors return shallow values; slices, schema pointers and nested
  data alias the document. Constructors use Validate/Marshal/Background Parse to
  produce separate mutable storage under defaults. Documentation.Markup is unused
  by Parse/Marshal. Source locations identify decoder positions after start tokens,
  not trusted resource identities or exact source-byte preservation.
- Builders copy top-level values but retain nested inputs until reuse ends. New is
  needed for named additions; nil/zero receivers fail rather than acquire resources.
  Build re-enters XML canonicalization despite the old parsing-free package claim.
- Compiler options copy validation names but borrow resolvers. Concurrent safety
  needs safe collaborators. Invocation-local traversal/cancellation stays separate
  from reusable configuration; callbacks cannot be forcibly preempted. URI identities
  are exact strings and absolute/fragment-free syntax is not an SSRF policy.
  WSDL URI revisits/cycles reuse resources; interface-inheritance cycles are rejected.
  Complete operation identity retains effective WSDL 1.1 input/output labels.
- Compiled Set storage is private; nested accessor collections are independently
  copied. Lookup rejects ambiguous operation names, rather than selecting an overload.
  XSD schema compilation/validation remains separately injected and version-pinned.
- MaxDiagnostics drops later messages without stopping validation. Marshal MaxBytes
  caps encoded output after prefix/model traversal. Codegen limits may be checked
  after accessor copies and model construction. Namespace maps copy inherited entries
  per tree node; lexical attribute count is not an aggregate namespace-copy cap.
  These are not universal heap ceilings or hard deadlines. Caller-mutated/built models
  and collaborator response allocation remain separate resource responsibilities.
- Current WSDL-DEC-006 security/resource consequence prose is clarified to these
  actual boundaries, without changing normative interpretation or selected behavior.
  Its prior digest remains in append-only typed history; only the changed current
  digest is appended using the authoritative Tools decision type and algorithm.
- Serialization chooses stable prefixes/components, preserves operation order and
  optional-property presence, and rejects unsafe raw extension XML. It is not W3C
  XML canonicalization or lexical byte identity. Version-specific unknown extensions
  remain data; typed adjuncts are removed from the opaque bag to avoid double emission.
- Compose rejects any duplicate top-level identity, even equal content. Lexical
  import/include deduplication does not invoke resolvers. Diff compares structural
  graphs (nil as empty) and retains explicit unknown/wire caveats, not conformance.
  Neither API has an independent context or report/input-memory budget.
- Default root/compile/composition error categories hide input text, but diagnostics,
  conflict lists and unwrapped causes are intentionally inspectable and sensitive.
  Internal errorprivacy wraps categories; it does not sanitize arbitrary causes.
- Three fuzz targets protect bounded XML/model/builder partitions; three benchmarks
  name the fixed parse/marshal/compile workload and excluded setup. Tests include
  resource admission, cancellation before publication, owner handoffs, partial write
  failures, exact limits, version-specific adjuncts, inherited/overloaded identity,
  raw extension refusal, caller-copy mutation and injected resolver/fixture boundaries.
  Existing tests/fixtures are reviewed, not rerun as a broad assurance campaign.

No generated/build-tagged/platform-specific Go files or generation/deprecation
markers were found. One public /v3 module has 20 implementation/helper and 51
assurance files in eight source directories. Fixture pins remain in
specification/manifest.tsv, interoperability.tsv and tooling.tsv; requirement and
assertion matrices, external/W3C corpora, module/API pins and dated release records
remain unchanged. Native provenance validation is offline only.

Unresolved future behavior decisions are not documentation gaps: nil-context
classification, strict pre-allocation namespace/model budgets, interruptible validation/
generation and arbitrary resolver calls, markup preservation and mutable raw-model
snapshot APIs need separately scoped requirements and tests. No rationale is invented
for those limitations. All current contracts receive structural/rendered review;
no new conformance, interoperability, performance or release claim is made.

## Objective Enforcement

The checksum-pinned public Tools v2.4.0 verifier now supplies
`make source-docs-check` as a prerequisite of the existing local and CI Make
routes. Ordinary shared CI selects strict source-documentation checking and
requires its result. Native prerequisites, operation budgets, security and
release controls are retained. Hosted exact-head evidence is recorded by the
coordination ledger; configuration alone does not establish green delivery.
