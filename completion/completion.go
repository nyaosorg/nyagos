package completion

import "bytes"
import "fmt"
import "os"
import "path"
import "strings"
import "unicode"

import "../box"
import "../conio"

func listUpFiles(str string) ([]string, error) {
	str = strings.Replace(strings.Replace(str, "\\", "/", -1), "\"", "", -1)
	var directory string
	str = strings.Replace(str, "\\", "/", -1)
	if strings.HasSuffix(str, "/") {
		directory = path.Join(str, ".")
	} else {
		directory = path.Dir(str)
	}
	cutprefix := 0
	if strings.HasPrefix(directory,"/") {
		wd,_ := os.Getwd()
		directory = wd[0:2] + directory
		cutprefix = 2
	}

	fd, err := os.Open(directory)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	files, err2 := fd.Readdir(-1)
	if err2 != nil {
		return nil, err2
	}
	commons := make([]string, 0)
	STR := strings.ToUpper(str)
	for _, node1 := range files {
		name := path.Join(directory, node1.Name())
		if node1.IsDir() {
			name += "/"
		}
		if cutprefix > 0 {
			name = name[2:]
		}
		NAME := strings.ToUpper(name)
		if strings.HasPrefix(NAME, STR) {
			commons = append(commons, name)
		}
	}
	return commons, nil
}

func KeyFuncCompletionList(this *conio.ReadLineBuffer) conio.KeyFuncResult {
	str, _ := this.CurrentWord()
	list, _ := listUpFiles(str)
	if list == nil {
		return conio.CONTINUE
	}
	fmt.Print("\n")
	box.Print(list, 80, os.Stdout)
	this.RepaintAll("$ ")
	return conio.CONTINUE
}

func getCommmon(list []string) string {
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

func KeyFuncCompletion(this *conio.ReadLineBuffer) conio.KeyFuncResult {
	str, wordStart := this.CurrentWord()

	slashToBackSlash := true
	firstFoundSlashPos := strings.IndexRune(str, '/')
	firstFoundBackSlashPos := strings.IndexRune(str, '\\')
	if firstFoundSlashPos >= 0 && firstFoundBackSlashPos >= 0 && firstFoundSlashPos < firstFoundBackSlashPos {
		slashToBackSlash = false
	}

	list, err := listUpFiles(str)
	if err != nil || len(list) <= 0 {
		return conio.CONTINUE
	}
	commonStr := getCommmon(list)
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
		var buffer bytes.Buffer
		buffer.WriteRune('"')
		buffer.WriteString(commonStr)
		if len(list) <= 1 {
			buffer.WriteRune('"')
		}
		commonStr = buffer.String()
	}
	if len(list) == 1 && ! strings.HasSuffix(commonStr,"/") && ! strings.HasSuffix(commonStr,"/\"") {
		commonStr += " "
	}
	if slashToBackSlash {
		commonStr = strings.Replace(commonStr, "/", "\\", -1)
	}
	if compareWithoutQuote(str, commonStr) {
		return KeyFuncCompletionList(this)
	}
	this.ReplaceAndRepaint(wordStart, commonStr)
	conio.PutRep('\a', 1)
	return conio.CONTINUE
}
