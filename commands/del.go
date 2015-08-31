package commands

import (
	"fmt"
	"os"
	"syscall"

	"../conio"
	. "../interpreter"
)

func cmd_del(cmd *Interpreter) (ErrorLevel, error) {
	n := len(cmd.Args)
	if n <= 1 {
		fmt.Fprintln(cmd.Stderr, "Usage: del   FILE(S)...")
		fmt.Fprintln(cmd.Stderr, "       erase FILE(S)...")
		return NOERROR, nil
	}
	all := false
	errorcount := ErrorLevel(0)
	for i := 1; i < n; i++ {
		path := cmd.Args[i]
		stat, statErr := os.Lstat(path)
		if statErr != nil {
			fmt.Fprintf(cmd.Stdout, "(%d/%d) %s: %s\n",
				i, n-1, path, statErr.Error())
			errorcount++
			continue
		}
		if mode := stat.Mode(); mode.IsDir() {
			fmt.Fprintf(cmd.Stdout, "(%d/%d) %s is directory and passed.\n",
				i, n-1, path)
			errorcount++
			continue
		}
		if all {
			fmt.Fprintf(cmd.Stdout, "(%d/%d) %s: Remove ", i, n-1, path)
		} else {
			fmt.Fprintf(cmd.Stdout,
				"(%d/%d) %s: Remove ? [Yes/No/All/Quit] ",
				i, n-1, path)
			ch := conio.GetCh()
			fmt.Fprintf(cmd.Stdout, "%c ", ch)
			switch ch {
			case 'q', 'Q':
				fmt.Fprintln(cmd.Stdout)
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
		err := syscall.Unlink(path)
		if err != nil {
			fmt.Printf("-> %s\n", err.Error())
			errorcount++
			continue
		}
		fmt.Println("-> done.")
	}
	return errorcount, nil
}
