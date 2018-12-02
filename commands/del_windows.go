package commands

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"unicode"

	"github.com/zetamatta/nyagos/dos"
)

func setReadonly(path string) error {
	perm, err := dos.GetFileAttributes(path)
	if err != nil {
		return err
	}
	return dos.SetFileAttributes(path, perm&^dos.FILE_ATTRIBUTE_READONLY)
}

func cmdDel(ctx context.Context, cmd Param) (int, error) {
	n := len(cmd.Args())
	if n <= 1 {
		fmt.Fprintln(cmd.Err(), "Usage: del   [/q] FILE(S)...")
		fmt.Fprintln(cmd.Err(), "       erase [/q] FILE(S)...")
		return 0, nil
	}
	all := false
	force := false
	errorcount := 0
	i := 1
	for _, arg1 := range cmd.Args()[1:] {
		if ctx != nil {
			select {
			case <-ctx.Done():
				return 0, nil
			default:
			}
		}
		if arg1 == "/q" {
			all = true
			n--
			continue
		}
		if arg1 == "/f" {
			force = true
			n--
			continue
		}
		path := arg1
		stat, err := os.Lstat(path)
		if _, ok := err.(*os.PathError); ok || os.IsNotExist(err) {
			fmt.Fprintf(cmd.Out(), "(%d/%d) %s: not found.\n", i, n-1, path)
			errorcount++
			continue
		}
		if err != nil {
			fmt.Fprintf(cmd.Out(), "(%d/%d) %s: %s\n", i, n-1, path, err)
			errorcount++
			continue
		}
		if mode := stat.Mode(); mode.IsDir() {
			fmt.Fprintf(cmd.Out(), "(%d/%d) %s is directory and passed.\n",
				i, n-1, path)
			errorcount++
			continue
		}
		if all {
			fmt.Fprintf(cmd.Out(), "(%d/%d) %s: Remove ", i, n-1, path)
		} else {
			fmt.Fprintf(cmd.Out(),
				"(%d/%d) %s: Remove ? [Yes/No/All/Quit] ",
				i, n-1, path)
			ch, err := getkey()
			if err != nil {
				return 1, err
			}
			if unicode.IsPrint(ch) {
				fmt.Fprintf(cmd.Out(), "%c ", ch)
			}
			switch ch {
			case 'q', 'Q', rune(0x03):
				fmt.Fprintln(cmd.Out())
				return errorcount, nil
			case 'y', 'Y':
				break
			case 'a', 'A':
				all = true
			default: // for 'n','N'
				fmt.Println("-> canceled")
				continue
			}
		}
		err = syscall.Unlink(path)
		if err != nil && force {
			if err1 := setReadonly(path); err1 == nil {
				err = syscall.Unlink(path)
			}
		}
		if err != nil {
			fmt.Printf("-> %s\n", err)
			errorcount++
			continue
		}
		fmt.Println("-> done.")
		i++
	}
	return errorcount, nil
}
