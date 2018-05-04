package shell

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zetamatta/go-mbcs"

	"github.com/zetamatta/nyagos/dos"
)

func readEnv(scan *bufio.Scanner, verbose io.Writer) (int, error) {
	errorlevel := -1
	for scan.Scan() {
		line, err := mbcs.AtoU(scan.Bytes())
		if err != nil {
			continue
		}
		line = strings.TrimSpace(line)
		eqlPos := strings.Index(line, "=")
		if eqlPos > 0 {
			left := line[:eqlPos]
			right := line[eqlPos+1:]
			if left == "ERRORLEVEL_" {
				value, err := strconv.ParseUint(right, 10, 32)
				if err != nil {
					fmt.Fprintf(verbose, "Could not read ERRORLEVEL(%s)\n", right)
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
		return errors.New("Could not load the new current directory")
	}
	if err := scan.Err(); err != nil {
		return err
	}
	line, err := mbcs.AtoU(scan.Bytes())
	if err != nil {
		return err
	}
	line = strings.TrimSpace(line)
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

	scan := bufio.NewScanner(fp)
	if err := readPwd(scan, verbose); err != nil {
		return -1, err
	}
	return readEnv(scan, verbose)
}

func callBatch(batch string,
	args []string,
	tmpfile string,
	verbose io.Writer,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer) (int, error) {
	params := []string{
		os.Getenv("COMSPEC"),
		"/C",
		batch,
	}
	fd, err := os.Create(batch)
	if err != nil {
		return 1, err
	}
	var writer *bufio.Writer
	if verbose != nil && verbose != ioutil.Discard {
		writer = bufio.NewWriter(io.MultiWriter(fd, verbose))
	} else {
		writer = bufio.NewWriter(fd)
	}
	fmt.Fprint(writer, "@call")
	for _, arg1 := range args {
		// UTF8 parameter to ANSI
		ansi, err := mbcs.UtoA(arg1)
		if err != nil {
			// println("utoa: " + err.Error())
			fd.Close()
			return -1, err
		}
		ansi = bytes.TrimSuffix(ansi, []byte{0})
		fmt.Fprintf(writer, " %s", ansi)
	}
	fmt.Fprintf(writer, "\r\n@set \"ERRORLEVEL_=%%ERRORLEVEL%%\"\r\n")

	// Sometimes %TEMP% has not ASCII letters.
	ansi, err := mbcs.UtoA(tmpfile)
	if err != nil {
		fd.Close()
		return -1, err
	}
	ansi = bytes.TrimSuffix(ansi, []byte{0})
	fmt.Fprintf(writer, "@(cd & set) > \"%s\"\r\n", ansi)
	fmt.Fprintf(writer, "@exit /b \"%%ERRORLEVEL_%%\"\r\n")
	writer.Flush()
	if err := fd.Close(); err != nil {
		return 1, err
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

func RawSource(args []string, verbose io.Writer, debug bool, stdin io.Reader, stdout io.Writer, stderr io.Writer) (int, error) {
	tempDir := os.TempDir()
	pid := os.Getpid()
	batch := filepath.Join(tempDir, fmt.Sprintf("nyagos-%d.cmd", pid))
	tmpfile := filepath.Join(tempDir, fmt.Sprintf("nyagos-%d.tmp", pid))

	errorlevel, err := callBatch(
		batch,
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
		defer os.Remove(batch)
	}

	if errorlevel, err = loadTmpFile(tmpfile, verbose); err != nil {
		if os.IsNotExist(err) {
			return 1, fmt.Errorf("%s: the batch file may use `exit` without `/b` option. Could not find the change of the environment variables", args[0])
		}
		return 1, err
	}
	// println("ERRORLEVEL=", errorlevel)
	if err != nil {
		return errorlevel, err
	}
	if errorlevel != 0 {
		return errorlevel, fmt.Errorf("exit status %d", errorlevel)
	}
	return 0, nil
}
