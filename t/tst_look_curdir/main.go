package main

import (
	"fmt"
	"os"
)

func main() {
	path, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		fmt.Println(path)
	}
}
