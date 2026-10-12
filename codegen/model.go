// Package codegen uses Build to project an immutable compiled WSDL graph into
// an owned mutable generation model. Collection limits cap accepted results,
// not all intermediate copying/allocation. There is no context/cancellation API.
// Naming, source emission, transport, credentials, files and workers remain
// application-owned; generation candidates are not executable clients.
package codegen

import (
	"errors"
	"fmt"

	wsdl "github.com/faustbrian/go-wsdl/v3"
	wsdlcompile "github.com/faustbrian/go-wsdl/v3/compile"
)

// ErrLimitExceeded identifies a generation model beyond configured bounds.
var ErrLimitExceeded = errors.New("wsdl codegen: resource limit exceeded")

const defaultModelLimit = 100000

// Limits bounds materialized generation-model collections.
type Limits struct {
	// MaxInterfaces caps returned interfaces count; zero selects 100,000, not a pre-allocation
	// heap limit.
	MaxInterfaces int
	// MaxOperations caps returned operations count; zero selects 100,000, not a pre-allocation
	// heap limit.
	MaxOperations int
	// MaxParts caps logical parts across Inputs/Outputs, excluding duplicated legacy
	// Input/Output views; zero selects 100,000. It is not a pre-allocation heap limit.
	MaxParts int
	// MaxFaults caps returned faults count; zero selects 100,000, not a pre-allocation heap
	// limit.
	MaxFaults int
	// MaxBindings caps returned bindings count; zero selects 100,000, not a pre-allocation
	// heap limit.
	MaxBindings int
	// MaxServices caps returned services count; zero selects 100,000, not a pre-allocation
	// heap limit.
	MaxServices int
	// MaxEndpoints caps returned endpoints count; zero selects 100,000, not a pre-allocation
	// heap limit.
	MaxEndpoints int
	// MaxTypes caps returned types count; zero selects 100,000, not a pre-allocation heap
	// limit.
	MaxTypes int
	// MaxElements caps returned elements count; zero selects 100,000, not a pre-allocation
	// heap limit.
	MaxElements int
}

// Options configures generation-model construction.
type Options struct {
	// Limits selects phase-specific collection budgets; zero fields select finite defaults.
	Limits Limits
}

// Part is one generated message field candidate.
type Part struct {
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Element is the expanded XML Schema element reference; zero means no named element.
	Element wsdl.QName
	// Type is the expanded XML Schema type reference; zero means no named type.
	Type wsdl.QName
}

// Message is one generated operation payload candidate.
type Message struct {
	// Label is the effective message-reference label, including WSDL default naming.
	Label string
	// Name is the expanded component identity; lexical prefixes are not identity.
	Name wsdl.QName
	// Element is the expanded XML Schema element reference; zero means no named element.
	Element wsdl.QName
	// ContentModel records the effective WSDL 2.0 payload model, not a wire encoder.
	ContentModel wsdl.MessageContentModel
	// Parts contains the WSDL 1.1 message parts in model order.
	Parts []Part
}

// Fault is one generated operation error candidate.
type Fault struct {
	// Name is the expanded component identity; lexical prefixes are not identity.
	Name wsdl.QName
	// Label retains the fault-reference/message name; omitted WSDL 2.0 labels may be empty.
	Label string
	// Direction identifies the input/output role of this fault or RPC parameter.
	Direction string
	// Message is the expanded WSDL 1.1 abstract message reference.
	Message wsdl.QName
	// Element is the expanded XML Schema element reference; zero means no named element.
	Element wsdl.QName
	// ContentModel records the effective WSDL 2.0 payload model, not a wire encoder.
	ContentModel wsdl.MessageContentModel
}

