package main

import "fmt"
import "os"

import "./builtincmd"
import "./completion"
import "./conio"
import "./history"
import "./interpreter"
import "./prompt"

func main() {
	conio.KeyMap['\t'] = completion.KeyFuncCompletion
	conio.ZeroMap[conio.K_UP] = history.KeyFuncHistoryUp
	conio.ZeroMap[conio.K_DOWN] = history.KeyFuncHistoryDown
	conio.KeyMap['P'&0x1F] = history.KeyFuncHistoryUp
	conio.KeyMap['N'&0x1F] = history.KeyFuncHistoryDown
	for {
		fmt.Printf("%s", prompt.Prompt(os.Getenv("PROMPT")))
		line, cont := conio.ReadLine()
		if cont == conio.ABORT {
			break
		}
		history.Push(line)
		whatToDo, err := interpreter.Interpret(line, builtincmd.Exec, nil)
		if err != nil {
			fmt.Println(err)
		}
		if whatToDo == interpreter.SHUTDOWN {
			break
		}
	}
}
