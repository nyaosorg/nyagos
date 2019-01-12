package shell

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/nodos"
)

func (cmd *Cmd) lookpath() string {
	return nodos.LookPath(LookCurdirOrder, cmd.args[0], "NYAGOSPATH")
}

func (cmd *Cmd) startProcess(ctx context.Context) (int, error) {
	if cmd.UseShellExecute {
		// GUI Application
		cmdline := makeCmdline(cmd.args[1:], cmd.rawArgs[1:])
		return 0, dos.ShellExecute("open", cmd.args[0], cmdline, "")
	}
	if closer, err := dos.ChangeConsoleMode(windows.Stdin, dos.ModeSet(0x7)); err == nil {
		defer closer()
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
	return startAndWaitProcess(ctx, cmd.args[0], cmd.args, procAttr)
}

func isGui(path string) bool {
	return dos.IsGui(path)
}
