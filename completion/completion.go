package completion

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"unicode"

	"../conio"
	"../readline"
	/* dbg "github.com/zetamatta/goutputdebugstring" */)

type CompletionList struct {
	AllLine string
	List    []string
	RawWord string // have quotation
	Word    string
	Pos     int
}

func listUpComplete(this *readline.Buffer) (*CompletionList, rune, error) {
	var err error
	rv := new(CompletionList)

	// environment completion.
	rv.AllLine = this.String()
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

	if rv.Pos > 0 {
		rv.List, err = listUpFiles(rv.Word[start:])
	} else {
		rv.List, err = listUpCommands(rv.Word[start:])
	}

	for i := 0; i < len(rv.List); i++ {
		rv.List[i] = rv.Word[:start] + rv.List[i]
	}
	rv, err = luaHook(this, rv)
	return rv, default_delimiter, err
}

func KeyFuncCompletionList(ctx context.Context, this *readline.Buffer) readline.Result {
	comp, _, err := listUpComplete(this)
	if comp == nil {
		return readline.CONTINUE
	}
	fmt.Print("\n")
	os.Stdout.Sync()
	if err != nil {
		fmt.Printf("(warning) %s\n", err.Error())
		os.Stderr.Sync()
	}
	conio.BoxPrint(ctx, comp.List, os.Stdout)
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
		var buffer bytes.Buffer
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
	return strings.HasSuffix(path, "\\") || strings.HasSuffix(path, "/")
}

func KeyFuncCompletion(this *readline.Buffer) readline.Result {
	comp, default_delimiter, err := listUpComplete(this)
	if comp.List == nil || len(comp.List) <= 0 {
		return readline.CONTINUE
	}

	slashToBackSlash := true
	firstFoundSlashPos := strings.IndexRune(comp.Word, '/')
	firstFoundBackSlashPos := strings.IndexRune(comp.Word, '\\')
	if firstFoundSlashPos >= 0 && (firstFoundBackSlashPos == -1 || firstFoundSlashPos < firstFoundBackSlashPos) {
		slashToBackSlash = false
	}

	commonStr := CommonPrefix(comp.List)
	quotechar := byte(0)
	if i := strings.IndexAny(comp.Word, readline.Delimiters); i >= 0 {
		quotechar = comp.Word[i]
	} else {
		for _, node := range comp.List {
			if strings.ContainsAny(node, " &") {
				quotechar = byte(default_delimiter)
				break
			}
		}
	}
	if quotechar != 0 {
		buffer := make([]byte, 0, 100)
		buffer = append(buffer, quotechar)
		buffer = append(buffer, commonStr...)
		if len(comp.List) == 1 && !endWithRoot(comp.List[0]) {
			buffer = append(buffer, quotechar)
		}
		commonStr = string(buffer)
	}
	if len(comp.List) == 1 && !endWithRoot(commonStr) {
		commonStr += " "
	}
	if slashToBackSlash {
		commonStr = strings.Replace(commonStr, "/", "\\", -1)
	}
	if comp.RawWord == commonStr {
		fmt.Print("\n")
		if err != nil {
			fmt.Printf("(warning) %s\n", err.Error())
		}
		conio.BoxPrint(nil, comp.List, os.Stdout)
		this.RepaintAll()
		return readline.CONTINUE
	}
	this.ReplaceAndRepaint(comp.Pos, commonStr)
	return readline.CONTINUE
}
