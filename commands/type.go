package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/zetamatta/go-mbcs"

	"github.com/zetamatta/nyagos/shell"
)

func cat(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		var text string
		if utf8.Valid(line) {
			text = string(line)
		} else {
			var err error
			text, err = mbcs.AtoU(line)
			if err != nil {
				text = err.Error()
			}
		}
		text = strings.Replace(text, "\xEF\xBB\xBF", "", 1)
		fmt.Fprintln(w, text)
	}
}

func cmd_type(ctx context.Context, cmd *shell.Cmd) (int, error) {
	if len(cmd.Args) <= 1 {
		cat(cmd.Stdin, cmd.Stdout)
	} else {
		for _, arg1 := range cmd.Args[1:] {
			r, err := os.Open(arg1)
			if err != nil {
				return 1, err
			}
			cat(r, cmd.Stdout)
			r.Close()
		}
	}
	return 0, nil
}
