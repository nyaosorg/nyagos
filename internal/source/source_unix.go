//go:build !windows
// +build !windows

package source

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func loadTmpFile(fname string, verbose io.Writer) (int, error) {
	fp, err := os.Open(fname)
	if err != nil {
		return -1, err
	}
	defer fp.Close()

	scan := bufio.NewScanner(fp)
	if err := readPwd(scan, verbose); err != nil {
		return -1, err
	}
	return readEnv(scan, verbose)
}

func (system *System) run() (int, error) {
	args := []string{
		"/bin/sh",
		"-c",
		system.Cmdline,
	}
	cmd := &exec.Cmd{
		Path:   "/bin/sh",
		Args:   args,
		Stdin:  system.Stdin,
		Stdout: system.Stdout,
		Stderr: system.Stderr,
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

	cmdline.WriteByte('.')

	if fullpath, err := filepath.Abs(strings.ReplaceAll(batch.Args[0], `"`, ``)); err == nil {
		fmt.Fprintf(&cmdline, ` "%s"`, fullpath)
	} else {
		cmdline.WriteByte(' ')
		cmdline.WriteString(batch.Args[0])
	}

	for _, arg1 := range batch.Args[1:] {
		cmdline.WriteByte(' ')
		cmdline.WriteString(arg1)
	}
	cmdline.WriteString(` ; (pwd ; env) > '`)
	cmdline.WriteString(tmpfile)
	cmdline.WriteString(`'`)

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
