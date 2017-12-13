package completion

import (
	"testing"
)

func TestSplitLikeShell(t *testing.T) {
	field := SplitLikeShell(`1 2 "3  4 5" 6 7`)
	if len(field) < 5 {
		t.Fatal("len error")
		return
	}
	if field[2] != `"3  4 5"` {
		t.Fatal("quote err:" + field[3])
		return
	}
	if field[3] != "6" {
		t.Fatal("6 error:" + field[4])
		return
	}
}
