// +build !windows

package shell

import (
	"os"
	"os/exec"
	"syscall"
)

func (cmd *Cmd) lookpath() string {
	path, err := exec.LookPath(cmd.args[0])
	if err != nil {
		return ""
	}
	return path
}

func (cmd *Cmd) startProcess() (int, error) {
	procAttr := &os.ProcAttr{
		Env:   os.Environ(),
		Files: []*os.File{cmd.Stdin, cmd.Stdout, cmd.Stderr},
	}
	process, err := os.StartProcess(cmd.args[0], cmd.args, procAttr)
	if err != nil {
		return 255, err
	}
	processState, err := process.Wait()
	if err != nil {
		return 254, err
	}
	if processState.Success() {
		return 0, nil
	}
	if t, ok := processState.Sys().(syscall.WaitStatus); ok {
		return t.ExitStatus(), nil
	}
	return 253, nil
}

func isGui(path string) bool {
	return false
}
