package commands_test

import (
	"context"
	"io"
	"testing"

	"github.com/nyaosorg/nyagos/internal/commands"
)

func TestLsMain(t *testing.T) {
	_, err := commands.Ls(context.Background(), []string{"."}, io.Discard, io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("ls .: %s", err.Error())
	}
	_, err = commands.Ls(context.Background(), []string{"-l", "ls_test.go"}, io.Discard, io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("ls *: %s", err.Error())
	}
}
