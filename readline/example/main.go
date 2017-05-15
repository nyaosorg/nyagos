package main

// Example to use readline

import (
	"context"
	"fmt"
	// "github.com/zetamatta/nyagos/readline"
	"../../readline"
)

func main() {
	readline1 := readline.Editor{
		Default: "AHAHA",
		Cursor:  3,
	}

	enter_status := 0
	readline.BindKeyClosure(readline.K_CTRL_P,
		 func(r *readline.Buffer) readline.Result {
			enter_status = -1
			return readline.ENTER
		})

	readline.BindKeyClosure(readline.K_CTRL_N,
		func(r *readline.Buffer) readline.Result {
			enter_status = +1
			return readline.ENTER
		})

	text, err := readline1.ReadLine(context.Background())
	fmt.Printf("ENTER_STATUS=%d\n", enter_status)
	if err != nil {
		fmt.Printf("ERR=%s\n", err.Error())
	} else {
		fmt.Printf("TEXT=%s\n", text)
	}
}
