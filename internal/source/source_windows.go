package source

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/nyaosorg/go-windows-mbcs"
	"github.com/nyaosorg/go-windows-netresource"
)

// loadTmpFile - read update the current-directory and environment-variables from tmp-file.
func loadTmpFile(fname string, verbose io.Writer) (int, error) {
	fp, err := os.Open(fname)
	if err != nil {
		return -1, err
	}
	defer fp.Close()

	scan := mbcs.NewFilter(fp, mbcs.ConsoleCP())
	if err := readPwd(scan, verbose); err != nil {
		return -1, err
	}
	return readEnv(scan, verbose)
}

func (system *System) run() (int, error) {
	if wd, err := os.Getwd(); err == nil && strings.HasPrefix(wd, `\\`) {
		netdrive, closer := netresource.UNCtoNetDrive(wd)
		defer closer(false, false)
		if netdrive != "" {
			if err := os.Chdir(netdrive); err == nil {
				defer os.Chdir(wd)
			}
		}
	}

	cmdexe := os.Getenv("COMSPEC")

	if cmdexe == "" {
		cmdexe = "cmd.exe"
	}

	var buffer strings.Builder
	buffer.WriteString(`/S /C "`)
	buffer.WriteString(system.Cmdline)
	buffer.WriteString(` "`)

	cmd := &exec.Cmd{
		Path:        cmdexe,
		Stdin:       system.Stdin,
		Stdout:      system.Stdout,
		Stderr:      system.Stderr,
		Env:         system.Env,
		SysProcAttr: &syscall.SysProcAttr{CmdLine: buffer.String()},
	}
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}
	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Start(); err != nil {
		return -1, err
	}
	if system.OnExec != nil && cmd.Process != nil {
		system.OnExec(cmd.Process.Pid)
	}
	if err := cmd.Wait(); err != nil {
		return -1, err
	}
	if system.OnDone != nil && cmd.Process != nil {
		system.OnDone(cmd.Process.Pid)
	}
	return cmd.ProcessState.ExitCode(), nil
}

func (batch *Batch) call(tmpfile string) (int, error) {
	var cmdline strings.Builder

	cmdline.WriteString(`call`)
	for _, arg1 := range batch.Args {
		cmdline.WriteByte(' ')
		cmdline.WriteString(arg1)
	}
	cmdline.WriteString(` & call set ERRORLEVEL_=%^ERRORLEVEL% & (cd & set) > "`)
	cmdline.WriteString(tmpfile)
	cmdline.WriteString(`"`)

	return System{
		Cmdline: cmdline.String(),
		Stdin:   batch.Stdin,
		Stdout:  batch.Stdout,
		Stderr:  batch.Stderr,
		Env:     batch.Env,
		OnExec:  batch.OnExec,
		OnDone:  batch.OnDone,
	}.Run()
}
