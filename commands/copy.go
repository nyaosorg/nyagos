package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"../conio"
	"../dos"
	"../interpreter"
)

type actionT struct {
	Do      func(src, dst string) error
	IsDirOk bool
}

func cmd_copy(cmd *interpreter.Interpreter) (interpreter.ErrorLevel, error) {
	return interpreter.CONTINUE,
		cmd_xxxx(cmd, actionT{
			func(src, dst string) error {
				return dos.Copy(src, dst, false)
			}, false})
}

func cmd_move(cmd *interpreter.Interpreter) (interpreter.ErrorLevel, error) {
	return interpreter.CONTINUE,
		cmd_xxxx(cmd, actionT{
			func(src, dst string) error {
				return dos.Move(src, dst)
			}, true})
}

func cmd_xxxx(cmd *interpreter.Interpreter, action actionT) error {
	switch len(cmd.Args) {
	case 0, 1, 2:
		fmt.Fprintf(cmd.Stderr,
			"Usage: %s SOURCE-FILENAME DESITINATE-FILENAME\n"+
				"       %s FILENAMES... DESINATE-DIRECTORY\n",
			cmd.Args[0], cmd.Args[0])
	case 3:
		src := cmd.Args[1]
		dst := cmd.Args[2]
		if !action.IsDirOk {
			fi, err := os.Stat(src)
			if err == nil && fi.Mode().IsDir() {
				fmt.Fprintf(cmd.Stderr, "%s is directory and passed.\n", src)
				return nil
			}
		}
		fi, err := os.Stat(dst)
		if err == nil && fi != nil && fi.Mode().IsDir() {
			dst = dos.Join(dst, filepath.Base(src))
			fi, err = os.Stat(dst)
		}
		fmt.Fprintf(cmd.Stderr, "%s -> %s\n", src, dst)
		if fi != nil && err == nil {
			fmt.Fprintf(cmd.Stderr, "%s: override? [Yes/No] ", dst)
			ch := conio.GetCh()
			fmt.Fprintf(cmd.Stderr, "%c\n", ch)
			if ch != 'y' && ch != 'Y' {
				return nil
			}
		}
		return action.Do(src, dst)
	default:
		all := false
		for i, n := 1, len(cmd.Args)-1; i < n; i++ {
			src := cmd.Args[i]
			dst := dos.Join(cmd.Args[n], filepath.Base(src))
			if !action.IsDirOk {
				fi, err := os.Stat(src)
				if err == nil && fi.Mode().IsDir() {
					fmt.Fprintf(cmd.Stderr, "%s is directory and passed.\n", src)
					continue
				}
			}

			fmt.Fprintf(cmd.Stderr, "%s -> %s\n", src, dst)
			if !all {
				fi, err := os.Stat(dst)
				if fi != nil && err == nil {
					fmt.Fprintf(cmd.Stderr,
						"%s: override? [Yes/No/All/Quit] ",
						dst)
					ch := conio.GetCh()
					fmt.Fprintf(cmd.Stderr, "%c\n", ch)
					switch ch {
					case 'y', 'Y':

					case 'a', 'A':
						all = true
					case 'q', 'Q':
						return nil
					default:
						continue
					}
				}
			}
			err := action.Do(src, dst)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
