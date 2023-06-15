package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nyaosorg/go-windows-shortcut"
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
	sAbs, err := filepath.Abs(s)
	if err != nil {
		return err
	}
	tAbs, err := filepath.Abs(t)
	if err != nil {
		return err
	}
	stat1, err := os.Stat(tAbs)
	if err == nil && stat1 != nil {
		if stat1.IsDir() {
			tAbs = filepath.Join(tAbs, filepath.Base(sAbs))
		} else {
			return fmt.Errorf("%s: file already exists", t)
		}
	}
	if !strings.EqualFold(filepath.Ext(tAbs), ".lnk") {
		tAbs = tAbs + ".lnk"
	}
	err = shortcut.Make(sAbs, tAbs, d)
	if err == nil {
		printShortcut(sAbs, tAbs, d, out)
	}
	return err
}

func cmdLnk(_ context.Context, cmd1 Param) (int, error) {
	switch len(cmd1.Args()) {
	case 0, 1:
		fmt.Fprintln(cmd1.Err(), "usage: lnk FILENAME SHORTCUT WORKING-DIR")
		return 0, nil
	case 2:
		fn := cmd1.Arg(1)
		if strings.ToLower(filepath.Ext(fn)) != ".lnk" {
			return 1, fmt.Errorf("%s: not shotcut file", fn)
		}
		if _, err := os.Stat(fn); err != nil {
			return 1, fmt.Errorf("%s: not exist", fn)
		}
		target, dir, err := shortcut.Read(fn)
		if err != nil {
			return 1, err
		}
		printShortcut(target, fn, dir, cmd1.Out())
	case 3:
		if err := makeShortcut(cmd1.Arg(1), cmd1.Arg(2), "", cmd1.Out()); err != nil {
			return 1, err
		}
	case 4:
		if err := makeShortcut(cmd1.Arg(1), cmd1.Arg(2), cmd1.Arg(3), cmd1.Out()); err != nil {
			return 1, err
		}
	}
	return 0, nil
}
