package functions

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/frame"
	"github.com/zetamatta/nyagos/readline"
)

var rxAnsiEscCode = regexp.MustCompile("\x1b[^a-zA-Z]*[a-zA-Z]")

func setTitle(s string) {
	fmt.Fprintf(readline.Console, "\x1B]0;%s\007", s)
}

func Prompt(args []interface{}) []interface{} {
	if len(args) >= 2 {
		setTitle(fmt.Sprint(args[1]))
	} else if wd, err := os.Getwd(); err == nil {
		if flag, _ := dos.IsElevated(); flag {
			setTitle("(Admin) - " + wd)
		} else {
			setTitle("NYAGOS - " + wd)
		}
	} else {
		if flag, _ := dos.IsElevated(); flag {
			setTitle("(Admin)")
		} else {
			setTitle("NYAGOS")
		}
	}
	var template string
	if len(args) >= 1 {
		template = fmt.Sprint(args[0])
	} else {
		template = "[too few arguments]"
	}
	text := frame.Format2Prompt(template)

	io.WriteString(readline.Console, text)

	text = rxAnsiEscCode.ReplaceAllString(text, "")
	lfPos := strings.LastIndex(text, "\n")
	if lfPos >= 0 {
		text = text[lfPos+1:]
	}
	return []interface{}{readline.GetStringWidth(text)}
}
