package main

import "fmt"
import "./conio"
import "./interpreter"

func main() {
	for {
		fmt.Print("$ ")
		line,cont := conio.ReadLine()
		if line == "exit" || cont == conio.ABORT {
			break
		}
		// fmt.Println(line)
		_, err := interpreter.Interpret(line)
		if err != nil {
			fmt.Println(err)
		}
	}
}
