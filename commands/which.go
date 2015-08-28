package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"../alias"
	"../dos"
	. "../interpreter"
)

func cmd_which(cmd *Interpreter) (ErrorLevel, error) {
	for _, name := range cmd.Args[1:] {
		if a, ok := alias.Table[strings.ToLower(name)]; ok {
			fmt.Fprintf(cmd.Stdout, "%s: aliased to %s\n", name, a.String())
			continue
		}
		if _, ok := BuildInCommand[name]; ok {
			fmt.Fprintf(cmd.Stdout, "%s: built-in command\n", name)
			continue
		}
		path, err := exec.LookPath(name)
		if err != nil {
			return CONTINUE, err
		}
		fmt.Fprintln(cmd.Stdout, dos.YenYen2Yen(path))
	}
	return CONTINUE, nil
}
