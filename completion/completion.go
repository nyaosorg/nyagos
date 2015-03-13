package completion

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"

	"../conio"
)

type CompletionList struct {
	AllLine string
	List    []string
	RawWord string // have quotation
	Word    string
	Pos     int
}

func listUpComplete(this *conio.Buffer) (*CompletionList, error) {
	var err error
	rv := CompletionList{}

	// environment completion.
	rv.AllLine = this.String()
	rv.List, rv.Pos, err = listUpEnv(rv.AllLine)
	if len(rv.List) > 0 && rv.Pos >= 0 && err == nil {
		rv.RawWord = rv.AllLine[rv.Pos:]
		rv.Word = rv.RawWord
		return &rv, nil
	}

	// filename or commandname completion
	rv.RawWord, rv.Pos = this.CurrentWord()
	rv.Word = strings.Replace(rv.RawWord, "\"", "", -1)
	if rv.Pos > 0 {
		rv.List, err = listUpFiles(rv.Word)
	} else {
		rv.List, err = listUpCommands(rv.Word)
	}
	return &rv, err
}

func KeyFuncCompletionList(this *conio.Buffer) conio.Result {
	comp, err := listUpComplete(this)
	if err != nil {
		return conio.CONTINUE
	}
	fmt.Print("\n")
	conio.BoxPrint(comp.List, os.Stdout)
	this.RepaintAll()
	return conio.CONTINUE
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

func KeyFuncCompletion(this *conio.Buffer) conio.Result {
	comp, err := listUpComplete(this)
	if err != nil || comp.List == nil || len(comp.List) <= 0 {
		return conio.CONTINUE
	}

	slashToBackSlash := true
	firstFoundSlashPos := strings.IndexRune(comp.Word, '/')
	firstFoundBackSlashPos := strings.IndexRune(comp.Word, '\\')
	if firstFoundSlashPos >= 0 && firstFoundBackSlashPos >= 0 && firstFoundSlashPos < firstFoundBackSlashPos {
		slashToBackSlash = false
	}

	commonStr := CommonPrefix(comp.List)
	needQuote := strings.ContainsRune(comp.Word, '"')
	if !needQuote {
		for _, node := range comp.List {
			if strings.ContainsAny(node, " &") {
				needQuote = true
				break
			}
		}
	}
	if needQuote {
		buffer := make([]byte, 0, 100)
		buffer = append(buffer, '"')
		buffer = append(buffer, commonStr...)
		if len(comp.List) == 1 && !endWithRoot(comp.List[0]) {
			buffer = append(buffer, '"')
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
		return KeyFuncCompletionList(this)
	}
	this.ReplaceAndRepaint(comp.Pos, commonStr)
	return conio.CONTINUE
}
