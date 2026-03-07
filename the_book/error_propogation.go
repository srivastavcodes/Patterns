package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

type Errors struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]any
}

func wrapErrorf(err error, format string, args ...any) Errors {
	return Errors{
		Inner:      err,
		Message:    fmt.Sprintf(format, args...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]any),
	}
}

func (e Errors) Error() string {
	return e.Message
}

// pretend this is a different module

type LowLevelError struct {
	error
}

func isGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelError{wrapErrorf(err, err.Error())}
	}
	return info.Mode().Perm()&0o100 == 0o100, nil
}

// another module which calls functions from the low-leve-module

type IntermediateError struct {
	error
}

func runJob(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGloballyExec(jobBinPath)
	if err != nil {
		// Here we're passing on errors from the low leve module. Because of our architectural
		// decision to consider errors passed from modules without wrapping them in our own
		// types a bug, this will cause us issues later.
		return err
	} else if !isExecutable {
		return wrapErrorf(nil, "job binary %s is not executable", jobBinPath)
	}
	return exec.Command(jobBinPath, "--id="+id).Run()
}

func runJobWithErrsWrapped(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGloballyExec(jobBinPath)
	if err != nil {
		// here we are wrapping the error that we didn't in the previous iteration of this
		// method.
		return &IntermediateError{wrapErrorf(
			err,
			"cannot run job %q: requisite binaries not available", id,
		)}
	} else if !isExecutable {
		return wrapErrorf(nil, "cannot run job %q: requisite binaries not available", id)
	}
	return exec.Command(jobBinPath, "--id="+id).Run()
}

// This is the top-level function that will call functions from the intermediate package.
// This is the user-facing part of the program.

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err) // this log-outs the full error in case someone needs to dig what happened.
	fmt.Printf("[%v] %v\n", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.LUTC)

	err := runJobWithErrsWrapped("1")
	if err != nil {
		msg := "There was an unexpected issue; please report this as a bug"
		// here we check to see if the error is of the expected type. If it is, we know
		// it's a well-crafted error, and we can simply pass its message onto the user.
		if _, ok := errors.AsType[*IntermediateError](err); ok {
			msg = err.Error()
		}
		// here we bind the log message and error together with an ID of 1, we can easily
		// make the ID increase monotonically or use a GUID to ensure a unique ID.
		handleError(1, err, msg)
	}
}
