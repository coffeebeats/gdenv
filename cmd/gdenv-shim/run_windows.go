//go:build windows

package main

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"unsafe"

	"golang.org/x/sys/windows"
)

/* ------------------------------ Function: run ----------------------------- */

// run starts Godot and waits for it, since Windows can't replace the shim's
// process with Godot's, and returns the code Godot exited with.
//
// NOTE: Godot runs in a job that kills it if the shim ends first, as it does
// when a caller terminates the shim; what Godot starts is left out of the job.
func run(binary string, args ...string) (int, error) {
	// Godot shares the shim's console and answers Ctrl+C itself, so the shim
	// waits for it to exit rather than ending first.
	signal.Notify(make(chan os.Signal, 1), os.Interrupt)

	cmd := exec.Command(binary, args...) //nolint:gosec,noctx

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return 0, err
	}

	bind(cmd.Process.Pid)

	err := cmd.Wait()

	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return exit.ExitCode(), nil
	}

	return 0, err
}

/* ----------------------------- Function: bind ----------------------------- */

// bind places the process pid in a job that kills it once the shim exits, and
// does nothing if it can't.
//
// NOTE: The job's handle is never closed, so the job lasts as long as the shim.
func bind(pid int) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return
	}

	var info windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION

	info.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE |
		windows.JOB_OBJECT_LIMIT_SILENT_BREAKAWAY_OK

	_, err = windows.SetInformationJobObject(
		job,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)),
	)
	if err != nil {
		_ = windows.CloseHandle(job)

		return
	}

	process, err := windows.OpenProcess(
		windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE,
		false,
		uint32(pid), //nolint:gosec
	)
	if err != nil {
		_ = windows.CloseHandle(job)

		return
	}

	defer windows.CloseHandle(process) //nolint:errcheck

	if windows.AssignProcessToJobObject(job, process) != nil {
		_ = windows.CloseHandle(job)
	}
}
