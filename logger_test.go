package unilog

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/blugnu/unilog/internal"
)

type MockAdapter struct {
	emitLevel       *Level
	emitString      *string
	emitCalled      *bool
	newEntryCalled  *bool
	withFieldCalled *bool
}

func (mock MockAdapter) Emit(level Level, s string) {
	*mock.emitLevel = level
	*mock.emitString = s
	*mock.emitCalled = true
}

func (mock MockAdapter) NewEntry() Adapter {
	*mock.newEntryCalled = true
	return mock
}

func (mock MockAdapter) WithField(name string, value any) Adapter {
	*mock.withFieldCalled = true
	return mock
}

func TestLogEmissions(t *testing.T) {
	// ARRANGE
	var (
		exitFnWasCalled bool
		exitCode        int

		emitLevel       Level
		emitString      string
		emitCalled      bool
		newEntryCalled  bool
		withFieldCalled bool

		enrichmentFuncsCalled bool
	)
	initSpies := func() {
		exitFnWasCalled = false
		exitCode = 0

		emitLevel = Level(-1)
		emitString = ""
		emitCalled = false
		newEntryCalled = false
		withFieldCalled = false

		enrichmentFuncsCalled = false
	}

	ofn := ExitFn
	defer func() { ExitFn = ofn }()

	ExitFn = func(code int) {
		exitCode = code
		exitFnWasCalled = true
	}

	sut := logger{Adapter: &MockAdapter{
		emitLevel:       &emitLevel,
		emitString:      &emitString,
		emitCalled:      &emitCalled,
		newEntryCalled:  &newEntryCalled,
		withFieldCalled: &withFieldCalled,
	}}

	// ACT
	t.Run("emit", func(t *testing.T) {
		testcases := []struct {
			name            string
			fn              func(string)
			message         string
			output          string
			callsExit       bool
			callsDecorators bool
			Level
		}{
			{name: "trace", fn: func(s string) { sut.Trace(s) }, Level: Trace, message: "test", callsExit: false},
			{name: "tracef", fn: func(s string) { sut.Tracef("formatted: %s", s) }, Level: Trace, message: "test", output: "formatted: test", callsExit: false},
			{name: "debug", fn: func(s string) { sut.Debug(s) }, Level: Debug, message: "test", callsExit: false},
			{name: "debugf", fn: func(s string) { sut.Debugf("formatted: %s", s) }, Level: Debug, message: "test", output: "formatted: test", callsExit: false},
			{name: "info", fn: func(s string) { sut.Info(s) }, Level: Info, message: "test", callsExit: false},
			{name: "infof", fn: func(s string) { sut.Infof("formatted: %s", s) }, Level: Info, message: "test", output: "formatted: test", callsExit: false},
			{name: "warn", fn: func(s string) { sut.Warn(s) }, Level: Warn, message: "test", callsExit: false},
			{name: "warnf", fn: func(s string) { sut.Warnf("formatted: %s", s) }, Level: Warn, message: "test", output: "formatted: test", callsExit: false},
			{name: "error", fn: func(s string) { sut.Error(errors.New(s)) }, Level: Error, message: "test", callsExit: false},
			{name: "errorf", fn: func(s string) { sut.Errorf("formatted: %s", errors.New(s)) }, Level: Error, message: "test", output: "formatted: test", callsExit: false},
			{name: "fatal", fn: func(s string) { sut.Fatal(s) }, Level: Fatal, message: "test", callsExit: true},
			{name: "fatalf", fn: func(s string) { sut.Fatalf("formatted: %s", errors.New(s)) }, Level: Fatal, message: "test", output: "formatted: test", callsExit: true},
			{name: "fatalerror", fn: func(s string) { sut.FatalError(errors.New(s)) }, Level: Fatal, message: "test", callsExit: true},
			{name: "withdecoration", fn: func(s string) {
				od := enrichmentFuncs
				defer func() { enrichmentFuncs = od }()
				EnrichWith(func(ctx context.Context, e Enricher) Entry {
					enrichmentFuncsCalled = true
					return e.(Entry)
				})
				sut.Info(s)
			}, Level: Info, message: "test", callsExit: false, callsDecorators: true},
		}
		for _, tc := range testcases {
			initSpies()

			t.Run(tc.name, func(t *testing.T) {
				// ACT
				tc.fn(tc.message)

				// ASSERT
				t.Run("level", func(t *testing.T) {
					wanted := tc.Level
					got := emitLevel
					if wanted != got {
						t.Errorf("wanted %v, got %v", wanted, got)
					}
				})

				t.Run("string", func(t *testing.T) {
					wanted := tc.output
					if wanted == "" {
						wanted = tc.message
					}
					got := emitString
					if wanted != got {
						t.Errorf("wanted %v, got %v", wanted, got)
					}
				})

				t.Run("calls decorators", func(t *testing.T) {
					wanted := tc.callsDecorators
					got := enrichmentFuncsCalled
					if wanted != got {
						t.Errorf("wanted %v, got %v", wanted, got)
					}
				})

				t.Run("calls exit fn", func(t *testing.T) {
					wanted := tc.callsExit
					got := exitFnWasCalled
					if wanted != got {
						t.Errorf("wanted %v, got %v", wanted, got)
					}
					if !tc.callsExit {
						return
					}

					t.Run("with exit code", func(t *testing.T) {
						wanted := 1
						got := exitCode
						if wanted != got {
							t.Errorf("wanted %v, got %v", wanted, got)
						}
					})
				})
			})
		}
	})
}

