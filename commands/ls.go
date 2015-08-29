package commands

import (
	"io"
	"os"

	"github.com/shiena/ansicolor"

	"../interpreter"
	"./ls"
)

func cmd_ls(cmd *interpreter.Interpreter) (interpreter.ErrorLevel, error) {
	var out io.Writer
	if cmd.Stdout == os.Stdout {
		out = ansicolor.NewAnsiColorWriter(cmd.Stdout)
	} else {
		out = cmd.Stdout
	}
	return interpreter.NOERROR, ls.Main(cmd.Args[1:], out)
}
