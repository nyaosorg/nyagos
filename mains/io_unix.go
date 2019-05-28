// +build !windows

package mains

import "os/exec"

func newCommand(command string) *exec.Cmd {
	return exec.Command("/bin/sh", "-c", command)
}
