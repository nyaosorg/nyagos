package main

import "fmt"
import "os"
import "io"

import "github.com/shiena/ansicolor"

import "./completion"
import "./conio"
import "./history"
import "./interpreter"
import "./option"
import "./prompt"

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
	argc := 0
	option.Parse(func() (string, bool) {
		argc++
		if argc < len(os.Args) {
			return os.Args[argc], true
		} else {
			return "", false
		}
	})

	for {
		line, cont := conio.ReadLine(
			func() {
				io.WriteString(ansiOut,
					prompt.Format2Prompt(os.Getenv("PROMPT")))
			})
		if cont == conio.ABORT {
			break
		}
		var isReplaced bool
		line, isReplaced = history.Replace(line)
		if isReplaced {
			os.Stdout.WriteString(line)
			os.Stdout.WriteString("\n")
		}
		history.Push(line)
		whatToDo, err := interpreter.Interpret(line, option.CommandHooks, nil)
		if err != nil {
			fmt.Println(err)
		}
		if whatToDo == interpreter.SHUTDOWN {
			break
		}
	}
}
