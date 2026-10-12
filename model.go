// Package wsdl parses and emits WSDL 1.1 (15 March 2001 Note) and WSDL 2.0
// (26 June 2007 Recommendations) description models using Wire XML and XSD.
// Start with Parse, then Validate for local semantics or package compile for
// an explicitly resolved graph. Parsing does no resource loading; SOAP/HTTP
// bindings are data, not service execution or endpoint authorization.
//
// Documents expose mutable nested models through shallow accessors. Keep all
// borrowed source/model data stable during operations and synchronize mutation.
// Limits are phase-specific, not universal heap or hard-deadline guarantees.
// Default wrapped errors hide input text; diagnostics and causes remain sensitive.
// No file/network acquisition, cgo or background workers are owned here.
package wsdl

import (
	"context"

	"github.com/faustbrian/go-wsdl/v3/internal/errorprivacy"
	xsd "github.com/faustbrian/go-xsd/v2"
)

// Version identifies a WSDL language version, or SOAP 1.1/1.2 in SOAP model
// fields. Version12 is not a supported WSDL document language.
type Version string

const (
	// Version11 identifies WSDL 1.1 documents.
	Version11 Version = "1.1"
	// Version20 identifies WSDL 2.0 documents.
	Version20 Version = "2.0"
)

// Definitions11 is the root of a WSDL 1.1 description.
type Definitions11 struct {
	// Name is the optional WSDL 1.1 definitions name, not a component QName.
	Name string
	// TargetNamespace is the namespace assigned to this document's top-level component
	// identities.
	TargetNamespace string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Imports contains external namespace references; parsing does not resolve resources.
	Imports []Import11
	// Types contains XML Schema declarations or names; schema semantics remain owned by XSD.
	Types *Types11
	// Messages contains WSDL 1.1 abstract message declarations.
	Messages []Message11
	// PortTypes contains WSDL 1.1 abstract interface declarations.
	PortTypes []PortType11
	// Bindings contains concrete protocol-description bindings, not runtime adapters.
	Bindings []Binding11
	// Services contains service declarations and their described endpoints.
	Services []Service11
	// Extensions contains preserved foreign elements; required semantics need explicit
	// understanding.
	Extensions []Extension
	// ExtensionAttributes contains expanded foreign attribute names and lexical values.
	ExtensionAttributes []ExtensionAttribute
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Types11 contains the XML Schema documents embedded in WSDL 1.1 types.
// XML Schema parsing and compilation remain owned by xsd.
type Types11 struct {
	Extensibility
	// Schemas contains inline XSD document pointers; shallow WSDL model access does not
	// clone their storage.
	Schemas []*xsd.Document
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Description20 is the root of a WSDL 2.0 description.
type Description20 struct {
	Extensibility
	// TargetNamespace is the namespace assigned to this document's top-level component
	// identities.
	TargetNamespace string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Imports contains external namespace references; parsing does not resolve resources.
	Imports []Import20
	// Includes contains same-namespace WSDL 2.0 references; parsing performs no load.
	Includes []Include20
	// Types contains XML Schema declarations or names; schema semantics remain owned by XSD.
	Types *Types20
	// Interfaces contains abstract interface declarations, not executable clients.
	Interfaces []Interface20
	// Bindings contains concrete protocol-description bindings, not runtime adapters.
	Bindings []Binding20
	// Services contains service declarations and their described endpoints.
	Services []Service20
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Types20 contains the XML Schema documents embedded in WSDL 2.0 types.
// XML Schema parsing and compilation remain owned by xsd.
type Types20 struct {
	Extensibility
	// Imports contains external namespace references; parsing does not resolve resources.
	Imports []xsd.SchemaReference
	// Schemas contains inline XSD document pointers; shallow WSDL model access does not
	// clone their storage.
	Schemas []*xsd.Document
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Document normally contains one version-specific model. Its zero value has no
// version and cannot be marshaled. Accessors return shallow copies: nested slices,
// pointers and schema documents still alias it. Concurrent reads require no mutation.
type Document struct {
	version       Version
	definitions11 *Definitions11
	description20 *Description20
}

var (
	canonicalMarshal  = Marshal
	canonicalParse    = Parse
	canonicalValidate = Validate
)

// NewDocument11 validates, marshals, reparses and revalidates a borrowed WSDL 1.1
// model into separately owned mutable storage. Keep nested input stable until return.
// Canonicalization uses default output/parse budgets and a Background context,
// may discard lexical/source-only details, and does not resolve external resources.
func NewDocument11(
	value Definitions11,
	validation ValidationOptions,
) (*Document, error) {
	return canonicalDocument(
		&Document{version: Version11, definitions11: &value},
		validation,
	)
}

// NewDocument20 validates, marshals, reparses and revalidates a borrowed WSDL 2.0
// model into separately owned mutable storage under default parse/output budgets.
// It uses Background, performs no resolution, and may discard source-only details.
func NewDocument20(
	value Description20,
	validation ValidationOptions,
) (*Document, error) {
	return canonicalDocument(
		&Document{version: Version20, description20: &value},
		validation,
	)
}

func canonicalDocument(
	document *Document,
	validation ValidationOptions,
) (*Document, error) {
	diagnostics := canonicalValidate(document, validation)
	if err := diagnostics.Err(); err != nil {
		return nil, errorprivacy.Wrap("wsdl: validate model", err)
	}
	payload, err := canonicalMarshal(document, MarshalOptions{})
	if err != nil {
		return nil, errorprivacy.Wrap("wsdl: canonicalize model", err)
	}
	canonical, err := canonicalParse(context.Background(), payload, ParseOptions{})
	if err != nil {
		return nil, errorprivacy.Wrap("wsdl: canonicalize model", err)
	}
	diagnostics = canonicalValidate(canonical, validation)
	if err := diagnostics.Err(); err != nil {
		return nil, errorprivacy.Wrap("wsdl: validate model", err)
	}
	return canonical, nil
}

// Description20 returns a shallow model copy, or zero/false for nil/wrong-version
// documents. Nested storage aliases the document; it is not an immutable snapshot.
func (d *Document) Description20() (Description20, bool) {
	if d == nil || d.description20 == nil {
		return Description20{}, false
	}
	return *d.description20, true
}

// Version reports the parsed WSDL version.
func (d *Document) Version() Version {
	if d == nil {
		return ""
	}
	return d.version
}

// Definitions11 returns a shallow model copy, or zero/false for nil/wrong-version
// documents. Nested storage aliases the document; callers synchronize mutations.
func (d *Document) Definitions11() (Definitions11, bool) {
	if d == nil || d.definitions11 == nil {
		return Definitions11{}, false
	}
	return *d.definitions11, true
}
