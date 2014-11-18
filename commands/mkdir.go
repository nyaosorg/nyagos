package commands

import "fmt"
import "os"
import "syscall"

import "../conio"
import "../dos"
import "../interpreter"

func cmd_mkdir(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	for _, arg1 := range cmd.Args[1:] {
		err := syscall.Mkdir(arg1, 0777)
		if err != nil {
			fmt.Fprintf(cmd.Stderr, "%s: %s\n", arg1, err)
		}
	}
	return interpreter.CONTINUE, nil
}

func cmd_rmdir(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	s_option := false
	message := "%s: Rmdir Are you sure? [Yes/No/Quit] "
	for _, arg1 := range cmd.Args[1:] {
		if arg1 == "/s" {
			s_option = true
			message = "%s : Delete Tree. Are you sure? [Yes/No/Quit] "
			continue
		}
		fi, err := os.Stat(arg1)
		if err != nil {
			fmt.Fprintf(cmd.Stderr, "%s: %s\n", arg1, err.Error())
			continue
		}
		if !fi.IsDir() {
			fmt.Fprintf(cmd.Stderr, "%s: not directory\n", arg1)
			continue
		}
		fmt.Fprintf(cmd.Stderr, message, arg1)
		ch := conio.GetCh()
		fmt.Fprintf(cmd.Stderr, "%c ", ch)
		switch ch {
		case 'y', 'Y':

		case 'q', 'Q':
			fmt.Fprintln(cmd.Stderr, "-> canceled all")
			return interpreter.CONTINUE, nil
		default:
			fmt.Fprintln(cmd.Stderr, "-> canceled")
			continue
		}
		if s_option {
			fmt.Fprintln(cmd.Stdout)
			err = dos.Truncate(arg1, cmd.Stdout)
		} else {
			err = syscall.Rmdir(arg1)
		}
		if err != nil {
			fmt.Fprintf(cmd.Stderr, "-> %s\n", err.Error())
		} else {
			fmt.Fprintln(cmd.Stderr, "-> done.")
		}
	}
	return interpreter.CONTINUE, nil
}
