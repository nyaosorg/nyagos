package completion

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/zetamatta/go-box/v2"

	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/texts"
)

type Element interface {
	String() string
	Display() string
}

type Element2 [2]string

func (s Element2) String() string  { return s[0] }
func (s Element2) Display() string { return s[1] }

type Element1 string

func (s Element1) String() string  { return string(s) }
func (s Element1) Display() string { return string(s) }

type List struct {
	AllLine string
	List    []Element
	RawWord string // have quotation
	Word    string
	Pos     int
	Field   []string
	Left    string
}

var UseSlash = false

func isTop(s string, indexes [][]int) bool {
	if len(indexes) < 1 {
		return true
	}
	if len(indexes) == 1 {
		return indexes[0][1] == len(s)
	}
	prev := s[indexes[len(indexes)-2][0]:indexes[len(indexes)-2][1]]
	return prev == ";" || prev == "|" || prev == "&"
}

type CustomCompleter interface {
	Complete(context.Context, []string) ([]Element, error)
	String() string
}

type customComplete struct {
	Func func(context.Context, []string) ([]Element, error)
	Name string
}

func (f customComplete) Complete(ctx context.Context, args []string) ([]Element, error) {
	return f.Func(ctx, args)
}

func (f customComplete) String() string {
	return f.Name
}

var CustomCompletion = map[string]CustomCompleter{
	"set":      &customComplete{Func: completionSet, Name: "Built-in `set` completer"},
	"cd":       &customComplete{Func: completionCd, Name: "Built-in `cd` completer"},
	"env":      &customComplete{Func: completionEnv, Name: "built-in `env` completer"},
	"which":    &customComplete{Func: completionWhich, Name: "built-in `which` completer"},
	"pushd":    &customComplete{Func: completionCd, Name: "Built-in `pushd` completer"},
	"rmdir":    &customComplete{Func: completionCd, Name: "Built-in `rmdir` completer"},
	"rd":       &customComplete{Func: completionCd, Name: "Built-in `rmdir` completer"},
	"kill":     &customComplete{Func: completionProcessName, Name: "Built-in `kill` completer"},
	"taskkill": &customComplete{Func: completionTaskKill, Name: "Built-in `taskkill` completer"},
}

func listUpComplete(ctx context.Context, this *readline.Buffer) (*List, rune, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var err error
	rv := &List{
		AllLine: this.String(),
		Left:    this.SubString(0, this.Cursor),
	}

	// environment completion.

	indexes := texts.SplitLikeShell(rv.Left)
	for _, p := range indexes {
		rv.Field = append(rv.Field, rv.Left[p[0]:p[1]])
	}
	rv.List, rv.Pos, err = listUpEnv(rv.AllLine)
	default_delimiter := rune(readline.Delimiters[0])
	if len(rv.List) > 0 && rv.Pos >= 0 && err == nil {
		rv.RawWord = rv.AllLine[rv.Pos:]
		rv.Word = rv.RawWord
		return rv, default_delimiter, nil
	}

	// filename or commandname completion
	rv.RawWord, rv.Pos = this.CurrentWord()
	found_delimter := false
	rv.Word = strings.Map(func(c rune) rune {
		if strings.ContainsRune(readline.Delimiters, c) {
			if !found_delimter {
				default_delimiter = c
			}
			found_delimter = true
			return -1
		} else {
			return c
		}
	}, rv.RawWord)

	start := strings.LastIndexAny(rv.Word, ";=") + 1

	replace := false
	if isTop(rv.Left, indexes) {
		rv.List, err = listUpCommands(ctx, rv.Word[start:])
	} else {
		args := make([]string, 0, len(rv.Field))
		for i, w := range rv.Field {
			if indexes[i][0] > this.Cursor {
				break
			}
			args = append(args, strings.Replace(w, `"`, ``, -1))
		}
		if len(indexes) <= 0 || indexes[len(indexes)-1][1] < this.Cursor {
			args = append(args, "")
		}

		if f, ok := CustomCompletion[strings.ToLower(args[0])]; ok {
			rv.List, err = f.Complete(ctx, args)
			if rv.List != nil && err == nil {
				replace = true
			} else {
				rv.List, err = listUpFiles(ctx, rv.Word[start:])
			}
		} else {
			rv.List, err = listUpFiles(ctx, rv.Word[start:])
		}
	}
	if err != nil {
		return rv, default_delimiter, err
	}
	if !replace {
		for i := 0; i < len(rv.List); i++ {
			rv.List[i] = Element2{
				rv.Word[:start] + rv.List[i].String(),
				rv.List[i].Display(),
			}
		}
	}
	for _, f := range HookToList {
		rv, err = f(ctx, this, rv)
		if err != nil {
			break
		}
	}
	return rv, default_delimiter, err
}

