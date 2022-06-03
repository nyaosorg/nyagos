package texts_test

import (
	"testing"

	"github.com/nyaosorg/nyagos/internal/texts"
)

func TestFirstWord(t *testing.T) {
	if value := texts.FirstWord("aaaa bbbb cccc"); value != "aaaa" {
		t.Error("Case-1: failed")
	}
	if value := texts.FirstWord("\"12 34\" bbb"); value != "\"12 34\"" {
		t.Error("Case-2: failed")
	}
}

func TestSplitLikeShellString(t *testing.T) {
	values := texts.SplitLikeShellString("\"a b\" bbb ccc \"1 2 3\" 'a  b' c")
	if len(values) != 6 {
		t.Error("Case-1: failed")
	}
	if values[0] != "\"a b\"" {
		t.Error("Case-2: failed")
	}
	if values[3] != "\"1 2 3\"" {
		t.Error("Case-3: failed")
	}
}
