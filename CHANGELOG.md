# Changelog

## Unreleased

### Changed

- Adopt the `go-library-tools` v1.4.0 schema-v2 cohesion contract, local
  `make cohesion` gate, and immutable reusable CI workflow without changing
  the WSDL API or runtime behavior.

- Resolve the existing `go-wire` and `go-xsd` v1.0.0 requirements from their
  public Go proxy and checksum-database identities after the shared bootstrap
  proxy was made fallback-only.

- Adopt checksum-pinned `go-library-tools` v1.2.0 and its immutable workflow
  so CI executes specification governance while keeping WSDL conformance,
  interoperability, mutation, and benchmark evidence in this repository.

### Specification Decisions

- WSDL11-DEC-001 sha256:d2f247795494ca3558fbc9474f453030b838be12a20d2e28a4f135bc932866e1: Govern [The 15 March 2001 Note is the WSDL 1.1 baseline](docs/specification-decisions.md).
- WSDL11-DEC-002 sha256:59ad4cff0b8f49e0dc7a034e782527fcfd19d54d421331c85905bfb6deef58ae: Govern [Operation child order determines operation style](docs/specification-decisions.md).
- WSDL11-DEC-003 sha256:9d4a3c6d8df3f63859cd013ca822e8a71d8c8d5a05f99038fc81b5decd0f1c0e: Govern [Binding extensions are description data, not runtime clients](docs/specification-decisions.md).
- WSDL11-DEC-004 sha256:58182321af2108468a86dc25080a3c70af5c2708c37795e8c1f2568aa30958f3: Govern [Overloaded operations use complete message identity](docs/specification-decisions.md).
- WSDL20-DEC-001 sha256:b48502d1aa2100c4ed2ee9295240db82f5816b30842b055f06776a6ef098201d: Govern [The 26 June 2007 Recommendations are the WSDL 2.0 baseline](docs/specification-decisions.md).
- WSDL20-DEC-002 sha256:676968bbfb774d5a979d9c2e0f6806951da6c83530adc9196eca85252024dbf7: Govern [Absent and explicitly defaulted values remain distinguishable](docs/specification-decisions.md).
- WSDL20-DEC-003 sha256:33bb6ac1db3a47601c5dc0f8722c216e2fabb2816892513c1fa00400dba0ba99: Govern [Unknown absolute message exchange patterns remain extensible](docs/specification-decisions.md).
- WSDL20-DEC-004 sha256:08c8510e8d8caf35559a3dab5b0e88c29d0fafc438e1418c9370489bfc7a66bc: Govern [RPC validation is split between WSDL and compiled XML Schema](docs/specification-decisions.md).
- WSDL20-DEC-005 sha256:82129b688cdf5ebb1b2dcf266e8d8ef52666d71f09a53dfd247dfea8d30455ab: Govern [Operation safety accepts both published spellings and emits the normative one](docs/specification-decisions.md).
- WSDL20-DEC-006 sha256:ce7fe8de75f4e9cf9cb8d88f2e230f02eba239e7b39eadaee56533cb70025d1a: Govern [Operation style validation is split between WSDL and XML Schema](docs/specification-decisions.md).
- WSDL-DEC-001 sha256:af518f4f24e5e8b3fcc406ee5c959e6530e4d8ad8c2fd8bf962ae56c439bf0a4: Govern [Normative prose outranks schemas, examples, fixtures, and peers](docs/specification-decisions.md).
- WSDL-DEC-002 sha256:609be5a66fcf307b118a81617b0dd55895da772ab19f1217181c31fd8c68d0d9: Govern [Parsing performs no external resolution](docs/specification-decisions.md).
- WSDL-DEC-003 sha256:42b75d33065fd5546fe9296a74917777b8ebd30126074e16767b6acc532d8812: Govern [DTDs, directives, and custom entity processing are rejected](docs/specification-decisions.md).
- WSDL-DEC-004 sha256:dbb7598bc1bfcedc75f2c32c9821ede43c4a0f94b09a0623cde77c67c9432ecb: Govern [Expanded QNames own identity and extensions remain opaque unless understood](docs/specification-decisions.md).
- WSDL-DEC-005 sha256:a29a8c7703e0a759efac5b623da9713351e222f28b2b9923d88dee3fed859ea6: Govern [Serialization is deterministic and semantic, not lexically preserving](docs/specification-decisions.md).
- WSDL-DEC-006 sha256:c69d6e6bc3f5a0c0bd2dc1844ac5547bf6fa8eda890ccce211e6da3d7ac5ee10: Govern [All parsing, compilation, validation, and output work is explicitly bounded](docs/specification-decisions.md).
- WSDL-DEC-007 sha256:a21e739d2a7560c3a1adf76450148de87b6aaaf635506b1988279e2bebe611ec: Govern [Conformance and interoperability claims remain attributable](docs/specification-decisions.md).
- WSDL-DEC-008 sha256:fc0330be7640d488fd01a99f39ea6c7c6f660ab7cb3b5dbc13ffa3891d99a15e: Govern [The package models descriptions and never owns service execution](docs/specification-decisions.md).

