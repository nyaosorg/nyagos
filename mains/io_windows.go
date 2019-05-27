package mains

import (
	"os/exec"
	"syscall"
)

func newCommand(command string) *exec.Cmd {
	xcmd := exec.Command("cmd.exe", "/S", "/C", string(command+` `))
	xcmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: `/S /C "` + string(command) + ` "`,
	}
	return xcmd
}
