package check

import (
	"fmt"
	"strings"
	"os/exec"
	"bytes"
	"syscall"
)

type ExecStruct struct {
	CmdInput	string
	CmdMain		string
	CmdArgs		[]string

	StdOut		string
	StdErr		string
	Output		string
	ExitCode	int
}

const defaultExitCode = 1

func Exec(input string) *ExecStruct {
	e := &ExecStruct{}
	e.CmdInput = input
	e.ExitCode = defaultExitCode

	// split command
 	cmd_parts := strings.Split(e.CmdInput, " ")
 	e.CmdMain = cmd_parts[0]
	e.CmdArgs = cmd_parts[1:len(cmd_parts)]
	 
	// prepare execution
	var stdout, stderr bytes.Buffer
	x := exec.Command(e.CmdMain, e.CmdArgs...)
	x.Stdout = &stdout
	x.Stderr = &stderr

	// execute
	err := x.Run()
	e.StdOut = stdout.String()
	e.StdErr = stderr.String()

	// get exit code
    if err != nil {
        if exitError, ok := err.(*exec.ExitError); ok {
            ws := exitError.Sys().(syscall.WaitStatus)
            e.ExitCode = ws.ExitStatus()
        } 
    } else {
        ws := x.ProcessState.Sys().(syscall.WaitStatus)
        e.ExitCode = ws.ExitStatus()
	}
	
	// set output
	if err == nil {
		e.Output = e.StdOut
	} else {
		if len(e.StdErr) != 0 {
			e.Output = e.StdErr
		} else {
			e.Output = err.Error()
		}
	}
	e.Output = strings.TrimSpace(e.Output)

	return e
}

func (e *ExecStruct) Dump() {
	fmt.Println("CmdInput: " + e.CmdInput)
	fmt.Println("CmdMain: " + e.CmdMain)
	fmt.Println("CmdArgs:", e.CmdArgs)
	fmt.Println("StdOut: " + e.StdOut)
	fmt.Println("StdErr: " + e.StdErr)
	fmt.Println("Output: " + e.Output)
	fmt.Println("ExitCode:", e.ExitCode)
}
