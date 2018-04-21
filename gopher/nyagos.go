package main

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
	"github.com/zetamatta/nyagos/shell"
)

type luaKeyT struct{}

var luaKey luaKeyT

type ScriptEngineForOptionImpl struct{}

func (this *ScriptEngineForOptionImpl) SetArg(args []string) {}

func (this *ScriptEngineForOptionImpl) RunFile(ctx context.Context, fname string) ([]byte, error) {
	L, ok := ctx.Value(luaKey).(Lua)
	if !ok {
		return nil, errors.New("Script is not supported.")
	}
	if save := L.Context(); save != nil {
		defer L.SetContext(save)
	}
	L.SetContext(ctx)
	return nil, L.DoFile(fname)
}

func (this *ScriptEngineForOptionImpl) RunString(ctx context.Context, code string) error {
	L, ok := ctx.Value(luaKey).(Lua)
	if !ok {
		return errors.New("Script is not supported.")
	}
	if save := L.Context(); save != nil {
		defer L.SetContext(save)
	}
	L.SetContext(ctx)
	return L.DoString(code)
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
	ctx := context.Background()

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

	langEngine := func(fname string) ([]byte, error) {
		ctxTmp := context.WithValue(ctx, shellKey, sh)
		L.SetContext(ctxTmp)
		return nil, L.DoFile(fname)
	}
	shellEngine := func(fname string) error {
		return sh.Source(ctx, fname)
	}

	script, err := frame.OptionParse(sh, &ScriptEngineForOptionImpl{})
	if err != nil {
		return err
	}

	if !isatty.IsTerminal(os.Stdin.Fd()) || script != nil {
		frame.SilentMode = true
	}

	if !frame.OptionNorc {
		if !frame.SilentMode {
			fmt.Printf("Nihongo Yet Another GOing Shell %s-%s by %s\n",
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
		constream := frame.NewCmdStreamConsole(
			func() (int, error) {
				functions.Prompt(
					[]interface{}{frame.Format2Prompt(os.Getenv("PROMPT"))})
				return 0, nil
			})
		stream1 = constream
		frame.DefaultHistory = constream.History
		ctx = context.WithValue(ctx, history.PackageId, constream.History)
	} else {
		stream1 = shell.NewCmdStreamFile(os.Stdin)
	}
	sh.ForEver(ctx, stream1)
	return nil
}
