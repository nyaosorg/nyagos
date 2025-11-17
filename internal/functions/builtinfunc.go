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

	"github.com/nyaosorg/go-box/v3"
	"github.com/nyaosorg/go-ttyadapter/tty8"
	"github.com/nyaosorg/go-windows-findfile"

	"github.com/nyaosorg/nyagos/internal/completion"
	"github.com/nyaosorg/nyagos/internal/config"
	"github.com/nyaosorg/nyagos/internal/defined"
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

func CmdChdir(param *Param) []any {
	args := param.Args
	if len(args) >= 1 {
		nodos.Chdir(fmt.Sprint(args[0]))
		return []any{true}
	}
	return []any{nil, "directory is required"}
}

func CmdBox(param *Param) []any {
	args := param.Args
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
	choice, err := box.SelectString(sources, true, param.Term)
	if err != nil {
		return []any{nil, err.Error()}
	}
	for _, s := range choice {
		values = append(values, s)
	}
	return values
}

func CmdGetwd(param *Param) []any {
	wd, err := os.Getwd()
	if err != nil {
		return []any{nil, err}
	}
	return []any{wd}
}

func CmdGetKey(*Param) []any {
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

func CmdGetKeys(*Param) []any {
	tty := &tty8.Tty{}
	if err := tty.Open(nil); err != nil {
		return []any{nil, err.Error()}
	}
	defer tty.Close()
	key, err := tty.GetKey()
	if err != nil {
		return []any{nil, err.Error()}
	}
	return []any{key}
}

func CmdGetViewWidth(*Param) []any {
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

func CmdPathJoin(param *Param) []any {
	args := param.Args
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

func CmdDirName(param *Param) []any {
	args := param.Args
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

func CmdAccess(param *Param) []any {
	args := param.Args
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

func CmdStat(param *Param) []any {
	args := param.Args
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

func CmdSetEnv(param *Param) []any {
	args := param.Args
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

func CmdGetEnv(param *Param) []any {
	args := param.Args
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

func CmdWhich(param *Param) []any {
	args := param.Args
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

func CmdGlob(param *Param) []any {
	args := param.Args
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

func CmdGetHistory(param *Param) []any {
	editor := param.Editor
	if editor == nil {
		return []any{}
	}
	history := editor.History
	if history == nil {
		return []any{}
	}
	args := param.Args
	if len(args) >= 1 {
		if n, ok := toNumber(args[len(args)-1]); ok {
			return []any{history.At(n)}
		}
	}
	return []any{history.Len()}
}

func CmdLenHistory(param *Param) []any {
	editor := param.Editor
	if editor == nil {
		return []any{}
	}
	history := editor.History
	if history == nil {
		return []any{}
	}
	return []any{history.Len()}
}

func CmdRawEval(param *Param) []any {
	argv := stackToSlice(param)
	cmd1 := exec.Command(argv[0], argv[1:]...)
	out, err := cmd1.Output()
	if err != nil {
		return []any{nil, err.Error()}
	}
	return []any{out}
}

func CmdCommonPrefix(param *Param) []any {
	args := param.Args
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

func CmdWriteSub(param *Param, out io.Writer) []any {
	args := param.Args
	if f, ok := out.(*os.File); ok && isatty.IsTerminal(f.Fd()) {
		cout := bufio.NewWriter(param.Term)
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

func CmdWrite(param *Param) []any {
	return CmdWriteSub(param, param.Out)
}

func CmdWriteErr(param *Param) []any {
	return CmdWriteSub(param, param.Err)
}

func CmdPrint(param *Param) []any {
	rc := CmdWrite(param)
	fmt.Fprintln(param.Out)
	return rc
}

func stackToSlice(param *Param) []string {
	argv := make([]string, 0, len(param.Args))
	for _, arg1 := range param.Args {
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

func GetOption(param *Param) []any {
	args := param.Args
	if len(args) < 2 {
		return []any{nil, "too few arguments"}
	}
	key := fmt.Sprint(args[1])
	if ptr, ok := config.Bools.Load(key); ok {
		return []any{ptr.Get()}
	}
	if ptr, ok := config.Strings.Load(key); ok {
		return []any{ptr.Get()}
	}
	return []any{nil, fmt.Sprintf("key: %s: not found", key)}
}

func SetOption(param *Param) []any {
	args := param.Args
	if len(args) < 3 {
		return []any{nil, "too few arguments"}
	}
	key := fmt.Sprint(args[1])
	if ptr, ok := config.Bools.Load(key); ok {
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
	if ptr, ok := config.Strings.Load(key); ok {
		if s, ok := args[2].(string); ok {
			ptr.Set(s)
			return []any{true}
		} else {
			return []any{nil, "not string"}
		}
	}
	return []any{nil, "key: %s: not found"}
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

func CmdBitAnd(param *Param) []any {
	args := param.Args
	return bitOperators(args, ^0, func(r, v int) int { return r & v })
}

func CmdBitOr(param *Param) []any {
	args := param.Args
	return bitOperators(args, 0, func(r, v int) int { return r | v })
}

func CmdBitXor(param *Param) []any {
	args := param.Args
	return bitOperators(args, 0, func(r, v int) int { return r ^ v })
}

func CmdFields(param *Param) []any {
	args := param.Args
	if len(args) <= 0 {
		return []any{nil}
	}
	fields := strings.Fields(fmt.Sprint(args[0]))
	return []any{fields}
}

func CmdEnvAdd(param *Param) []any {
	args := param.Args
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

func CmdEnvDel(param *Param) (result []any) {
	args := param.Args
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

func CmdCompleteForFiles(param *Param) []any {
	args := param.Args
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

func CmdSetNextLine(param *Param) []any {
	editor := param.Editor
	if editor == nil {
		return []any{nil, "can not find the current editor"}
	}

	args := param.Args
	var buffer strings.Builder
	if len(args) > 0 {
		for {
			fmt.Fprint(&buffer, args[0])
			args = args[1:]
			if len(args) <= 0 {
				editor.Default = buffer.String()
				break
			}
			buffer.WriteByte(' ')
		}
	}
	return []any{true}
}
