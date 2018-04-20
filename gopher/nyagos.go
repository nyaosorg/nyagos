package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/mattn/go-isatty"
	"github.com/yuin/gopher-lua"

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
	L, ok := ctx.Value(luaKey).(*lua.LState)
	if !ok {
		return nil, errors.New("Script is not supported.")
	}
	return nil, L.DoFile(fname)
}

func (this *ScriptEngineForOptionImpl) RunString(ctx context.Context, code string) error {
	L, ok := ctx.Value(luaKey).(*lua.LState)
	if !ok {
		return errors.New("Script is not supported.")
	}
	return L.DoString(code)
}

type luaWrapper struct {
	*lua.LState
}

func (this *luaWrapper) Clone(ctx context.Context) (context.Context, shell.CloneCloser, error) {
	newL := Clone(this.LState)
	ctx = context.WithValue(ctx, luaKey, newL)
	return ctx, &luaWrapper{newL}, nil
}

func (this *luaWrapper) Close() error {
	this.LState.Close()
	return nil
}

func Main() error {
	sh := shell.New()
	defer sh.Close()

	L := lua.NewState()
	defer L.Close()
	sh.SetTag(&luaWrapper{L})

	ctx := context.Background()

	langEngine := func(fname string) ([]byte, error) {
		return nil, L.DoFile(fname)
	}
	shellEngine := func(fname string) error {
		fd, err := os.Open(fname)
		if err != nil {
			return err
		}
		stream1 := shell.NewCmdStreamFile(fd)
		_, err = sh.Loop(ctx, stream1)
		fd.Close()
		if err == io.EOF {
			return nil
		} else {
			return err
		}
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
