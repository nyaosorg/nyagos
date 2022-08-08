//go:build !windows
// +build !windows

package shell

import (
	"context"
	"os"
	"os/exec"
)

func (cmd *Cmd) lookpath() string {
	path, err := exec.LookPath(cmd.args[0])
	if err != nil {
		return ""
	}
	return path
}

func (cmd *Cmd) startProcess(ctx context.Context) (int, error) {
	procAttr := &os.ProcAttr{
		Env:   cmd.DumpEnv(),
		Files: []*os.File{cmd.Stdio[0], cmd.Stdio[1], cmd.Stdio[2]},
	}
	return startAndWaitProcess(ctx, cmd.args[0], cmd.args, procAttr, cmd.OnBackExec, cmd.OnBackDone)
}

func isGui(path string) bool {
	return false
}
