package main

import (
	conio ".."
	"fmt"
)

func test1(message, want_text string, want_result bool) {
	text, result := conio.ReadLinePromptStr(message)
	fmt.Printf("  -> text='%v' result='%v'\n", text, result)
	if text != want_text {
		fmt.Printf("got '%v' want '%v'", text, want_text)
	}
	if result != want_result {
		fmt.Printf("got '%v' want '%v'", result, want_result)
	}

}

func main() {
	fmt.Println("--- Test for ReadLinePromptStr ---")
	test1("Please Type 'foo'+Enter>", "foo", true)
	test1("Please Type Anything+CTRL-C>", "", true)
	test1("Please Type CTRL-D only>", "", false)
}
