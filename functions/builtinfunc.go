package functions

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/mattn/go-colorable"
	"github.com/mattn/msgbox"

	"github.com/zetamatta/go-box"
	"github.com/zetamatta/go-findfile"
	"github.com/zetamatta/go-getch"
	"github.com/zetamatta/go-mbcs"

	"github.com/zetamatta/nyagos/commands"
	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/defined"
	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/frame"
	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/shell"
)

type any_t = interface{}

const TooFewArguments = "Too few arguments"

func toStr(arr []any_t, n int) string {
	if n < len(arr) {
		if defined.DBG {
			println(fmt.Sprint(arr[n]))
		}
		return fmt.Sprint(arr[n])
	} else {
		if defined.DBG {
			println("''")
		}
		return ""
	}
}

func CmdElevated([]any_t) []any_t {
	flag, _ := dos.IsElevated()
	return []any_t{flag}
}

func CmdChdir(args []any_t) []any_t {
	if len(args) >= 1 {
		dos.Chdir(fmt.Sprint(args[0]))
		return []any_t{true}
	}
	return []any_t{nil, "directory is required"}
}

func CmdBox(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, TooFewArguments}
	}
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
	return []any_t{box.Choice(sources, colorable.NewColorableStdout())}
}

func CmdResetCharWidth(args []any_t) []any_t {
	readline.ResetCharWidth()
	return []any_t{}
}

func CmdNetDriveToUNC(args []any_t) []any_t {
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

func CmdShellExecute(args []any_t) []any_t {
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

func CmdGetwd(args []any_t) []any_t {
	wd, err := os.Getwd()
	if err == nil {
		return []any_t{wd}
	} else {
		return []any_t{nil, err}
	}
}

func CmdGetKey(args []any_t) []any_t {
	keycode, scancode, shiftstatus := getch.Full()
	return []any_t{keycode, scancode, shiftstatus}
}

func CmdGetViewWidth(args []any_t) []any_t {
	width, height := box.GetScreenBufferInfo().ViewSize()
	return []any_t{width, height}
}

func CmdPathJoin(args []any_t) []any_t {
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

func CmdAccess(args []any_t) []any_t {
	if len(args) < 2 {
		return []any_t{nil, "nyagos.access requilres two arguments"}
	}
	path := fmt.Sprint(args[0])
	mode, mode_ok := args[1].(int)
	if !mode_ok {
		return []any_t{nil, "mode value must be interger"}
	}
	if defined.DBG {
		fmt.Fprintf(os.Stderr, "given mode==%o\n", mode)
	}
	fi, err := os.Stat(path)

	var result bool
	if err != nil || fi == nil {
		result = false
	} else {
		if defined.DBG {
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

func CmdStat(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, TooFewArguments}
	}
	path := fmt.Sprint(args[0])
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

func CmdSetEnv(args []any_t) []any_t {
	if len(args) < 2 {
		return []any_t{nil, TooFewArguments}
	}
	name := fmt.Sprint(args[len(args)-2])
	value := fmt.Sprint(args[len(args)-1])
	if args[len(args)-1] != nil && len(value) > 0 {
		os.Setenv(name, value)
	} else {
		os.Unsetenv(name)
	}
	return []any_t{true}
}

func CmdGetEnv(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, TooFewArguments}
	}
	name := fmt.Sprint(args[len(args)-1])
	value, ok := shell.OurGetEnv(name)
	if ok && len(value) > 0 {
		return []any_t{value}
	} else {
		return []any_t{nil}
	}
}

func CmdAtoU(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, TooFewArguments}
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

func CmdUtoA(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, TooFewArguments}
	}
	utf8 := fmt.Sprint(args[0])
	bin, err := mbcs.UtoA(utf8)
	if err != nil {
		return []any_t{nil, err}
	}
	if len(bin) >= 1 {
		// trim the last zero byte from SJIS string
		return []any_t{bin[:len(bin)-1], nil}
	} else {
		return []any_t{"", nil}
	}
}

func CmdWhich(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, TooFewArguments}
	}
	name := fmt.Sprint(args[0])
	path := dos.LookPath(name, "NYAGOSPATH")
	if path != "" {
		return []any_t{path}
	} else {
		return []any_t{nil, name + ": Path not found"}
	}
}

func CmdGlob(args []any_t) []any_t {
	result := make([]string, 0)
	for _, arg1 := range args {
		wildcard := fmt.Sprint(arg1)
		list, err := findfile.Glob(wildcard)
		if list == nil || err != nil {
			result = append(result, wildcard)
		} else {
			result = append(result, list...)
		}
	}
	sort.StringSlice(result).Sort()
	return []any_t{result}
}

func CmdGetHistory(args []any_t) []any_t {
	if frame.DefaultHistory == nil {
		return []any_t{}
	}
	if len(args) >= 1 {
		if n, ok := args[len(args)-1].(int); ok {
			return []any_t{frame.DefaultHistory.At(n)}
		}
	}
	return []any_t{frame.DefaultHistory.Len()}
}

