package mains

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/mattn/go-isatty"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/frame"
	"github.com/zetamatta/nyagos/history"
	"github.com/zetamatta/nyagos/lua"
	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/shell"
)

var rxAnsiEscCode = regexp.MustCompile("\x1b[^a-zA-Z]*[a-zA-Z]")

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
	text := frame.Format2Prompt(template)

	fmt.Fprint(readline.Console, text)

	text = rxAnsiEscCode.ReplaceAllString(text, "")
	lfPos := strings.LastIndex(text, "\n")
	if lfPos >= 0 {
		text = text[lfPos+1:]
	}
	L.PushInteger(lua.Integer(readline.GetStringWidth(text)))
	return 1
}

var prompt_hook lua.Object = lua.TGoFunction(nyagosPrompt)

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

var luaFilter lua.Object = lua.TNil{}

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

type luaWrapper struct {
	lua.Lua
}

func (this *luaWrapper) Clone() (shell.CloneCloser, error) {
	L := this.Lua
	newL, err := NewLua()
	if err != nil {
		return nil, err
	}
	err = L.CloneTo(newL)
	if err != nil {
		return nil, err
	}
	return &luaWrapper{newL}, nil
}

func (this *luaWrapper) Close() error {
	return this.Lua.Close()
}

type MainStream struct {
	shell.Stream
	L lua.Lua
}

func (this *MainStream) ReadLine(ctx context.Context) (context.Context, string, error) {
	ctx = context.WithValue(ctx, lua.PackageId, this.L)
	ctx = context.WithValue(ctx, history.PackageId, default_history)
	ctx, line, err := this.Stream.ReadLine(ctx)
	if err != nil {
		return ctx, "", err
	}
	return ctx, doLuaFilter(this.L, line), nil
}

type ScriptEngineForOptionImpl struct {
	L  lua.Lua
	Sh *shell.Shell
}

func (this *ScriptEngineForOptionImpl) SetArg(args []string) {
	setLuaArg(this.L, args)
}

func (this *ScriptEngineForOptionImpl) RunFile(fname string) ([]byte, error) {
	return runLua(this.Sh, this.L, fname)
}

func (this *ScriptEngineForOptionImpl) RunString(code string) error {
	if err := this.L.LoadString(code); err != nil {
		return err
	}
	this.L.Call(0, 0)
	return nil
}

func optionParseLua(sh *shell.Shell, L lua.Lua) (func() error, error) {
	e := &ScriptEngineForOptionImpl{Sh: sh, L: L}
	return frame.OptionParse(sh, e)
}

func Main() error {
	// for issue #155 & #158
	lua.NG_UPVALUE_NAME["prompter"] = struct{}{}

	// Lua extension
	L, err := NewLua()
	if err != nil {
		return err
	}
	defer L.Close()

	sh := shell.New()
	sh.SetTag(&luaWrapper{L})
	defer sh.Close()

	langEngine := func(fname string) ([]byte, error) {
		return runLua(sh, L, fname)
	}
	shellEngine := func(fname string) error {
		fd, err := os.Open(fname)
		if err != nil {
			return err
		}
		stream1 := frame.NewCmdStreamFile(fd)
		_, err = sh.Loop(stream1)
		fd.Close()
		if err == io.EOF {
			return nil
		} else {
			return err
		}
	}

	script, err := optionParseLua(sh, L)
	if err != nil {
		return err
	}

	if !isatty.IsTerminal(os.Stdin.Fd()) || script != nil {
		frame.SilentMode = true
	}

	if !frame.OptionNorc {
		if !frame.SilentMode {
			fmt.Printf("Nihongo Yet Another GOing Shell %s-%s by %s & Lua 5.3\n",
				frame.VersionOrStamp(),
				runtime.GOARCH,
				runtime.Version())
			fmt.Println("(c) 2014-2018 NYAOS.ORG <http://www.nyaos.org>")
		}
		if err := frame.LoadScripts(shellEngine, langEngine); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	if script != nil {
		if err := script(); err != nil {
			if err != io.EOF {
				return err
			} else {
				return nil
			}
		}
	}

	backupHistory := default_history
	defer func() {
		default_history = backupHistory
	}()

	var stream1 shell.Stream
	if isatty.IsTerminal(os.Stdin.Fd()) {
		constream := frame.NewCmdStreamConsole(
			func() (int, error) { return printPrompt(L) })
		stream1 = constream
		default_history = constream.History
	} else {
		stream1 = frame.NewCmdStreamFile(os.Stdin)
	}

	for {
		_, err = sh.Loop(&MainStream{stream1, L})
		if err == io.EOF {
			return err
		}
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
