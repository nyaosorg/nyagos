package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"../dos"
	"../interpreter"
)

func cmd_source(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	args := cmd.Args
	verbose := false
	if len(args) >= 2 && args[1] == "-v" {
		verbose = true
		args = args[1:]
	}
	if len(cmd.Args) < 2 {
		return interpreter.CONTINUE, nil
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
		return interpreter.CONTINUE, err
	}
	defer os.Remove(envTxtPath)
	defer os.Remove(pwdTxtPath)

	fp, err := os.Open(envTxtPath)
	if err != nil {
		return interpreter.CONTINUE, err
	}
	defer fp.Close()

	for {
		line, lineErr := dos.ReadAnsiLine(fp)
		if lineErr != nil {
			break
		}
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
		return interpreter.CONTINUE, err2
	}
	defer fp2.Close()
	line, lineErr := dos.ReadAnsiLine(fp2)
	if lineErr == nil {
		os.Chdir(line)
	}
	return interpreter.CONTINUE, nil
}
