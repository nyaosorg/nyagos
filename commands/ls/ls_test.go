package ls

import (
	"os"

	"testing"
)

func TestMain(t *testing.T) {
	err := Main([]string{"."}, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("ls .: %s", err.Error())
	}
	err = Main([]string{"-l","ls_test.go"}, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("ls *: %s", err.Error())
	}
}
