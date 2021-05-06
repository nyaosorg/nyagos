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
