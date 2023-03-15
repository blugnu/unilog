package unilog

import (
	"context"
	"testing"
)

func TestDecorator(t *testing.T) {
	// ARRANGE
	oef := enrichmentFuncs
	defer func() { enrichmentFuncs = oef }()

	f := func(ctx context.Context, e Enricher) Entry { return e.(Entry) }

	// ACT
	if len(oef) != 0 {
		t.Fatal("`decorators` is not empty")
	}
	EnrichWith(f)

	// ASSERT
	wanted := 1
	got := len(enrichmentFuncs)
	if wanted != got {
		t.Errorf("wanted %v, got %v", wanted, got)
	}
}
