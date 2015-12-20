package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"../conio"
	"../dos"
	. "../interpreter"
)

type actionT struct {
	Do      func(src, dst string) error
	IsDirOk bool
}

func cmd_copy(cmd *Interpreter) (ErrorLevel, error) {
	return cmd_xxxx(cmd, actionT{
		func(src, dst string) error {
			return dos.Copy(src, dst, false)
		}, false})
}

func cmd_move(cmd *Interpreter) (ErrorLevel, error) {
	return cmd_xxxx(cmd, actionT{
		func(src, dst string) error {
			return dos.Move(src, dst)
		}, true})
}

func cmd_ln(cmd *Interpreter) (ErrorLevel, error) {
	return cmd_xxxx(cmd, actionT{
		func(src, dst string) error {
			return os.Link(src, dst)
		}, false})
}

func cmd_xxxx(cmd *Interpreter, action actionT) (ErrorLevel, error) {
	if len(cmd.Args) <= 2 {
		fmt.Fprintf(cmd.Stderr,
			"Usage: %s [/y] SOURCE-FILENAME DESITINATE-FILENAME\n"+
				"       %s [/y] FILENAMES... DESINATE-DIRECTORY\n",
			cmd.Args[0], cmd.Args[0])
		return NOERROR, nil
	}
	fi, err := os.Stat(cmd.Args[len(cmd.Args)-1])
	isDir := err == nil && fi.Mode().IsDir()
	all := false
	for i, n := 1, len(cmd.Args)-1; i < n; i++ {
		if cmd.Args[i] == "/y" {
			all = true
			continue
		}
		src := cmd.Args[i]
		dst := cmd.Args[n]
		if isDir {
			dst = dos.Join(dst, filepath.Base(src))
		}
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
					return NOERROR, nil
				default:
					continue
				}
			}
		}
		err := action.Do(src, dst)
		if err != nil {
			return ErrorLevel(1), err
		}
	}
	return NOERROR, nil
}
