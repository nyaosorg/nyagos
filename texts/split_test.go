package texts_test

import (
	"testing"

	"github.com/nyaosorg/nyagos/texts"
)

func TestSplitLikeShell(t *testing.T) {
	s := `1 2 "3 \" 4 5" 6 7`
	indexes := texts.SplitLikeShell(s)
	if len(indexes) < 5 {
		t.Fatal("len error")
		return
	}
	field := make([]string, 0, len(indexes))
	for _, p := range indexes {
		field = append(field, s[p[0]:p[1]])
	}
	if field[2] != `"3 \" 4 5"` {
		t.Fatal("quote err:" + field[2])
		return
	}
	if field[3] != "6" {
		t.Fatal("6 error:" + field[4])
		return
	}
}
