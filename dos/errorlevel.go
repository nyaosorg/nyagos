package dos

import (
	"os"
	"os/exec"
	"syscall"
)

func getErrorLevel(processState *os.ProcessState) (int, bool) {
	if processState == nil {
		return 255, false
	} else if processState.Success() {
		return 0, true
	} else if t, ok := processState.Sys().(syscall.WaitStatus); ok {
		return t.ExitStatus(), true
	} else {
		return 255, false
	}
}

func GetErrorLevel(cmd *exec.Cmd) (int, bool) {
	return getErrorLevel(cmd.ProcessState)
}
