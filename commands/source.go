package commands

import "bufio"
import "path/filepath"
import "fmt"
import "os"
import "os/exec"
import "strings"

import "../interpreter"

func cmd_source(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	if len(cmd.Args) < 2 {
		return interpreter.CONTINUE, nil
	}
	envTxtPath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("nyagos-%d.tmp", os.Getpid()))

	args := []string{
		os.Getenv("COMSPEC"),
		"/C",
	}
	for _, v := range cmd.Args[1:] {
		args = append(args,
			strings.Replace(
				strings.Replace(
					strings.Replace(v, " ", "^ ", -1), "(", "^(", -1),
				")", "^)", -1))
	}
	args = append(args, "&", "set", ">", envTxtPath)

	cmd2 := exec.Cmd{Path: args[0], Args: args}
	if err := cmd2.Run(); err != nil {
		return interpreter.CONTINUE, err
	}
	fp, err := os.Open(envTxtPath)
	if err != nil {
		return interpreter.CONTINUE, err
	}
	defer os.Remove(envTxtPath)
	defer fp.Close()

	for scr := bufio.NewScanner(fp); scr.Scan(); {
		line := scr.Text()
		eqlPos := strings.Index(line, "=")
		if eqlPos > 0 {
			os.Setenv(line[:eqlPos], line[eqlPos+1:])
		}
	}
	return interpreter.CONTINUE, nil
}
