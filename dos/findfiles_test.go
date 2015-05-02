package dos

import (
	"fmt"
	"testing"
)

func TestFindFirst(t *testing.T) {
	fd, err := FindFirst("*")
	if err != nil {
		return
	}
	defer fd.Close()
	for ; err == nil; err = fd.FindNext() {
		fmt.Print(fd.Name())
		if fd.IsDir() {
			fmt.Print("/")
		}
		fmt.Println()
	}
}
