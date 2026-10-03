package wsdl

import (
	"errors"
	"testing"
)

func TestOrdinaryModelAndMarshalErrorPrivacy(t *testing.T) {
	if payload, err := Marshal(&Document{}, MarshalOptions{}); payload != nil || err == nil || err.Error() != "wsdl: marshal failed" || errors.Unwrap(err) == nil {
		t.Fatalf("Marshal(empty model) = %q, %v", payload, err)
	}
	originalMarshal := canonicalMarshal
	t.Cleanup(func() { canonicalMarshal = originalMarshal })
	cause := errors.New("ordinary model collaborator unavailable")
	canonicalMarshal = func(*Document, MarshalOptions) ([]byte, error) { return nil, cause }
	document, err := NewDocument20(Description20{TargetNamespace: "urn:catalog"}, ValidationOptions{})
	if document != nil || err == nil || err.Error() != "wsdl: canonicalize model" || !errors.Is(err, cause) || errors.Unwrap(err) != cause {
		t.Fatalf("NewDocument20() = %v, %v", document, err)
	}
}
