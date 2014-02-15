package main

import "fmt"
import "os/exec"
import "./conio"
import "./interpreter"

func hook(cmd *exec.Cmd, IsBackground bool) (interpreter.WhatToDoAfterCmd, error) {
	if cmd.Args[0] == "exit" {
		return interpreter.SHUTDOWN, nil
	} else {
		return interpreter.THROUGH, nil
	}
}

func main() {
	for {
		fmt.Print("$ ")
		line, cont := conio.ReadLine()
		if cont == conio.ABORT {
			break
		}
		whatToDo, err := interpreter.Interpret(line, hook)
		if err != nil {
			fmt.Println(err)
		}
		if whatToDo == interpreter.SHUTDOWN {
			break
		}
	}
}
