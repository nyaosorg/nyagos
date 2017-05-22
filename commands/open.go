package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/zetamatta/go-getch"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/shell"
)

func open1(fname string, out io.Writer) {
	err1 := dos.ShellExecute("open", fname, "", "")
	if err1 != nil {
		fmt.Fprintf(out, "%s: %s\n", fname, err1.Error())
		truepath := dos.TruePath(fname)
		err2 := dos.ShellExecute("open", truepath, "", "")
		if err2 != nil {
			fmt.Fprintf(out, "%s: %s\n", truepath, err2.Error())
		}
	}
}

func cmd_open(ctx context.Context, cmd *shell.Cmd) (int, error) {
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
