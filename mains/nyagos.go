// +build !vanilla

package mains

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/yuin/gopher-lua"

	"github.com/zetamatta/go-windows-consoleicon"

	"github.com/zetamatta/nyagos/alias"
	"github.com/zetamatta/nyagos/commands"
	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/frame"
	"github.com/zetamatta/nyagos/functions"
	"github.com/zetamatta/nyagos/shell"
)

type luaKeyT struct{}

var luaKey luaKeyT

type ScriptEngineForOptionImpl struct {
	L  Lua
	Sh *shell.Shell
}

func (this *ScriptEngineForOptionImpl) SetArg(args []string) {
	if L := this.L; L != nil {
		table := L.NewTable()
		for i, arg1 := range args {
			L.SetTable(table, lua.LNumber(i), lua.LString(arg1))
		}
		L.SetGlobal("arg", table)
	}
}

func DoFileExceptForAtmarkLines(L *lua.LState, fname string) error {
	fd, err := os.Open(fname)
	if err != nil {
		return err
	}
	reader, writer := io.Pipe()
	go func() {
		scan := bufio.NewScanner(fd)
		for scan.Scan() {
			line := scan.Text()
			if len(line) > 0 && line[0] == '@' {
				line = ""
			}
			fmt.Fprintln(writer, line)
		}
		writer.Close()
		fd.Close()
	}()
	f, err := L.Load(reader, fname)
	reader.Close()
	if err != nil {
		return err
	}
	L.Push(f)
	return L.PCall(0, 0, nil)
}

func (this *ScriptEngineForOptionImpl) RunFile(ctx context.Context, fname string) ([]byte, error) {
	L, ok := ctx.Value(luaKey).(Lua)
	if !ok {
		return nil, errors.New("Script is not supported.")
	}
	defer setContext(L, getContext(L))
	setContext(L, ctx)
	return nil, luaRedirect(ctx, os.Stdin, os.Stdout, os.Stderr, L, func() error {
		return DoFileExceptForAtmarkLines(L, fname)
	})
}

func (this *ScriptEngineForOptionImpl) RunString(ctx context.Context, code string) error {
	L, ok := ctx.Value(luaKey).(Lua)
	if !ok {
		return errors.New("Script is not supported.")
	}
	ctx = context.WithValue(ctx, shellKey, this.Sh)
	defer setContext(L, getContext(L))
	setContext(L, ctx)
	return luaRedirect(ctx, os.Stdin, os.Stdout, os.Stderr, L, func() error {
		return L.DoString(code)
	})
}

type luaWrapper struct {
	Lua
}

func (this *luaWrapper) Clone(ctx context.Context) (context.Context, shell.CloneCloser, error) {
	newL, err := Clone(this.Lua)
	if err != nil {
		return nil, nil, err
	}
	ctx = context.WithValue(ctx, luaKey, newL)
	return ctx, &luaWrapper{newL}, nil
}

func (this *luaWrapper) Close() error {
	this.Lua.Close()
	return nil
}

func Main() error {
	disableColors := colorable.EnableColorsStdout(nil)
	defer disableColors()

	if clean, err := consoleicon.SetFromExe(); err == nil {
		defer clean(true)
	}
	ctx := context.Background()

	completion.HookToList = append(completion.HookToList, luaHookForComplete)

	L, err := NewLua()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		ctx = context.WithValue(ctx, luaKey, L)
		defer L.Close()
	}

	sh := shell.New()
	if L != nil {
		sh.SetTag(&luaWrapper{L})
	}
	defer sh.Close()
	sh.Console = colorable.NewColorableStdout()
	ctx = context.WithValue(ctx, shellKey, sh)

	langEngine := func(fname string) ([]byte, error) {
		ctxTmp := context.WithValue(ctx, shellKey, sh)
		defer setContext(L, getContext(L))
		setContext(L, ctxTmp)
		return nil, L.DoFile(fname)
	}
	shellEngine := func(fname string) error {
		return sh.Source(ctx, fname)
	}

	alias.LineFilter = func(ctx context.Context, line string) string {
		if L, ok := ctx.Value(luaKey).(Lua); ok {
			return luaLineFilter(ctx, L, line)
		} else {
			return line
		}
	}

	script, err := frame.OptionParse(ctx, sh, &ScriptEngineForOptionImpl{L: L, Sh: sh})
	if err != nil {
		return err
	}

	if !isatty.IsTerminal(os.Stdin.Fd()) || script != nil {
		frame.SilentMode = true
	}

	if !frame.OptionNorc {
		if !frame.SilentMode {
			frame.Title()
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
	if !commands.ReadStdinAsFile && isatty.IsTerminal(os.Stdin.Fd()) {
		constream := frame.NewCmdStreamConsole(
			func() (int, error) {
				if L != nil {
					return printPrompt(ctx, sh, L)
				} else {
					functions.Prompt(
						&functions.Param{
							Args: []interface{}{frame.Format2Prompt(os.Getenv("PROMPT"))},
							In:   os.Stdin,
							Out:  os.Stdout,
							Err:  os.Stderr,
							Term: colorable.NewColorableStdout(),
						})
					return 0, nil
				}
			})
		stream1 = constream
		frame.DefaultHistory = constream.History
		sh.History = constream.History
		ctx = context.WithValue(ctx, shellKey, sh)
	} else {
		stream1 = shell.NewCmdStreamFile(os.Stdin)
	}
	if L != nil {
		sh.ForEver(ctx, &luaFilterStream{Stream: stream1, L: L})
	} else {
		sh.ForEver(ctx, stream1)
	}
	return nil
}
