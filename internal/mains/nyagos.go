//go:build !vanilla
// +build !vanilla

package mains

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"golang.org/x/text/transform"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/yuin/gopher-lua"

	"github.com/nyaosorg/nyagos/internal/alias"
	"github.com/nyaosorg/nyagos/internal/completion"
	"github.com/nyaosorg/nyagos/internal/config"
	"github.com/nyaosorg/nyagos/internal/frame"
	"github.com/nyaosorg/nyagos/internal/functions"
	"github.com/nyaosorg/nyagos/internal/shell"
)

type luaKeyT struct{}

var luaKey luaKeyT

type _LuaCallBack struct {
	Lua
}

type _ScriptEngineForOptionImpl struct {
	L  Lua
	Sh *shell.Shell
}

func (impl *_ScriptEngineForOptionImpl) SetArg(args []string) {
	if L := impl.L; L != nil {
		table := L.NewTable()
		for i, arg1 := range args {
			L.SetTable(table, lua.LNumber(i), lua.LString(arg1))
		}
		L.SetGlobal("arg", table)
	}
}

// doFileExceptAtMarkLines reads and executes a Lua script file, but ignores
// lines starting with `@` until a line not starting with `@` is found.
// This function is used in `nyagos --norc --lua-file FILENAME.CMD`
func doFileExceptAtMarkLines(L *lua.LState, fname string) error {
	fd, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fd.Close()

	reader := transform.NewReader(fd, &_AtShebangFilter{})
	f, err := L.Load(reader, fname)
	if err != nil {
		return err
	}
	L.Push(f)
	return L.PCall(0, 0, nil)
}

func (*_ScriptEngineForOptionImpl) RunFile(ctx context.Context, fname string) ([]byte, error) {
	L, ok := ctx.Value(luaKey).(Lua)
	if !ok {
		return nil, errors.New("script is not supported")
	}
	defer setContext(getContext(L), L)
	setContext(ctx, L)
	return nil, luaRedirect(ctx, os.Stdin, os.Stdout, os.Stderr, L, func() error {
		return doFileExceptAtMarkLines(L, fname)
	})
}

func (impl *_ScriptEngineForOptionImpl) RunString(ctx context.Context, code string) error {
	L, ok := ctx.Value(luaKey).(Lua)
	if !ok {
		return errors.New("script is not supported")
	}
	defer setContext(getContext(L), L)
	setContext(ctx, L)
	return luaRedirect(ctx, os.Stdin, os.Stdout, os.Stderr, L, func() error {
		return L.DoString(code)
	})
}

type luaWrapper struct {
	Lua
	Env *env
}

func (lw *luaWrapper) Clone(ctx context.Context) (context.Context, shell.CloneCloser, error) {
	newL, err := Clone(lw.Lua, lw.Env)
	if err != nil {
		return nil, nil, err
	}
	ctx = context.WithValue(ctx, luaKey, newL)
	return ctx, &luaWrapper{Lua: newL, Env: lw.Env}, nil
}

func (lw *luaWrapper) Close() error {
	lw.Lua.Close()
	return nil
}

func warningOnly(err error) error {
	fmt.Fprintln(os.Stderr, strings.TrimSpace(err.Error()))
	return nil
}

// Run is the entry of this package.
func Run(fsys fs.FS) error {
	ctx := context.Background()

	sh := shell.New()
	defer sh.Close()

	env := &env{Shell: sh, Env: &functions.Env{Value: sh}}

	L, err := NewLua(env)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		ctx = context.WithValue(ctx, luaKey, L)
		defer L.Close()
	}

	completion.HookToList = append(completion.HookToList, (&_LuaCallBack{Lua: L}).luaHookForComplete)

	if L != nil {
		sh.SetTag(&luaWrapper{Lua: L, Env: env})
	}
	sh.Console = colorable.NewColorableStdout()

	defer setContext(getContext(L), L)

	alias.LineFilter = func(ctx context.Context, line string) string {
		if L, ok := ctx.Value(luaKey).(Lua); ok {
			return luaLineFilter(ctx, L, line)
		}
		return line
	}

	script, err := frame.OptionParse(ctx, sh, &_ScriptEngineForOptionImpl{L: L, Sh: sh})
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
		if embed, err := fs.Sub(fsys, "embed"); err != nil {
			warningOnly(err)
		} else if err := frame.LoadScriptsFs(L, embed, warningOnly); err != nil {
			return err
		}
		if err := frame.LoadScripts(L, warningOnly); err != nil {
			return err
		}
	}

	if script != nil {
		if err := script(ctx); err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
	}

	var stream1 shell.Stream
	if !config.ReadStdinAsFile && isatty.IsTerminal(os.Stdin.Fd()) {
		constream := frame.NewCmdStreamConsole(
			func(w io.Writer) (int, error) {
				if L != nil {
					return printPrompt(ctx, sh, L, w)
				}
				(&functions.Env{Value: sh}).Prompt(
					&functions.Param{
						Args: []interface{}{frame.Format2Prompt(os.Getenv("PROMPT"))},
						In:   os.Stdin,
						Out:  os.Stdout,
						Err:  os.Stderr,
						Term: colorable.NewColorableStdout(),
					})
				return 0, nil
			})
		stream1 = constream
		sh.History = constream.History
	} else {
		stream1 = shell.NewCmdStreamFile(os.Stdin, nil)
	}
	if L != nil {
		return sh.ForEver(ctx, &luaFilterStream{Stream: stream1, L: L})
	}
	return sh.ForEver(ctx, stream1)
}
