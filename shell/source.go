package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

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

func callBatch(
	args []string,
	tmpfile string,
	verbose io.Writer,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer) (int, error) {

	var cmdline strings.Builder

	cmdline.WriteString("call")
	for _, arg1 := range args {
		cmdline.WriteByte(' ')
		cmdline.WriteString(arg1)
	}
	cmdline.WriteString(` & set "ERRORLEVEL_=!ERRORLEVEL!" & (cd & set) > "`)
	cmdline.WriteString(tmpfile)
	cmdline.WriteString(`"`)

	backup := os.Getenv("NYAGOSCMDLINE")
	os.Setenv("NYAGOSCMDLINE", cmdline.String())
	defer os.Setenv("NYAGOSCMDLINE", backup)

	params := []string{
		os.Getenv("COMSPEC"),
		"/V:ON",
		"/C",
		"%NYAGOSCMDLINE%",
	}

	cmd := exec.Cmd{
		Path:   params[0],
		Args:   params,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
	if err := cmd.Run(); err != nil {
		return 1, err
	}
	errorlevel, errorlevelOk := dos.GetErrorLevel(&cmd)
	if !errorlevelOk {
		errorlevel = 255
	}
	return errorlevel, nil
}

// RawSource calls the batchfiles and load the changed variable the batchfile has done.
func RawSource(args []string, verbose io.Writer, debug bool, stdin io.Reader, stdout io.Writer, stderr io.Writer) (int, error) {
	tempDir := os.TempDir()
	pid := os.Getpid()
	tmpfile := filepath.Join(tempDir, fmt.Sprintf("nyagos-%d.tmp", pid))

	errorlevel, err := callBatch(
		args,
		tmpfile,
		verbose,
		stdin,
		stdout,
		stderr)

	if err != nil {
		return errorlevel, err
	}

	if !debug {
		defer os.Remove(tmpfile)
	}

	if errorlevel, err = loadTmpFile(tmpfile, verbose); err != nil {
		if os.IsNotExist(err) {
			return 1, fmt.Errorf("%s: the batch file may use `exit` without `/b` option. Could not find the change of the environment variables", args[0])
		}
		return 1, err
	}
	if err != nil {
		return errorlevel, err
	}
	if errorlevel != 0 {
		return errorlevel, fmt.Errorf("exit status %d", errorlevel)
	}
	return 0, nil
}
