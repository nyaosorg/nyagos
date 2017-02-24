package commands

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/zetamatta/go-getch"

	"../dos"
)

func open1(fname string, out io.Writer) {
	err1 := dos.ShellExecute("open", fname, "", "")
	if err1 != nil {
		fmt.Fprintln(out, err1.Error())
		err2 := dos.ShellExecute("open", dos.TruePath(fname), "", "")
		if err2 != nil {
			fmt.Fprintln(out, err2.Error())
		}
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
