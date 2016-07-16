package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zetamatta/go-mbcs"

	"../dos"
)

func cmd_source(cmd *exec.Cmd) (int, error) {
	args := cmd.Args
	verbose := false
	if len(args) >= 2 && args[1] == "-v" {
		verbose = true
		args = args[1:]
	}
	if len(cmd.Args) < 2 {
		return 255, nil
	}
	tempDir := os.TempDir()
	pid := os.Getpid()
	batchPath := filepath.Join(tempDir, fmt.Sprintf("nyagos-%d.cmd", pid))
	envTxtPath := filepath.Join(tempDir, fmt.Sprintf("nyagos-%d.tmp", pid))
	pwdTxtPath := filepath.Join(tempDir, fmt.Sprintf("nyagos_%d.tmp", pid))

	params := []string{
		os.Getenv("COMSPEC"),
		"/C",
		batchPath,
	}
	batchFd, batchFd_err := os.Create(batchPath)
	if batchFd_err != nil {
		return -1, batchFd_err
	}
	var batchWriter io.Writer
	if verbose {
		batchWriter = io.MultiWriter(batchFd, cmd.Stdout)
	} else {
		batchWriter = batchFd
	}
	fmt.Fprint(batchWriter, "@call")
	for _, v := range args[1:] {
		if strings.ContainsRune(v, ' ') {
			fmt.Fprintf(batchWriter, " \"%s\"", v)
		} else {
			fmt.Fprintf(batchWriter, " %s", v)
		}
	}
	fmt.Fprintf(batchWriter, "\n@set \"ERRORLEVEL_=%%ERRORLEVEL%%\"\n")
	fmt.Fprintf(batchWriter, "@set > \"%s\"\n", envTxtPath)
	fmt.Fprintf(batchWriter, "@cd > \"%s\"\n", pwdTxtPath)
	fmt.Fprintf(batchWriter, "@exit /b \"%%ERRORLEVEL_%%\"\n")
	batchFd.Close()
	defer os.Remove(batchPath)

	cmd2 := exec.Cmd{Path: params[0], Args: params}
	if err := cmd2.Run(); err != nil {
		return 1, err
	}
	errorlevel, errorlevelOk := dos.GetErrorLevel(&cmd2)
	if !errorlevelOk {
		errorlevel = 255
	}
	defer os.Remove(envTxtPath)
	defer os.Remove(pwdTxtPath)

	fp, err := os.Open(envTxtPath)
	if err != nil {
		return 1, err
	}
	defer fp.Close()

	br := bufio.NewReader(fp)
	for {
		lineB, readErr := br.ReadBytes(byte('\n'))
		if readErr != nil {
			if readErr != io.EOF {
				fmt.Fprintf(cmd.Stderr, "%s: %s (environment-readline error)\n", envTxtPath, readErr.Error())
			}
			break
		}
		line, atouErr := mbcs.AtoU(lineB)
		if atouErr != nil {
			fmt.Fprintf(cmd.Stderr, "%s: %s(environment-ansi-to-unicode error)\n", envTxtPath, atouErr.Error())
			continue
		}
		line = strings.TrimSpace(line)
		eqlPos := strings.Index(line, "=")
		if eqlPos > 0 {
			left := line[:eqlPos]
			right := line[eqlPos+1:]
			if left != "ERRORLEVEL_" {
				if verbose {
					fmt.Fprintf(cmd.Stdout, "%s=%s\n", left, right)
				}
				os.Setenv(left, right)
			}
		}
	}

	fp2, err2 := os.Open(pwdTxtPath)
	if err2 != nil {
		return 1, err2
	}
	defer fp2.Close()
	br2 := bufio.NewReader(fp2)
	lineB, lineErr := br2.ReadBytes(byte('\n'))
	if lineErr != nil {
		return 1, errors.New("source : could not get current-directory")
	}
	line, err := mbcs.AtoU(lineB)
	if err == nil {
		line = strings.TrimSpace(line)
		if verbose {
			fmt.Fprintf(cmd.Stdout, "cd \"%s\"\n", line)
		}
		os.Chdir(line)
	}
	return errorlevel, nil
}
