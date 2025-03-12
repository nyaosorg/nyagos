package functions

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/mattn/go-tty"

	"github.com/nyaosorg/go-box/v2"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-windows-findfile"

	"github.com/nyaosorg/nyagos/internal/completion"
	"github.com/nyaosorg/nyagos/internal/config"
	"github.com/nyaosorg/nyagos/internal/defined"
	"github.com/nyaosorg/nyagos/internal/frame"
	"github.com/nyaosorg/nyagos/internal/nodos"
	"github.com/nyaosorg/nyagos/internal/shell"
)

func toNumber(value any) (int, bool) {
	if f, ok := value.(float64); ok {
		return int(f), true
	}
	return 0, false
}

const TooFewArguments = "Too few arguments"

func toStr(arr []any, n int) string {
	if n < len(arr) {
		if defined.DBG {
			println(fmt.Sprint(arr[n]))
		}
		return fmt.Sprint(arr[n])
	}
	if defined.DBG {
		println("''")
	}
	return ""
}

func CmdChdir(args []any) []any {
	if len(args) >= 1 {
		nodos.Chdir(fmt.Sprint(args[0]))
		return []any{true}
	}
	return []any{nil, "directory is required"}
}

func CmdBox(this *Param) []any {
	args := this.Args
	if len(args) < 1 {
		return []any{nil, TooFewArguments}
	}
	t, ok := args[0].(map[any]any)
	if !ok {
		return []any{nil, "Not a table"}
	}
	if len(t) == 0 {
		return []any{}
	}
	sources := make([]string, 0, len(t))
	for i, _i := 1, len(t); i <= _i; i++ {
		if val, ok := t[i]; ok {
			sources = append(sources, fmt.Sprint(val))
		}
	}
	values := make([]any, 0)
	choice, err := box.SelectString(sources, true, this.Term)
	if err != nil {
		return []any{nil, err.Error()}
	}
	for _, s := range choice {
		values = append(values, s)
	}
	return values
}

func CmdResetCharWidth(args []any) []any {
	readline.ResetCharWidth()
	return []any{}
}

func CmdGetwd(args []any) []any {
	wd, err := os.Getwd()
	if err != nil {
		return []any{nil, err}
	}
	return []any{wd}
}

func CmdGetKey(args []any) []any {
	tty1, err := tty.Open()
	if err != nil {
		return []any{nil, err.Error()}
	}
	defer tty1.Close()
	for {
		r, err := tty1.ReadRune()
		if err != nil {
			return []any{nil, err.Error()}
		}
		if r != 0 {
			return []any{r, 0, 0}
		}
	}
}

func CmdGetKeys(args []any) []any {
	tty1, err := tty.Open()
	if err != nil {
		return []any{nil, err.Error()}
	}
	defer tty1.Close()
	key, err := readline.GetKey(tty1)
	if err != nil {
		return []any{nil, err.Error()}
	}
	return []any{key}
}

func CmdGetViewWidth(args []any) []any {
	tty1, err := tty.Open()
	if err != nil {
		return []any{nil, err.Error()}
	}
	defer tty1.Close()
	width, height, err := tty1.Size()
	if err != nil {
		return []any{nil, err.Error()}
	}
	return []any{width, height}
}

var rxEnv = regexp.MustCompile("%[^%]+%")

func expandEnv(str string) string {
	if (len(str) >= 2 && str[0] == '~' && os.IsPathSeparator(str[1])) || str == "~" {
		home := os.Getenv("HOME")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		str = home + str[1:]
	}
	return rxEnv.ReplaceAllStringFunc(str, func(s string) string {
		name := s[1 : len(s)-1]
		return os.Getenv(name)
	})
}

func CmdPathJoin(args []any) []any {
	if len(args) < 1 {
		return []any{""}
	}
	path := expandEnv(fmt.Sprint(args[0]))
	for i, _i := 1, len(args); i < _i; i++ {
		sub := expandEnv(fmt.Sprint(args[i]))
		path = filepath.Join(path, sub)
	}
	return []any{path}
}

func CmdDirName(args []any) []any {
	if len(args) < 1 {
		return []any{""}
	}
	result := []any{}
	for _, arg1 := range args {
		if s, ok := arg1.(string); ok {
			result = append(result, any(filepath.Dir(s)))
		} else {
			result = append(result, any(""))
		}
	}
	return result
}

