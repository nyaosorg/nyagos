package main

import (
	"bufio"
	"context"
	"fmt"

	"github.com/mattn/go-colorable"

	"github.com/zetamatta/nyagos/readline"
)

func main() {
	editor := readline.Editor{
		Default: "InitialValue",
		Cursor:  3,
		Writer:  bufio.NewWriter(colorable.NewColorableStdout()),
	}
	text, err := editor.ReadLine(context.Background())

	if err != nil {
		fmt.Printf("ERR=%s\n", err.Error())
	} else {
		fmt.Printf("TEXT=%s\n", text)
	}
}
