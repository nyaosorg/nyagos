package ls

import (
	"context"
	"os"

	"testing"
)

func TestMain(t *testing.T) {
	err := Main(context.Background(), []string{"."}, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("ls .: %s", err.Error())
	}
	err = Main(context.Background(), []string{"-l", "ls_test.go"}, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("ls *: %s", err.Error())
	}
}
