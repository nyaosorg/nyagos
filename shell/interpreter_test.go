package shell

import (
	"context"
	"fmt"
	"testing"
)

func TestInterpret(t *testing.T) {
	ctx := context.Background()
	_, err := New().Interpret(ctx, "ls | cat -n > hogehoge")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("no error")
	}
}

func TestMain(t *testing.T) {
	in := []string{`ahahaha ihhihi`, `foo bar`, `"foo bar"`}
	out := makeCmdline(in, in)
	tst := `"ahahaha ihhihi" "foo bar" "\"foo bar\""`
	if out != tst {
		t.Fatalf(`Fail "%s" != "%s"`, out, tst)
	}
}
