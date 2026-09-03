# Builders and composition

`builder.New11` and `builder.New20` provide zero-I/O programmatic construction.
Add methods reject duplicate top-level component names. `Build` delegates to
the canonical document constructors, so invalid references or binding
properties cannot be emitted as a successful document.

Builder inputs, including their nested slices and maps, remain caller-owned and
borrowed for as long as the builder may be reused. Callers must keep that
storage stable between builder calls and must not mutate it concurrently with
builder use. Each successful `Build` returns a document that owns a separate
mutable snapshot, so later caller mutation cannot change an already built
document.

`compose.Merge` combines same-version documents with the same target namespace.
It sorts components and extensions deterministically, deduplicates imports,
includes, and schemas, and returns a structured `ConflictError` for components
with the same identity but different content. Inputs remain caller-owned.

Composition does not resolve imports and does not silently merge different
namespaces or WSDL versions.
