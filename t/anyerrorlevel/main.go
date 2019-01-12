package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) <= 1 {
		os.Exit(0)
	}
	val, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[1], err.Error())
		os.Exit(0)
	}
	fmt.Printf("os.Exit(%d)\n",val)
	os.Exit(val)
}
