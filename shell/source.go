package shell

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type scannerT interface {
	Scan() bool
	Err() error
	Text() string
}

func readEnv(scan scannerT, verbose io.Writer) (int, error) {
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

func readPwd(scan scannerT, verbose io.Writer) error {
	if !scan.Scan() {
		if err := scan.Err(); err != nil {
			return err
		}
		return io.EOF
	}
	line := strings.TrimSpace(scan.Text())
	if verbose != nil {
		fmt.Fprintf(verbose, "cd \"%s\"\n", line)
	}
	os.Chdir(line)
	return nil
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

func (source Source) Call() (int, error) {
	tempDir := os.TempDir()
	pid := os.Getpid()
	tmpfile := filepath.Join(tempDir, fmt.Sprintf("nyagos-%d-%d.tmp", pid, rand.Int()))

	errorlevel, err := source.callBatch(tmpfile)

	if err != nil {
		return errorlevel, err
	}

	if !source.Debug {
		defer os.Remove(tmpfile)
	}

	if errorlevel, err = loadTmpFile(tmpfile, source.Verbose); err != nil {
		if os.IsNotExist(err) {
			return 1, fmt.Errorf("%s: the batch file may use `exit` without `/b` option. Could not find the change of the environment variables", source.Args[0])
		}
		return 1, err
	}
	return errorlevel, err
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

type CmdExe struct {
	Cmdline string
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Env     []string
	OnExec  func(int)
	OnDone  func(int)
}

func (cmdExe CmdExe) Run() (int, error) {
	return cmdExe.run()
}
