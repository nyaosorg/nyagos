package shell

import (
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/zetamatta/nyagos/dos"
)

func (cmd *Cmd) startProcess() (int, error) {
	if cmd.UseShellExecute {
		// GUI Application
		cmdline := makeCmdline(cmd.args[1:], cmd.rawArgs[1:])
		return 0, dos.ShellExecute("open", cmd.args[0], cmdline, "")
	}
	if UseSourceRunBatch {
		lowerName := strings.ToLower(cmd.args[0])
		if strings.HasSuffix(lowerName, ".cmd") || strings.HasSuffix(lowerName, ".bat") {
			rawargs := cmd.RawArgs()
			args := make([]string, len(rawargs))
			args[0] = encloseWithQuote(cmd.args[0])
			for i, end := 1, len(rawargs); i < end; i++ {
				args[i] = rawargs[i]
			}
			// Batch files
			return RawSource(args, ioutil.Discard, false, cmd.Stdin, cmd.Stdout, cmd.Stderr)
		}
	}

	cmdline := makeCmdline(cmd.args, cmd.rawArgs)

	procAttr := &os.ProcAttr{
		Env:   os.Environ(),
		Files: []*os.File{cmd.Stdin, cmd.Stdout, cmd.Stderr},
		Sys:   &syscall.SysProcAttr{CmdLine: cmdline},
	}

	process, err := os.StartProcess(cmd.args[0], cmd.args[1:], procAttr)
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
	return dos.IsGui(path)
}
