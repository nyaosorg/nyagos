package functions

import (
	"fmt"
	"github.com/nyaosorg/nyagos/internal/frame"
	"io"
	"os"
	"regexp"
)

var rxAnsiEscCode = regexp.MustCompile("\x1b[^a-zA-Z]*[a-zA-Z]")

func setTitle(w io.Writer, s string) {
	fmt.Fprintf(w, "\x1B]0;%s\007", s)
}

// Prompt is the body of the lua-function `nyagos.default_prompt`
func Prompt(param *Param) []interface{} {
	return []interface{}{PromptCore(param.Term, param.Args...)}
}

// PromptCore prints prompt-str(args[0]) to console.
func PromptCore(console io.Writer, args ...interface{}) string {
	if len(args) >= 2 {
		setTitle(console, fmt.Sprint(args[1]))
	} else if wd, err := os.Getwd(); err == nil {
		if flag := isElevated(); flag {
			setTitle(console, "(Admin) - "+wd)
		} else {
			setTitle(console, "NYAGOS - "+wd)
		}
	} else {
		if flag := isElevated(); flag {
			setTitle(console, "(Admin)")
		} else {
			setTitle(console, "NYAGOS")
		}
	}
	var template string
	if len(args) >= 1 {
		template = fmt.Sprint(args[0])
	} else {
		template = "[too few arguments]"
	}
	return frame.Format2Prompt(template)
}
