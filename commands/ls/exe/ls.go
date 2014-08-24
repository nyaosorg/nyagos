package main

import "fmt"
import "os"

import ".."

func main() {
	err := ls.Main(os.Args[1:], os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}
