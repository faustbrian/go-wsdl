package wsdl

const (
	// NamespaceWSDL11 is the namespace of WSDL 1.1 core elements.
	NamespaceWSDL11 = "http://schemas.xmlsoap.org/wsdl/"
	// NamespaceXMLSchema is the namespace of XML Schema 1.0 and 1.1.
	NamespaceXMLSchema = "http://www.w3.org/2001/XMLSchema"
	// NamespaceSOAP11Binding is the WSDL 1.1 SOAP 1.1 extension namespace.
	NamespaceSOAP11Binding = "http://schemas.xmlsoap.org/wsdl/soap/"
	// NamespaceSOAP12Binding is the WSDL 1.1 SOAP 1.2 extension namespace.
	NamespaceSOAP12Binding = "http://schemas.xmlsoap.org/wsdl/soap12/"
	// NamespaceHTTPBinding is the WSDL 1.1 HTTP extension namespace.
	NamespaceHTTPBinding = "http://schemas.xmlsoap.org/wsdl/http/"
	// NamespaceMIMEBinding is the WSDL 1.1 MIME extension namespace.
	NamespaceMIMEBinding = "http://schemas.xmlsoap.org/wsdl/mime/"
)

// QName is an expanded XML qualified name.
type QName struct {
	// Namespace is the namespace URI used in expanded identity, not a lexical prefix.
	Namespace string
	// Local is the local NCName portion of an expanded XML name, without a prefix.
	Local string
}

