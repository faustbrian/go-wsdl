package wsdl

const (
	// NamespaceWSDL20 is the namespace of WSDL 2.0 core elements.
	NamespaceWSDL20 = "http://www.w3.org/ns/wsdl"
	// NamespaceWSDL20SOAP is the WSDL 2.0 SOAP binding namespace.
	NamespaceWSDL20SOAP = "http://www.w3.org/ns/wsdl/soap"
	// NamespaceWSDL20HTTP is the WSDL 2.0 HTTP binding namespace.
	NamespaceWSDL20HTTP = "http://www.w3.org/ns/wsdl/http"
	// NamespaceWSDL20Extensions contains the WSDL 2.0 predefined extensions.
	NamespaceWSDL20Extensions = "http://www.w3.org/ns/wsdl-extensions"
	// NamespaceWSDL20RPC contains the WSDL 2.0 RPC style extension.
	NamespaceWSDL20RPC = "http://www.w3.org/ns/wsdl/rpc"
)

// MessageContentModel identifies how a WSDL 2.0 message payload is described.
type MessageContentModel string

const (
	// MessageContentAny identifies the #any message payload marker.
	MessageContentAny MessageContentModel = "#any"
	// MessageContentNone identifies a message with no payload.
	MessageContentNone MessageContentModel = "#none"
	// MessageContentOther identifies unspecified payload content, including an absent element attribute.
	MessageContentOther MessageContentModel = "#other"
	// MessageContentElement identifies payload described by a concrete Element QName.
	MessageContentElement MessageContentModel = "#element"
)

// RPCDirection identifies one WSDL 2.0 RPC signature parameter direction.
type RPCDirection string

const (
	// RPCDirectionIn identifies an RPC input parameter.
	RPCDirectionIn RPCDirection = "#in"
	// RPCDirectionOut identifies an RPC output parameter.
	RPCDirectionOut RPCDirection = "#out"
	// RPCDirectionInOut identifies a parameter shared across RPC directions.
	RPCDirectionInOut RPCDirection = "#inout"
	// RPCDirectionReturn identifies the RPC return value parameter.
	RPCDirectionReturn RPCDirection = "#return"
)

// RPCSignatureParameter20 is one qualified parameter and direction pair.
type RPCSignatureParameter20 struct {
	// Name is the expanded component identity; lexical prefixes are not identity.
	Name QName
	// Direction identifies the input/output role of this fault or RPC parameter.
	Direction RPCDirection
}

// MessageExchangePattern identifies the message flow of a WSDL 2.0 operation.
type MessageExchangePattern string

const (
	// MEPInOnly identifies the recognized input-only message-exchange pattern.
	MEPInOnly MessageExchangePattern = "http://www.w3.org/ns/wsdl/in-only"
	// MEPRobustInOnly identifies input with an optional fault response.
	MEPRobustInOnly MessageExchangePattern = "http://www.w3.org/ns/wsdl/robust-in-only"
	// MEPInOut identifies the recognized input/output request-response pattern.
	MEPInOut MessageExchangePattern = "http://www.w3.org/ns/wsdl/in-out"
	// MEPInOptionalOut identifies input with optional output.
	MEPInOptionalOut MessageExchangePattern = "http://www.w3.org/ns/wsdl/in-opt-out"
	// MEPOutOnly identifies the recognized additional output-only pattern.
	MEPOutOnly MessageExchangePattern = "http://www.w3.org/ns/wsdl/out-only"
	// MEPRobustOutOnly identifies output with an optional fault response.
	MEPRobustOutOnly MessageExchangePattern = "http://www.w3.org/ns/wsdl/robust-out-only"
	// MEPOutIn identifies the recognized output/input pattern.
	MEPOutIn MessageExchangePattern = "http://www.w3.org/ns/wsdl/out-in"
	// MEPOutOptionalIn identifies output with optional input.
	MEPOutOptionalIn MessageExchangePattern = "http://www.w3.org/ns/wsdl/out-opt-in"
)

const (
	// StyleIRI identifies the WSDL 2.0 IRI operation style.
	StyleIRI = "http://www.w3.org/ns/wsdl/style/iri"
	// StyleMultipart identifies the WSDL 2.0 multipart operation style.
	StyleMultipart = "http://www.w3.org/ns/wsdl/style/multipart"
	// StyleRPC identifies the WSDL 2.0 RPC operation style.
	StyleRPC = "http://www.w3.org/ns/wsdl/style/rpc"
)

