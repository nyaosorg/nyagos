package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/zetamatta/go-getch"

	"../alias"
	"../commands"
	"../completion"
	"../conio"
	"../dos"
	"../history"
	"../interpreter"
	"../lua"
	"../readline"
)

var rxAnsiEscCode = regexp.MustCompile("\x1b[^a-zA-Z]*[a-zA-Z]")

var stamp string
var commit string
var version string
var ansiOut io.Writer

func nyagosPrompt(L lua.Lua) int {
	title, title_err := L.ToString(2)
	if title_err == nil && title != "" {
		conio.SetTitle(title)
	} else if wd, wdErr := os.Getwd(); wdErr == nil {
		conio.SetTitle("NYAGOS - " + wd)
	} else {
		conio.SetTitle("NYAGOS")
	}
	template, err := L.ToString(1)
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
	L.PushInteger(lua.Integer(conio.GetStringWidth(text)))
	return 1
}

var prompt_hook lua.Pushable = lua.TGoFunction{nyagosPrompt}

func printPrompt(this *readline.LineEditor) (int, error) {
	L := NewNyagosLua()
	defer L.Close()
	L.Push(prompt_hook)

	if !L.IsFunction(-1) {
		L.Pop(1)
		return 0, nil
	}
	L.PushString(os.Getenv("PROMPT"))
	if err := L.Call(1, 1); err != nil {
		return 0, err
	}
	length, lengthErr := L.ToInteger(-1)
	L.Pop(1)
	if lengthErr == nil {
		return length, nil
	} else {
		return 0, fmt.Errorf("nyagos.prompt: return-value(length) is invalid: %s", lengthErr.Error())
	}
}

func when_panic() {
	err := recover()
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, "************ Panic Occured. ***********")
	fmt.Fprintln(os.Stderr, err)
	debug.PrintStack()
	fmt.Fprintln(os.Stderr, "*** Please copy these error message ***")
	fmt.Fprintln(os.Stderr, "*** And hit ENTER key to quit.      ***")
	var dummy [1]byte
	os.Stdin.Read(dummy[:])
}

var luaFilter lua.Pushable = lua.TNil{}

func itprCloneHook(this *interpreter.Interpreter) error {
	LL := NewNyagosLua()
	this.Tag = LL
	this.OnClone = itprCloneHook
	this.Closers = append(this.Closers, LL)
	return nil
}

func NewCmdStreamFile(f *os.File) func() (string, error) {
	breader := bufio.NewReader(os.Stdin)
	return func() (string, error) {
		line, err := breader.ReadString('\n')
		if err != nil {
			return "", err
		}
		line = strings.TrimRight(line, "\r\n")
		return line, nil
	}
}

func NewCmdStreamConsole(it *interpreter.Interpreter) func() (string, error) {
	readline.DefaultEditor.Prompt = printPrompt
	readline.DefaultEditor.Tag = it
	return readline.DefaultEditor.ReadLine
}

var optionK = flag.String("k", "", "like `cmd /k`")
var optionC = flag.String("c", "", "like `cmd /c`")
var optionF = flag.String("f", "", "run lua script")
var optionE = flag.String("e", "", "run inline-lua-code")

func main() {
	defer when_panic()

	flag.Parse()

	interpreter.SetHook(func(it *interpreter.Interpreter) (int, bool, error) {
		rc, done, err := commands.Exec(&it.Cmd)
		return rc, done, err
	})
	completion.AppendCommandLister(commands.AllNames)
	completion.AppendCommandLister(alias.AllNames)

	dos.CoInitializeEx(0, dos.COINIT_MULTITHREADED)
	defer dos.CoUninitialize()

	getch.DisableCtrlC()

	completion := readline.KeyGoFuncT{F: completion.KeyFuncCompletion}

	if err := readline.BindKeySymbolFunc(readline.K_CTRL_I, "COMPLETE", &completion); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	// ANSI Escape Sequence Support
	ansiOut = colorable.NewColorableStdout()

	commands.Init()
	alias.Init()

	// Lua extension
	L := NewNyagosLua()
	defer L.Close()

	if !isatty.IsTerminal(os.Stdin.Fd()) || *optionC != "" || *optionF != "" || *optionE != "" {
		silentmode = true
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

	it := interpreter.New()
	it.Tag = L
	it.OnClone = itprCloneHook
	it.Closers = append(it.Closers, L)

	if !optionParse(it, L) {
		return
	}

	var command_stream func() (string, error)
	if isatty.IsTerminal(os.Stdin.Fd()) {
		command_stream = NewCmdStreamConsole(it)
	} else {
		command_stream = NewCmdStreamFile(os.Stdin)
	}

	for {
		history_count := readline.DefaultEditor.HistoryLen()

		line, err := command_stream()
		if err != nil {
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, err.Error())
			}
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
		if readline.DefaultEditor.HistoryLen() > history_count {
			fd, err := os.OpenFile(histPath, os.O_APPEND, 0600)
			if err != nil && os.IsNotExist(err) {
				// print("create ", histPath, "\n")
				fd, err = os.Create(histPath)
			}
			if err == nil {
				fmt.Fprintln(fd, line)
				fd.Close()
			} else {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		} else {
			readline.DefaultEditor.HistoryResetPointer()
		}

		stackPos := L.GetTop()
		L.Push(luaFilter)
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

		_, err = it.Interpret(line)
		if err != nil {
			if err == io.EOF {
				break
			}
			if err1, ok := err.(interpreter.AlreadyReportedError); ok {
				if err1.Err == io.EOF {
					break
				}
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}
