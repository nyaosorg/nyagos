package completion

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"../conio"
	"../dos"
)

func isExecutable(path string) bool {
	return dos.IsExecutableSuffix(filepath.Ext(path))
}

func KeyFuncCompletionList(this *conio.Buffer) conio.Result {
	str, pos := this.CurrentWord()
	var list []string
	if pos > 0 {
		list, _ = listUpFiles(str)
	} else {
		list, _ = listUpCommands(str)
	}
	if list == nil {
		return conio.CONTINUE
	}
	fmt.Print("\n")
	conio.BoxPrint(list, os.Stdout)
	this.RepaintAll()
	return conio.CONTINUE
}

func GetCommon(list []string) string {
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

func compareWithoutQuote(this string, that string) bool {
	return strings.Replace(this, "\"", "", -1) == strings.Replace(that, "\"", "", -1)
}

func completeXXX(this *conio.Buffer, listUpper func(string) ([]string, int, error)) bool {
	allstring := this.String()
	matches, lastPercentPos, listUpErr := listUpper(allstring)
	if listUpErr != nil {
		fmt.Fprintln(os.Stderr, listUpErr.Error())
		return false
	}
	if matches == nil || len(matches) <= 0 {
		return false
	}
	if len(matches) == 1 { // one matches.
		this.ReplaceAndRepaint(lastPercentPos, matches[0])
		return true
	}
	// more than one match.
	commonStr := GetCommon(matches)
	originStr := allstring[lastPercentPos:]
	if commonStr != originStr {
		this.ReplaceAndRepaint(lastPercentPos, commonStr)
	} else {
		// no difference -> listing.
		fmt.Println()
		conio.BoxPrint(matches, os.Stdout)
		this.RepaintAll()
	}
	return true
}

func KeyFuncCompletion(this *conio.Buffer) conio.Result {
	if completeXXX(this, listUpEnv) {
		return conio.CONTINUE
	}
	str, wordStart := this.CurrentWord()

	slashToBackSlash := true
	firstFoundSlashPos := strings.IndexRune(str, '/')
	firstFoundBackSlashPos := strings.IndexRune(str, '\\')
	if firstFoundSlashPos >= 0 && firstFoundBackSlashPos >= 0 && firstFoundSlashPos < firstFoundBackSlashPos {
		slashToBackSlash = false
	}

	var list []string
	var err error
	if wordStart > 0 {
		list, err = listUpFiles(str)
	} else {
		list, err = listUpCommands(str)
	}
	if err != nil || list == nil || len(list) <= 0 {
		return conio.CONTINUE
	}
	commonStr := GetCommon(list)
	needQuote := strings.ContainsRune(str, '"')
	if !needQuote {
		for _, node := range list {
			if strings.ContainsRune(node, ' ') {
				needQuote = true
				break
			}
		}
	}
	if needQuote {
		buffer := make([]byte, 0, 100)
		buffer = append(buffer, '"')
		buffer = append(buffer, commonStr...)
		if len(list) <= 1 {
			buffer = append(buffer, '"')
		}
		commonStr = string(buffer)
	}
	if len(list) == 1 && !strings.HasSuffix(commonStr, "/") && !strings.HasSuffix(commonStr, "/\"") {
		commonStr += " "
	}
	if slashToBackSlash {
		commonStr = strings.Replace(commonStr, "/", "\\", -1)
	}
	if compareWithoutQuote(str, commonStr) {
		return KeyFuncCompletionList(this)
	}
	this.ReplaceAndRepaint(wordStart, commonStr)
	return conio.CONTINUE
}
