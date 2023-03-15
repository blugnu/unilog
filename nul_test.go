package unilog

import (
	"reflect"
	"testing"

	log "github.com/blugnu/go-logspy"
)

func TestNulAdapter(t *testing.T) {
	// ARRANGE
	sut := &nulAdapter{}

	testcases := []struct {
		name string
		fn   func(string)
	}{
		{name: "debug", fn: func(s string) { sut.Emit(Debug, s) }},
		{name: "info", fn: func(s string) { sut.Emit(Info, s) }},
		{name: "warn", fn: func(s string) { sut.Emit(Warn, s) }},
		{name: "error", fn: func(s string) { sut.Emit(Error, s) }},
		{name: "debug and error", fn: func(s string) { sut.Emit(Debug, s); sut.Emit(Error, s) }},
		{name: "withfield", fn: func(s string) {
			a := sut
			b := sut.WithField("field", "data")

			t.Run("returns SAME logger", func(t *testing.T) {
				wanted := true
				got := a == b
				if wanted != got {
					t.Errorf("wanted %v, got %v", wanted, got)
				}
			})

			a.Emit(Info, s)
			b.Emit(Info, s)
		}},
		{name: "newentry", fn: func(s string) {
			a := sut
			b := sut.NewEntry()

			t.Run("returns SAME logger", func(t *testing.T) {
				wanted := true
				got := a == b
				if wanted != got {
					t.Errorf("wanted %v, got %v", wanted, got)
				}
			})

			a.Emit(Info, s)
			b.Emit(Info, s)
		}},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			defer log.Reset()

			// ACT
			tc.fn("entry text")

			// ASSERT
			wanted := ""
			got := log.String()
			if wanted != got {
				t.Errorf("\nwanted %q\ngot    %q", wanted, got)
			}
		})
	}
}

func TestNul(t *testing.T) {
	// ACT
	result := Nul()

	// ASSERT
	wanted := &logger{
		Context: nil,
		Adapter: &nulAdapter{},
	}
	got := result
	if !reflect.DeepEqual(wanted, got) {
		t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
	}
}
