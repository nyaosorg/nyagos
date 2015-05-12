package conio

import (
	"fmt"
	"testing"
)

func TestQuotedFirstWord(t *testing.T) {
	fmt.Println("--- splitq ---")

	if value := QuotedFirstWord("aaaa bbbb cccc"); value != "aaaa" {
		t.Error("Case-1: failed")
	} else {
		fmt.Println(value)
	}
	if value := QuotedFirstWord("\"12 34\" bbb"); value != "\"12 34\"" {
		t.Error("Case-2: failed")
	} else {
		fmt.Println(value)
	}
}

func TestSplitQ(t *testing.T) {
	fmt.Println("*** Test SplitQ() ***")
	values := SplitQ("\"a b\" bbb ccc \"1 2 3\" 'a  b' c")
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
