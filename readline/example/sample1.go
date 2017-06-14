package main

import (
	"context"
	"fmt"
	"github.com/zetamatta/nyagos/readline"
)

func main() {
	editor := readline.Editor{
		Default: "InitialValue",
		Cursor:  3,
	}
	text, err := editor.ReadLine(context.Background())

	if err != nil {
		fmt.Printf("ERR=%s\n", err.Error())
	} else {
		fmt.Printf("TEXT=%s\n", text)
	}
}