// Import20 references a WSDL 2.0 description in another namespace.
type Import20 struct {
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

// Include20 references a WSDL 2.0 description in the same namespace.
type Include20 struct {
	Extensibility
	// Location retains the lexical URI/address text, not a resolved or authorized resource.
	Location string
	// URI is the resource identity, not a filesystem path or network authorization decision.
	URI string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Source records where the lexical URI/address was decoded.
	Source Location
}

// InterfaceFault20 declares a fault available to interface operations.
type InterfaceFault20 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Element is the expanded XML Schema element reference; zero means no named element.
	Element QName
	// MessageContentModel records the effective element/content marker; absence defaults to
	// #other.
	MessageContentModel MessageContentModel
	// MessageContentModelSet tracks explicit MessageContentModel presence independently of its effective value.
	MessageContentModelSet bool
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// InterfaceMessageReference20 describes one operation message.
type InterfaceMessageReference20 struct {
	Extensibility
	// MessageLabel records the lexical message label; semantic validation applies MEP
	// defaults.
	MessageLabel string
	// Element is the expanded XML Schema element reference; zero means no named element.
	Element QName
	// MessageContentModel records the effective element/content marker; absence defaults to
	// #other.
	MessageContentModel MessageContentModel
	// MessageContentModelSet tracks explicit MessageContentModel presence independently of its effective value.
	MessageContentModelSet bool
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// InterfaceFaultReference20 connects an operation to an interface fault.
type InterfaceFaultReference20 struct {
	Extensibility
	// Ref is the expanded component reference in this binding/operation's symbol space.
	Ref QName
	// MessageLabel records the lexical message label; semantic validation applies MEP
	// defaults.
	MessageLabel string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// InterfaceOperation20 describes one abstract WSDL 2.0 operation.
type InterfaceOperation20 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Pattern is the message-exchange pattern IRI; unknown absolute IRIs remain extensible.
	Pattern MessageExchangePattern
	// Style records the operation/binding style, not an application serialization policy.
	Style []string
	// Safe records the description's operation-safety assertion, not verified endpoint
	// behavior.
	Safe bool
	// SafeSet tracks explicit Safe presence independently of its effective value.
	SafeSet bool
	// RPCSignature contains ordered expanded parameter names and RPC directions.
	RPCSignature []RPCSignatureParameter20
	// RPCSignatureSet tracks explicit RPCSignature presence independently of its effective value.
	RPCSignatureSet bool
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Inputs contains input message references; retain order and explicit labels.
	Inputs []InterfaceMessageReference20
	// Outputs contains output message references; retain order and explicit labels.
	Outputs []InterfaceMessageReference20
	// Input and Output retain the first message for source compatibility.
	Input *InterfaceMessageReference20
	// Output is a legacy first-message view; nonempty Outputs takes precedence.
	Output *InterfaceMessageReference20
	// InFaults contains incoming fault references; direction is part of their identity.
	InFaults []InterfaceFaultReference20
	// OutFaults contains outgoing fault references; direction is part of their identity.
	OutFaults []InterfaceFaultReference20
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Interface20 is a named WSDL 2.0 abstract interface.
type Interface20 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Extends contains expanded parent-interface references, not inherited copies.
	Extends []QName
	// StyleDefault contains operation-style IRIs inherited when operation Style is empty.
	StyleDefault []string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Faults contains declared or referenced faults in this component's vocabulary.
	Faults []InterfaceFault20
	// Operations contains operation declarations or references in this component.
	Operations []InterfaceOperation20
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// SOAPModule20 declares a SOAP module required by a binding component.
type SOAPModule20 struct {
	Extensibility
	// Ref is the SOAP module IRI, not an expanded component QName.
	Ref string
	// Required records an explicit extension/module/header requirement, not package support
	// for it.
	Required bool
	// RequiredSet tracks explicit Required presence independently of its effective value.
	RequiredSet bool
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// SOAPHeader20 declares an element serialized as a SOAP header block.
type SOAPHeader20 struct {
	Extensibility
	// Element is the expanded XML Schema element reference; zero means no named element.
	Element QName
	// MustUnderstand records the declared SOAP header requirement, not runtime processing.
	MustUnderstand bool
	// MustUnderstandSet tracks explicit MustUnderstand presence independently of its effective value.
	MustUnderstandSet bool
	// Required records an explicit extension/module/header requirement, not package support
	// for it.
	Required bool
	// RequiredSet tracks explicit Required presence independently of its effective value.
	RequiredSet bool
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// SOAPBinding20 contains SOAP properties attached to a binding.
type SOAPBinding20 struct {
	// Version identifies the applicable language or protocol version of this component.
	Version string
	// VersionSet tracks explicit Version presence independently of its effective value.
	VersionSet bool
	// Protocol is the described SOAP underlying-protocol IRI.
	Protocol string
	// ProtocolSet tracks explicit Protocol presence independently of its effective value.
	ProtocolSet bool
	// MEPDefault is the default SOAP message-exchange protocol IRI for operations.
	MEPDefault string
	// MEPDefaultSet tracks explicit MEPDefault presence independently of its effective value.
	MEPDefaultSet bool
	// Modules contains SOAP module requirements as opaque description metadata.
	Modules []SOAPModule20
}

// SOAPFaultBinding20 contains SOAP properties attached to a binding fault.
type SOAPFaultBinding20 struct {
	// Code is the concrete SOAP fault-code QName when CodeAny is false.
	Code QName
	// CodeAny distinguishes the SOAP #any code marker from a concrete Code QName.
	CodeAny bool
	// CodeSet tracks explicit Code presence independently of its effective value.
	CodeSet bool
	// Subcodes contains ordered SOAP fault subcode QNames when SubcodesAny is false.
	Subcodes []QName
	// SubcodesAny distinguishes the SOAP #any subcodes marker from a QName list.
	SubcodesAny bool
	// SubcodesSet tracks explicit Subcodes presence independently of its effective value.
	SubcodesSet bool
	// Modules contains SOAP module requirements as opaque description metadata.
	Modules []SOAPModule20
	// Headers contains described HTTP or SOAP header fields for this message/fault.
	Headers []SOAPHeader20
}

// SOAPOperationBinding20 contains SOAP properties attached to an operation.
type SOAPOperationBinding20 struct {
	// MEP records an operation-specific SOAP message-exchange protocol IRI.
	MEP string
	// MEPSet tracks explicit MEP presence independently of its effective value.
	MEPSet bool
	// Action records the described SOAP action; no request is sent.
	Action string
	// ActionSet tracks explicit Action presence independently of its effective value.
	ActionSet bool
	// Modules contains SOAP module requirements as opaque description metadata.
	Modules []SOAPModule20
}

// SOAPMessageBinding20 contains SOAP properties attached to a message reference.
type SOAPMessageBinding20 struct {
	// Modules contains SOAP module requirements as opaque description metadata.
	Modules []SOAPModule20
	// Headers contains described HTTP or SOAP header fields for this message/fault.
	Headers []SOAPHeader20
}

// SOAPFaultReferenceBinding20 contains SOAP modules attached to a fault reference.
type SOAPFaultReferenceBinding20 struct {
	// Modules contains SOAP module requirements as opaque description metadata.
	Modules []SOAPModule20
}

// HTTPHeader20 declares one HTTP header field used by a binding message.
type HTTPHeader20 struct {
	Extensibility
	// Name is the HTTP field name token, not an XML component QName.
	Name string
	// Type is the expanded XML Schema type reference; zero means no named type.
	Type QName
	// Required records an explicit extension/module/header requirement, not package support
	// for it.
	Required bool
	// RequiredSet tracks explicit Required presence independently of its effective value.
	RequiredSet bool
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// HTTPBinding20 contains HTTP properties attached to a binding.
type HTTPBinding20 struct {
	// MethodDefault records the binding's default HTTP method; presence remains separate.
	MethodDefault string
	// MethodDefaultSet tracks explicit MethodDefault presence independently of its effective value.
	MethodDefaultSet bool
	// Version identifies the applicable language or protocol version of this component.
	Version string
	// VersionSet tracks explicit Version presence independently of its effective value.
	VersionSet bool
	// QueryParameterSeparatorDefault records the default HTTP query separator, not a URI
	// encoder.
	QueryParameterSeparatorDefault string
	// QueryParameterSeparatorDefaultSet tracks explicit QueryParameterSeparatorDefault presence independently of its effective value.
	QueryParameterSeparatorDefaultSet bool
	// ContentEncodingDefault records default HTTP content-coding text, not installed
	// compression.
	ContentEncodingDefault string
	// ContentEncodingDefaultSet tracks explicit ContentEncodingDefault presence independently of its effective value.
	ContentEncodingDefaultSet bool
	// DefaultTransferCoding records default HTTP transfer-coding text, not a stream encoder.
	DefaultTransferCoding string
	// DefaultTransferCodingSet tracks explicit DefaultTransferCoding presence independently of its effective value.
	DefaultTransferCodingSet bool
	// Cookies records the HTTP cookie-use declaration, not a cookie store.
	Cookies bool
	// CookiesSet tracks explicit Cookies presence independently of its effective value.
	CookiesSet bool
}

// HTTPFaultBinding20 contains HTTP properties attached to a binding fault.
type HTTPFaultBinding20 struct {
	// Code is the described HTTP fault status-code text.
	Code string
	// CodeSet tracks explicit Code presence independently of its effective value.
	CodeSet bool
	// ContentEncoding records HTTP content-coding text, not an installed decoder.
	ContentEncoding string
	// ContentEncodingSet tracks explicit ContentEncoding presence independently of its effective value.
	ContentEncodingSet bool
	// TransferCoding records HTTP transfer-coding text, not a transport implementation.
	TransferCoding string
	// TransferCodingSet tracks explicit TransferCoding presence independently of its effective value.
	TransferCodingSet bool
	// Headers contains described HTTP or SOAP header fields for this message/fault.
	Headers []HTTPHeader20
}

// HTTPOperationBinding20 contains HTTP properties attached to an operation.
type HTTPOperationBinding20 struct {
	// Location retains the lexical URI/address text, not a resolved or authorized resource.
	Location string
	// LocationSet tracks explicit Location presence independently of its effective value.
	LocationSet bool
	// Method records the operation HTTP method without issuing requests.
	Method string
	// MethodSet tracks explicit Method presence independently of its effective value.
	MethodSet bool
	// InputSerialization records the described input Internet media type.
	InputSerialization string
	// InputSerializationSet tracks explicit InputSerialization presence independently of its effective value.
	InputSerializationSet bool
	// OutputSerialization records the described output Internet media type.
	OutputSerialization string
	// OutputSerializationSet tracks explicit OutputSerialization presence independently of its effective value.
	OutputSerializationSet bool
	// FaultSerialization records the described fault Internet media type.
	FaultSerialization string
	// FaultSerializationSet tracks explicit FaultSerialization presence independently of its effective value.
	FaultSerializationSet bool
	// QueryParameterSeparator records the operation-specific HTTP query separator.
	QueryParameterSeparator string
	// QueryParameterSeparatorSet tracks explicit QueryParameterSeparator presence independently of its effective value.
	QueryParameterSeparatorSet bool
	// ContentEncodingDefault records default HTTP content-coding text, not installed
	// compression.
	ContentEncodingDefault string
	// ContentEncodingDefaultSet tracks explicit ContentEncodingDefault presence independently of its effective value.
	ContentEncodingDefaultSet bool
	// DefaultTransferCoding records default HTTP transfer-coding text, not a stream encoder.
	DefaultTransferCoding string
	// DefaultTransferCodingSet tracks explicit DefaultTransferCoding presence independently of its effective value.
	DefaultTransferCodingSet bool
	// IgnoreUncited records the HTTP ignoreUncited declaration for message parts.
	IgnoreUncited bool
	// IgnoreUncitedSet tracks explicit IgnoreUncited presence independently of its effective value.
	IgnoreUncitedSet bool
}

// HTTPMessageBinding20 contains HTTP properties attached to a message reference.
type HTTPMessageBinding20 struct {
	// ContentEncoding records HTTP content-coding text, not an installed decoder.
	ContentEncoding string
	// ContentEncodingSet tracks explicit ContentEncoding presence independently of its effective value.
	ContentEncodingSet bool
	// TransferCoding records HTTP transfer-coding text, not a transport implementation.
	TransferCoding string
	// TransferCodingSet tracks explicit TransferCoding presence independently of its effective value.
	TransferCodingSet bool
	// Headers contains described HTTP or SOAP header fields for this message/fault.
	Headers []HTTPHeader20
}

// HTTPFaultReferenceBinding20 contains HTTP transfer properties on a fault reference.
type HTTPFaultReferenceBinding20 struct {
	// TransferCoding records HTTP transfer-coding text, not a transport implementation.
	TransferCoding string
	// TransferCodingSet tracks explicit TransferCoding presence independently of its effective value.
	TransferCodingSet bool
}

// HTTPEndpoint20 contains HTTP authentication properties on an endpoint.
type HTTPEndpoint20 struct {
	// AuthenticationScheme records an authentication scheme name, not credentials or an
	// authenticator.
	AuthenticationScheme string
	// AuthenticationSchemeSet tracks explicit AuthenticationScheme presence independently of its effective value.
	AuthenticationSchemeSet bool
	// AuthenticationRealm records advertised realm text, not trusted authentication state.
	AuthenticationRealm string
	// AuthenticationRealmSet tracks explicit AuthenticationRealm presence independently of its effective value.
	AuthenticationRealmSet bool
}

// BindingFault20 binds one WSDL 2.0 interface fault.
type BindingFault20 struct {
	Extensibility
	// Ref is the expanded component reference in this binding/operation's symbol space.
	Ref QName
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// SOAP is optional SOAP extension description data; no envelope runtime is owned.
	SOAP *SOAPFaultBinding20
	// HTTP is optional HTTP extension description data; no HTTP client is owned.
	HTTP *HTTPFaultBinding20
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// BindingMessageReference20 binds one message in an interface operation.
type BindingMessageReference20 struct {
	Extensibility
	// MessageLabel records the lexical message label; semantic validation applies MEP
	// defaults.
	MessageLabel string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// SOAP is optional SOAP extension description data; no envelope runtime is owned.
	SOAP *SOAPMessageBinding20
	// HTTP is optional HTTP extension description data; no HTTP client is owned.
	HTTP *HTTPMessageBinding20
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// BindingFaultReference20 binds one fault in an interface operation.
type BindingFaultReference20 struct {
	Extensibility
	// Ref is the expanded component reference in this binding/operation's symbol space.
	Ref QName
	// MessageLabel records the lexical message label; semantic validation applies MEP
	// defaults.
	MessageLabel string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// SOAP is optional SOAP extension description data; no envelope runtime is owned.
	SOAP *SOAPFaultReferenceBinding20
	// HTTP is optional HTTP extension description data; no HTTP client is owned.
	HTTP *HTTPFaultReferenceBinding20
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// BindingOperation20 binds one WSDL 2.0 interface operation.
type BindingOperation20 struct {
	Extensibility
	// Ref is the expanded component reference in this binding/operation's symbol space.
	Ref QName
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// SOAP is optional SOAP extension description data; no envelope runtime is owned.
	SOAP *SOAPOperationBinding20
	// HTTP is optional HTTP extension description data; no HTTP client is owned.
	HTTP *HTTPOperationBinding20
	// Inputs contains input message references; retain order and explicit labels.
	Inputs []BindingMessageReference20
	// Outputs contains output message references; retain order and explicit labels.
	Outputs []BindingMessageReference20
	// InFaults contains incoming fault references; direction is part of their identity.
	InFaults []BindingFaultReference20
	// OutFaults contains outgoing fault references; direction is part of their identity.
	OutFaults []BindingFaultReference20
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Binding20 binds an interface to a concrete protocol.
type Binding20 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Interface is the expanded abstract interface reference; zero may mean unspecified.
	Interface QName
	// Type is the absolute protocol binding-type IRI, not an XML Schema type.
	Type string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// SOAP is optional SOAP extension description data; no envelope runtime is owned.
	SOAP *SOAPBinding20
	// HTTP is optional HTTP extension description data; no HTTP client is owned.
	HTTP *HTTPBinding20
	// Faults contains declared or referenced faults in this component's vocabulary.
	Faults []BindingFault20
	// Operations contains operation declarations or references in this component.
	Operations []BindingOperation20
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Endpoint20 is one concrete WSDL 2.0 service endpoint.
type Endpoint20 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Binding is the expanded concrete binding reference, not a transport instance.
	Binding QName
	// Address is endpoint-description text; it is never contacted or trusted as an allowlist.
	Address string
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// HTTP is optional HTTP extension description data; no HTTP client is owned.
	HTTP *HTTPEndpoint20
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Service20 groups endpoints implementing one interface.
type Service20 struct {
	Extensibility
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Interface is the expanded abstract interface reference; zero may mean unspecified.
	Interface QName
	// Documentation is optional description text; nil omits its documentation element.
	Documentation *Documentation
	// Endpoints contains described endpoint metadata without owning connections.
	Endpoints []Endpoint20
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

func interfaceInputs20(value InterfaceOperation20) []InterfaceMessageReference20 {
	if len(value.Inputs) > 0 {
		return value.Inputs
	}
	if value.Input != nil {
		return []InterfaceMessageReference20{*value.Input}
	}
	return nil
}

func interfaceOutputs20(value InterfaceOperation20) []InterfaceMessageReference20 {
	if len(value.Outputs) > 0 {
		return value.Outputs
	}
	if value.Output != nil {
		return []InterfaceMessageReference20{*value.Output}
	}
	return nil
}
