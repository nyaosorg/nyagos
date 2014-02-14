package interpreter

import "testing"
import "fmt"

func TestInterpret(t *testing.T) {
	_, err := Interpret("ls.exe | cat.exe -n > hogehoge")
	fmt.Println(err)
}
