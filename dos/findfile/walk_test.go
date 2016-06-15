package findfile

import (
	"fmt"
	"testing"
)

func TestWalk(t *testing.T) {
	Walk("*", func(fd *FileInfo) bool {
		fmt.Print(fd.Name())
		if fd.IsDir() {
			fmt.Print("/")
		}
		fmt.Println()
		return true
	})
}
