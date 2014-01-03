package main

import "../conio"
import "fmt"

func main() {
	fmt.Print("conio.ReadLine>")
	result := conio.ReadLine()
	fmt.Println("Result=" + result)
}
