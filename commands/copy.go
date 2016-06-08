package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"../conio/getch"
	"../dos"
	. "../interpreter"
)

func cmd_copy(cmd *Interpreter) (ErrorLevel, error) {
	return cmd_xxxx(cmd.Args,
		cmd.Stderr,
		func(src, dst string) error {
			return dos.Copy(src, dst, false)
		},
		false)
}

func cmd_move(cmd *Interpreter) (ErrorLevel, error) {
	return cmd_xxxx(cmd.Args,
		cmd.Stderr,
		func(src, dst string) error {
			return dos.Move(src, dst)
		},
		true)
}

func cmd_ln(cmd *Interpreter) (ErrorLevel, error) {
	if len(cmd.Args) >= 2 && cmd.Args[1] == "-s" {
		args := make([]string, 0, len(cmd.Args)-1)
		args = append(args, cmd.Args[0])
		args = append(args, cmd.Args[2:]...)
		return cmd_xxxx(
			args,
			cmd.Stderr,
			func(src, dst string) error {
				return os.Symlink(src, dst)
			},
			true)
	} else {
		return cmd_xxxx(
			cmd.Args,
			cmd.Stderr,
			func(src, dst string) error {
				return os.Link(src, dst)
			},
			false)
	}
}

func cmd_xxxx(args []string,
	out io.Writer,
	action func(src, dst string) error,
	isDirOk bool) (ErrorLevel, error) {
	if len(args) <= 2 {
		fmt.Fprintf(out,
			"Usage: %s [/y] SOURCE-FILENAME DESITINATE-FILENAME\n"+
				"       %s [/y] FILENAMES... DESINATE-DIRECTORY\n",
			args[0], args[0])
		return NOERROR, nil
	}
	fi, err := os.Stat(args[len(args)-1])
	isDir := err == nil && fi.Mode().IsDir()
	all := false
	for i, n := 1, len(args)-1; i < n; i++ {
		if args[i] == "/y" {
			all = true
			continue
		}
		src := args[i]
		dst := args[n]
		if isDir {
			dst = dos.Join(dst, filepath.Base(src))
		}
		if !isDirOk {
			fi, err := os.Stat(src)
			if err == nil && fi.Mode().IsDir() {
				fmt.Fprintf(out, "%s is directory and passed.\n", src)
				continue
			}
		}

		fmt.Fprintf(out, "%s -> %s\n", src, dst)
		if !all {
			fi, err := os.Stat(dst)
			if fi != nil && err == nil {
				fmt.Fprintf(out,
					"%s: override? [Yes/No/All/Quit] ",
					dst)
				ch := getch.Rune()
				fmt.Fprintf(out, "%c\n", ch)
				switch ch {
				case 'y', 'Y':

				case 'a', 'A':
					all = true
				case 'q', 'Q':
					return NOERROR, nil
				default:
					continue
				}
			}
		}
		err := action(src, dst)
		if err != nil {
			if i == n-1 {
				return ErrorLevel(1), err
			}
			fmt.Fprintf(out, "%s\nContinue? [Yes/No] ", err.Error())
			ch := getch.Rune()
			fmt.Fprintf(out, "%c\n", ch)
			switch ch {
			case 'y', 'Y':

			default:
				return NOERROR, nil
			}
		}
	}
	return NOERROR, nil
}
