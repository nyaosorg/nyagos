package main

import "../conio"
import "fmt"

func main() {
	fmt.Print("conio.ReadLine>")
	result, rc := conio.ReadLine()
	fmt.Println("Result=", result)
	fmt.Println("Rc=", rc)
}