func CmdLenHistory(args []any_t) []any_t {
	if frame.DefaultHistory == nil {
		return []any_t{}
	}
	return []any_t{frame.DefaultHistory.Len()}
}

func CmdMsgBox(args []any_t) []any_t {
	var message string
	title := "nyagos"
	if len(args) >= 1 {
		message = fmt.Sprint(args[0])
	}
	if len(args) >= 2 {
		title = fmt.Sprint(args[1])
	}
	msgbox.Show(0, message, title, msgbox.OK)
	return []any_t{}
}

func CmdRawEval(args []any_t) []any_t {
	argv := make([]string, 0, len(args))
	for _, s := range args {
		argv = append(argv, fmt.Sprint(s))
	}
	cmd1 := exec.Command(argv[0], argv[1:]...)
	out, err := cmd1.Output()
	if err != nil {
		return []any_t{nil, err.Error()}
	} else {
		return []any_t{out}
	}
}

func CmdSetRuneWidth(args []any_t) []any_t {
	if len(args) < 2 {
		return []any_t{nil, "too few aruments"}
	}
	char, ok := args[0].(int)
	if !ok {
		return []any_t{nil, "not a number"}
	}
	width, ok := args[1].(int)
	if !ok {
		return []any_t{nil, "not a number"}
	}
	readline.SetCharWidth(rune(char), width)
	return []any_t{true}
}

func CmdCommonPrefix(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, "too few arguments"}
	}
	list := []string{}

	table, ok := args[0].(map[any_t]any_t)
	if !ok {
		return []any_t{nil, "not a table"}
	}
	for _, val := range table {
		list = append(list, fmt.Sprint(val))
	}
	return []any_t{completion.CommonPrefix(list)}
}

func CmdWriteSub(args []any_t, out io.Writer) []any_t {
	if f, ok := out.(*os.File); ok {
		cout := bufio.NewWriter(colorable.NewColorable(f))
		defer cout.Flush()
		out = cout
	}
	for i, arg1 := range args {
		if i > 0 {
			io.WriteString(out, "\t")
		}
		var str string
		if arg1 == nil {
			str = "nil"
		} else {
			switch v := arg1.(type) {
			case bool:
				if v {
					str = "true"
				} else {
					str = "false"
				}
			default:
				str = fmt.Sprint(v)
			}
		}
		io.WriteString(out, str)
	}
	return []any_t{true}
}

func CmdWrite(this *Param) []any_t {
	return CmdWriteSub(this.Args, this.Out)
}

func CmdWriteErr(this *Param) []any_t {
	return CmdWriteSub(this.Args, this.Err)
}

func CmdPrint(this *Param) []any_t {
	rc := CmdWrite(this)
	fmt.Fprintln(this.Out)
	return rc
}

func CmdRawExec(this *Param) []any_t {
	argv := make([]string, 0, len(this.Args))
	for _, arg1 := range this.Args {
		argv = append(argv, fmt.Sprint(arg1))
	}
	xcmd := exec.Command(argv[0], argv[1:]...)
	xcmd.Stdin = this.In
	xcmd.Stdout = this.Out
	xcmd.Stderr = this.Err
	err := xcmd.Run()
	errorlevel, errorlevelOk := dos.GetErrorLevel(xcmd)
	if !errorlevelOk {
		errorlevel = 255
	}
	if err != nil {
		fmt.Fprintln(xcmd.Stderr, err.Error())
		return []any_t{errorlevel, err.Error()}
	} else {
		return []any_t{errorlevel}
	}
}

func GetOption(args []any_t) []any_t {
	if len(args) < 2 {
		return []any_t{nil, "too few arguments"}
	}
	key := fmt.Sprint(args[1])
	ptr, ok := commands.BoolOptions[key]
	if !ok {
		return []any_t{nil, fmt.Sprintf("key: %s: not found", key)}
	}
	return []any_t{*ptr}
}

func SetOption(args []any_t) []any_t {
	if len(args) < 3 {
		return []any_t{nil, "too few arguments"}
	}
	key := fmt.Sprint(args[1])
	ptr, ok := commands.BoolOptions[key]
	if !ok || ptr == nil {
		return []any_t{nil, "key: %s: not found"}
	}
	val := args[2]
	*ptr = (val != nil && val != false && val == "")
	return []any_t{true}
}

func CmdBitAnd(args []any_t) []any_t {
	result := ^0
	for _, arg1tmp := range args {
		if arg1, ok := arg1tmp.(int); ok {
			result &= arg1
		} else {
			return []any_t{nil, fmt.Sprintf("%s : not a number", arg1tmp)}
		}
	}
	return []any_t{result}
}

func CmdBitOr(args []any_t) []any_t {
	result := 0
	for _, arg1tmp := range args {
		if arg1, ok := arg1tmp.(int); ok {
			result |= arg1
		} else {
			return []any_t{nil, fmt.Sprintf("%s : not a number", arg1tmp)}
		}
	}
	return []any_t{result}
}
