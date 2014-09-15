package main

import "fmt"
import "os"
import "os/signal"
import "path/filepath"
import "regexp"
import "strings"

import "github.com/shiena/ansicolor"

import "./alias"
import "./commands"
import "./completion"
import "./conio"
import "./dos"
import "./history"
import "./interpreter"
import "./lua"

func signalOff() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for _ = range ch {

		}
	}()
}

var rxAnsiEscCode = regexp.MustCompile("\x1b[^a-zA-Z]*[a-zA-Z]")

var stamp string
var commit string

func main() {
	signalOff()

	// KeyBind += completion Module
	historyUp := conio.KeyGoFuncT{history.KeyFuncHistoryUp}
	historyDown := conio.KeyGoFuncT{history.KeyFuncHistoryDown}
	completion := conio.KeyGoFuncT{completion.KeyFuncCompletion}

	conio.KeyMap['\t'] = &completion
	conio.ZeroMap[conio.K_UP] = &historyUp
	conio.ZeroMap[conio.K_DOWN] = &historyDown
	conio.KeyMap['P'&0x1F] = &historyUp
	conio.KeyMap['N'&0x1F] = &historyDown

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
	L.Pop(1)
	defer L.Close()

	// Parameter Parsing
	argc := 0
	OptionParse(func() (string, bool) {
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
				text := Format2Prompt(os.Getenv("PROMPT"))
				fmt.Fprint(ansiOut, text)
				text = rxAnsiEscCode.ReplaceAllString(text, "")
				lfPos := strings.LastIndex(text, "\n")
				if lfPos >= 0 {
					text = text[lfPos+1:]
				}
				return conio.GetStringWidth(text)
			})
		if cont == conio.ABORT {
			break
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
					line = L.ToString(-1)
				}
			}
		}
		L.SetTop(stackPos)

		var isReplaced bool
		line, isReplaced = history.Replace(line)
		if isReplaced {
			fmt.Fprintln(os.Stdout, line)
		}
		history.Push(line)
		whatToDo, err := interpreter.Interpret(line, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if whatToDo == interpreter.SHUTDOWN {
			break
		}
	}
}
