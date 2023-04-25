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
	"unicode/utf8"

	"golang.org/x/term"

	"github.com/mattn/go-isatty"

	"github.com/nyaosorg/go-box/v2"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-windows-findfile"

	"github.com/nyaosorg/nyagos/internal/commands"
	"github.com/nyaosorg/nyagos/internal/completion"
	"github.com/nyaosorg/nyagos/internal/defined"
	"github.com/nyaosorg/nyagos/internal/frame"
	"github.com/nyaosorg/nyagos/internal/nodos"
	"github.com/nyaosorg/nyagos/internal/shell"
)

type anyT = interface{}

func toNumber(value anyT) (int, bool) {
	if f, ok := value.(float64); ok {
		return int(f), true
	}
	return 0, false
}

const TooFewArguments = "Too few arguments"

func toStr(arr []anyT, n int) string {
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

func CmdChdir(args []anyT) []anyT {
	if len(args) >= 1 {
		nodos.Chdir(fmt.Sprint(args[0]))
		return []anyT{true}
	}
	return []anyT{nil, "directory is required"}
}

func CmdBox(this *Param) []anyT {
	args := this.Args
	if len(args) < 1 {
		return []anyT{nil, TooFewArguments}
	}
	t, ok := args[0].(map[anyT]anyT)
	if !ok {
		return []anyT{nil, "Not a table"}
	}
	if len(t) == 0 {
		return []anyT{}
	}
	sources := make([]string, 0, len(t))
	for i, _i := 1, len(t); i <= _i; i++ {
		if val, ok := t[i]; ok {
			sources = append(sources, fmt.Sprint(val))
		}
	}
	values := make([]anyT, 0)
	for _, s := range box.ChoiceMulti(sources, this.Term) {
		values = append(values, s)
	}
	return values
}

func CmdResetCharWidth(args []anyT) []anyT {
	readline.ResetCharWidth()
	return []anyT{}
}

func CmdGetwd(args []anyT) []anyT {
	wd, err := os.Getwd()
	if err != nil {
		return []anyT{nil, err}
	}
	return []anyT{wd}
}

func CmdGetKey(args []anyT) []anyT {
	stdin := int(os.Stdin.Fd())
	state, err := term.MakeRaw(stdin)
	if err != nil {
		return []anyT{nil, err.Error()}
	}
	defer term.Restore(stdin, state)

	for {
		var buffer [256]byte

		n, err := os.Stdin.Read(buffer[:])
		if err != nil {
			return []anyT{nil, err.Error()}
		}
		key := buffer[:n]
		for len(key) > 0 {
			r, size := utf8.DecodeRune(key)
			if r != 0 {
				return []anyT{r, 0, 0}
			}
			key = key[size:]
		}
	}
}

func CmdGetViewWidth(args []anyT) []anyT {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return []anyT{nil, err.Error()}
	}
	return []anyT{width, height}
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

func CmdPathJoin(args []anyT) []anyT {
	if len(args) < 1 {
		return []anyT{""}
	}
	path := expandEnv(fmt.Sprint(args[0]))
	for i, _i := 1, len(args); i < _i; i++ {
		sub := expandEnv(fmt.Sprint(args[i]))
		path = filepath.Join(path, sub)
	}
	return []anyT{path}
}

func CmdDirName(args []anyT) []anyT {
	if len(args) < 1 {
		return []anyT{""}
	}
	result := []anyT{}
	for _, arg1 := range args {
		if s, ok := arg1.(string); ok {
			result = append(result, anyT(filepath.Dir(s)))
		} else {
			result = append(result, anyT(""))
		}
	}
	return result
}

func CmdAccess(args []anyT) []anyT {
	if len(args) < 2 {
		return []anyT{nil, "nyagos.access requilres two arguments"}
	}
	path := fmt.Sprint(args[0])
	mode, ok := toNumber(args[1])
	if !ok {
		return []anyT{nil, "mode value must be interger"}
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
	return []anyT{result}
}

func CmdStat(args []anyT) []anyT {
	if len(args) < 1 {
		return []anyT{nil, TooFewArguments}
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
		return []anyT{nil, statErr}
	}
	if stat == nil {
		return []anyT{nil, fmt.Errorf("%s: failed to stat", path)}
	}
	t := stat.ModTime()
	return []anyT{
		map[string]anyT{
			"name":  stat.Name(),
			"size":  stat.Size(),
			"isdir": stat.IsDir(),
			"mtime": map[string]anyT{
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

func CmdSetEnv(args []anyT) []anyT {
	if len(args) < 2 {
		return []anyT{nil, TooFewArguments}
	}
	name := fmt.Sprint(args[len(args)-2])
	value := fmt.Sprint(args[len(args)-1])
	if args[len(args)-1] != nil && len(value) > 0 {
		os.Setenv(name, value)
	} else {
		os.Unsetenv(name)
	}
	return []anyT{true}
}

func CmdGetEnv(args []anyT) []anyT {
	if len(args) < 1 {
		return []anyT{nil, TooFewArguments}
	}
	name := fmt.Sprint(args[len(args)-1])
	value, ok := shell.OurGetEnv(name)
	if ok && len(value) > 0 {
		return []anyT{value}
	}
	return []anyT{nil}
}

func CmdWhich(args []anyT) []anyT {
	if len(args) < 1 {
		return []anyT{nil, TooFewArguments}
	}
	name := fmt.Sprint(args[0])
	path := nodos.LookPath(shell.LookCurdirOrder, name, "NYAGOSPATH")
	if path != "" {
		return []anyT{path}
	}
	return []anyT{nil, name + ": Path not found"}
}

func CmdGlob(args []anyT) []anyT {
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
	return []anyT{result}
}

func CmdGetHistory(args []anyT) []anyT {
	if frame.DefaultHistory == nil {
		return []anyT{}
	}
	if len(args) >= 1 {
		if n, ok := toNumber(args[len(args)-1]); ok {
			return []anyT{frame.DefaultHistory.At(n)}
		}
	}
	return []anyT{frame.DefaultHistory.Len()}
}

func CmdLenHistory(args []anyT) []anyT {
	if frame.DefaultHistory == nil {
		return []anyT{}
	}
	return []anyT{frame.DefaultHistory.Len()}
}

func CmdRawEval(this *Param) []anyT {
	argv := stackToSlice(this)
	cmd1 := exec.Command(argv[0], argv[1:]...)
	out, err := cmd1.Output()
	if err != nil {
		return []anyT{nil, err.Error()}
	}
	return []anyT{out}
}

func CmdSetRuneWidth(args []anyT) []anyT {
	if len(args) < 2 {
		return []anyT{nil, "too few aruments"}
	}
	char, ok := toNumber(args[0])
	if !ok {
		return []anyT{nil, "not a number"}
	}
	width, ok := toNumber(args[1])
	if !ok {
		return []anyT{nil, "not a number"}
	}
	readline.SetCharWidth(rune(char), width)
	return []anyT{true}
}

func CmdCommonPrefix(args []anyT) []anyT {
	if len(args) < 1 {
		return []anyT{nil, "too few arguments"}
	}
	list := []string{}

	table, ok := args[0].(map[anyT]anyT)
	if !ok {
		return []anyT{nil, "not a table"}
	}
	for _, val := range table {
		list = append(list, fmt.Sprint(val))
	}
	return []anyT{completion.CommonPrefix(list)}
}

func CmdWriteSub(this *Param, out io.Writer) []anyT {
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
	return []anyT{true}
}

func CmdWrite(this *Param) []anyT {
	return CmdWriteSub(this, this.Out)
}

func CmdWriteErr(this *Param) []anyT {
	return CmdWriteSub(this, this.Err)
}

func CmdPrint(this *Param) []anyT {
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

func GetOption(args []anyT) []anyT {
	if len(args) < 2 {
		return []anyT{nil, "too few arguments"}
	}
	key := fmt.Sprint(args[1])
	ptr, ok := commands.BoolOptions.Load(key)
	if !ok {
		return []anyT{nil, fmt.Sprintf("key: %s: not found", key)}
	}
	return []anyT{*ptr.V}
}

func SetOption(args []anyT) []anyT {
	if len(args) < 3 {
		return []anyT{nil, "too few arguments"}
	}
	key := fmt.Sprint(args[1])
	ptr, ok := commands.BoolOptions.Load(key)
	if !ok || ptr == nil {
		return []anyT{nil, "key: %s: not found"}
	}
	val := args[2]
	if val == nil {
		*ptr.V = false
	} else if s, ok := val.(string); ok && s == "" {
		*ptr.V = false
	} else if b, ok := val.(bool); ok {
		*ptr.V = b
	} else {
		*ptr.V = true
	}
	return []anyT{true}
}

func bitOperators(args []anyT, result int, f func(int, int) int) []anyT {
	for _, arg1tmp := range args {
		if arg1, ok := toNumber(arg1tmp); ok {
			result = f(result, arg1)
		} else {
			return []anyT{nil, fmt.Sprintf("%s : not a number", arg1tmp)}
		}
	}
	return []anyT{result}
}

func CmdBitAnd(args []anyT) []anyT {
	return bitOperators(args, ^0, func(r, v int) int { return r & v })
}

func CmdBitOr(args []anyT) []anyT {
	return bitOperators(args, 0, func(r, v int) int { return r | v })
}

func CmdBitXor(args []anyT) []anyT {
	return bitOperators(args, 0, func(r, v int) int { return r ^ v })
}

func CmdFields(args []anyT) []anyT {
	if len(args) <= 0 {
		return []anyT{nil}
	}
	fields := strings.Fields(fmt.Sprint(args[0]))
	return []anyT{fields}
}

func CmdEnvAdd(args []anyT) []anyT {
	if len(args) >= 1 {
		list := make([]string, 1, len(args))
		name := strings.ToUpper(fmt.Sprint(args[0]))
		list[0] = os.Getenv(name)
		for _, s := range args[1:] {
			list = append(list, expandEnv(fmt.Sprint(s)))
		}
		os.Setenv(name, nodos.JoinList(list...))
	}
	return []anyT{}
}

func CmdEnvDel(args []anyT) (result []anyT) {
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

func CmdCompleteForFiles(args []anyT) []anyT {
	if len(args) < 1 {
		return []anyT{nil, errors.New("too few arguments")}
	}
	if s, ok := args[0].(string); ok {
		elements, err := completion.ListUpFiles(
			context.TODO(),
			completion.DoNotUncCompletion,
			s)
		if err != nil {
			return []anyT{nil, err.Error()}
		}
		result := make([]string, len(elements))
		for i := 0; i < len(elements); i++ {
			result[i] = elements[i].String()
		}
		return []anyT{result}
	}
	return []anyT{nil, errors.New("invalid arguments")}
}
