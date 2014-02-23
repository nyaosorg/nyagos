package main

import "fmt"
import "os"

import "./completion"
import "./conio"
import "./internalcmd"
import "./interpreter"

func main() {
	conio.KeyMap['\t'] = completion.KeyFuncCompletion
	for {
		wd, _ := os.Getwd()
		fmt.Printf("[%s]\n$ ", wd)
		line, cont := conio.ReadLine()
		if cont == conio.ABORT {
			break
		}
		whatToDo, err := interpreter.Interpret(line, internalcmd.Exec)
		if err != nil {
			fmt.Println(err)
		}
		if whatToDo == interpreter.SHUTDOWN {
			break
		}
	}
}
