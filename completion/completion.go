package completion

import "bytes"
import "fmt"
import "os"
import "path"
import "strings"

import "../box"
import "../conio"

func listUpFiles(str string) ([]string, error) {
	directory := path.Dir(str)
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
	STR := strings.Replace(strings.ToUpper(str), "\\", "/", -1)
	for _, node1 := range files {
		name := path.Join(directory, node1.Name())
		NAME := strings.Replace(strings.ToUpper(name), "\\", "/", -1)
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

func KeyFuncCompletion(this *conio.ReadLineBuffer) conio.KeyFuncResult {
	str, pos := this.CurrentWord()
	list, err := listUpFiles(str)
	if err == nil {
		if len(list) == 1 {
			str = list[0]
			if strings.ContainsRune(str, ' ') {
				var buffer bytes.Buffer
				buffer.WriteRune('"')
				buffer.WriteString(str)
				buffer.WriteRune('"')
				str = buffer.String()
			}
			this.ReplaceAndRepaint(pos, str)
			for _, ch := range str {
				conio.PutRep(ch, 1)
				this.Cursor++
			}
		} else {
			return KeyFuncCompletionList(this)
		}
	}
	return conio.CONTINUE
}
