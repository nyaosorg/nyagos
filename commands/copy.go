package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/zetamatta/go-getch"

	"../dos"
)

type copymove_t struct {
	*exec.Cmd
	Action  func(src, dst string) error
	IsDirOk bool
}

func cmd_copy(cmd *exec.Cmd) (int, error) {
	return copymove_t{
		Cmd: cmd,
		Action: func(src, dst string) error {
			return dos.Copy(src, dst, false)
		},
	}.Run()
}

func cmd_move(cmd *exec.Cmd) (int, error) {
	return copymove_t{
		Cmd:     cmd,
		Action:  dos.Move,
		IsDirOk: true,
	}.Run()
}

func cmd_ln(cmd *exec.Cmd) (int, error) {
	if len(cmd.Args) >= 2 && cmd.Args[1] == "-s" {
		args := make([]string, 0, len(cmd.Args)-1)
		args = append(args, cmd.Args[0])
		args = append(args, cmd.Args[2:]...)
		return copymove_t{
			Cmd:     cmd,
			Action:  os.Symlink,
			IsDirOk: true,
		}.Run()
	} else {
		return copymove_t{
			Cmd:    cmd,
			Action: os.Link,
		}.Run()
	}
}

func (this copymove_t) Run() (int, error) {
	if len(this.Args) <= 2 {
		fmt.Fprintf(this.Stderr,
			"Usage: %s [/y] SOURCE-FILENAME DESITINATE-FILENAME\n"+
				"       %s [/y] FILENAMES... DESINATE-DIRECTORY\n",
			this.Args[0], this.Args[0])
		return 0, nil
	}
	fi, err := os.Stat(this.Args[len(this.Args)-1])
	isDir := err == nil && fi.Mode().IsDir()
	all := false
	for i, src := range this.Args[1 : len(this.Args)-1] {
		if src == "/y" {
			all = true
			continue
		}
		dst := this.Args[len(this.Args)-1]
		if isDir {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		if !this.IsDirOk {
			fi, err := os.Stat(src)
			if err == nil && fi.Mode().IsDir() {
				fmt.Fprintf(this.Stderr, "%s is directory and passed.\n", src)
				continue
			}
		}

		fmt.Fprintf(this.Stderr, "%s -> %s\n", src, dst)
		if !all {
			fi, err := os.Stat(dst)
			if fi != nil && err == nil {
				fmt.Fprintf(this.Stderr,
					"%s: override? [Yes/No/All/Quit] ",
					dst)
				ch := getch.Rune()
				fmt.Fprintf(this.Stderr, "%c\n", ch)
				switch ch {
				case 'y', 'Y':

				case 'a', 'A':
					all = true
				case 'q', 'Q':
					return 0, nil
				default:
					continue
				}
			}
		}
		err := this.Action(src, dst)
		if err != nil {
			if i == len(this.Args)-1 {
				return 1, err
			}
			fmt.Fprintf(this.Stderr, "%s\nContinue? [Yes/No] ", err.Error())
			ch := getch.Rune()
			fmt.Fprintf(this.Stderr, "%c\n", ch)
			if ch != 'y' && ch != 'Y' {
				return 0, nil
			}
		}
	}
	return 0, nil
}
