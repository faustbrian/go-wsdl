package compile

import (
	"context"
	"errors"
	"testing"

	wsdl "github.com/faustbrian/go-wsdl/v2"
)

// phaseCancellationContext cancels a real context after an Err snapshot, so
// the next owned observation sees ordinary cancellation, not a synthetic error.
type phaseCancellationContext struct {
	context.Context
	observations int
	cancelAfter  int
	cancel       context.CancelFunc
}

func (ctx *phaseCancellationContext) Err() error {
	err := ctx.Context.Err()
	ctx.observations++
	if ctx.observations == ctx.cancelAfter && ctx.cancel != nil {
		ctx.cancel()
	}
	return err
}

func TestOwnedCompilePhasesRefuseCompletionCancellation(t *testing.T) {
	document, err := wsdl.NewDocument20(wsdl.Description20{
		TargetNamespace: "urn:ordinary",
	}, wsdl.ValidationOptions{})
	if err != nil {
		t.Fatal(err)
	}
	compiler, err := New(Options{})
	if err != nil {
		t.Fatal(err)
	}
	newState := func() *compileState {
		return &compileState{compiler: compiler, resources: map[string]*resourceDocument{
			"urn:ordinary": {document: document},
		}}
	}
	for _, phase := range []struct {
		name           string
		published      bool
		cancelSnapshot int
		run            func(*compileState, context.Context) (bool, error)
	}{
		// Snapshot positions belong to this fixed empty-model fixture, excluding
		// each owner's completion guard. Do not derive them from that guard:
		// removing a guard must not recalibrate cancellation onto an earlier step.
		{"resolution", false, 3, func(state *compileState, ctx context.Context) (bool, error) {
			return false, state.resolveDocument(ctx, "urn:ordinary", 1)
		}},
		{"set", true, 21, func(state *compileState, ctx context.Context) (bool, error) {
			set, err := state.buildSet(ctx)
			return set != nil, err
		}},
		{"schema-free", false, 2, func(state *compileState, ctx context.Context) (bool, error) {
			set, err := state.compileSchemas(ctx, []string{"urn:ordinary"})
			return set != nil, err
		}},
	} {
		t.Run(phase.name, func(t *testing.T) {
			baseline := &phaseCancellationContext{Context: context.Background()}
			published, err := phase.run(newState(), baseline)
			if err != nil || published != phase.published || baseline.observations < 2 {
				t.Fatal("valid control must finish the owned phase")
			}
			underlying, cancel := context.WithCancel(context.Background())
			defer cancel()
			ctx := &phaseCancellationContext{
				Context: underlying, cancelAfter: phase.cancelSnapshot, cancel: cancel,
			}
			published, err = phase.run(newState(), ctx)
			if underlying.Err() != context.Canceled {
				t.Fatal("fixture must cancel the real context before phase completion")
			}
			if published || !errors.Is(err, context.Canceled) {
				t.Fatal("owned phase must refuse completion and preserve cancellation")
			}
		})
	}
}

func TestInheritanceCancellationNeverConsumesAnInvalidOperationIndex(t *testing.T) {
	parent := wsdl.QName{Namespace: "urn:ordinary", Local: "Parent"}
	child := wsdl.QName{Namespace: "urn:ordinary", Local: "Child"}
	for _, localOperation := range []bool{false, true} {
		newInterfaces := func() []Interface {
			values := []Interface{
				{Name: parent, Operations: []Operation{{Name: "Inherited"}}},
				{Name: child, Extends: []wsdl.QName{parent}},
			}
			if localOperation {
				values[1].Operations = []Operation{{Name: "Local"}}
			}
			return values
		}
		baseline := &phaseCancellationContext{Context: context.Background()}
		values := newInterfaces()
		added, err := expandInterfaceInheritance(&compilationWalk{ctx: baseline}, values)
		if err != nil || added != 1 || baseline.observations < 2 {
			t.Fatal("valid control must inherit one operation")
		}
		// Cover cancellation between each pair of owned observations on this
		// tiny graph. Counts bound the fixture; assertions concern the result.
		for snapshot := 1; snapshot < baseline.observations; snapshot++ {
			underlying, cancel := context.WithCancel(context.Background())
			ctx := &phaseCancellationContext{
				Context: underlying, cancelAfter: snapshot, cancel: cancel,
			}
			added, err := expandInterfaceInheritance(&compilationWalk{ctx: ctx}, newInterfaces())
			cancel()
			if added != 0 || !errors.Is(err, context.Canceled) {
				t.Fatal("cancellation must return refusal rather than consume an invalid index")
			}
		}
	}
}
