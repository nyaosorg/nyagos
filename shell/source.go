package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zetamatta/go-mbcs"

	"github.com/zetamatta/nyagos/dos"
)

func loadEnvFile(fname string, verbose io.Writer) error {
	fp, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fp.Close()

	scan := bufio.NewScanner(fp)
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
			if left != "ERRORLEVEL_" {
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
	return scan.Err()
}

func loadPwdFile(fname string, verbose io.Writer) error {
	fp, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fp.Close()
	scan := bufio.NewScanner(fp)
	if !scan.Scan() {
		return fmt.Errorf("Could not load the new current directory from %s", fname)
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

func callBatch(batch string,
	args []string,
	env string,
	pwd string,
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
	if verbose != nil {
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
		// chop last '\0'
		if ansi[len(ansi)-1] == 0 {
			ansi = ansi[:len(ansi)-1]
		}
		fmt.Fprintf(writer, " %s", ansi)
	}
	fmt.Fprintf(writer, "\n@set \"ERRORLEVEL_=%%ERRORLEVEL%%\"\n")
	fmt.Fprintf(writer, "@set > \"%s\"\n", env)
	fmt.Fprintf(writer, "@cd > \"%s\"\n", pwd)
	fmt.Fprintf(writer, "@exit /b \"%%ERRORLEVEL_%%\"\n")
	writer.Flush()
	if err := fd.Close(); err != nil {
		return 1, err
	}

	cmd2 := exec.Cmd{
		Path:   params[0],
		Args:   params,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
	if err := cmd2.Run(); err != nil {
		return 1, err
	}
	errorlevel, errorlevelOk := dos.GetErrorLevel(&cmd2)
	if !errorlevelOk {
		errorlevel = 255
	}
	return errorlevel, nil
}

func RawSource(args []string, verbose io.Writer, debug bool, stdin io.Reader, stdout io.Writer, stderr io.Writer) (int, error) {
	tempDir := os.TempDir()
	pid := os.Getpid()
	batch := filepath.Join(tempDir, fmt.Sprintf("nyagos-%d.cmd", pid))
	env := filepath.Join(tempDir, fmt.Sprintf("nyagos-%d.tmp", pid))
	pwd := filepath.Join(tempDir, fmt.Sprintf("nyagos_%d.tmp", pid))

	errorlevel, err := callBatch(
		batch,
		args,
		env,
		pwd,
		verbose,
		stdin,
		stdout,
		stderr)

	if err != nil {
		return -1, err
	}

	if !debug {
		defer os.Remove(env)
		defer os.Remove(pwd)
		defer os.Remove(batch)
	}

	if err != nil {
		return errorlevel, err
	}

	if err := loadEnvFile(env, verbose); err != nil {
		return 1, err
	}

	if err := loadPwdFile(pwd, verbose); err != nil {
		return 1, err
	}
	return errorlevel, err
}

func Source(args []string, verbose io.Writer, debug bool, stdin io.Reader, stdout io.Writer, stderr io.Writer) (int, error) {
	rawArgs := make([]string, 0, len(args))
	for _, s := range args {
		if strings.ContainsAny(s, " \r\n\v\t\f<>&|") {
			rawArgs = append(rawArgs, fmt.Sprintf("\"%s\"", s))
		} else {
			rawArgs = append(rawArgs, s)
		}
	}
	return RawSource(rawArgs, verbose, debug, stdin, stdout, stderr)
}
