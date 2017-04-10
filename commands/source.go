package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zetamatta/go-mbcs"

	"../dos"
	"../shell"
)

func load_envfile(fname string, verbose io.Writer) error {
	fp, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fp.Close()

	scan := bufio.NewScanner(fp)
	for scan.Scan() {
		line, err := mbcs.AtoU(scan.Bytes())
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		eqlPos := strings.Index(line, "=")
		if eqlPos > 0 {
			left := line[:eqlPos]
			right := line[eqlPos+1:]
			if left != "ERRORLEVEL_" {
				if verbose != nil {
					fmt.Fprintf(verbose, "%s=%s\n", left, right)
				}
				os.Setenv(left, right)
			}
		}
	}
	if err := scan.Err(); err != nil {
		return err
	}
	return nil
}

func load_pwdfile(fname string, verbose io.Writer) error {
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

func call_batch(batch string, args []string, env string, pwd string, verbose io.Writer) (int, error) {
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
	for _, v := range args {
		if strings.ContainsRune(v, ' ') {
			fmt.Fprintf(writer, " \"%s\"", v)
		} else {
			fmt.Fprintf(writer, " %s", v)
		}
	}
	fmt.Fprintf(writer, "\n@set \"ERRORLEVEL_=%%ERRORLEVEL%%\"\n")
	fmt.Fprintf(writer, "@set > \"%s\"\n", env)
	fmt.Fprintf(writer, "@cd > \"%s\"\n", pwd)
	fmt.Fprintf(writer, "@exit /b \"%%ERRORLEVEL_%%\"\n")
	writer.Flush()
	fd.Close()

	cmd2 := exec.Cmd{Path: params[0], Args: params}
	if err := cmd2.Run(); err != nil {
		return 1, err
	}
	errorlevel, errorlevelOk := dos.GetErrorLevel(&cmd2)
	if !errorlevelOk {
		errorlevel = 255
	}
	return errorlevel, nil
}

func cmd_source(ctx context.Context, cmd *shell.Cmd) (int, error) {
	var verbose io.Writer
	args := make([]string, 0, len(cmd.Args))
	debug := false
	for _, arg1 := range cmd.Args[1:] {
		switch arg1 {
		case "-v":
			verbose = cmd.Stderr
		case "-d":
			debug = true
		default:
			args = append(args, arg1)
		}
	}
	if len(cmd.Args) <= 0 {
		return 255, nil
	}

	tempDir := os.TempDir()
	pid := os.Getpid()
	batch := filepath.Join(tempDir, fmt.Sprintf("nyagos-%d.cmd", pid))
	env := filepath.Join(tempDir, fmt.Sprintf("nyagos-%d.tmp", pid))
	pwd := filepath.Join(tempDir, fmt.Sprintf("nyagos_%d.tmp", pid))

	errorlevel, err := call_batch(batch, args, env, pwd, verbose)

	if !debug {
		defer os.Remove(env)
		defer os.Remove(pwd)
		defer os.Remove(batch)
	}

	if err != nil {
		return errorlevel, err
	}

	if err := load_envfile(env, verbose); err != nil {
		return 1, err
	}

	if err := load_pwdfile(pwd, verbose); err != nil {
		return 1, err
	}

	return errorlevel, nil
}
