package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zetamatta/nyagos/dos"
)

func printShortcut(s, t, d string, out io.Writer) {
	fmt.Fprintf(out, "    %s\n--> %s", s, t)
	if d != "" {
		fmt.Fprintf(out, "(%s)\n", d)
	} else {
		fmt.Fprintln(out)
	}
}

func makeShortcut(s, t, d string, out io.Writer) error {
	s_, err := filepath.Abs(s)
	if err != nil {
		return err
	}
	t_, err := filepath.Abs(t)
	if err != nil {
		return err
	}
	stat1, err := os.Stat(t_)
	if err == nil && stat1 != nil {
		if stat1.IsDir() {
			t_ = filepath.Join(t_, filepath.Base(s_))
		} else {
			return fmt.Errorf("%s: file already exists", t)
		}
	}
	if !strings.EqualFold(filepath.Ext(t_), ".lnk") {
		t_ = t_ + ".lnk"
	}
	err = dos.MakeShortcut(s_, t_, d)
	if err == nil {
		printShortcut(s_, t_, d, out)
	}
	return err
}

func cmdLnk(_ context.Context, cmd1 Param) (int, error) {
	switch len(cmd1.Args()) {
	case 0, 1:
		fmt.Fprintln(cmd1.Err(), "usage: lnk FILENAME SHORTCUT WORKING-DIR")
		return 0, nil
	case 2:
		target, dir, err := dos.ReadShortcut(cmd1.Arg(1))
		if err != nil {
			return 1, err
		}
		printShortcut(target, cmd1.Arg(1), dir, cmd1.Out())
		break
	case 3:
		if err := makeShortcut(cmd1.Arg(1), cmd1.Arg(2), "", cmd1.Out()); err != nil {
			return 1, err
		}
		break
	case 4:
		if err := makeShortcut(cmd1.Arg(1), cmd1.Arg(2), cmd1.Arg(3), cmd1.Out()); err != nil {
			return 1, err
		}
		break
	}
	return 0, nil
}