func TestLogWithField(t *testing.T) {
	// ARRANGE
	a := &logger{Adapter: &nulAdapter{}}

	// ACT
	b := a.WithField("field", "value")

	// ASSERT
	t.Run("returns new logger", func(t *testing.T) {
		wanted := true
		got := a != b
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})
}

func TestLogger_WithContext(t *testing.T) {
	// ARRANGE
	newEntryCalled := false

	ctx := context.Background()
	adapter := MockAdapter{
		newEntryCalled: &newEntryCalled,
	}
	sut := &logger{ctx, adapter}

	// ACT
	log := sut.WithContext(ctx)

	// ASSERT
	wanted := &logger{ctx, adapter}
	got := log
	if !reflect.DeepEqual(wanted, got) {
		t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
	}
}

func TestLogger_NewEntry(t *testing.T) {
	// ARRANGE
	newEntryCalled := false

	ctx := context.Background()
	adapter := MockAdapter{
		newEntryCalled: &newEntryCalled,
	}
	sut := &logger{ctx, adapter}

	// ACT
	log := sut.NewEntry()

	// ASSERT
	wanted := &logger{ctx, adapter}
	got := log
	if !reflect.DeepEqual(wanted, got) {
		t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
	}
}

func TestFromContext(t *testing.T) {
	t.Run("when context does not contain a Logger", func(t *testing.T) {
		// ARRANGE
		ctx := context.Background()

		// ACT
		entry := FromContext(ctx)

		// ASSERT
		wanted := (Entry)(nil)
		got := entry
		if wanted != got {
			t.Errorf("wanted %#v, got %#v", wanted, got)
		}
	})

	t.Run("when context contains a Logger", func(t *testing.T) {
		// ARRANGE
		ctx := context.Background()
		ctx = context.WithValue(ctx, internal.LoggerKey, nul)

		// ACT
		entry := FromContext(ctx)

		// ASSERT
		wanted := nul.WithContext(ctx).(*logger)
		got := entry.(*logger)
		if *wanted != *got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})
}

func TestUsingAdapter(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	adapter := &nulAdapter{}

	// ACT
	got := UsingAdapter(ctx, adapter).(*logger)

	// ASSERT
	wanted := &logger{ctx, adapter}
	if *wanted != *got {
		t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
	}
}
