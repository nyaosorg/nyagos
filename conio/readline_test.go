package conio

import (
	"fmt"
	"testing"
)

func test1(t *testing.T, message, want_text string, want_result bool) {
	text, result := ReadLinePromptStr(message)
	fmt.Printf("  -> text='%v' result='%v'\n", text, result)
	if text != want_text {
		t.Errorf("got '%v' want '%v'", text, want_text)
	}
	if result != want_result {
		t.Errorf("got '%v' want '%v'", result, want_result)
	}

}

func TestReadLinePromptStr(t *testing.T) {
	fmt.Println("--- Test for ReadLinePromptStr ---")
	test1(t, "Please Type 'foo'+Enter>", "foo", true)
	test1(t, "Please Type Anything+CTRL-C>", "", true)
	test1(t, "Please Type CTRL-D only>", "", false)
}
