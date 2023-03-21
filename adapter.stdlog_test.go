package unilog

import (
	"context"
	"fmt"
	stdlog "log"
	"reflect"
	"testing"

	log "github.com/blugnu/go-logspy"
)

func TestStdLogAdapter(t *testing.T) {
	// ARRANGE
	stdlog.SetOutput(log.Sink())
	stdlog.SetFlags(0) // clear all flags so that we can test only the output produced by LogAdapter
	sut := &stdlogAdapter{}

	testcases := []struct {
		name   string
		fn     func(string, ...any)
		args   []any
		output string
	}{
		{name: "debug", fn: func(s string, args ...any) { sut.Emit(Debug, s) }, args: []any{"entry text"}, output: "DEBUG: entry text"},
		{name: "debugf", fn: func(format string, args ...any) { sut.Emit(Debug, fmt.Sprintf(format, args...)) }, args: []any{"entry %d", 1}, output: "DEBUG: entry 1"},
		{name: "info", fn: func(s string, args ...any) { sut.Emit(Info, s) }, args: []any{"entry text"}, output: "INFO: entry text"},
		{name: "infof", fn: func(format string, args ...any) { sut.Emit(Info, fmt.Sprintf(format, args...)) }, args: []any{"entry %d", 1}, output: "INFO: entry 1"},
		{name: "warn", fn: func(s string, args ...any) { sut.Emit(Warn, s) }, args: []any{"entry text"}, output: "WARN: entry text"},
		{name: "warnf", fn: func(format string, args ...any) { sut.Emit(Warn, fmt.Sprintf(format, args...)) }, args: []any{"entry %d", 1}, output: "WARN: entry 1"},
		{name: "error", fn: func(s string, args ...any) { sut.Emit(Error, s) }, args: []any{"entry text"}, output: "ERROR: entry text"},
		{name: "debug and error", fn: func(s string, args ...any) { sut.Emit(Debug, s); sut.Emit(Error, s) }, args: []any{"entry"}, output: "DEBUG: entry\nERROR: entry"},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			defer log.Reset()
			// ARRANGE
			s := tc.args[0].(string)
			a := []any{}
			if len(tc.args) > 1 {
				a = append(a, tc.args[1:]...)
			}

			// ACT
			tc.fn(s, a...)

			// ASSERT
			wanted := tc.output
			got := log.String()
			if !log.Contains(wanted) {
				t.Errorf("\nwanted %q\ngot    %q", wanted, got)
			}
		})
	}
}

func TestLogAdapterWithFields(t *testing.T) {
	// ARRANGE
	stdlog.SetOutput(log.Sink())
	stdlog.SetFlags(0) // clear all flags so that we can test only the output produced by LogAdapter
	adapter := &stdlogAdapter{}
	sut := adapter.WithField("fieldname", "data")
	sut = sut.WithField("name space", "da ta")

	testcases := []struct {
		name   string
		fn     func(string, ...any)
		args   []any
		output string
	}{
		{name: "debug", fn: func(s string, args ...any) { sut.Emit(Debug, s) }, args: []any{"entry text"}, output: "fieldname=data \"name space\"=\"da ta\" DEBUG: entry text\n"},
		{name: "debugf", fn: func(format string, args ...any) { sut.Emit(Debug, fmt.Sprintf(format, args...)) }, args: []any{"entry %d", 1}, output: "fieldname=data \"name space\"=\"da ta\" DEBUG: entry 1\n"},
		{name: "info", fn: func(s string, args ...any) { sut.Emit(Info, s) }, args: []any{"entry text"}, output: "fieldname=data \"name space\"=\"da ta\" INFO: entry text\n"},
		{name: "infof", fn: func(format string, args ...any) { sut.Emit(Info, fmt.Sprintf(format, args...)) }, args: []any{"entry %d", 1}, output: "fieldname=data \"name space\"=\"da ta\" INFO: entry 1\n"},
		{name: "warn", fn: func(s string, args ...any) { sut.Emit(Warn, s) }, args: []any{"entry text"}, output: "fieldname=data \"name space\"=\"da ta\" WARN: entry text\n"},
		{name: "warnf", fn: func(format string, args ...any) { sut.Emit(Warn, fmt.Sprintf(format, args...)) }, args: []any{"entry %d", 1}, output: "fieldname=data \"name space\"=\"da ta\" WARN: entry 1\n"},
		{name: "error", fn: func(fmt string, args ...any) { sut.Emit(Error, fmt) }, args: []any{"entry text"}, output: "fieldname=data \"name space\"=\"da ta\" ERROR: entry text\n"},
		{name: "debug and error", fn: func(s string, args ...any) { sut.Emit(Debug, s); sut.Emit(Error, s) }, args: []any{"entry"}, output: "fieldname=data \"name space\"=\"da ta\" DEBUG: entry\nfieldname=data \"name space\"=\"da ta\" ERROR: entry\n"},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			defer log.Reset()
			// ARRANGE
			s := tc.args[0].(string)
			a := []any{}
			if len(tc.args) > 1 {
				a = append(a, tc.args[1:]...)
			}

			// ACT
			tc.fn(s, a...)

			// ASSERT
			wanted := tc.output
			got := log.String()
			if wanted != got {
				t.Errorf("\nwanted %q\ngot    %q", wanted, got)
			}
		})
	}
}

func TestUsingStdLog(t *testing.T) {
	// ACT
	ctx := context.Background()
	result := StdLog()

	// ASSERT
	wanted := UsingAdapter(ctx, &stdlogAdapter{fields: map[string]any{}})

	got := result
	if !reflect.DeepEqual(wanted, got) {
		t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
	}
}
