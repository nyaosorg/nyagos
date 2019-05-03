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
		Files: []*os.File{cmd.Stdin, cmd.Stdout, cmd.Stderr},
	}
	return startAndWaitProcess(ctx, cmd.args[0], cmd.args, procAttr, cmd.report)
}

func isGui(path string) bool {
	return false
}