// Operation is one generated client or server method candidate.
type Operation struct {
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Pattern is the message-exchange pattern IRI; unknown absolute IRIs remain extensible.
	Pattern string
	// Style records the operation/binding style, not an application serialization policy.
	Style string
	// Styles contains effective WSDL 2.0 style IRIs after interface defaults.
	Styles []string
	// Safe records the description's operation-safety assertion, not verified endpoint
	// behavior.
	Safe bool
	// RPCSignature contains ordered expanded parameter names and RPC directions.
	RPCSignature []wsdl.RPCSignatureParameter20
	// RPCSignatureSet tracks explicit RPCSignature presence independently of its effective value.
	RPCSignatureSet bool
	// Inputs contains input message references; retain order and explicit labels.
	Inputs []Message
	// Outputs contains output message references; retain order and explicit labels.
	Outputs []Message
	// Input is the optional first input message view; it does not imply an input exists.
	Input *Message
	// Output is the optional first output message view; it does not imply an output exists.
	Output *Message
	// Faults contains declared or referenced faults in this component's vocabulary.
	Faults []Fault
}

// Interface groups generated operation candidates.
type Interface struct {
	// Name is the expanded component identity; lexical prefixes are not identity.
	Name wsdl.QName
	// Extends contains expanded parent-interface references, not inherited copies.
	Extends []wsdl.QName
	// Operations contains operation declarations or references in this component.
	Operations []Operation
}

// Binding preserves the protocol association used by generated adapters.
type Binding struct {
	// Name is the expanded component identity; lexical prefixes are not identity.
	Name wsdl.QName
	// Interface is the expanded abstract interface reference; zero may mean unspecified.
	Interface wsdl.QName
	// Type records the WSDL 2.0 protocol binding IRI; WSDL 1.1 may leave it empty.
	Type string
	// OperationReferences retains complete bound-operation identity, including WSDL 1.1
	// overload labels.
	OperationReferences []wsdlcompile.OperationReference
}

// Endpoint preserves generated client endpoint metadata without transport behavior.
type Endpoint struct {
	// Name is the local component identifier in its owning symbol space.
	Name string
	// Binding is the expanded concrete binding reference, not a transport instance.
	Binding wsdl.QName
	// Address is endpoint-description text; it is never contacted or trusted as an allowlist.
	Address string
}

// Service groups generated endpoint candidates.
type Service struct {
	// Name is the expanded component identity; lexical prefixes are not identity.
	Name wsdl.QName
	// Interface is the expanded abstract interface reference; zero may mean unspecified.
	Interface wsdl.QName
	// Endpoints contains described endpoint metadata without owning connections.
	Endpoints []Endpoint
}

// Model is an owned deterministic input for optional source generators.
type Model struct {
	// Interfaces contains abstract interface declarations, not executable clients.
	Interfaces []Interface
	// Bindings contains concrete protocol-description bindings, not runtime adapters.
	Bindings []Binding
	// Services contains service declarations and their described endpoints.
	Services []Service
	// Types contains expanded named simple/complex XSD type identities.
	Types []wsdl.QName
	// Elements contains expanded global XML Schema element names available to generators.
	Elements []wsdl.QName
}

