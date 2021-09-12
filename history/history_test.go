package history_test

import (
	"strings"
	"testing"
	"time"

	"github.com/nyaosorg/nyagos/history"
)

type historyT struct {
	List []string
}

func (h *historyT) Len() int {
	return len(h.List)
}

func (h *historyT) At(n int) string {
	return h.List[n]
}

func (h *historyT) Push(line string) {
	h.List = append(h.List, line)
}

func TestReplace(t *testing.T) {
	testdata := &historyT{
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
}

func newLine(text string) *history.Line {
	return &history.Line{
		Text:  text,
		Dir:   ".",
		Stamp: time.Now(),
		Pid:   0,
	}
}

func TestExpandMacro(t *testing.T) {
	var buffer strings.Builder

	history.ExpandMacro(&buffer, strings.NewReader("^"), newLine("aaa bbb ccc"))
	if buffer.String() != "bbb" {
		t.Fail()
		return
	}

	buffer.Reset()
	history.ExpandMacro(&buffer, strings.NewReader("$"), newLine("aaa bbb ccc ddd"))
	if buffer.String() != "ddd" {
		t.Fail()
		return
	}

	buffer.Reset()
	history.ExpandMacro(&buffer, strings.NewReader(":1"), newLine(`aaa "b bb" ccc ddd`))
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
	hisObj := &history.Container{}
	hisObj.LoadViaReader(strings.NewReader(source))
	if hisObj.Len() != 5 ||
		hisObj.At(0) != "aaaa" ||
		hisObj.At(1) != "aaaa" ||
		hisObj.At(2) != "bbbb" ||
		hisObj.At(3) != "bbbb" ||
		hisObj.At(4) != "cccc" {

		var buffer strings.Builder
		for i := 0; i < hisObj.Len(); i++ {
			if i > 0 {
				buffer.WriteString(",")
			}
			buffer.WriteString(hisObj.At(i))
		}
		t.Fatalf("Failed: [%s]", buffer.String())
	}
}

// func TestSaveToWriter(t *testing.T) {
// 	hisObj := &Container{
// 		[]Line{
// 			{ Text:"aaaa" },
// 			{ Text:"bbbb" },
// 			{ Text:"aaaa" },
// 			{ Text:"dddd" },
// 			{ Text:"eeee" },
// 		},
// 	}
// 	var buffer strings.Builder
// 	hisObj.SaveViaWriter(&buffer)
// 	if buffer.String() != "bbbb\naaaa\ndddd\neeee\n" {
// 		println(buffer.String())
// 		t.Fail()
// 	}
// }
