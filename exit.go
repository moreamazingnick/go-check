package check

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/mitchellh/go-ps"
)

// AllowExit lets you disable the call to os.Exit() in ExitXxx() functions of this package.
//
// This should be used carefully and most likely only for testing.
var AllowExit = true

// PrintStack prints the error stack when recovering from a panic with CatchPanic()
var PrintStack = true

// Exitf prints the plugin output using formatting and exits the program.
//
// Output is the formatting string, and the rest of the arguments help adding values.
//
// Also see fmt package: https://golang.org/pkg/fmt
func Exitf(rc int, output string, args ...interface{}) {
	ExitRaw(rc, fmt.Sprintf(output, args...))
}

// ExitRaw prints the plugin output with the state prefixed and exits the program.
//
// Example:
//	OK - everything is fine
func ExitRaw(rc int, output ...string) {
	text := StatusText(rc) + " -"

	for _, s := range output {
		text += " " + s
	}

	text += "\n"

	_, _ = os.Stdout.WriteString(text)

	BaseExit(rc)
}

// Exit prints the plugin output and exits the program
//
// Deprecated, please use Exitf or ExitRaw.
func Exit(rc int, output string, args ...interface{}) {
	Exitf(rc, output, args...)
}

func BaseExit(rc int) {
	if AllowExit {
		os.Exit(rc)
	} else {
		_, _ = os.Stdout.WriteString("would exit with code " + strconv.Itoa(rc) + "\n")
	}
}

// ExitError exists with an Unknown state while reporting the error
func ExitError(err error) {
	Exitf(Unknown, "%s (%T)", err.Error(), err)
}

// CatchPanic is a general function for defer, to capture any panic that occurred during runtime of a check
//
// The function will recover from the condition and exit with a proper UNKNOWN status, while showing error
// and the call stack.
func CatchPanic() {
	ppid := os.Getppid()
	if parent, err := ps.FindProcess(ppid); err == nil {
		if parent.Executable() == "dlv" {
			// seems to be a debugger, don't continue with recover
			return
		}
	}

	if r := recover(); r != nil {
		output := fmt.Sprint("Golang encountered a panic: ", r)
		if PrintStack {
			output += "\n\n" + string(debug.Stack())
		}

		ExitRaw(Unknown, output)
	}
}
