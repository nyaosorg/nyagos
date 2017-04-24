package main

// Example to use readline

import (
	"context"
	"fmt"
	"github.com/zetamatta/nyagos/readline"
	//"../../readline"
)

func main() {
	readline1 := readline.Editor{}
	text, err := readline1.ReadLine(context.Background())
	if err != nil {
		fmt.Printf("ERR=%s\n", err.Error())
	} else {
		fmt.Printf("TEXT=%s\n", text)
	}
}
