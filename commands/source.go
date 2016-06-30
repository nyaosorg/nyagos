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
	envTxtPath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("nyagos-%d.tmp", os.Getpid()))
	pwdTxtPath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("nyagos_%d.tmp", os.Getpid()))

	params := []string{
		os.Getenv("COMSPEC"),
		"/C",
	}
	for _, v := range args[1:] {
		params = append(params,
			strings.Replace(
				strings.Replace(
					strings.Replace(v, " ", "^ ", -1), "(", "^(", -1),
				")", "^)", -1))
	}
	params = append(params,
		"&", "set", ">", envTxtPath,
		"&", "cd", ">", pwdTxtPath)

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
			if verbose {
				fmt.Fprintf(cmd.Stdout, "%s=%s\n", left, right)
			}
			os.Setenv(left, right)
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
		os.Chdir(line)
	}
	return errorlevel, nil
}
