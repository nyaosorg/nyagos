package shell

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/zetamatta/go-texts/mbcs"

	"github.com/zetamatta/nyagos/dos"
)

// loadTmpFile - read update the current-directory and environment-variables from tmp-file.
func loadTmpFile(fname string, verbose io.Writer) (int, error) {
	fp, err := os.Open(fname)
	if err != nil {
		return -1, err
	}
	defer fp.Close()

	scan := bufio.NewScanner(mbcs.NewAtoUReader(fp, mbcs.ConsoleCP()))
	if err := readPwd(scan, verbose); err != nil {
		return -1, err
	}
	return readEnv(scan, verbose)
}

func (this *CmdExe) run() (int, error) {
	if wd, err := os.Getwd(); err == nil && strings.HasPrefix(wd, `\\`) {
		netdrive, closer := dos.UNCtoNetDrive(wd)
		defer closer()
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
	buffer.WriteString(this.Cmdline)
	buffer.WriteString(` "`)

	cmd := exec.Cmd{
		Path:        cmdexe,
		Stdin:       this.Stdin,
		Stdout:      this.Stdout,
		Stderr:      this.Stderr,
		Env:         this.Env,
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
	if this.OnExec != nil && cmd.Process != nil {
		this.OnExec(cmd.Process.Pid)
	}
	if err := cmd.Wait(); err != nil {
		return -1, err
	}
	if this.OnDone != nil && cmd.Process != nil {
		this.OnDone(cmd.Process.Pid)
	}
	return cmd.ProcessState.ExitCode(), nil
}

func (this *Source) callBatch(tmpfile string) (int, error) {
	var cmdline strings.Builder

	cmdline.WriteString(`call`)
	for _, arg1 := range this.Args {
		cmdline.WriteByte(' ')
		cmdline.WriteString(arg1)
	}
	cmdline.WriteString(` & call set ERRORLEVEL_=%^ERRORLEVEL% & (cd & set) > "`)
	cmdline.WriteString(tmpfile)
	cmdline.WriteString(`"`)

	return CmdExe{
		Cmdline: cmdline.String(),
		Stdin:   this.Stdin,
		Stdout:  this.Stdout,
		Stderr:  this.Stderr,
		Env:     this.Env,
		OnExec:  this.OnExec,
		OnDone:  this.OnDone,
	}.Run()
}
