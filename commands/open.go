package commands

import (
	"context"
	"fmt"
	"github.com/zetamatta/go-getch"
	"io"
	"os"
	"os/exec"
	"strings"

	"../dos"
)

func open1(fname string, out io.Writer) {
	p := strings.Split(fname, ":")
	if p[0] != "http" && p[0] != "https" {
		if _, err := os.Stat(fname); err != nil {
			fmt.Fprintln(out, err.Error())
			return
		}
	}
	if err := dos.ShellExecute("open", dos.TruePath(fname), "", ""); err != nil {
		fmt.Fprintln(out, err.Error())
	}
}

func cmd_open(ctx context.Context, cmd *exec.Cmd) (int, error) {
	switch len(cmd.Args) {
	case 1:
		open1(".", cmd.Stderr)
	case 2:
		open1(cmd.Args[1], cmd.Stderr)
	default:
		for _, arg := range cmd.Args[1:] {
			fmt.Fprintf(cmd.Stderr, "open: %s ? [y/n/q] ", arg)
			ch := getch.Rune()
			fmt.Fprintf(cmd.Stderr, "%c\n", ch)
			if ch == 'q' {
				break
			} else if ch != 'y' && ch != 'Y' {
				continue
			}
			open1(arg, cmd.Stderr)
		}
	}
	return 0, nil
}
