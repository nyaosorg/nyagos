package main

import "fmt"
import "os"
import "path/filepath"
import "regexp"
import "strings"

import "github.com/shiena/ansicolor"

import "./alias"
import "./commands"
import "./completion"
import "./conio"
import "./conio/readline"
import "./dos"
import "./history"
import "./interpreter"
import "./lua"

var rxAnsiEscCode = regexp.MustCompile("\x1b[^a-zA-Z]*[a-zA-Z]")

var stamp string
var commit string
var version string

func main() {
	conio.DisableCtrlC()

	historyUp := readline.KeyGoFuncT{history.KeyFuncHistoryUp}
	historyDown := readline.KeyGoFuncT{history.KeyFuncHistoryDown}
	completion := readline.KeyGoFuncT{completion.KeyFuncCompletion}

	if err := readline.BindKeySymbolFunc(readline.K_CTRL_I, "COMPLETE", &completion); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if err := readline.BindKeySymbolFunc(readline.K_UP, "PREVIOUS_HISTORY", &historyUp); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if err := readline.BindKeySymbolFunc(readline.K_DOWN, "NEXT_HISTORY", &historyDown); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if err := readline.BindKeySymbol(readline.K_CTRL_P, "PREVIOUS_HISTORY"); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if err := readline.BindKeySymbol(readline.K_CTRL_N, "NEXT_HISTORY"); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	// ANSI Escape Sequence Support
	ansiOut := ansicolor.NewAnsiColorWriter(os.Stdout)

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
	L.Pop(1)
	defer L.Close()

	// Parameter Parsing
	argc := 0
	cont := OptionParse(L, func() (string, bool) {
		argc++
		if argc < len(os.Args) {
			return os.Args[argc], true
		} else {
			return "", false
		}
	})
	if !cont {
		return
	}

	appData := filepath.Join(os.Getenv("APPDATA"), "NYAOS_ORG")
	os.Mkdir(appData, 0777)
	histPath := filepath.Join(appData, "nyagos.history")
	history.Load(histPath)
	defer history.Save(histPath)

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
		line, cont := readline.ReadLine(
			func() int {
				text := Format2Prompt(os.Getenv("PROMPT"))
				fmt.Fprint(ansiOut, text)
				text = rxAnsiEscCode.ReplaceAllString(text, "")
				lfPos := strings.LastIndex(text, "\n")
				if lfPos >= 0 {
					text = text[lfPos+1:]
				}
				return conio.GetStringWidth(text)
			})
		if cont == readline.ABORT {
			break
		}

		var isReplaced bool
		line, isReplaced = history.Replace(line)
		if isReplaced {
			fmt.Fprintln(os.Stdout, line)
		}
		if line == "" {
			continue
		}
		if line != history.LastHistory() {
			history.Push(line)
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
