package commands_test

import (
	"context"
	"os"
	"testing"

	"github.com/zetamatta/nyagos/commands"
)

func TestLsMain(t *testing.T) {
	_, err := commands.Ls(context.Background(), []string{"."}, os.Stdout, os.Stderr, os.Stdout)
	if err != nil {
		t.Fatalf("ls .: %s", err.Error())
	}
	_, err = commands.Ls(context.Background(), []string{"-l", "ls_test.go"}, os.Stdout, os.Stderr, os.Stdout)
	if err != nil {
		t.Fatalf("ls *: %s", err.Error())
	}
}
