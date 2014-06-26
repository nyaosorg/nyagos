package main

import "bufio"
import "fmt"
import "os"
import "path/filepath"

import "github.com/shiena/ansicolor"

import "./exename"
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

	appData := filepath.Join(os.Getenv("APPDATA"), "NYAOS_ORG")
	os.Mkdir(appData, 0777)
	histPath := filepath.Join(appData, "nyagos.history")
	history.Load(histPath)
	defer history.Save(histPath)

	exeFolder := filepath.Dir(exename.Query())
	rcPath := filepath.Join(exeFolder, "nyagos.rc")
	if fd, err := os.Open(rcPath); err == nil {
		defer fd.Close()
		scr := bufio.NewScanner(fd)
		for scr.Scan() {
			interpreter.Interpret(scr.Text(), option.CommandHooks, nil)
		}
	}

	for {
		line, cont := conio.ReadLine(
			func() {
				fmt.Fprint(ansiOut, prompt.Format2Prompt(os.Getenv("PROMPT")))
			})
		if cont == conio.ABORT {
			break
		}
		var isReplaced bool
		line, isReplaced = history.Replace(line)
		if isReplaced {
			fmt.Fprintln(os.Stdout, line)
		}
		history.Push(line)
		whatToDo, err := interpreter.Interpret(line, option.CommandHooks, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if whatToDo == interpreter.SHUTDOWN {
			break
		}
	}
}
