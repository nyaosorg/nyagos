package functions

import (
	"bufio"
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

	"github.com/zetamatta/go-box/v2"
	"github.com/zetamatta/go-findfile"

	"github.com/zetamatta/nyagos/commands"
	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/defined"
	"github.com/zetamatta/nyagos/frame"
	"github.com/zetamatta/nyagos/nodos"
	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/shell"
)

type any_t = interface{}

func toNumber(value any_t) (int, bool) {
	if f, ok := value.(float64); ok {
		return int(f), true
	}
	return 0, false
}

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

func CmdChdir(args []any_t) []any_t {
	if len(args) >= 1 {
		nodos.Chdir(fmt.Sprint(args[0]))
		return []any_t{true}
	}
	return []any_t{nil, "directory is required"}
}

func CmdBox(this *Param) []any_t {
	args := this.Args
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
	values := make([]any_t, 0)
	for _, s := range box.ChoiceMulti(sources, this.Term) {
		values = append(values, s)
	}
	return values
}

func CmdResetCharWidth(args []any_t) []any_t {
	readline.ResetCharWidth()
	return []any_t{}
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
	tty1, err := tty.Open()
	if err != nil {
		return []any_t{nil, err.Error()}
	}
	defer tty1.Close()
	for {
		r, err := tty1.ReadRune()
		if err != nil {
			return []any_t{nil, err.Error()}
		}
		if r != 0 {
			return []any_t{r, 0, 0}
		}
	}
}

func CmdGetViewWidth(args []any_t) []any_t {
	tty1, err := tty.Open()
	if err != nil {
		return []any_t{nil, err.Error()}
	}
	defer tty1.Close()
	width, height, err := tty1.Size()
	if err != nil {
		return []any_t{nil, err.Error()}
	}
	return []any_t{width, height}
}

var rxEnv = regexp.MustCompile("%[^%]+%")

func expandEnv(str string) string {
	if len(str) >= 2 && str[0] == '~' && (str[1] == '\\' || str[1] == '/') {
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

func CmdPathJoin(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{""}
	}
	path := expandEnv(fmt.Sprint(args[0]))
	for i, i_ := 1, len(args); i < i_; i++ {
		sub := expandEnv(fmt.Sprint(args[i]))
		path = filepath.Join(path, sub)
	}
	return []any_t{path}
}

func CmdDirName(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{""}
	}
	result := []any_t{}
	for _, arg1 := range args {
		if s, ok := arg1.(string); ok {
			result = append(result, any_t(filepath.Dir(s)))
		} else {
			result = append(result, any_t(""))
		}
	}
	return result
}

func CmdAccess(args []any_t) []any_t {
	if len(args) < 2 {
		return []any_t{nil, "nyagos.access requilres two arguments"}
	}
	path := fmt.Sprint(args[0])
	mode, mode_ok := toNumber(args[1])
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

func CmdWhich(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, TooFewArguments}
	}
	name := fmt.Sprint(args[0])
	path := nodos.LookPath(shell.LookCurdirOrder, name, "NYAGOSPATH")
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
		if n, ok := toNumber(args[len(args)-1]); ok {
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

func CmdRawEval(this *Param) []any_t {
	argv := stackToSlice(this)
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
	char, ok := toNumber(args[0])
	if !ok {
		return []any_t{nil, "not a number"}
	}
	width, ok := toNumber(args[1])
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

func CmdWriteSub(this *Param, out io.Writer) []any_t {
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
	return []any_t{true}
}

func CmdWrite(this *Param) []any_t {
	return CmdWriteSub(this, this.Out)
}

func CmdWriteErr(this *Param) []any_t {
	return CmdWriteSub(this, this.Err)
}

func CmdPrint(this *Param) []any_t {
	rc := CmdWrite(this)
	fmt.Fprintln(this.Out)
	return rc
}

func stackToSlice(this *Param) []string {
	argv := make([]string, 0, len(this.Args))
	for _, arg1 := range this.Args {
		if table, ok := arg1.(map[interface{}]interface{}); ok {
			for i := 0; i < len(table); i++ {
				if _, ok = table[i]; ok {
					argv = append(argv, fmt.Sprint(table[i]))
				}
			}
		} else {
			argv = append(argv, fmt.Sprint(arg1))
		}
	}
	return argv
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
	return []any_t{*ptr.V}
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
	if val == nil {
		*ptr.V = false
	} else if s, ok := val.(string); ok && s == "" {
		*ptr.V = false
	} else if b, ok := val.(bool); ok {
		*ptr.V = b
	} else {
		*ptr.V = true
	}
	return []any_t{true}
}

func bitOperators(args []any_t, result int, f func(int, int) int) []any_t {
	for _, arg1tmp := range args {
		if arg1, ok := toNumber(arg1tmp); ok {
			result = f(result, arg1)
		} else {
			return []any_t{nil, fmt.Sprintf("%s : not a number", arg1tmp)}
		}
	}
	return []any_t{result}
}

func CmdBitAnd(args []any_t) []any_t {
	return bitOperators(args, ^0, func(r, v int) int { return r & v })
}

func CmdBitOr(args []any_t) []any_t {
	return bitOperators(args, 0, func(r, v int) int { return r | v })
}

func CmdBitXor(args []any_t) []any_t {
	return bitOperators(args, 0, func(r, v int) int { return r ^ v })
}

func CmdFields(args []any_t) []any_t {
	if len(args) <= 0 {
		return []any_t{nil}
	}
	fields := strings.Fields(fmt.Sprint(args[0]))
	return []any_t{fields}
}

func CmdEnvAdd(args []any_t) []any_t {
	if len(args) >= 1 {
		list := make([]string, 1, len(args))
		name := strings.ToUpper(fmt.Sprint(args[0]))
		list[0] = os.Getenv(name)
		for _, s := range args[1:] {
			list = append(list, expandEnv(fmt.Sprint(s)))
		}
		os.Setenv(name, nodos.JoinList(list...))
	}
	return []any_t{}
}

func CmdEnvDel(args []any_t) []any_t {
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
			if !doRemove {
				newlist = append(newlist, e)
			}
		}
		os.Setenv(name, strings.Join(newlist, string(os.PathListSeparator)))
	}
	return []any_t{}
}
