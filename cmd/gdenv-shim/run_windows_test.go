//go:build windows

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

/* -------------------------------------------------------------------------- */
/*                              Function: TestRun                             */
/* -------------------------------------------------------------------------- */

// TestRun covers the exit codes the shim passes back from the process it runs,
// which is this test binary exiting through 'TestHelperProcess'.
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
			args := helperArgs(strconv.Itoa(tc.code))

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
/*                             Function: TestBind                             */
/* -------------------------------------------------------------------------- */

// TestBind covers the job 'bind' places a process in, which ends the process
// once the job's handle is closed, as it is when the shim exits.
func TestBind(t *testing.T) {
	t.Parallel()

	// Given: A process which runs until it's ended.
	cmd := exec.Command(os.Args[0], helperArgs("sleep")...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = cmd.Process.Kill() })

	// Given: The process has started sleeping.
	line, err := bufio.NewReader(stdout).ReadString('\n')
	if err != nil || strings.TrimSpace(line) != "ready" {
		t.Fatalf("helper: got %q (err: %v), want %q", line, err, "ready")
	}

	// Given: The process is bound to a job.
	job, err := bind(cmd.Process.Pid)
	if err != nil {
		t.Fatalf("err: got %v, want nil", err)
	}

	// When: The job's handle is closed, as it is when the shim exits.
	if err := windows.CloseHandle(job); err != nil {
		t.Fatal(err)
	}

	// Then: The process ends.
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("process outlived its job")
	}
}

/* -------------------------------------------------------------------------- */
/*                         Function: TestHelperProcess                        */
/* -------------------------------------------------------------------------- */

// TestHelperProcess stands in for Godot when a test runs this binary with an
// argument after '--'. A number makes it exit with that code, and 'sleep' makes
// it run until it's ended; run as an ordinary test, it does nothing.
func TestHelperProcess(t *testing.T) {
	t.Parallel()

	args := flag.Args()
	if len(args) != 1 {
		return
	}

	if args[0] == "sleep" {
		fmt.Println("ready")
		time.Sleep(time.Minute)

		return
	}

	code, err := strconv.Atoi(args[0])
	if err != nil {
		t.Fatal(err)
	}

	os.Exit(code)
}

/* -------------------------------------------------------------------------- */
/*                            Function: helperArgs                            */
/* -------------------------------------------------------------------------- */

// helperArgs returns the arguments that run this test binary as
// 'TestHelperProcess' with arg.
func helperArgs(arg string) []string {
	return []string{"-test.run=^TestHelperProcess$", "--", arg}
}
