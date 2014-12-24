package dos

import (
	"os"
	"strings"
	"testing"
)

func TestFixPathCase(t *testing.T) {
	path1, err1 := os.Getwd()
	if err1 != nil {
		t.Errorf("os.Getwd(): %v", err1)
	}
	path1 = strings.ToUpper(path1)
	path2, err2 := CorrectPathCase(path1)
	if err2 != nil {
		t.Errorf("CorrectPathCase: %v", err2)
	}
	println(path1, "->", path2)
}
