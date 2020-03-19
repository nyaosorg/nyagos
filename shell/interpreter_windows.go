package shell

import (
	"context"
	"os"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"

	"github.com/zetamatta/go-windows-su"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/nodos"
)

func encloseWithQuote(fullpath string) string {
	if strings.ContainsRune(fullpath, ' ') {
		var f strings.Builder
		f.WriteByte('"')
		f.WriteString(fullpath)
		f.WriteByte('"')
		return f.String()
	} else {
		return fullpath
	}
}

func (cmd *Cmd) lookpath() string {
	return nodos.LookPath(LookCurdirOrder, cmd.args[0], "NYAGOSPATH")
}

func (cmd *Cmd) startProcess(ctx context.Context) (int, error) {
	if cmd.UseShellExecute {
		// GUI Application
		cmdline := makeCmdline(cmd.args[1:], cmd.rawArgs[1:])
		pid, err := su.ShellExecute("open", cmd.args[0], cmdline, "")
		if err == nil && pid != 0 && cmd.OnBackExec != nil {
			cmd.OnBackExec(pid)
			if cmd.OnBackDone != nil {
				if process, err := os.FindProcess(pid); err == nil {
					go func(f func(int)) {
						process.Wait()
						f(pid)
					}(cmd.OnBackDone)
				}
			}
		}
		return 0, err
	}
	if closer, err := dos.ChangeConsoleMode(windows.Stdin,
		dos.ModeSet(
			windows.ENABLE_PROCESSED_INPUT|
				windows.ENABLE_LINE_INPUT|
				windows.ENABLE_ECHO_INPUT)); err == nil {
		defer closer()
	}
	if closer, err := dos.ChangeConsoleMode(windows.Stdout); err == nil {
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
			return Source{
				Args:   args,
				Stdin:  cmd.Stdio[0],
				Stdout: cmd.Stdio[1],
				Stderr: cmd.Stdio[2],
				Env:    cmd.DumpEnv(),
				OnExec: cmd.OnBackExec,
				OnDone: cmd.OnBackDone,
			}.Call()
		}
	}

	cmdline := makeCmdline(cmd.args, cmd.rawArgs)

	procAttr := &os.ProcAttr{
		Env:   cmd.DumpEnv(),
		Files: cmd.Stdio[:],
		Sys:   &syscall.SysProcAttr{CmdLine: cmdline},
	}
	return startAndWaitProcess(ctx, cmd.args[0], cmd.args, procAttr, cmd.OnBackExec, cmd.OnBackDone)
}

func isGui(path string) bool {
	return dos.IsGui(path)
}
