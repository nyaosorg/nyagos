package dbg

import (
	"os"
	"testing"
)

func TestTypeName(t *testing.T) {
	_, err := os.Stat("NOT_FOUND_FILE")
	s := TypeName(err)
	if s != "*os.PathError" {
		t.Fatalf("TypeName == '%s'",s)
		return
	}
}
