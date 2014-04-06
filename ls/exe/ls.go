package main

import "os"
import ".."

func main() {
	ls.Main(os.Args[1:], os.Stdout)
}
