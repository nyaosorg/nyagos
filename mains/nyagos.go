package mains

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/mattn/go-isatty"

	"github.com/zetamatta/nyagos/frame"
	"github.com/zetamatta/nyagos/functions"
	"github.com/zetamatta/nyagos/history"
	"github.com/zetamatta/nyagos/lua"
	"github.com/zetamatta/nyagos/shell"
)

var noLuaEngineErr = errors.New("no lua engine")

var prompt_hook lua.Object = lua.TGoFunction(lua2cmd(functions.Prompt))

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

type luaWrapper struct {
	lua.Lua
}

func (this *luaWrapper) Clone(ctx context.Context) (context.Context, shell.CloneCloser, error) {
	L := this.Lua
	newL, err := NewLua()
	if err != nil {
		return ctx, nil, err
	}
	err = L.CloneTo(newL)
	if err != nil {
		return ctx, nil, err
	}
	ctx = context.WithValue(ctx, lua.PackageId, newL)
	return ctx, &luaWrapper{newL}, nil
}

func (this *luaWrapper) Close() error {
	return this.Lua.Close()
}

type LuaFilterStream struct {
	shell.Stream
	L lua.Lua
}

func (this *LuaFilterStream) ReadLine(ctx context.Context) (context.Context, string, error) {
	ctx, line, err := this.Stream.ReadLine(ctx)
	if err != nil {
		return ctx, "", err
	}

	L := this.L

	stackPos := L.GetTop()
	defer L.SetTop(stackPos)

	L.Push(luaFilter)
	if !L.IsFunction(-1) {
		return ctx, line, nil
	}
	L.PushString(line)
	err = L.Call(1, 1)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ctx, line, nil
	}
	if !L.IsString(-1) {
		return ctx, line, nil
	}
	newLine, err := L.ToString(-1)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ctx, line, nil
	}
	return ctx, newLine, nil
}

type ScriptEngineForOptionImpl struct {
	L  lua.Lua
	Sh *shell.Shell
}

func (this *ScriptEngineForOptionImpl) SetArg(args []string) {
	if this.L != 0 {
		setLuaArg(this.L, args)
	}
}

func (this *ScriptEngineForOptionImpl) RunFile(fname string) ([]byte, error) {
	if this.L != 0 {
		return runLua(this.Sh, this.L, fname)
	} else {
		return nil, noLuaEngineErr
	}
}

func (this *ScriptEngineForOptionImpl) RunString(code string) error {
	if this.L == 0 {
		return noLuaEngineErr
	}
	if err := this.L.LoadString(code); err != nil {
		return err
	}
	this.L.Call(0, 0)
	return nil
}

func optionParseLua(sh *shell.Shell, L lua.Lua) (func(context.Context) error, error) {
	e := &ScriptEngineForOptionImpl{Sh: sh, L: L}
	return frame.OptionParse(sh, e)
}

func Main() error {
	// for issue #155 & #158
	lua.NG_UPVALUE_NAME["prompter"] = struct{}{}

	// Lua extension
	L, err := NewLua()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		defer L.Close()
	}

	sh := shell.New()
	if L != 0 {
		sh.SetTag(&luaWrapper{L})
	}
	defer sh.Close()

	ctx := context.Background()

	langEngine := func(fname string) ([]byte, error) {
		if L != 0 {
			return runLua(sh, L, fname)
		} else {
			return nil, nil
		}
	}
	shellEngine := func(fname string) error {
		fd, err := os.Open(fname)
		if err != nil {
			return err
		}
		stream1 := frame.NewCmdStreamFile(fd)
		_, err = sh.Loop(ctx, stream1)
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
			fmt.Printf("Nihongo Yet Another GOing Shell %s-%s by %s",
				frame.VersionOrStamp(),
				runtime.GOARCH,
				runtime.Version())
			if L != 0 {
				fmt.Print(" & Lua 5.3")
			}
			fmt.Println("\n(c) 2014-2018 NYAOS.ORG <http://www.nyaos.org>")
		}
		if err := frame.LoadScripts(shellEngine, langEngine); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	if script != nil {
		if err := script(ctx); err != nil {
			if err != io.EOF {
				return err
			} else {
				return nil
			}
		}
	}

	var stream1 shell.Stream
	if isatty.IsTerminal(os.Stdin.Fd()) {
		constream := frame.NewCmdStreamConsole(func() (int, error) {
			if L != 0 {
				return printPrompt(L)
			} else {
				functions.Prompt(
					[]interface{}{frame.Format2Prompt(os.Getenv("PROMPT"))})
				return 0, nil
			}
		})
		stream1 = constream
		frame.DefaultHistory = constream.History
		ctx = context.WithValue(ctx, history.PackageId, constream.History)
	} else {
		stream1 = frame.NewCmdStreamFile(os.Stdin)
	}
	if L != 0 {
		ctx = context.WithValue(ctx, lua.PackageId, L)
		sh.ForEver(ctx, &LuaFilterStream{stream1, L})
	} else {
		sh.ForEver(ctx, stream1)
	}
	return nil
}
