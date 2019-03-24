package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/nodos"
)

type copyMoveT struct {
	Param
	Action  func(src, dst string) error
	IsDirOk bool
}

func cmdCopy(ctx context.Context, cmd Param) (int, error) {
	return copyMoveT{
		Param: cmd,
		Action: func(src, dst string) error {
			return nodos.Copy(src, dst, false)
		},
	}.Run(ctx, cmd.Args())
}

func cmdMove(ctx context.Context, cmd Param) (int, error) {
	return copyMoveT{
		Param:   cmd,
		Action:  nodos.Move,
		IsDirOk: true,
	}.Run(ctx, cmd.Args())
}

func cmdLn(ctx context.Context, cmd Param) (int, error) {
	if len(cmd.Args()) >= 2 && cmd.Arg(1) == "-s" {
		args := make([]string, 0, len(cmd.Args())-1)
		args = append(args, cmd.Arg(0))
		args = append(args, cmd.Args()[2:]...)
		return copyMoveT{
			Param:   cmd,
			Action:  os.Symlink,
			IsDirOk: true,
		}.Run(ctx, args)
	}
	return copyMoveT{
		Param: cmd,
		Action: func(src, dst string) error {
			if stat, err := os.Stat(src); err != nil {
				return err
			} else if stat.IsDir() {
				if fullpath, err := filepath.Abs(src); err != nil {
					return err
				} else {
					return dos.CreateJunction(dst, fullpath)
				}
			} else {
				return os.Link(src, dst)
			}
		},
		IsDirOk: true,
	}.Run(ctx, cmd.Args())
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

func (cm copyMoveT) Run(ctx context.Context, args []string) (int, error) {
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
		fmt.Fprintf(cm.Err(),
			"Usage: %s [/y] SOURCE-FILENAME DESITINATE-FILENAME\n"+
				"       %s [/y] FILENAMES... DESINATE-DIRECTORY\n",
			cm.Arg(0), cm.Arg(0))
		return 0, nil
	}

	_dst := args[len(args)-1]
	if strings.ToLower(filepath.Ext(_dst)) == ".lnk" {
		if __dst, _, err := nodos.ReadShortcut(_dst); err == nil {
			_dst = __dst
		}
	}

	isDir := judgeDir(_dst)
	srcs := args[0 : len(args)-1]
	for i, src := range srcs {
		dst := _dst
		if isDir {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		if !cm.IsDirOk {
			fi, err := os.Stat(src)
			if err == nil && fi.Mode().IsDir() {
				fmt.Fprintf(cm.Err(), "%s is directory and passed.\n", src)
				continue
			}
		}

		fmt.Fprintf(cm.Err(), "%s -> %s\n", src, dst)
		if !all {
			fi, err := os.Stat(dst)
			if fi != nil && err == nil {
				fmt.Fprintf(cm.Err(),
					"%s: override? [Yes/No/All/Quit] ",
					dst)
				ch, err := getkey()
				if err != nil {
					return 1, err
				}
				if unicode.IsPrint(ch) {
					fmt.Fprintf(cm.Err(), "%c\n", ch)
				} else {
					io.WriteString(cm.Err(), "\n")
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
				fmt.Fprintln(cm.Err(), "^C")
				return 0, nil
			default:
			}
		}
		err := cm.Action(src, dst)
		if err != nil {
			if i >= len(srcs)-1 {
				return 1, err
			}
			fmt.Fprintf(cm.Err(), "%s\nContinue? [Yes/No] ", err.Error())
			ch, err := getkey()
			if err != nil {
				return 1, err
			}
			if unicode.IsPrint(ch) {
				fmt.Fprintf(cm.Err(), "%c\n", ch)
			} else {
				io.WriteString(cm.Err(), "\n")
			}
			if ch != 'y' && ch != 'Y' {
				return 0, nil
			}
		}
	}
	return 0, nil
}