func toComplete(source []Element) []string {
	result := make([]string, len(source))
	for key, val := range source {
		result[key] = val.String()
	}
	return result
}

func toDisplay(source []Element) []string {
	result := make([]string, len(source))
	for key, val := range source {
		result[key] = val.Display()
	}
	return result
}

func KeyFuncCompletionList(ctx context.Context, this *readline.Buffer) readline.Result {
	comp, _, err := listUpComplete(ctx, this)
	if err != nil {
		fmt.Fprintf(this.Out, "\n%s\n", err)
		this.RepaintAll()
		return readline.CONTINUE
	}
	if comp == nil {
		return readline.CONTINUE
	}
	this.Out.WriteByte('\n')
	if err != nil {
		fmt.Fprintf(this.Out, "(warning) %s\n", err.Error())
	}
	box.Print(ctx, toDisplay(comp.List), this.Out)
	this.RepaintAll()
	return readline.CONTINUE
}

func CommonPrefix(list []string) string {
	if len(list) < 1 {
		return ""
	}
	common := list[0]
	var cr, fr *strings.Reader
	for _, f := range list[1:] {
		if cr != nil {
			cr.Reset(common)
		} else {
			cr = strings.NewReader(common)
		}
		if fr != nil {
			fr.Reset(f)
		} else {
			fr = strings.NewReader(f)
		}
		i := 0
		var buffer strings.Builder
		for {
			ch, _, err := cr.ReadRune()
			if err != nil {
				break
			}
			fh, _, err := fr.ReadRune()
			if err != nil || unicode.ToUpper(ch) != unicode.ToUpper(fh) {
				break
			}
			buffer.WriteRune(ch)
			i++
		}
		common = buffer.String()
	}
	return common
}

func endWithRoot(path string) bool {
	return len(path) >= 1 && os.IsPathSeparator(path[len(path)-1])
}

func KeyFuncCompletion(ctx context.Context, this *readline.Buffer) readline.Result {
	comp, default_delimiter, err := listUpComplete(ctx, this)
	if err != nil {
		fmt.Fprintf(this.Out, "\n%s\n", err)
		this.RepaintAll()
		return readline.CONTINUE
	}
	if comp.List == nil || len(comp.List) <= 0 {
		return readline.CONTINUE
	}

	complete_list := toComplete(comp.List)
	commonStr := CommonPrefix(complete_list)
	quotechar := byte(0)
	if i := strings.IndexAny(comp.Word, readline.Delimiters); i >= 0 {
		quotechar = comp.Word[i]
	} else {
		for _, node := range complete_list {
			if strings.ContainsAny(node, " &!") {
				quotechar = byte(default_delimiter)
				break
			}
		}
	}
	if quotechar != 0 {
		var buffer strings.Builder
		buffer.Grow(len(commonStr) + 3)
		if len(commonStr) >= 2 && commonStr[0] == '~' && (os.IsPathSeparator(commonStr[1]) || unicode.IsLetter(rune(commonStr[1]))) {
			buffer.WriteString(commonStr[:1])
			buffer.WriteByte(quotechar)
			buffer.WriteString(commonStr[1:])
		} else {
			buffer.WriteByte(quotechar)
			buffer.WriteString(commonStr)
		}
		if len(comp.List) == 1 && !endWithRoot(comp.List[0].String()) {
			buffer.WriteByte(quotechar)
		}
		commonStr = buffer.String()
	}
	if len(comp.List) == 1 && !endWithRoot(commonStr) && !strings.HasSuffix(commonStr, `%`) {
		commonStr += " "
	}
	if comp.RawWord == commonStr {
		this.Out.WriteByte('\n')
		if err != nil {
			fmt.Fprintf(this.Out, "(warning) %s\n", err.Error())
		}
		box.Print(nil, toDisplay(comp.List), this.Out)
		this.RepaintAll()
		return readline.CONTINUE
	}
	this.ReplaceAndRepaint(comp.Pos, commonStr)
	return readline.CONTINUE
}
