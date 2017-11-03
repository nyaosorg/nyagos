package mains

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zetamatta/go-box"
	"github.com/zetamatta/go-findfile"
	"github.com/zetamatta/go-getch"
	"github.com/zetamatta/go-mbcs"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/ifdbg"
	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/shell"
)

type any_t = interface{}

func toStr(arr []any_t, n int) string {
	if n < len(arr) {
		if ifdbg.DBG {
			println(fmt.Sprint(arr[n]))
		}
		return fmt.Sprint(arr[n])
	} else {
		if ifdbg.DBG {
			println("''")
		}
		return ""
	}
}

func cmdElevated([]any_t) []any_t {
	flag, _ := dos.IsElevated()
	return []any_t{flag}
}

func cmdChdir(args []any_t) []any_t {
	if len(args) >= 1 {
		dos.Chdir(fmt.Sprint(args[0]))
		return []any_t{true}
	}
	return []any_t{nil, errors.New("directory is required")}
}

func cmdBox(args []any_t) []any_t {
	t, ok := args[0].(map[any_t]any_t)
	if !ok {
		return []any_t{nil, "Not a table"}
	}
	if len(t) == 0 {
		return []any_t{}
	}
	sources := make([]string, 0, len(t))
	for i, i_ := 1, len(t); i <= i_; i++ {
		if val, ok := t[i]; ok {
			sources = append(sources, fmt.Sprint(val))
		}
	}
	return []any_t{box.Choice(sources, readline.Console)}
}

func cmdResetCharWidth(args []any_t) []any_t {
	readline.ResetCharWidth()
	return []any_t{}
}

func cmdNetDriveToUNC(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{}
	}
	path, ok := args[0].(string)
	if !ok {
		return []any_t{path}
	}
	unc := dos.NetDriveToUNC(path)
	return []any_t{unc}
}

func cmdShellExecute(args []any_t) []any_t {
	err := dos.ShellExecute(
		toStr(args, 0),
		dos.TruePath(toStr(args, 1)),
		toStr(args, 2),
		toStr(args, 3))
	if err != nil {
		return []any_t{nil, err}
	} else {
		return []any_t{true}
	}
}

func cmdGetwd(args []any_t) []any_t {
	wd, err := os.Getwd()
	if err == nil {
		return []any_t{wd}
	} else {
		return []any_t{nil, err}
	}
}

func cmdGetKey(args []any_t) []any_t {
	keycode, scancode, shiftstatus := getch.Full()
	return []any_t{keycode, scancode, shiftstatus}
}

func cmdGetViewWidth(args []any_t) []any_t {
	width, height := box.GetScreenBufferInfo().ViewSize()
	return []any_t{width, height}
}

func cmdPathJoin(args []any_t) []any_t {
	if len(args) < 0 {
		return []any_t{""}
	}
	path := fmt.Sprint(args[0])
	for i, i_ := 1, len(args); i < i_; i++ {
		sub := fmt.Sprint(args[i])
		path = filepath.Join(path, sub)
	}
	return []any_t{path}
}

func cmdAccess(args []any_t) []any_t {
	if len(args) < 2 {
		return []any_t{nil, "nyagos.access requilres two arguments"}
	}
	path := fmt.Sprint(args[0])
	mode, mode_ok := args[1].(int)
	if !mode_ok {
		return []any_t{nil, "mode value must be interger"}
	}
	if ifdbg.DBG {
		fmt.Fprintf(os.Stderr, "given mode==%o\n", mode)
	}
	fi, err := os.Stat(path)

	var result bool
	if err != nil || fi == nil {
		result = false
	} else {
		if ifdbg.DBG {
			fmt.Fprintf(os.Stderr, "file mode==%o\n", fi.Mode().Perm())
		}
		switch {
		case mode == 0:
			result = true
		case (mode & 1) != 0: // X_OK
		case (mode & 2) != 0: // W_OK
			result = ((fi.Mode().Perm() & 0200) != 0)
		case (mode & 4) != 0: // R_OK
			result = ((fi.Mode().Perm() & 0400) != 0)
		}
	}
	return []any_t{result}
}

func cmdStat(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, errors.New("fee arguments")}
	}
	path := toStr(args, 0)
	var stat os.FileInfo
	var path_ string
	if len(path) > 0 && path[len(path)-1] == '\\' {
		path_ = filepath.Join(path, ".")
	} else {
		path_ = path
	}
	statErr := findfile.Walk(path_, func(f *findfile.FileInfo) bool {
		stat = f
		return false
	})
	if statErr != nil {
		return []any_t{nil, statErr}
	}
	if stat == nil {
		return []any_t{nil, fmt.Errorf("%s: failed to stat", path)}
	}
	t := stat.ModTime()
	return []any_t{
		map[string]any_t{
			"name":  stat.Name(),
			"size":  stat.Size(),
			"isdir": stat.IsDir(),
			"mtime": map[string]any_t{
				"year":   t.Year(),
				"month":  t.Month(),
				"day":    t.Day(),
				"hour":   t.Hour(),
				"minute": t.Minute(),
				"second": t.Second(),
			},
		},
	}
}

func cmdSetEnv(args []any_t) []any_t {
	if len(args) < 2 {
		return []any_t{nil, "too few arguments"}
	}
	name := fmt.Sprint(args[len(args)-2])
	value := fmt.Sprint(args[len(args)-1])
	if len(value) > 0 {
		os.Setenv(name, value)
	} else {
		os.Unsetenv(name)
	}
	return []any_t{true}
}

func cmdGetEnv(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, "too few arguments"}
	}
	name := fmt.Sprint(args[len(args)-1])
	value, ok := shell.OurGetEnv(name)
	if ok && len(value) > 0 {
		return []any_t{value}
	} else {
		return []any_t{nil}
	}
}

func cmdAtoU(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, "too few arguments"}
	}
	if s, ok := args[0].(string); ok {
		if val, err := mbcs.AtoU([]byte(s)); err == nil {
			return []any_t{val}
		} else {
			return []any_t{nil, err}
		}
	} else {
		return []any_t{fmt.Sprint(args[0])}
	}
}
