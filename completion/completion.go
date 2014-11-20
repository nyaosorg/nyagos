package completion

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"../conio"
	"../conio/readline"
	"../dos"
)

func isExecutable(path string) bool {
	return dos.IsExecutableSuffix(filepath.Ext(path))
}

var RxEnvironPattern = regexp.MustCompile("%[^%]+%")
var RxTilda = regexp.MustCompile("^~[/\\\\]")

func listUpFiles(str string) ([]string, error) {
	str = RxEnvironPattern.ReplaceAllStringFunc(str, func(p string) string {
		if len(p) == 2 {
			return "%"
		} else if val := os.Getenv(p[1 : len(p)-1]); val != "" {
			return val
		} else {
			return p
		}
	})
	str = RxTilda.ReplaceAllStringFunc(str, func(p string) string {
		home := os.Getenv("HOME")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		if home != "" {
			return home + "\\"
		} else {
			return p
		}
	})
	str = strings.Replace(strings.Replace(str, "\\", "/", -1), "\"", "", -1)
	var directory string
	str = strings.Replace(str, "\\", "/", -1)
	if strings.HasSuffix(str, "/") {
		directory = path.Join(str, ".")
	} else {
		directory = path.Dir(str)
	}
	cutprefix := 0
	if strings.HasPrefix(directory, "/") {
		wd, _ := os.Getwd()
		directory = wd[0:2] + directory
		cutprefix = 2
	}

	fd, fdErr := os.Open(directory)
	if fdErr != nil {
		return nil, fdErr
	}
	defer fd.Close()
	files, filesErr := fd.Readdir(-1)
	if filesErr != nil {
		return nil, filesErr
	}
	commons := make([]string, 0)
	STR := strings.ToUpper(str)
	for _, node1 := range files {
		name := path.Join(directory, node1.Name())
		if attr, attrErr := dos.GetFileAttributes(name); attrErr == nil && (attr&dos.FILE_ATTRIBUTE_HIDDEN) != 0 {
			continue
		}
		if node1.IsDir() {
			name += "/"
		}
		if cutprefix > 0 {
			name = name[2:]
		}
		nameUpr := strings.ToUpper(name)
		if strings.HasPrefix(nameUpr, STR) {
			commons = append(commons, name)
		}
	}
	return commons, nil
}

func listUpCommands(str string) ([]string, error) {
	listTmp, listErr := listUpFiles(str)
	if listErr != nil {
		return nil, listErr
	}
	list := make([]string, 0)
	for _, fname := range listTmp {
		if strings.HasSuffix(fname, "/") || strings.HasSuffix(fname, "\\") || isExecutable(fname) {
			list = append(list, fname)
		}
	}
	pathEnv := os.Getenv("PATH")
	dirList := strings.Split(pathEnv, ";")
	strUpr := strings.ToUpper(str)
	for _, dir1 := range dirList {
		dirHandle, dirErr := os.Open(dir1)
		if dirErr != nil {
			continue
		}
		defer dirHandle.Close()
		files, filesErr := dirHandle.Readdir(0)
		if filesErr != nil {
			continue
		}
		for _, file1 := range files {
			name1Upr := strings.ToUpper(file1.Name())
			if !strings.HasPrefix(name1Upr, strUpr) {
				continue
			}
			if file1.IsDir() {
				continue
			}
			name := file1.Name()
			if isExecutable(name) {
				list = append(list, path.Base(name))
			}
		}
	}
	// remove dupcalites
	uniq := make([]string, 0)
	lastone := ""
	for _, cur := range list {
		if cur != lastone {
			uniq = append(uniq, cur)
		}
		lastone = cur
	}
	return uniq, nil
}
func KeyFuncCompletionList(this *readline.Buffer) readline.Result {
	str, pos := this.CurrentWord()
	var list []string
	if pos > 0 {
		list, _ = listUpFiles(str)
	} else {
		list, _ = listUpCommands(str)
	}
	if list == nil {
		return readline.CONTINUE
	}
	fmt.Print("\n")
	conio.BoxPrint(list, os.Stdout)
	this.RepaintAll()
	return readline.CONTINUE
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

func KeyFuncCompletion(this *readline.Buffer) readline.Result {
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
		return readline.CONTINUE
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
	return readline.CONTINUE
}
