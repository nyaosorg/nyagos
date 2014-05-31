package main

import "fmt"
import "os"
import "io"
import "os/exec"

import "github.com/shiena/ansicolor"

import "./alias"
import "./builtincmd"
import "./completion"
import "./conio"
import "./history"
import "./interpreter"
import "./prompt"

func commandHooks(cmd *exec.Cmd, IsBackground bool) (interpreter.WhatToDoAfterCmd, error) {
	status, _ := alias.Hook(cmd, IsBackground)
	if status != interpreter.THROUGH {
		return status, nil
	}
	return builtincmd.Exec(cmd, IsBackground)
}

func main() {
	conio.KeyMap['\t'] = completion.KeyFuncCompletion
	conio.ZeroMap[conio.K_UP] = history.KeyFuncHistoryUp
	conio.ZeroMap[conio.K_DOWN] = history.KeyFuncHistoryDown
	conio.KeyMap['P'&0x1F] = history.KeyFuncHistoryUp
	conio.KeyMap['N'&0x1F] = history.KeyFuncHistoryDown
	ansiOut := ansicolor.NewAnsiColorWriter(os.Stdout)
	for {
		io.WriteString(ansiOut, prompt.Format2Prompt(os.Getenv("PROMPT")))
		line, cont := conio.ReadLine()
		if cont == conio.ABORT {
			break
		}
		history.Push(line)
		whatToDo, err := interpreter.Interpret(line, commandHooks, nil)
		if err != nil {
			fmt.Println(err)
		}
		if whatToDo == interpreter.SHUTDOWN {
			break
		}
	}
}
