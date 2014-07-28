package commands

import "bufio"
import "path/filepath"
import "fmt"
import "os"
import "os/exec"
import "strings"

import "../interpreter"

func cmd_source(cmd *exec.Cmd) (interpreter.NextT, error) {
	if len(cmd.Args) < 2 {
		return interpreter.CONTINUE, nil
	}
	envTxtPath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("nyagos-%d.tmp", os.Getpid()))

	args := make([]string, 2)
	args[0] = os.Getenv("COMSPEC")
	args[1] = "/C"
	for _, v := range cmd.Args[1:] {
		args = append(args,
			strings.Replace(
				strings.Replace(
					strings.Replace(v, " ", "^ ", -1), "(", "^(", -1),
				")", "^)", -1))
	}
	args = append(args, "&")
	args = append(args, "set")
	args = append(args, ">")
	args = append(args, envTxtPath)

	var cmd2 exec.Cmd
	cmd2.Path = args[0]
	cmd2.Args = args
	cmd2.Env = nil
	cmd2.Dir = ""
	if err := cmd2.Run(); err != nil {
		return interpreter.CONTINUE, err
	}
	fp, err := os.Open(envTxtPath)
	if err != nil {
		return interpreter.CONTINUE, err
	}
	defer os.Remove(envTxtPath)
	defer fp.Close()

	scr := bufio.NewScanner(fp)
	for scr.Scan() {
		line := scr.Text()
		eqlPos := strings.Index(line, "=")
		if eqlPos > 0 {
			os.Setenv(line[:eqlPos], line[eqlPos+1:])
		}
	}
	return interpreter.CONTINUE, nil
}
