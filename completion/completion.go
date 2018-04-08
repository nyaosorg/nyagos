package completion

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/zetamatta/go-box"

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

func listUpComplete(ctx context.Context, this *readline.Buffer) (*List, rune, error) {
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

	if isTop(rv.Left, indexes) {
		rv.List, err = listUpCommands(ctx, rv.Word[start:])
	} else {
		rv.List, err = listUpFiles(ctx, rv.Word[start:])
	}

	for i := 0; i < len(rv.List); i++ {
		rv.List[i] = Element2{
			rv.Word[:start] + rv.List[i].String(),
			rv.List[i].Display(),
		}
	}
	for _, f := range HookToList {
		rv, err = f(this, rv)
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
	if comp == nil {
		return readline.CONTINUE
	}
	fmt.Fprint(readline.Console, "\n")
	if err != nil {
		fmt.Fprintf(readline.Console, "(warning) %s\n", err.Error())
	}
	box.Print(ctx, toDisplay(comp.List), readline.Console)
	this.RepaintAll()
	return readline.CONTINUE
}

func CommonPrefix(list []string) string {
	if len(list) < 1 {
		return ""
	}
	common := list[0]
	for _, f := range list[1:] {
		cr := strings.NewReader(common)
		fr := strings.NewReader(f)
		i := 0
		var buffer strings.Builder
		for {
			ch, _, cerr := cr.ReadRune()
			fh, _, ferr := fr.ReadRune()
			if cerr != nil || ferr != nil || unicode.ToUpper(ch) != unicode.ToUpper(fh) {
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
	if comp.List == nil || len(comp.List) <= 0 {
		return readline.CONTINUE
	}

	slashToBackSlash := true
	firstFoundSlashPos := strings.IndexRune(comp.Word, '/')
	firstFoundBackSlashPos := strings.IndexRune(comp.Word, os.PathSeparator)
	if UseSlash {
		slashToBackSlash = false
		if firstFoundBackSlashPos >= 0 && (firstFoundSlashPos == -1 || firstFoundBackSlashPos < firstFoundSlashPos) {
			slashToBackSlash = true
		}
	} else {
		if firstFoundSlashPos >= 0 && (firstFoundBackSlashPos == -1 || firstFoundSlashPos < firstFoundBackSlashPos) {
			slashToBackSlash = false
		}
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
		if len(commonStr) >= 2 && commonStr[0] == '~' && os.IsPathSeparator(commonStr[1]) {
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
	if slashToBackSlash {
		commonStr = filepath.FromSlash(commonStr)
	}
	if comp.RawWord == commonStr {
		fmt.Fprint(readline.Console, "\n")
		if err != nil {
			fmt.Fprintf(readline.Console, "(warning) %s\n", err.Error())
		}
		box.Print(nil, toDisplay(comp.List), readline.Console)
		this.RepaintAll()
		return readline.CONTINUE
	}
	this.ReplaceAndRepaint(comp.Pos, commonStr)
	return readline.CONTINUE
}