func CmdAccess(args []any) []any {
	if len(args) < 2 {
		return []any{nil, "nyagos.access requilres two arguments"}
	}
	path := fmt.Sprint(args[0])
	mode, ok := toNumber(args[1])
	if !ok {
		return []any{nil, "mode value must be interger"}
	}
	if defined.DBG {
		fmt.Fprintf(os.Stderr, "given mode==%o\n", mode)
	}
	fi, err := os.Stat(expandEnv(path))

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
	return []any{result}
}

func CmdStat(args []any) []any {
	if len(args) < 1 {
		return []any{nil, TooFewArguments}
	}
	path := expandEnv(fmt.Sprint(args[0]))
	var stat os.FileInfo
	var _path string
	if len(path) > 0 && path[len(path)-1] == '\\' {
		_path = filepath.Join(path, ".")
	} else {
		_path = path
	}
	statErr := findfile.Walk(_path, func(f *findfile.FileInfo) bool {
		stat = f
		return false
	})
	if statErr != nil {
		return []any{nil, statErr}
	}
	if stat == nil {
		return []any{nil, fmt.Errorf("%s: failed to stat", path)}
	}
	t := stat.ModTime()
	return []any{
		map[string]any{
			"name":  stat.Name(),
			"size":  stat.Size(),
			"isdir": stat.IsDir(),
			"mtime": map[string]any{
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

func CmdSetEnv(args []any) []any {
	if len(args) < 2 {
		return []any{nil, TooFewArguments}
	}
	name := fmt.Sprint(args[len(args)-2])
	value := fmt.Sprint(args[len(args)-1])
	if args[len(args)-1] != nil && len(value) > 0 {
		os.Setenv(name, value)
	} else {
		os.Unsetenv(name)
	}
	return []any{true}
}

func CmdGetEnv(args []any) []any {
	if len(args) < 1 {
		return []any{nil, TooFewArguments}
	}
	name := fmt.Sprint(args[len(args)-1])
	value, ok := shell.OurGetEnv(name)
	if ok && len(value) > 0 {
		return []any{value}
	}
	return []any{nil}
}

func CmdWhich(args []any) []any {
	if len(args) < 1 {
		return []any{nil, TooFewArguments}
	}
	name := fmt.Sprint(args[0])
	path := nodos.LookPath(shell.LookCurdirOrder, name, "NYAGOSPATH")
	if path != "" {
		return []any{path}
	}
	return []any{nil, name + ": Path not found"}
}

func CmdGlob(args []any) []any {
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
	return []any{result}
}

func CmdGetHistory(args []any) []any {
	if frame.DefaultHistory == nil {
		return []any{}
	}
	if len(args) >= 1 {
		if n, ok := toNumber(args[len(args)-1]); ok {
			return []any{frame.DefaultHistory.At(n)}
		}
	}
	return []any{frame.DefaultHistory.Len()}
}

func CmdLenHistory(args []any) []any {
	if frame.DefaultHistory == nil {
		return []any{}
	}
	return []any{frame.DefaultHistory.Len()}
}

func CmdRawEval(this *Param) []any {
	argv := stackToSlice(this)
	cmd1 := exec.Command(argv[0], argv[1:]...)
	out, err := cmd1.Output()
	if err != nil {
		return []any{nil, err.Error()}
	}
	return []any{out}
}

func CmdSetRuneWidth(args []any) []any {
	if len(args) < 2 {
		return []any{nil, "too few aruments"}
	}
	char, ok := toNumber(args[0])
	if !ok {
		return []any{nil, "not a number"}
	}
	width, ok := toNumber(args[1])
	if !ok {
		return []any{nil, "not a number"}
	}
	readline.SetCharWidth(rune(char), width)
	return []any{true}
}

func CmdCommonPrefix(args []any) []any {
	if len(args) < 1 {
		return []any{nil, "too few arguments"}
	}
	list := []string{}

	table, ok := args[0].(map[any]any)
	if !ok {
		return []any{nil, "not a table"}
	}
	for _, val := range table {
		list = append(list, fmt.Sprint(val))
	}
	return []any{completion.CommonPrefix(list)}
}

func CmdWriteSub(this *Param, out io.Writer) []any {
	args := this.Args
	if f, ok := out.(*os.File); ok && isatty.IsTerminal(f.Fd()) {
		cout := bufio.NewWriter(this.Term)
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
	return []any{true}
}

func CmdWrite(this *Param) []any {
	return CmdWriteSub(this, this.Out)
}

func CmdWriteErr(this *Param) []any {
	return CmdWriteSub(this, this.Err)
}

func CmdPrint(this *Param) []any {
	rc := CmdWrite(this)
	fmt.Fprintln(this.Out)
	return rc
}

func stackToSlice(this *Param) []string {
	argv := make([]string, 0, len(this.Args))
	for _, arg1 := range this.Args {
		if table, ok := arg1.(map[interface{}]interface{}); ok {
			// Support both {0..(n-1)} and {1..n}
			for i := 0; i <= len(table); i++ {
				if _, ok = table[i]; ok { // check out of range here
					argv = append(argv, fmt.Sprint(table[i]))
				}
			}
		} else {
			argv = append(argv, fmt.Sprint(arg1))
		}
	}
	return argv
}

func GetOption(args []any) []any {
	if len(args) < 2 {
		return []any{nil, "too few arguments"}
	}
	key := fmt.Sprint(args[1])
	ptr, ok := config.Bools.Load(key)
	if !ok {
		return []any{nil, fmt.Sprintf("key: %s: not found", key)}
	}
	return []any{ptr.Get()}
}

func SetOption(args []any) []any {
	if len(args) < 3 {
		return []any{nil, "too few arguments"}
	}
	key := fmt.Sprint(args[1])
	ptr, ok := config.Bools.Load(key)
	if !ok || ptr == nil {
		return []any{nil, "key: %s: not found"}
	}
	val := args[2]
	if val == nil {
		ptr.Set(false)
	} else if s, ok := val.(string); ok && s == "" {
		ptr.Set(false)
	} else if b, ok := val.(bool); ok {
		ptr.Set(b)
	} else {
		ptr.Set(true)
	}
	return []any{true}
}

func bitOperators(args []any, result int, f func(int, int) int) []any {
	for _, arg1tmp := range args {
		if arg1, ok := toNumber(arg1tmp); ok {
			result = f(result, arg1)
		} else {
			return []any{nil, fmt.Sprintf("%s : not a number", arg1tmp)}
		}
	}
	return []any{result}
}

func CmdBitAnd(args []any) []any {
	return bitOperators(args, ^0, func(r, v int) int { return r & v })
}

func CmdBitOr(args []any) []any {
	return bitOperators(args, 0, func(r, v int) int { return r | v })
}

func CmdBitXor(args []any) []any {
	return bitOperators(args, 0, func(r, v int) int { return r ^ v })
}

func CmdFields(args []any) []any {
	if len(args) <= 0 {
		return []any{nil}
	}
	fields := strings.Fields(fmt.Sprint(args[0]))
	return []any{fields}
}

func CmdEnvAdd(args []any) []any {
	if len(args) >= 1 {
		list := make([]string, 1, len(args))
		name := strings.ToUpper(fmt.Sprint(args[0]))
		list[0] = os.Getenv(name)
		for _, s := range args[1:] {
			list = append(list, expandEnv(fmt.Sprint(s)))
		}
		os.Setenv(name, nodos.JoinList(list...))
	}
	return []any{}
}

func CmdEnvDel(args []any) (result []any) {
	if len(args) >= 1 {
		name := strings.ToUpper(fmt.Sprint(args[0]))
		list := filepath.SplitList(os.Getenv(name))
		newlist := make([]string, 0, len(list))

		for _, e := range list {
			E := strings.ToUpper(e)
			doRemove := false
			for _, substr := range args[1:] {
				if strings.Contains(E, strings.ToUpper(fmt.Sprint(substr))) {
					doRemove = true
					break
				}
			}
			if doRemove {
				result = append(result, e)
			} else {
				newlist = append(newlist, e)
			}
		}
		os.Setenv(name, strings.Join(newlist, string(os.PathListSeparator)))
	}
	return
}

func CmdCompleteForFiles(args []any) []any {
	if len(args) < 1 {
		return []any{nil, errors.New("too few arguments")}
	}
	if s, ok := args[0].(string); ok {
		elements, err := completion.ListUpFiles(
			context.TODO(),
			completion.DoNotUncCompletion,
			s)
		if err != nil {
			return []any{nil, err.Error()}
		}
		result := make([]string, len(elements))
		for i := 0; i < len(elements); i++ {
			result[i] = elements[i].String()
		}
		return []any{result}
	}
	return []any{nil, errors.New("invalid arguments")}
}
