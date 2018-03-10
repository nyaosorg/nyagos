package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/zetamatta/nyagos/dos"
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

func cmdOpen(ctx context.Context, cmd Param) (int, error) {
	switch len(cmd.Args()) {
	case 1:
		wd, err := os.Getwd()
		if err != nil {
			open1(".", cmd.Err())
		} else {
			open1(wd, cmd.Err())
		}
	case 2:
		open1(cmd.Arg(1), cmd.Err())
	default:
		fmt.Fprintln(cmd.Err(), "open: ambiguous shellexecute")
	}
	return 0, nil
}
