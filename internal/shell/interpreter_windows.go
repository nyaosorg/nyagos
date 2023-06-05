package shell

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"

	"github.com/nyaosorg/go-windows-su"

	"github.com/nyaosorg/nyagos/internal/nodos"
	"github.com/nyaosorg/nyagos/internal/source"
)

func encloseWithQuote(fullpath string) string {
	if strings.ContainsRune(fullpath, ' ') {
		var f strings.Builder
		f.WriteByte('"')
		f.WriteString(fullpath)
		f.WriteByte('"')
		return f.String()
	}
	return fullpath
}

func (cmd *Cmd) lookpath() string {
	return nodos.LookPath(LookCurdirOrder, cmd.args[0], "NYAGOSPATH")
}

func (cmd *Cmd) startProcess(ctx context.Context) (int, error) {
	if cmd.UseShellExecute {
		// GUI Application
		cmdline := makeCmdline(cmd.rawArgs[1:])
		truepath := cmd.args[0]
		if _truepath, err := filepath.Abs(truepath); err == nil {
			truepath = _truepath
		}
		if _truepath, err := filepath.EvalSymlinks(truepath); err == nil {
			truepath = _truepath
		}
		pid, err := su.ShellExecute("open", truepath, cmdline, "")
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
	if closer, err := nodos.ChangeConsoleMode(windows.Stdin,
		nodos.ModeSet(
			windows.ENABLE_PROCESSED_INPUT|
				windows.ENABLE_LINE_INPUT|
				windows.ENABLE_ECHO_INPUT)); err == nil {
		defer closer()
	}
	if closer, err := nodos.ChangeConsoleMode(windows.Stdout); err == nil {
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
			return source.Batch{
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

	cmdline := makeCmdline(cmd.rawArgs)

	procAttr := &os.ProcAttr{
		Env:   cmd.DumpEnv(),
		Files: cmd.Stdio[:],
		Sys:   &syscall.SysProcAttr{CmdLine: cmdline},
	}
	return startAndWaitProcess(ctx, cmd.args[0], cmd.args, procAttr, cmd.OnBackExec, cmd.OnBackDone)
}

func isGui(path string) bool {
	return nodos.IsGui(path)
}
