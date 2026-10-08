package wsdl

import (
	"context"
	"errors"
	"testing"
)

// completionContext cancels a real context immediately after an Err snapshot.
// It models cancellation between the parser's last observation and publication.
type completionContext struct {
	context.Context
	observations int
	cancelAfter  int
	cancel       context.CancelFunc
}

func (ctx *completionContext) Err() error {
	err := ctx.Context.Err()
	ctx.observations++
	if ctx.observations == ctx.cancelAfter && ctx.cancel != nil {
		ctx.cancel()
	}
	return err
}

func TestParseDoesNotPublishDocumentAfterCompletionCancellation(t *testing.T) {
	for _, source := range []string{
		`<definitions xmlns="http://schemas.xmlsoap.org/wsdl/"/>`,
		`<description xmlns="http://www.w3.org/ns/wsdl" targetNamespace="urn:example"/>`,
	} {
		t.Run(source, func(t *testing.T) {
			// Calibrate only the private parser, not the public completion guard.
			baseline := &completionContext{Context: context.Background()}
			document, err := parse(baseline, []byte(source), ParseOptions{})
			if err != nil || document == nil || baseline.observations == 0 {
				t.Fatal("valid document must finish the private parser")
			}
			underlying, cancel := context.WithCancel(context.Background())
			defer cancel()
			ctx := &completionContext{
				Context: underlying, cancelAfter: baseline.observations, cancel: cancel,
			}
			document, err = Parse(ctx, []byte(source), ParseOptions{})
			if underlying.Err() != context.Canceled {
				t.Fatal("fixture must cancel the real context at parser completion")
			}
			if document != nil || !errors.Is(err, context.Canceled) {
				t.Fatal("completion cancellation must suppress the document and preserve its cause")
			}
			if err.Error() != "wsdl: parse failed" {
				t.Fatal("completion cancellation must retain diagnostic privacy")
			}
		})
	}
}
