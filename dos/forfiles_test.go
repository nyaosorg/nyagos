package dos

import (
	"fmt"
	"testing"
)

func TestFindFirst(t *testing.T) {
	ForFiles("*", func(fd *FileInfo) bool {
		fmt.Print(fd.Name())
		if fd.IsDir() {
			fmt.Print("/")
		}
		fmt.Println()
		return true
	})
}
