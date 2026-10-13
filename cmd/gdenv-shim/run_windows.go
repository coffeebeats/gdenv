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

// run starts Godot, waits for it, and returns the code it exited with, since
// Windows can't replace the shim's process with Godot's.
//
// NOTE: Godot runs in a job that kills it if the shim ends first, as when a
// caller terminates the shim; what Godot starts is left out of the job.
func run(binary string, args ...string) (int, error) {
	// Godot answers Ctrl+C on the shared console, so the shim catches it and
	// waits; 'signal.Ignore' would let Windows' default handler end the shim.
	signal.Notify(make(chan os.Signal, 1), os.Interrupt)

	cmd := exec.Command(binary, args...) //nolint:gosec,noctx

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return 0, err
	}

	// Closing the job after Godot exits ends nothing. If 'bind' fails, Godot runs
	// unbound.
	if job, err := bind(cmd.Process.Pid); err == nil {
		defer windows.CloseHandle(job) //nolint:errcheck
	}

	err := cmd.Wait()
	if exit, ok := errors.AsType[*exec.ExitError](err); ok {
		return exit.ExitCode(), nil
	}

	return 0, err
}

/* ----------------------------- Function: bind ----------------------------- */

// bind places the process 'pid' in a job that ends it once the job's handle
// is closed, and returns that handle.
func bind(pid int) (windows.Handle, error) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return 0, err
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

		return 0, err
	}

	process, err := windows.OpenProcess(
		windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE,
		false,
		uint32(pid), //nolint:gosec
	)
	if err != nil {
		_ = windows.CloseHandle(job)

		return 0, err
	}

	defer windows.CloseHandle(process) //nolint:errcheck

	err = windows.AssignProcessToJobObject(job, process)
	if err != nil {
		_ = windows.CloseHandle(job)

		return 0, err
	}

	return job, nil
}
