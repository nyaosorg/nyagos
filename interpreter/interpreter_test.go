package interpreter

import (
	"fmt"
	"testing"
)

func TestInterpret(t *testing.T) {
	_, err := New().Interpret("ls.exe | cat.exe -n > hogehoge")
	fmt.Println(err)
}