### Documentation

- Record the public module identity review for the WSDL module's two owned
  dependencies.

- Publish the module's family, capabilities, ownership, lifecycle, supported
  environments, package-selection guidance, and immutable ecosystem links.

- Remove the archived monorepo documentation link; package guidance remains in
  the repository-owned documentation.

## 1.0.0 - 2026-08-25

### Changed

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Documentation

- Link the package README to the repository-wide Golib documentation portal.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-wsdl` identity while preserving its documented API and behavior.
- Validate benchmark output with the standard shell toolchain so clean Linux
  CI runners do not require an undeclared ripgrep installation.
- Consolidate WSDL 1.1, WSDL 2.0, and cross-version interpretation policies in
  a canonical decision register with stable identifiers, normative sources,
  consequences, and executable evidence.
- Delegate local mutation checks to the canonical exact-100 repository runner
  instead of a reduced package-local efficacy threshold.
- Require owned sibling modules at local `v0.0.0`; clean external consumers
  pin each module to an exact main pseudo-version.

### Distribution

- Include the canonical MIT licence in the independently published module.

### Compatibility

- Preserve the provenance-pinned W3C SOAP encoding schema byte-for-byte in
  source archives and clean clones.
- Check temporary interoperability fixture cleanup and normalize XML Schema
  extension detection to satisfy the canonical strict lint contract.
- Added a pinned module export baseline so incompatible public API changes
  fail the canonical repository gate.
- Separate W3C fixture conformance from Apache Woden interoperability so both
  results remain attributable.

- Refresh owned-module checksums against the final consolidated archives.
- Normalize standalone module metadata against the canonical owned dependency
  graph, including complete checksums for clean consumer resolution.
- Refresh the canonical XSD checksum after its API compatibility tooling was
  standardized.
- Add bounded WSDL 1.1 and WSDL 2.0 parsing and deterministic serialization.
- Add SOAP, HTTP, MIME, extension, validation, and XML Schema integration.
- Preserve the WSDL 2.0 `wsdlx:safe` operation property through parsing,
  compilation, generation, serialization, and semantic comparison.
- Add typed WSDL 2.0 `wrpc:signature` support and XSD-backed RPC style
  validation.
- Add XSD-backed WSDL 2.0 IRI and multipart operation-style validation.
- Add explicit resolution, immutable compilation, builders, composition,
  semantic diff, and a separate bounded code-generation model.
- Add provenance matrices, security documentation, fuzzing, benchmarks,
  coverage, mutation, race, and CI gates.
- Add provenance-pinned accepted W3C WSDL 2.0 interoperability fixtures.
- Add licensed SoapUI, `dotnet-svcutil`, Apache CXF, and DHL WSDL 1.1
  interoperability fixtures with executable round-trip evidence.
- Preserve descendant-local namespaces in embedded schemas and reuse canonical
  schema prefixes during deterministic WSDL serialization.
- Add a pinned Apache Woden Java gate for the W3C WSDL 2.0 fixture corpus.
- Run the Woden gate in a digest-pinned Eclipse Temurin container instead of
  depending on an untracked host Java installation.
