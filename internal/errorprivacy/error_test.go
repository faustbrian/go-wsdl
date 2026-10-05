package errorprivacy

import (
	"errors"
	"testing"
)

func TestWrapPreservesCauseWithoutDefaultDisclosure(t *testing.T) {
	if Wrap("wsdl: ordinary failure", nil) != nil {
		t.Fatal("successful operation acquired an error")
	}
	cause := errors.New("ordinary detail")
	err := Wrap("wsdl: ordinary failure", cause)
	if err.Error() != "wsdl: ordinary failure" || !errors.Is(err, cause) || errors.Unwrap(err) != cause {
		t.Fatal("safe category or inspectable cause changed")
	}
}