// Build rejects nil sets/negative limits and returns independent mutable model
// storage in compiled order. Zero limits select 100,000 entries per collection.
// Some accessor copies and collection construction precede count checks; limits
// are not pre-allocation heap ceilings. No partial model is returned on failure.
func Build(set *wsdlcompile.Set, options Options) (*Model, error) {
	if set == nil {
		return nil, errors.New("wsdl codegen: compiled set is nil")
	}
	limits, err := normalizeLimits(options.Limits)
	if err != nil {
		return nil, err
	}
	model := &Model{}
	// Set accessors return owned copies before generation-specific count checks.
	// These limits constrain accepted models, not every intermediate allocation.
	interfaces := set.Interfaces()
	if err := within("interfaces", len(interfaces), limits.MaxInterfaces); err != nil {
		return nil, err
	}
	operations, parts, faults := 0, 0, 0
	for _, value := range interfaces {
		item := Interface{
			Name: value.Name, Extends: append([]wsdl.QName(nil), value.Extends...),
		}
		for _, operation := range value.Operations {
			operations++
			for _, message := range operation.Inputs {
				parts += len(message.Parts)
			}
			for _, message := range operation.Outputs {
				parts += len(message.Parts)
			}
			faults += len(operation.Faults)
			item.Operations = append(item.Operations, copyOperation(operation))
		}
		model.Interfaces = append(model.Interfaces, item)
	}
	counts := []struct {
		name           string
		count, maximum int
	}{
		{name: "operations", count: operations, maximum: limits.MaxOperations},
		{name: "parts", count: parts, maximum: limits.MaxParts},
		{name: "faults", count: faults, maximum: limits.MaxFaults},
	}
	for _, value := range counts {
		if err := within(value.name, value.count, value.maximum); err != nil {
			return nil, err
		}
	}
	bindings := set.Bindings()
	if err := within("bindings", len(bindings), limits.MaxBindings); err != nil {
		return nil, err
	}
	for _, value := range bindings {
		model.Bindings = append(model.Bindings, Binding{
			Name: value.Name, Interface: value.Interface, Type: value.Type,
			OperationReferences: append(
				[]wsdlcompile.OperationReference(nil), value.OperationReferences...,
			),
		})
	}
	services := set.Services()
	if err := within("services", len(services), limits.MaxServices); err != nil {
		return nil, err
	}
	endpointCount := 0
	for _, value := range services {
		service := Service{Name: value.Name, Interface: value.Interface}
		for _, endpoint := range value.Endpoints {
			endpointCount++
			service.Endpoints = append(service.Endpoints, Endpoint{
				Name: endpoint.Name, Binding: endpoint.Binding, Address: endpoint.Address,
			})
		}
		model.Services = append(model.Services, service)
	}
	if err := within("endpoints", endpointCount, limits.MaxEndpoints); err != nil {
		return nil, err
	}
	if schemas := set.Schemas(); schemas != nil {
		for _, value := range append(schemas.SimpleTypeNames(), schemas.ComplexTypeNames()...) {
			model.Types = append(model.Types, wsdl.QName{Namespace: value.Namespace, Local: value.Local})
		}
		for _, value := range schemas.ElementNames() {
			model.Elements = append(model.Elements, wsdl.QName{Namespace: value.Namespace, Local: value.Local})
		}
	}
	if err := within("types", len(model.Types), limits.MaxTypes); err != nil {
		return nil, err
	}
	if err := within("elements", len(model.Elements), limits.MaxElements); err != nil {
		return nil, err
	}
	return model, nil
}

func normalizeLimits(value Limits) (Limits, error) {
	values := []*int{
		&value.MaxInterfaces, &value.MaxOperations, &value.MaxParts,
		&value.MaxFaults, &value.MaxBindings, &value.MaxServices,
		&value.MaxEndpoints, &value.MaxTypes, &value.MaxElements,
	}
	for _, limit := range values {
		if *limit < 0 {
			return Limits{}, errors.New("wsdl codegen: limits must not be negative")
		}
		if *limit == 0 {
			*limit = defaultModelLimit
		}
	}
	return value, nil
}

func within(name string, count, maximum int) error {
	if count > maximum {
		return fmt.Errorf("%w: %s exceed %d", ErrLimitExceeded, name, maximum)
	}
	return nil
}

func copyOperation(value wsdlcompile.Operation) Operation {
	result := Operation{
		Name: value.Name, Pattern: value.Pattern, Style: value.Style,
		Styles: append([]string(nil), value.Styles...), Safe: value.Safe,
		RPCSignature:    append([]wsdl.RPCSignatureParameter20(nil), value.RPCSignature...),
		RPCSignatureSet: value.RPCSignatureSet,
		Input:           copyMessage(value.Input), Output: copyMessage(value.Output),
	}
	for index := range value.Inputs {
		result.Inputs = append(result.Inputs, *copyMessage(&value.Inputs[index]))
	}
	for index := range value.Outputs {
		result.Outputs = append(result.Outputs, *copyMessage(&value.Outputs[index]))
	}
	for _, fault := range value.Faults {
		result.Faults = append(result.Faults, Fault{
			Name: fault.Name, Label: fault.Label, Direction: fault.Direction,
			Message: fault.Message, Element: fault.Element,
			ContentModel: fault.ContentModel,
		})
	}
	return result
}

func copyMessage(value *wsdlcompile.Message) *Message {
	if value == nil {
		return nil
	}
	result := &Message{
		Label: value.Label, Name: value.Name, Element: value.Element,
		ContentModel: value.ContentModel,
	}
	for _, part := range value.Parts {
		result.Parts = append(result.Parts, Part{
			Name: part.Name, Element: part.Element, Type: part.Type,
		})
	}
	return result
}
