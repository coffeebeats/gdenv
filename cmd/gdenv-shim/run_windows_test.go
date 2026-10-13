//go:build windows

package main

import (
	"flag"
	"os"
	"strconv"
	"testing"
)

/* -------------------------------------------------------------------------- */
/*                              Function: TestRun                             */
/* -------------------------------------------------------------------------- */

// TestRun covers the exit codes the shim passes back from the process it runs,
// which is this test binary exiting through 'TestRunHelper'.
func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code int
	}{
		{name: "passes back success", code: 0},
		{name: "passes back a failure", code: 3},
		{name: "passes back an interrupted exit", code: 0xC000013A},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Given: A process which exits with the test case's code.
			args := []string{"-test.run=^TestRunHelper$", "--", strconv.Itoa(tc.code)}

			// When: The shim runs it.
			got, err := run(os.Args[0], args...)

			// Then: The shim returns the process's code without an error.
			if err != nil {
				t.Fatalf("err: got %v, want nil", err)
			}

			if got != tc.code {
				t.Fatalf("exit code: got %d, want %d", got, tc.code)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                           Function: TestRunHelper                          */
/* -------------------------------------------------------------------------- */

// TestRunHelper exits with the code following '--' on the command line, which
// only 'TestRun' passes; run as an ordinary test, it does nothing.
func TestRunHelper(t *testing.T) {
	t.Parallel()

	args := flag.Args()
	if len(args) != 1 {
		return
	}

	code, err := strconv.Atoi(args[0])
	if err != nil {
		t.Fatal(err)
	}

	os.Exit(code)
}
