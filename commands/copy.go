package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"unicode"

	"github.com/zetamatta/go-getch"

	"github.com/zetamatta/nyagos/dos"
)

type copymove_t struct {
	Param
	Action  func(src, dst string) error
	IsDirOk bool
}

func cmdCopy(ctx context.Context, cmd Param) (int, error) {
	return copymove_t{
		Param: cmd,
		Action: func(src, dst string) error {
			return dos.Copy(src, dst, false)
		},
	}.Run(ctx, cmd.Args())
}

func cmdMove(ctx context.Context, cmd Param) (int, error) {
	return copymove_t{
		Param:   cmd,
		Action:  dos.Move,
		IsDirOk: true,
	}.Run(ctx, cmd.Args())
}

func cmdLn(ctx context.Context, cmd Param) (int, error) {
	if len(cmd.Args()) >= 2 && cmd.Arg(1) == "-s" {
		args := make([]string, 0, len(cmd.Args())-1)
		args = append(args, cmd.Arg(0))
		args = append(args, cmd.Args()[2:]...)
		return copymove_t{
			Param:   cmd,
			Action:  os.Symlink,
			IsDirOk: true,
		}.Run(ctx, args)
	} else {
		return copymove_t{
			Param:  cmd,
			Action: os.Link,
		}.Run(ctx, cmd.Args())
	}
}

var rxDir = regexp.MustCompile(`[\\/:]\.{0,2}$`)

func judgeDir(path string) bool {
	if path == "." || path == ".." || rxDir.MatchString(path) {
		return true
	}
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.Mode().IsDir()
}

func (this copymove_t) Run(ctx context.Context, args []string) (int, error) {
	all := false
	args = args[1:]
	for {
		if len(args) >= 1 && (args[0] == "/y" || args[0] == "/Y") {
			all = true
			args = args[1:]
		} else if len(args) >= 1 && (args[0] == "/-y" || args[0] == "/-Y") {
			all = false
			args = args[1:]
		} else {
			break
		}
	}
	if len(args) < 2 {
		fmt.Fprintf(this.Err(),
			"Usage: %s [/y] SOURCE-FILENAME DESITINATE-FILENAME\n"+
				"       %s [/y] FILENAMES... DESINATE-DIRECTORY\n",
			this.Arg(0), this.Arg(0))
		return 0, nil
	}
	isDir := judgeDir(args[len(args)-1])
	srcs := args[0 : len(args)-1]
	for i, src := range srcs {
		if getch.IsCtrlCPressed() {
			fmt.Fprintln(this.Err(), "^C")
			return 0, nil
		}
		dst := args[len(args)-1]
		if isDir {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		if !this.IsDirOk {
			fi, err := os.Stat(src)
			if err == nil && fi.Mode().IsDir() {
				fmt.Fprintf(this.Err(), "%s is directory and passed.\n", src)
				continue
			}
		}

		fmt.Fprintf(this.Err(), "%s -> %s\n", src, dst)
		if !all {
			fi, err := os.Stat(dst)
			if fi != nil && err == nil {
				fmt.Fprintf(this.Err(),
					"%s: override? [Yes/No/All/Quit] ",
					dst)
				ch := getch.Rune()
				if unicode.IsPrint(ch) {
					fmt.Fprintf(this.Err(), "%c\n", ch)
				} else {
					io.WriteString(this.Err(), "\n")
				}
				switch ch {
				case 'y', 'Y':

				case 'a', 'A':
					all = true
				case 'q', 'Q', rune(0x03):
					return 0, nil
				default:
					continue
				}
			}
		}
		if ctx != nil {
			select {
			case <-ctx.Done():
				return 0, nil
			default:
			}
		}
		err := this.Action(src, dst)
		if err != nil {
			if i >= len(srcs)-1 {
				return 1, err
			}
			fmt.Fprintf(this.Err(), "%s\nContinue? [Yes/No] ", err.Error())
			ch := getch.Rune()
			if unicode.IsPrint(ch) {
				fmt.Fprintf(this.Err(), "%c\n", ch)
			} else {
				io.WriteString(this.Err(), "\n")
			}
			if ch != 'y' && ch != 'Y' {
				return 0, nil
			}
		}
	}
	return 0, nil
}
