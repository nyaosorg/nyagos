package shell

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/zetamatta/go-texts/mbcs"

	"github.com/zetamatta/nyagos/dos"
)

func readEnv(scan *bufio.Scanner, verbose io.Writer) (int, error) {
	errorlevel := -1
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		eqlPos := strings.Index(line, "=")
		if eqlPos > 0 {
			left := line[:eqlPos]
			right := line[eqlPos+1:]
			if left == "ERRORLEVEL_" {
				value, err := strconv.ParseInt(right, 10, 32)
				if err != nil {
					if verbose != nil {
						fmt.Fprintf(verbose, "Could not read ERRORLEVEL(%s)\n", right)
					}
				} else {
					errorlevel = int(value)
				}
			} else {
				orig := os.Getenv(left)
				if verbose != nil {
					fmt.Fprintf(verbose, "%s=%s\n", left, right)
				}
				if orig != right {
					// fmt.Fprintf(os.Stderr, "%s:=%s\n", left, right)
					os.Setenv(left, right)
				}
			}
		}
	}
	return errorlevel, scan.Err()
}

func readPwd(scan *bufio.Scanner, verbose io.Writer) error {
	if !scan.Scan() {
		if err := scan.Err(); err != nil {
			return err
		} else {
			return io.EOF
		}
	}
	line := strings.TrimSpace(scan.Text())
	if verbose != nil {
		fmt.Fprintf(verbose, "cd \"%s\"\n", line)
	}
	os.Chdir(line)
	return nil
}

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

type CmdExe struct {
	Cmdline string
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Env     []string
	OnExec  func(int)
	OnDone  func(int)
}

func (this CmdExe) Call() (int, error) {

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

type Source struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Env     []string
	OnExec  func(int)
	OnDone  func(int)
	Args    []string
	Verbose io.Writer
	Debug   bool
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
	}.Call()
}

// RawSource calls the batchfiles and load the changed variable the batchfile has done.
func RawSource(args []string, verbose io.Writer, debug bool, stdin io.Reader, stdout, stderr io.Writer, env []string) (int, error) {
	return Source{
		Args:    args,
		Verbose: verbose,
		Debug:   debug,
		Stdin:   stdin,
		Stdout:  stdout,
		Stderr:  stderr,
		Env:     env,
	}.Call()
}

func (this Source) Call() (int, error) {
	tempDir := os.TempDir()
	pid := os.Getpid()
	tmpfile := filepath.Join(tempDir, fmt.Sprintf("nyagos-%d-%d.tmp", pid, rand.Int()))

	errorlevel, err := this.callBatch(tmpfile)

	if err != nil {
		return errorlevel, err
	}

	if !this.Debug {
		defer os.Remove(tmpfile)
	}

	if errorlevel, err = loadTmpFile(tmpfile, this.Verbose); err != nil {
		if os.IsNotExist(err) {
			return 1, fmt.Errorf("%s: the batch file may use `exit` without `/b` option. Could not find the change of the environment variables", this.Args[0])
		}
		return 1, err
	}
	return errorlevel, err
}
