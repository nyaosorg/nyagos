package main

import "fmt"
import "os"

import "./completion"
import "./conio"
import "./internalcmd"
import "./interpreter"
import "./history"

func main() {
	conio.KeyMap['\t'] = completion.KeyFuncCompletion
	conio.ZeroMap[conio.K_UP] = history.KeyFuncHistoryUp
	conio.ZeroMap[conio.K_DOWN] = history.KeyFuncHistoryDown
	conio.KeyMap['P'&0x1F] = history.KeyFuncHistoryUp
	conio.KeyMap['N'&0x1F] = history.KeyFuncHistoryDown
	for {
		wd, _ := os.Getwd()
		fmt.Printf("[%s]\n$ ", wd)
		line, cont := conio.ReadLine()
		if cont == conio.ABORT {
			break
		}
		history.Push(line)
		whatToDo, err := interpreter.Interpret(line, internalcmd.Exec)
		if err != nil {
			fmt.Println(err)
		}
		if whatToDo == interpreter.SHUTDOWN {
			break
		}
	}
}
