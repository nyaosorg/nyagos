package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shiena/ansicolor"

	"./alias"
	"./commands"
	"./completion"
	"./conio"
	"./dos"
	"./history"
	"./interpreter"
	"./lua"
)

var rxAnsiEscCode = regexp.MustCompile("\x1b[^a-zA-Z]*[a-zA-Z]")

var stamp string
var commit string
var version string
var ansiOut io.Writer

func nyagosPrompt(L *lua.Lua) int {
	template, err := L.ToString(-1)
	if err != nil {
		template = "[" + err.Error() + "]"
	}
	text := Format2Prompt(template)
	fmt.Fprint(ansiOut, text)
	text = rxAnsiEscCode.ReplaceAllString(text, "")
	lfPos := strings.LastIndex(text, "\n")
	if lfPos >= 0 {
		text = text[lfPos+1:]
	}
	L.PushInteger(conio.GetStringWidth(text))
	return 1
}

func main() {
	conio.DisableCtrlC()

	historyUp := conio.KeyGoFuncT{history.KeyFuncHistoryUp}
	historyDown := conio.KeyGoFuncT{history.KeyFuncHistoryDown}
	completion := conio.KeyGoFuncT{completion.KeyFuncCompletion}

	if err := conio.BindKeySymbolFunc(conio.K_CTRL_I, "COMPLETE", &completion); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if err := conio.BindKeySymbolFunc(conio.K_UP, "PREVIOUS_HISTORY", &historyUp); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if err := conio.BindKeySymbolFunc(conio.K_DOWN, "NEXT_HISTORY", &historyDown); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if err := conio.BindKeySymbol(conio.K_CTRL_P, "PREVIOUS_HISTORY"); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if err := conio.BindKeySymbol(conio.K_CTRL_N, "NEXT_HISTORY"); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	// ANSI Escape Sequence Support
	ansiOut = ansicolor.NewAnsiColorWriter(os.Stdout)

	commands.Init()
	alias.Init()

	// Lua extension
	L := lua.New()
	L.OpenLibs()
	SetLuaFunctions(L)
	L.GetGlobal("nyagos")
	L.PushString(stamp)
	L.SetField(-2, "stamp")
	L.PushString(commit)
	L.SetField(-2, "commit")
	L.PushString(version)
	L.SetField(-2, "version")
	L.PushGoFunction(nyagosPrompt)
	L.SetField(-2, "prompt")
	L.Pop(1)
	defer L.Close()

	if !optionParse(L) {
		return
	}

	appData := filepath.Join(os.Getenv("APPDATA"), "NYAOS_ORG")
	os.Mkdir(appData, 0777)
	histPath := filepath.Join(appData, "nyagos.history")
	history.Load(histPath)
	history.Save(histPath) // cut over max-line

	exeName, exeNameErr := dos.GetModuleFileName()
	if exeNameErr != nil {
		fmt.Fprintln(os.Stderr, exeNameErr)
	}
	exeFolder := filepath.Dir(exeName)
	nyagos_lua := filepath.Join(exeFolder, "nyagos.lua")
	if _, err := os.Stat(nyagos_lua); err == nil {
		err := L.Source(nyagos_lua)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	for {
		wd, wdErr := os.Getwd()
		if wdErr == nil {
			conio.SetTitle("NYAOS - " + wd)
		} else {
			conio.SetTitle("NYAOS - " + wdErr.Error())
		}
		line, cont := conio.ReadLine(
			func() int {
				L.GetGlobal("nyagos")
				L.GetField(-1, "prompt")
				L.Remove(-2)
				L.PushString(os.Getenv("PROMPT"))
				L.Call(1, 1)
				length, lengthErr := L.ToInteger(-1)
				if lengthErr == nil {
					return length
				} else {
					fmt.Fprintf(os.Stderr,
						"nyagos.prompt: Length invalid: %s\n",
						lengthErr.Error())
					return 0
				}
			})
		if cont == conio.ABORT {
			break
		}

		var isReplaced bool
		line, isReplaced = history.Replace(line)
		if isReplaced {
			fmt.Fprintln(os.Stdout, line)
		}
		if line == "" {
			history.ResetPointer()
			continue
		}
		if line != history.LastHistory() {
			history.Push(line)
			fd, err := os.OpenFile(histPath, os.O_APPEND, 0600)
			if err == nil {
				fmt.Fprintln(fd, line)
				fd.Close()
			}
		} else {
			history.ResetPointer()
		}

		stackPos := L.GetTop()
		L.GetGlobal("nyagos")
		L.GetField(-1, "filter")
		if L.IsFunction(-1) {
			L.PushString(line)
			err := L.Call(1, 1)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				if L.IsString(-1) {
					line, err = L.ToString(-1)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
				}
			}
		}
		L.SetTop(stackPos)

		whatToDo, err := interpreter.New().Interpret(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if whatToDo == interpreter.SHUTDOWN {
			break
		}
	}
}
