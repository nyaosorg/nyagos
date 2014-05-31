package main

import "fmt"
import "os"
import "io"
import "os/exec"
import "strings"

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
	// KeyBind += completion Module
	conio.KeyMap['\t'] = completion.KeyFuncCompletion
	conio.ZeroMap[conio.K_UP] = history.KeyFuncHistoryUp
	conio.ZeroMap[conio.K_DOWN] = history.KeyFuncHistoryDown
	conio.KeyMap['P'&0x1F] = history.KeyFuncHistoryUp
	conio.KeyMap['N'&0x1F] = history.KeyFuncHistoryDown

	// ANSI Escape Sequence Support
	ansiOut := ansicolor.NewAnsiColorWriter(os.Stdout)

	// Parameter Parsing
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i][0] == '-' {
			for _, o := range os.Args[i][1:] {
				switch o {
				case 'a':
					i++
					if i < len(os.Args) {
						equation := os.Args[i]
						equationArray := strings.SplitN(equation, "=", 2)
						if len(equationArray) >= 2 {
							alias.Table[strings.ToLower(equationArray[0])] =
								equationArray[1]
						} else {
							delete(alias.Table, strings.ToLower(
								equationArray[0]))
						}
					}
				case 'c', 'k':
					i++
					if i < len(os.Args) {
						interpreter.Interpret(os.Args[i], commandHooks, nil)
					}
					if o == 'c' {
						return
					}
				}
			}
		}
	}

	for {
		line, cont := conio.ReadLine(
			func() {
				io.WriteString(ansiOut,
					prompt.Format2Prompt(os.Getenv("PROMPT")))
			})
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
