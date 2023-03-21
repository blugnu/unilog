package unilog

import "testing"

func TestLevelsString(t *testing.T) {
	testcases := []struct {
		name  string
		level Level
	}{
		{name: "Trace", level: Trace},
		{name: "Debug", level: Debug},
		{name: "Info", level: Info},
		{name: "Warn", level: Warn},
		{name: "Error", level: Error},
		{name: "Fatal", level: Fatal},
		{name: "<invalid (-1)>", level: Level(-1)},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			got := tc.level.String()

			// ASSERT
			wanted := tc.name
			if wanted != got {
				t.Errorf("wanted %v, got %v", wanted, got)
			}
		})
	}
}
