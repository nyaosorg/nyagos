package completion

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"

	"../conio"
	"../lua"
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
	L, Lok := this.Session.Tag.(*lua.Lua)
	if !Lok {
		panic("conio.LineEditor.Tag is not *lua.Lua")
	}
	L.GetGlobal("nyagos")
	L.GetField(-1, "completion_hook")
	L.Remove(-2) // remove nyagos-table
	if L.IsFunction(-1) {
		L.NewTable()
		L.PushString(rv.RawWord)
		L.SetField(-2, "rawword")
		L.Push(rv.Pos + 1)
		L.SetField(-2, "pos")
		L.PushString(rv.AllLine)
		L.SetField(-2, "text")
		L.PushString(rv.Word)
		L.SetField(-2, "word")
		L.NewTable()
		for key, val := range rv.List {
			L.Push(key)
			L.PushString(val)
			L.SetTable(-3)
		}
		L.SetField(-2, "list")
		if err := L.Call(1, 1); err != nil {
			fmt.Println(err)
		}
		if L.IsTable(-1) {
			list := make([]string, 0, len(rv.List)+32)
			wordUpr := strings.ToUpper(rv.Word)
			for i := 1; true; i++ {
				L.Push(i)
				L.GetTable(-2)
				str, strErr := L.ToString(-1)
				L.Pop(1)
				if strErr != nil || str == "" {
					break
				}
				strUpr := strings.ToUpper(str)
				if strings.HasPrefix(strUpr, wordUpr) {
					list = append(list, str)
				}
			}
			if len(list) > 0 {
				rv.List = list
			}
		}
	}
	L.Pop(1) // remove something not function or result-table
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
		fmt.Print("\n")
		conio.BoxPrint(comp.List, os.Stdout)
		this.RepaintAll()
		return conio.CONTINUE
	}
	this.ReplaceAndRepaint(comp.Pos, commonStr)
	return conio.CONTINUE
}
