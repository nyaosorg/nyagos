package commands

import "fmt"
import "os"
import "syscall"

import "../interpreter"
import "../conio"

func cmd_del(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	n := len(cmd.Args)
	all := false
	for i := 1; i < n; i++ {
		path := cmd.Args[i]
		stat, statErr := os.Stat(path)
		if statErr != nil {
			fmt.Fprintf(cmd.Stdout, "(%d/%d) %s: %s\n",
				i, n-1, path, statErr.Error())
			continue
		}
		if mode := stat.Mode(); mode.IsDir() {
			fmt.Fprintf(cmd.Stdout, "(%d/%d) %s is directory and passed.\n",
				i, n-1, path)
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
				return interpreter.CONTINUE, nil
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
			return interpreter.CONTINUE, err
		}
		fmt.Println("-> done.")
	}
	return interpreter.CONTINUE, nil
}
