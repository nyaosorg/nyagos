package history

import (
	"bytes"
	"strings"
	"testing"
)

type history_t struct {
	List []string
}

func (this *history_t) Len() int {
	return len(this.List)
}

func (this *history_t) At(n int) string {
	return this.List[n]
}

func (this *history_t) Push(line string) {
	this.List = append(this.List, line)
}

func (this *history_t) Replace(line string) {
	this.List[len(this.List)-1] = line
}

func TestReplace(t *testing.T) {
	testdata := &history_t{
		List: []string{
			"aaa0 aaa1 aaa2",
			"bbb0 bbb1 bbb2",
			"ccc0 ccc1 ccc2",
		},
	}

	if testdata.Len() != 3 {
		t.Fail()
	}

	if testdata.At(1) != "bbb0 bbb1 bbb2" {
		t.Fail()
	}

	testdata.Push("xxxxx")
	if testdata.At(testdata.Len()-1) != "xxxxx" {
		t.Fail()
	}
	if testdata.Len() != 4 {
		t.Fail()
	}

	testdata.Replace("yyyyy")
	if testdata.At(testdata.Len()-1) != "yyyyy" {
		t.Fail()
	}
	if testdata.Len() != 4 {
		t.Fail()
	}
}

func TestInsertHistory(t *testing.T) {
	var buffer bytes.Buffer

	InsertHistory(&buffer, strings.NewReader("^"), "aaa bbb ccc")
	if buffer.String() != "bbb" {
		t.Fail()
		return
	}

	buffer.Reset()
	InsertHistory(&buffer, strings.NewReader("$"), "aaa bbb ccc ddd")
	if buffer.String() != "ddd" {
		t.Fail()
		return
	}

	buffer.Reset()
	InsertHistory(&buffer, strings.NewReader(":1"), `aaa "b bb" ccc ddd`)
	if buffer.String() != `"b bb"` {
		t.Fail()
		return
	}
}

func TestLoadFromReader(t *testing.T) {
	source := `aaaa
aaaa
bbbb
bbbb
cccc`
	hisObj := &history_t{[]string{}}
	LoadFromReader(strings.NewReader(source), hisObj)
	if hisObj.Len() != 3 || hisObj.At(0) != "aaaa" ||
		hisObj.At(1) != "bbbb" || hisObj.At(2) != "cccc" {

		t.Fail()
	}
}
