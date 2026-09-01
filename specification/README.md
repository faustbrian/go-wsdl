# Specification provenance

`manifest.tsv` pins the stable WSDL 1.1 Note, WSDL 2.0 Recommendation,
adjuncts, additional message exchange patterns, and published schemas used by
this project. Dated W3C URLs are preferred. Namespace documents are snapshots:
their recorded digest and size define the reviewed input.

The two requirement matrices are independent because WSDL 1.1 and WSDL 2.0
have different component models and normative status. A row is `implemented`
only when its evidence names an executable test or gate. `partial` and
`missing` rows are not conformance claims.

`assertions/wsdl-2.0.tsv` inventories all 84 Core and 110 Adjuncts assertion
identifiers from the Recommendations and maps each to a requirement row.
`assertions/wsdl-1.1.tsv` inventories the 23 distinct normative-keyword
statements in the 2001 Note by section and normalized-text digest because the
Note does not assign assertion identifiers. The provenance gate rejects
duplicates, missing rows, count drift, malformed sources, and unknown groups.

The published WSDL 2.0 errata page reports no substantive or editorial errata.
The WSDL 1.1 Note has no W3C errata document; implementation decisions record
the reviewed prose/schema discrepancies instead of implying corrections that
were never published.

`discrepancies-and-extensions.md` records those discrepancies and enumerates
every modeled foreign-attribute and foreign-element extension boundary.
The canonical [`docs/specification-decisions.md`](../docs/specification-decisions.md)
register records the selected observable behavior, alternatives, consequences,
and executable evidence. Version-specific paths under `decisions/` are stable
compatibility pointers to that register.

Run `make provenance` for offline structural checks. Set `VERIFY_REMOTE=1` to
download and hash every pinned resource. Parsing and compilation never use
these URLs implicitly.

`interoperability.tsv` pins every external corpus file by producer, license,
source revision, local path, digest, size, and URL. Offline provenance checks
verify local bytes; remote verification checks the upstream objects too.
`tooling.tsv` pins downloaded interoperability tools and their dependencies by
version, license, digest, size, and immutable artifact URL.

## Decision conformance matrix

| Decision | Title | Evidence |
| --- | --- | --- |
| WSDL11-DEC-001 | The 15 March 2001 Note is the WSDL 1.1 baseline | `TestParseRecognizesWSDL11Definitions`, `TestParseWSDL11CoreAndSOAPDescription`, `TestExternalWSDL11InteroperabilityCorpus` |
| WSDL11-DEC-002 | Operation child order determines operation style | `TestOperationStyle11RecognizesEveryMessageOrder`, `TestParseWSDL11PreservesSolicitResponseOperationOrder`, `TestWSDL11SolicitResponseSerializesOutputBeforeInput` |
| WSDL11-DEC-003 | Binding extensions are description data, not runtime clients | `TestParseWSDL11HTTPAndMIMEBindings`, `TestValidateWSDL11BindingProtocolProperties`, `TestSOAP12HeadersFaultsAndActionPresenceRoundTrip` |
| WSDL11-DEC-004 | Overloaded operations use complete message identity | `TestValidateWSDL11ResolvesOverloadedBindingOperation`, `TestValidateWSDL11RejectsAmbiguousBindingOperation`, `TestCompilerPreservesWSDL11OverloadedOperationIdentity` |
| WSDL20-DEC-001 | The 26 June 2007 Recommendations are the WSDL 2.0 baseline | `TestParseRecognizesWSDL20Description`, `TestValidateWSDL20PredefinedMessageExchangePatterns`, `TestAcceptedW3CWSDL20FixturesParseCompileAndRoundTrip` |
| WSDL20-DEC-002 | Absent and explicitly defaulted values remain distinguishable | `TestWSDL20MessageContentModelsRoundTrip`, `TestWSDL20OperationSafetyRoundTripsWithPresence`, `TestCompareMessageCoversPresencePropertiesAndParts` |
| WSDL20-DEC-003 | Unknown absolute message exchange patterns remain extensible | `TestValidateWSDL20PredefinedMessageExchangePatterns`, `TestValidateWSDL20CustomPatternMessageLabels`, `TestWSDL20InitialMessageFollowsPredefinedMEPDirection` |
| WSDL20-DEC-004 | RPC validation is split between WSDL and compiled XML Schema | `TestValidateWSDL20RPCStyleRules`, `TestCompilerRejectsInvalidWSDL20RPCSchemas`, `TestCompilerPreservesWSDL20RPCSignature` |
| WSDL20-DEC-005 | Operation safety accepts both published spellings and emits the normative one | `TestWSDL20OperationSafetyRoundTripsWithPresence`, `TestWSDL20RejectsInvalidBooleanAtEveryAdjunctBoundary`, `TestCompareDetectsWSDL20OperationSafetyChange` |
| WSDL20-DEC-006 | Operation style validation is split between WSDL and XML Schema | `TestValidateWSDL20IRIAndMultipartStylesRequireElementInitialMessage`, `TestCompilerValidatesWSDL20IRIAndMultipartStyleSchemas`, `TestIRISimpleTypeRulesCoverInlineAndBuiltInTypes` |
| WSDL-DEC-001 | Normative prose outranks schemas, examples, fixtures, and peers | `TestAcceptedW3CWSDL20FixturesParseCompileAndRoundTrip`, `TestExternalWSDL11InteroperabilityCorpus`, `TestParseRejectsMalformedWSDL20AtEveryNestedDecoderBoundary` |
| WSDL-DEC-002 | Parsing performs no external resolution | `TestParseWSDL11ImportResolvesURIWithoutLoading`, `TestParseWSDL20ImportAndIncludeResolveURIsWithoutLoading`, `TestCompilerDefaultsToDeniedResolution` |
| WSDL-DEC-003 | DTDs, directives, and custom entity processing are rejected | `TestParseRejectsGeneralXMLAndResourceBoundaries`, `TestSerializationHelpersRejectUnboundNamesAndUnsafeRawXML` |
| WSDL-DEC-004 | Expanded QNames own identity and extensions remain opaque unless understood | `TestParsePreservesWSDL11ExtensionElementsAndAttributes`, `TestWSDL20ExtensionsRoundTripAcrossComponents`, `TestValidateRequiresExplicitExtensionUnderstanding` |
| WSDL-DEC-005 | Serialization is deterministic and semantic, not lexically preserving | `TestMarshalWSDL11IsDeterministicAndRoundTrips`, `TestMarshalWSDL20IsDeterministicAndRoundTrips`, `TestSerializerAssignsTargetAndSchemaPreferredPrefixesForBothVersions` |
| WSDL-DEC-006 | All parsing, compilation, validation, and output work is explicitly bounded | `TestParseEnforcesComponentLimitsBeforeModelConstruction`, `TestCompilerAcceptsGraphsAtEveryExactResourceLimit`, `TestMarshalLimitAppliesAtEveryOutputBoundary` |
| WSDL-DEC-007 | Conformance and interoperability claims remain attributable | `TestExternalWSDL11InteroperabilityCorpus`, `TestAcceptedW3CWSDL20FixturesParseCompileAndRoundTrip` |
| WSDL-DEC-008 | The package models descriptions and never owns service execution | `TestBuildCreatesOwnedDeterministicGenerationModel`, `TestCompilerUsesOnlyInjectedSchemaResolver`, `TestMemoryReturnsOwnedResourceCopies` |
