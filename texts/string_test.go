package texts_test

import (
	"fmt"
	"testing"

	"github.com/zetamatta/nyagos/texts"

)

func TestFirstWord(t *testing.T) {
	fmt.Println("--- FirstWord ---")

	if value := texts.FirstWord("aaaa bbbb cccc"); value != "aaaa" {
		t.Error("Case-1: failed")
	} else {
		fmt.Println(value)
	}
	if value := texts.FirstWord("\"12 34\" bbb"); value != "\"12 34\"" {
		t.Error("Case-2: failed")
	} else {
		fmt.Println(value)
	}
}

func TestSplitLikeShellString(t *testing.T) {
	fmt.Println("*** Test SplitLikeShellString() ***")
	values := texts.SplitLikeShellString("\"a b\" bbb ccc \"1 2 3\" 'a  b' c")
	for key, val := range values {
		fmt.Printf("[%d]=\"%s\"\n", key, val)
	}
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