// ExtensionAttribute preserves an attribute outside the owning WSDL
// vocabulary as an expanded name and lexical value.
type ExtensionAttribute struct {
	// Name is the expanded component identity; lexical prefixes are not identity.
	Name QName
	// Value is the attribute's lexical text, not an interpreted extension value.
	Value string
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Extension preserves one foreign-namespace extension element.
type Extension struct {
	// Name is the expanded component identity; lexical prefixes are not identity.
	Name QName
	// Attributes contains preserved expanded attribute names and lexical values.
	Attributes []ExtensionAttribute
	// Required records an explicit extension/module/header requirement, not package support
	// for it.
	Required bool
	// RequiredSet tracks explicit Required presence independently of its effective value.
	RequiredSet bool
	// XML contains opaque extension XML replayed through the serializer's XML safety checks.
	XML []byte
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Extensibility preserves foreign attributes and child elements attached to
// a WSDL component.
type Extensibility struct {
	// Extensions contains preserved foreign elements; required semantics need explicit
	// understanding.
	Extensions []Extension
	// ExtensionAttributes contains expanded foreign attribute names and lexical values.
	ExtensionAttributes []ExtensionAttribute
}

// Location identifies an XML source position.
type Location struct {
	// SystemID is the caller's diagnostic source identifier; it can contain sensitive data.
	SystemID string
	// Line is the decoder-reported source line after the opening element token.
	Line int
	// Column is the decoder-reported column after the opening element token.
	Column int
	// Offset is the decoder byte offset after the opening element token.
	Offset int64
}

// Documentation preserves a WSDL documentation element.
type Documentation struct {
	// Language records xml:lang text; it does not perform language negotiation.
	Language string
	// Content contains trimmed direct character data; nested markup is not a preserved XML
	// payload.
	Content string
	// Markup is not populated by Parse or emitted by Marshal; do not rely on markup round
	// trips.
	Markup string
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Import11 references a WSDL 1.1 document for another namespace. URI is the
// resolved identity; parsing never loads it.
type Import11 struct {
	Extensibility
	// Namespace is the namespace URI used in expanded identity, not a lexical prefix.
	Namespace string
	// Location retains the lexical URI/address text, not a resolved or authorized resource.
	Location string
	// URI is the resource identity, not a filesystem path or network authorization decision.
	URI string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Source records where the lexical URI/address was decoded.
	Source Location
}

// Part11 is one abstract WSDL 1.1 message part.
type Part11 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Element is the expanded XML Schema element reference; zero means no named element.
	Element QName
	// Type is the expanded XML Schema type reference; zero means no named type.
	Type QName
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Message11 is a named WSDL 1.1 abstract message.
type Message11 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Parts contains the WSDL 1.1 message parts in model order.
	Parts []Part11
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// OperationMessage11 references a message used by an abstract operation.
type OperationMessage11 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Message is the expanded WSDL 1.1 abstract message reference.
	Message QName
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// OperationStyle11 identifies the message ordering of a WSDL 1.1 operation.
type OperationStyle11 string

const (
	// OperationStyleOneWay records an input-only WSDL 1.1 operation.
	OperationStyleOneWay OperationStyle11 = "one-way"
	// OperationStyleRequestResponse records input followed by output.
	OperationStyleRequestResponse OperationStyle11 = "request-response"
	// OperationStyleSolicitResponse records output followed by input; order is semantic.
	OperationStyleSolicitResponse OperationStyle11 = "solicit-response"
	// OperationStyleNotification records an output-only WSDL 1.1 operation.
	OperationStyleNotification OperationStyle11 = "notification"
)

// Operation11 is a WSDL 1.1 abstract operation.
type Operation11 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Style records the operation/binding style, not an application serialization policy.
	Style OperationStyle11
	// ParameterOrder retains the WSDL 1.1 RPC parameter-order names.
	ParameterOrder []string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Input is the optional first input message view; it does not imply an input exists.
	Input *OperationMessage11
	// Output is the optional first output message view; it does not imply an output exists.
	Output *OperationMessage11
	// Faults contains declared or referenced faults in this component's vocabulary.
	Faults []OperationMessage11
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// PortType11 is a named WSDL 1.1 abstract interface.
type PortType11 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Operations contains operation declarations or references in this component.
	Operations []Operation11
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// SOAPStyle is the SOAP RPC or document binding style.
type SOAPStyle string

const (
	// SOAPStyleDocument records document-style SOAP binding description data.
	SOAPStyleDocument SOAPStyle = "document"
	// SOAPStyleRPC records RPC-style SOAP binding description data.
	SOAPStyleRPC SOAPStyle = "rpc"
)

// SOAPUse is the SOAP encoded or literal message use.
type SOAPUse string

const (
	// SOAPUseLiteral describes literal XML message use, not a runtime encoder.
	SOAPUseLiteral SOAPUse = "literal"
	// SOAPUseEncoded describes SOAP encoding use without supplying an encoding runtime.
	SOAPUseEncoded SOAPUse = "encoded"
)

// SOAPBinding11 configures a WSDL 1.1 SOAP binding.
type SOAPBinding11 struct {
	// Version identifies the applicable language or protocol version of this component.
	Version Version
	// Style records the operation/binding style, not an application serialization policy.
	Style SOAPStyle
	// StyleSet tracks explicit Style presence independently of its effective value.
	StyleSet bool
	// Transport is the SOAP transport IRI, not an acquired transport.
	Transport string
	// TransportSet tracks explicit Transport presence independently of its effective value.
	TransportSet bool
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// SOAPOperation11 configures a bound SOAP operation.
type SOAPOperation11 struct {
	// Version identifies the applicable language or protocol version of this component.
	Version Version
	// Action records the described SOAP action; no request is sent.
	Action string
	// ActionSet tracks explicit Action presence independently of its effective value.
	ActionSet bool
	// ActionRequired records whether the description requires a SOAP action.
	ActionRequired bool
	// ActionRequiredSet tracks explicit ActionRequired presence independently of its effective value.
	ActionRequiredSet bool
	// Style records the operation/binding style, not an application serialization policy.
	Style SOAPStyle
	// StyleSet tracks explicit Style presence independently of its effective value.
	StyleSet bool
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// SOAPBody11 configures a SOAP message body.
type SOAPBody11 struct {
	// Version identifies the applicable language or protocol version of this component.
	Version Version
	// Use selects literal or encoded SOAP use as description data.
	Use SOAPUse
	// UseSet tracks explicit Use presence independently of its effective value.
	UseSet bool
	// Namespace is the namespace URI used in expanded identity, not a lexical prefix.
	Namespace string
	// NamespaceSet tracks explicit Namespace presence independently of its effective value.
	NamespaceSet bool
	// EncodingStyle contains described encoding-style IRIs, not installed codecs.
	EncodingStyle []string
	// EncodingStyleSet tracks explicit EncodingStyle presence independently of its effective value.
	EncodingStyleSet bool
	// Parts selects abstract message part names; PartsSet preserves omitted versus explicit
	// lists.
	Parts []string
	// PartsSet tracks explicit Parts presence independently of its effective value.
	PartsSet bool
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// SOAPHeaderFault11 binds a fault related to one SOAP header block.
type SOAPHeaderFault11 struct {
	// Version identifies the applicable language or protocol version of this component.
	Version Version
	// Message is the expanded WSDL 1.1 abstract message reference.
	Message QName
	// Part names the WSDL 1.1 abstract message part referenced by this binding.
	Part string
	// Use selects literal or encoded SOAP use as description data.
	Use SOAPUse
	// UseSet tracks explicit Use presence independently of its effective value.
	UseSet bool
	// Namespace is the namespace URI used in expanded identity, not a lexical prefix.
	Namespace string
	// NamespaceSet tracks explicit Namespace presence independently of its effective value.
	NamespaceSet bool
	// EncodingStyle contains described encoding-style IRIs, not installed codecs.
	EncodingStyle []string
	// EncodingStyleSet tracks explicit EncodingStyle presence independently of its effective value.
	EncodingStyleSet bool
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// SOAPHeader11 binds one message part to a SOAP header block.
type SOAPHeader11 struct {
	// Version identifies the applicable language or protocol version of this component.
	Version Version
	// Message is the expanded WSDL 1.1 abstract message reference.
	Message QName
	// Part names the WSDL 1.1 abstract message part referenced by this binding.
	Part string
	// Use selects literal or encoded SOAP use as description data.
	Use SOAPUse
	// UseSet tracks explicit Use presence independently of its effective value.
	UseSet bool
	// Namespace is the namespace URI used in expanded identity, not a lexical prefix.
	Namespace string
	// NamespaceSet tracks explicit Namespace presence independently of its effective value.
	NamespaceSet bool
	// EncodingStyle contains described encoding-style IRIs, not installed codecs.
	EncodingStyle []string
	// EncodingStyleSet tracks explicit EncodingStyle presence independently of its effective value.
	EncodingStyleSet bool
	// HeaderFaults contains faults associated with this SOAP header declaration.
	HeaderFaults []SOAPHeaderFault11
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// SOAPFault11 binds one WSDL fault to a SOAP fault.
type SOAPFault11 struct {
	// Version identifies the applicable language or protocol version of this component.
	Version Version
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Use selects literal or encoded SOAP use as description data.
	Use SOAPUse
	// UseSet tracks explicit Use presence independently of its effective value.
	UseSet bool
	// Namespace is the namespace URI used in expanded identity, not a lexical prefix.
	Namespace string
	// NamespaceSet tracks explicit Namespace presence independently of its effective value.
	NamespaceSet bool
	// EncodingStyle contains described encoding-style IRIs, not installed codecs.
	EncodingStyle []string
	// EncodingStyleSet tracks explicit EncodingStyle presence independently of its effective value.
	EncodingStyleSet bool
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// HTTPBinding11 configures the verb for a WSDL 1.1 HTTP binding.
type HTTPBinding11 struct {
	// Verb records the WSDL 1.1 HTTP binding verb without executing it.
	Verb string
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// HTTPOperation11 configures the relative URI template of an HTTP operation.
type HTTPOperation11 struct {
	// Location retains the lexical URI/address text, not a resolved or authorized resource.
	Location string
	// Source records where the lexical URI/address was decoded.
	Source Location
}

// HTTPMessageMode identifies WSDL 1.1 HTTP input serialization.
type HTTPMessageMode string

const (
	// HTTPURLEncoded describes WSDL 1.1 URL-encoded input binding.
	HTTPURLEncoded HTTPMessageMode = "urlEncoded"
	// HTTPURLReplacement describes WSDL 1.1 URL-template replacement input binding.
	HTTPURLReplacement HTTPMessageMode = "urlReplacement"
)

// HTTPMessage11 configures WSDL 1.1 HTTP input serialization.
type HTTPMessage11 struct {
	// Mode selects URL-encoded or URL-replacement input description semantics.
	Mode HTTPMessageMode
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// MIMEContent11 binds one message part to an Internet media type.
type MIMEContent11 struct {
	// Part names the WSDL 1.1 abstract message part referenced by this binding.
	Part string
	// Type records the bound Internet media type.
	Type string
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// MIMEXML11 binds one message part as XML.
type MIMEXML11 struct {
	// Part names the WSDL 1.1 abstract message part referenced by this binding.
	Part string
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// MIMEPart11 is one part of a multipart/related binding.
type MIMEPart11 struct {
	// Contents contains media-type bindings for WSDL 1.1 message parts.
	Contents []MIMEContent11
	// XML contains WSDL 1.1 MIME XML message-part bindings.
	XML []MIMEXML11
	// SOAPBody is optional SOAP body description data; nil omits the extension.
	SOAPBody *SOAPBody11
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// MIMEMultipart11 describes a multipart/related message.
type MIMEMultipart11 struct {
	// Parts contains multipart/related part-description values.
	Parts []MIMEPart11
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// MIMEMessage11 contains MIME extensions applied to a bound message.
type MIMEMessage11 struct {
	// Contents contains media-type bindings for WSDL 1.1 message parts.
	Contents []MIMEContent11
	// XML contains WSDL 1.1 MIME XML message-part bindings.
	XML []MIMEXML11
	// Multipart contains multipart/related descriptions, not a MIME encoder.
	Multipart []MIMEMultipart11
}

// BindingMessage11 applies concrete extensions to an operation message.
type BindingMessage11 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// SOAPBody is optional SOAP body description data; nil omits the extension.
	SOAPBody *SOAPBody11
	// SOAPHeaders contains described SOAP header bindings in model order.
	SOAPHeaders []SOAPHeader11
	// SOAPFault is optional concrete SOAP fault-description data.
	SOAPFault *SOAPFault11
	// HTTP is optional HTTP extension description data; no HTTP client is owned.
	HTTP *HTTPMessage11
	// MIME is optional MIME binding-description data; nil omits those extensions.
	MIME *MIMEMessage11
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// BindingOperation11 binds one abstract WSDL 1.1 operation.
type BindingOperation11 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// SOAP is optional SOAP extension description data; no envelope runtime is owned.
	SOAP *SOAPOperation11
	// HTTP is optional HTTP extension description data; no HTTP client is owned.
	HTTP *HTTPOperation11
	// Input is the optional first input message view; it does not imply an input exists.
	Input *BindingMessage11
	// Output is the optional first output message view; it does not imply an output exists.
	Output *BindingMessage11
	// Faults contains declared or referenced faults in this component's vocabulary.
	Faults []BindingMessage11
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Binding11 binds an abstract port type to a concrete protocol and format.
type Binding11 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Type is the expanded WSDL 1.1 port-type reference.
	Type QName
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// SOAP is optional SOAP extension description data; no envelope runtime is owned.
	SOAP *SOAPBinding11
	// HTTP is optional HTTP extension description data; no HTTP client is owned.
	HTTP *HTTPBinding11
	// Operations contains operation declarations or references in this component.
	Operations []BindingOperation11
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// SOAPAddress11 is the endpoint address supplied by a SOAP binding.
type SOAPAddress11 struct {
	// Version identifies the applicable language or protocol version of this component.
	Version Version
	// Location retains the lexical URI/address text, not a resolved or authorized resource.
	Location string
	// Source records where the lexical URI/address was decoded.
	Source Location
}

// HTTPAddress11 is the endpoint address supplied by an HTTP binding.
type HTTPAddress11 struct {
	// Location retains the lexical URI/address text, not a resolved or authorized resource.
	Location string
	// Source records where the lexical URI/address was decoded.
	Source Location
}

// Port11 is one named endpoint in a WSDL 1.1 service.
type Port11 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Binding is the expanded concrete binding reference, not a transport instance.
	Binding QName
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// SOAPAddress is an optional described SOAP endpoint address, not a live endpoint.
	SOAPAddress *SOAPAddress11
	// HTTPAddress is an optional described HTTP endpoint address, not a live endpoint.
	HTTPAddress *HTTPAddress11
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Service11 groups related WSDL 1.1 endpoints.
type Service11 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Ports contains WSDL 1.1 named endpoint declarations.
	Ports []Port11
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}
