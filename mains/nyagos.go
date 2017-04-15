package mains

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/zetamatta/go-getch"

	"../alias"
	"../commands"
	"../completion"
	"../dos"
	"../history"
	"../lua"
	"../readline"
	"../shell"
)

var rxAnsiEscCode = regexp.MustCompile("\x1b[^a-zA-Z]*[a-zA-Z]")

var Stamp string
var Commit string
var Version string

func setTitle(s string) {
	fmt.Fprintf(readline.Console, "\x1B]0;%s\007", s)
}

func nyagosPrompt(L lua.Lua) int {
	title, title_err := L.ToString(2)
	if title_err == nil && title != "" {
		setTitle(title)
	} else if wd, wdErr := os.Getwd(); wdErr == nil {
		if flag, _ := dos.IsElevated(); flag {
			setTitle("(Admin) - " + wd)
		} else {
			setTitle("NYAGOS - " + wd)
		}
	} else {
		if flag, _ := dos.IsElevated(); flag {
			setTitle("(Admin)")
		} else {
			setTitle("NYAGOS")
		}
	}
	template, err := L.ToString(1)
	if err != nil {
		template = "[" + err.Error() + "]"
	}
	text := Format2Prompt(template)

	fmt.Fprint(readline.Console, text)

	text = rxAnsiEscCode.ReplaceAllString(text, "")
	lfPos := strings.LastIndex(text, "\n")
	if lfPos >= 0 {
		text = text[lfPos+1:]
	}
	L.PushInteger(lua.Integer(readline.GetStringWidth(text)))
	return 1
}

var prompt_hook lua.Pushable = lua.TGoFunction(nyagosPrompt)

func printPrompt(L lua.Lua) (int, error) {
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

var luaFilter lua.Pushable = lua.TNil{}

func itprCloneHook(this *shell.Cmd) error {
	LL, err := NewNyagosLua()
	if err != nil {
		return err
	}
	this.Tag = LL
	this.OnClone = itprCloneHook
	this.Closers = append(this.Closers, LL)
	return nil
}

var optionK = flag.String("k", "", "like `cmd /k`")
var optionC = flag.String("c", "", "like `cmd /c`")
var optionF = flag.String("f", "", "run lua script")
var optionE = flag.String("e", "", "run inline-lua-code")

var appdatapath_ string

func AppDataDir() string {
	if appdatapath_ == "" {
		appdatapath_ = filepath.Join(os.Getenv("APPDATA"), "NYAOS_ORG")
		os.Mkdir(appdatapath_, 0777)
	}
	return appdatapath_
}

var default_history *history.Container

func doLuaFilter(L lua.Lua, line string) string {
	stackPos := L.GetTop()
	defer L.SetTop(stackPos)

	L.Push(luaFilter)
	if !L.IsFunction(-1) {
		return line
	}
	L.PushString(line)
	err := L.Call(1, 1)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return line
	}
	if !L.IsString(-1) {
		return line
	}
	line2, err2 := L.ToString(-1)
	if err2 != nil {
		fmt.Fprintln(os.Stderr, err2)
		return line
	}
	return line2
}

func Main() error {
	// for issue #155 & #158
	lua.NG_UPVALUE_NAME["prompter"] = struct{}{}

	flag.Parse()

	shell.SetHook(func(ctx context.Context, it *shell.Cmd) (int, bool, error) {
		rc, done, err := commands.Exec(ctx, it)
		return rc, done, err
	})
	completion.AppendCommandLister(commands.AllNames)
	completion.AppendCommandLister(alias.AllNames)
	completion.HookToList = append(completion.HookToList, luaHookForComplete)

	dos.CoInitializeEx(0, dos.COINIT_MULTITHREADED)
	defer dos.CoUninitialize()

	getch.DisableCtrlC()

	commands.Init()
	alias.Init()

	// Lua extension
	L, err := NewNyagosLua()
	if err != nil {
		return err
	}
	defer L.Close()

	if !isatty.IsTerminal(os.Stdin.Fd()) || *optionC != "" || *optionF != "" || *optionE != "" {
		silentmode = true
	}

	it := shell.New()
	it.Tag = L
	it.OnClone = itprCloneHook
	it.Closers = append(it.Closers, L)

	if err := loadScripts(it, L); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	if !optionParse(it, L) {
		return nil
	}

	var command_reader func(context.Context) (string, error)
	if isatty.IsTerminal(os.Stdin.Fd()) {
		command_reader, default_history = NewCmdStreamConsole(
			func() (int, error) { return printPrompt(L) })
	} else {
		command_reader = NewCmdStreamFile(os.Stdin)
	}

	sigint := make(chan os.Signal, 1)
	defer close(sigint)
	quit := make(chan struct{}, 1)
	defer close(quit)

	for {
		ctx, cancel := context.WithCancel(context.Background())
		ctx = context.WithValue(ctx, "lua", L)
		ctx = context.WithValue(ctx, "history", default_history)

		line, err := command_reader(ctx)

		if err != nil {
			cancel()
			return err
		}
		line = doLuaFilter(L, line)
		signal.Notify(sigint, os.Interrupt)

		go func(sigint_ chan os.Signal, quit_ chan struct{}, cancel_ func()) {
			for {
				select {
				case <-sigint_:
					cancel_()
					<-quit
					return
				case <-quit:
					cancel_()
					return
				}
			}
		}(sigint, quit, cancel)
		_, err = it.InterpretContext(ctx, line)
		signal.Stop(sigint)
		quit <- struct{}{}

		if err != nil {
			if err == io.EOF {
				break
			}
			if err1, ok := err.(shell.AlreadyReportedError); ok {
				if err1.Err == io.EOF {
					break
				}
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
	return nil
}
