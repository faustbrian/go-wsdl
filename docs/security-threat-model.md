# WSDL security threat model

Model version: 2.0.0. Reviewed scope: public v3.0.0, 2026-10-08.
Owner: Brian Faust, WSDL maintainer. The signed v3.0.0 release is published;
actual clean public consumption includes XSD v2 composition.

## Assets and boundaries

Untrusted XML, names, URIs, extensions and embedded schemas enter `Parse`.
Caller-created models enter validation, serialization, composition and code
generation. Import graphs enter `compile`; injected WSDL and XSD resolvers are
trusted application collaborators, not library-owned network implementations.
Wire v3.0.0 owns XML token/encoding boundaries; published XSD v2.0.0 owns schema
compilation. WSDL does not authenticate users, execute services, generate
executable clients, manage credentials or perform cryptographic operations.

Assets are process availability, model integrity, immutable compiled state and
confidential input-derived diagnostics. Parsing opens no network connection or
file and performs no implicit import resolution. A nil compiler resolver denies
loads. No hidden worker, global registry or environment lookup is authorized.

## Owned controls and evidence

`ParseOptions` bounds document bytes, XML depth, elements, attributes, text,
schemas, imports and model collections. Both XML passes apply the normalized
document-byte quota to charset conversion. DTDs/directives and custom entity
processing are rejected. `MarshalOptions` bounds emitted bytes. Compiler
`Limits` bounds documents, graph depth, references, cumulative bytes and
components; `SchemaLimits` delegates explicit schema bounds. Code-generation
`Limits` bounds each materialized model collection. Defaults are finite;
applications must choose deployment-appropriate limits rather than treat them
as a process-wide memory or CPU reservation.

Owned conversion, graph and reference loops observe the caller context and
refuse successful partial publication after observed cancellation. Compilers
keep invocation state separate so cancellation does not poison later reuse.
Returned models do not retain the parse context. Standard-library sorting,
escaping, model validation, resolver callbacks and XSD calls are cooperative
boundaries, not preemptive interruption or a cancellation-latency promise.

Default library-produced errors are categorical. Explicit diagnostics, conflict
details and unwrapped causes can contain untrusted data and are an opt-in
inspection boundary. This is a breaking v2 text contract, not a public v1 fix.
Applications must match categories/types rather than error strings.

Ordinary regression coverage includes XML and charset admission, hostile model
limits, contextual refusal and completion, resolver fallback, compiler reuse,
privacy, semantic round trips and independently owned output. CI retains native
mutation and scanner gates. These controls and tests do not by themselves prove
a public release or an arbitrary application's I/O policy.

## Conditional residual risks

### Published v3 supplier boundary

Public v3.0.0 adopts XSD v2 named types; the previous model's v2 scope remains
historical. Both direct inline parsers and
compiler root/import first parses forward caller-selected finite namespace and
model-copy allowances. The inline schema resolver keeps XSD v2's independent
finite constructor defaults; graph allowances do not silently enlarge them.
Zero namespace and model-copy allowances select finite defaults of 1,000,000
entries and 64 MiB; negative allowances are invalid. Inline resolver defaults
are independently limited to 256 resources, 64 MiB cumulative identity/content
bytes and 16 MiB per resource, inclusive. These work allowances do not measure
exact heap use or bound prior caller allocations.
The maintainer owns adoption, regression coverage and public dependency/release
qualification; revisit on any supplier or parser-boundary change. Applications
still own bounded cooperative custom resolvers and trusted diagnostic use.

| Risk | Owner and rationale | Mitigation and review condition |
| --- | --- | --- |
| Blocking or excessive collaborator work | Application owner; Go cannot preempt synchronous trusted resolver, XSD or standard-library operations. | Use bounded cooperative implementations and deployment budgets; revisit when a supplier or callback contract changes. |
| SSRF, filesystem access or credential forwarding in a supplied resolver | Application resolver owner; external access is explicit and denied without a resolver. | Allowlist resources, bound redirects/bytes/time, validate resolved addresses and avoid credential forwarding; review each new resolver or transport. |
| Sensitive opt-in diagnostics | Application observability owner; detail is explicitly inspectable, not implicitly printed by WSDL wrappers. | Authorize/redact before logging, rendering or tracing; revisit when adding an error or diagnostic surface. |
| Mutation of caller-owned models | Application owner; mutable models/builders require documented single-owner use, unlike immutable compiled lookup. | Avoid concurrent mutation and shared mutable aliases; revisit when changing cloning, builder or model ownership. |
| Supplier and publication compromise | Maintainer; public Wire/XSD and pinned tooling are distinct trust boundaries. | Verify source/dependency pins, scanners and signed tags/assets before publication, then actual clean public consumption before qualifying the release; revisit on any source, dependency or security advisory change. |

These are conditional use obligations, not acceptance of a known High finding.
Unknown scanner or release outcomes remain unqualified. Report suspected defects
through the [private reporting policy](../SECURITY.md); do not place credentials
or private reporter information in source, fixtures or public CI artifacts.
