package commands

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/zetamatta/go-findfile"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/shell"
)

func getbit(c byte) uint32 {
	switch c {
	case 'r', 'R':
		return dos.FILE_ATTRIBUTE_READONLY
	case 'h', 'H':
		return dos.FILE_ATTRIBUTE_HIDDEN
	case 's', 'S':
		return dos.FILE_ATTRIBUTE_SYSTEM
	case 'a', 'A':
		return dos.FILE_ATTRIBUTE_ARCHIVE
	default:
		return 0
	}
}

func bit2flg(bits uint32, flag uint32, r rune) rune {
	if (bits & flag) != 0 {
		return r
	} else {
		return ' '
	}
}

func globfile(pattern string) (result []string) {
	findfile.Walk(pattern, func(f *findfile.FileInfo) bool {
		if !f.IsDir() {
			one := filepath.Join(filepath.Dir(pattern), f.Name())
			result = append(result, one)
		}
		return true
	})
	return
}

func cmd_attrib(ctx context.Context, cmd *shell.Cmd) (int, error) {
	var set_bits uint32 = 0
	var reset_bits uint32 = 0
	files := make([]string, 0, len(cmd.Args)-1)

	for _, arg1 := range cmd.Args[1:] {
		if len(arg1) == 2 && arg1[0] == '+' {
			bits := getbit(arg1[1])
			if bits != 0 {
				set_bits |= bits
				continue
			}
		}
		if len(arg1) == 2 && arg1[0] == '-' {
			bits := getbit(arg1[1])
			if bits != 0 {
				reset_bits |= bits
				continue
			}
		}
		if arg1s := globfile(arg1); arg1s != nil && len(arg1s) > 0 {
			files = append(files, arg1s...)
		} else {
			files = append(files, arg1)
		}
	}
	if len(files) <= 0 {
		files = globfile(`.\*`)
		if files == nil {
			files = []string{}
		}
	}
	sort.Strings(files)
	for _, arg1 := range files {
		bits, err := dos.GetFileAttributes(arg1)
		if err != nil {
			return 1, err
		}
		if set_bits == 0 && reset_bits == 0 {
			fullpath, err := filepath.Abs(arg1)
			if err != nil {
				fullpath = arg1
			}
			fmt.Fprintf(cmd.Stdout, "%c  %c%c%c       %s\n",
				bit2flg(bits, dos.FILE_ATTRIBUTE_ARCHIVE, 'A'),
				bit2flg(bits, dos.FILE_ATTRIBUTE_SYSTEM, 'S'),
				bit2flg(bits, dos.FILE_ATTRIBUTE_HIDDEN, 'H'),
				bit2flg(bits, dos.FILE_ATTRIBUTE_READONLY, 'R'),
				fullpath)
		} else {
			bits = (bits | set_bits) &^ reset_bits
			err = dos.SetFileAttributes(arg1, bits)
			if err != nil {
				return 2, err
			}
		}
	}
	return 0, nil
}
