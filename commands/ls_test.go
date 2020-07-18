package commands

import (
	"context"
	"os"

	"testing"
)

func TestLsMain(t *testing.T) {
	err := lsMain(context.Background(), []string{"."}, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("ls .: %s", err.Error())
	}
	err = lsMain(context.Background(), []string{"-l", "ls_test.go"}, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("ls *: %s", err.Error())
	}
}
