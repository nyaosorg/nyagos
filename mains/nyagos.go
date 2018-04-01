package mains

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/mattn/go-isatty"
	"github.com/zetamatta/go-getch"

	"github.com/zetamatta/nyagos/alias"
	"github.com/zetamatta/nyagos/commands"
	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/history"
	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/shell"
)

var default_history *history.Container

type MainStream struct {
	shell.Stream
}

func (this *MainStream) ReadLine(ctx context.Context) (context.Context, string, error) {
	ctx = context.WithValue(ctx, history.PackageId, default_history)
	ctx, line, err := this.Stream.ReadLine(ctx)
	if err != nil {
		return ctx, "", err
	}
	return ctx, line, nil
}

type ScriptEngineForOptionImpl struct{}

func (this *ScriptEngineForOptionImpl) SetArg(args []string) {}

func (this *ScriptEngineForOptionImpl) RunFile(fname string) ([]byte, error) {
	println("Script is not supported.")
	return nil, nil
}

func (this *ScriptEngineForOptionImpl) RunString(code string) error {
	println("Script is not supported.")
	return nil
}

func Main() error {
	shell.SetHook(func(ctx context.Context, it *shell.Cmd) (int, bool, error) {
		rc, done, err := commands.Exec(ctx, it)
		return rc, done, err
	})
	completion.AppendCommandLister(commands.AllNames)
	completion.AppendCommandLister(alias.AllNames)

	dos.CoInitializeEx(0, dos.COINIT_MULTITHREADED)
	defer dos.CoUninitialize()

	getch.DisableCtrlC()

	alias.Init()

	sh := shell.New()

	langEngine := func(fname string) ([]byte, error) {
		return nil, nil
	}
	shellEngine := func(fname string) error {
		fd, err := os.Open(fname)
		if err != nil {
			return err
		}
		stream1 := NewCmdStreamFile(fd)
		_, err = sh.Loop(stream1)
		fd.Close()
		if err == io.EOF {
			return nil
		} else {
			return err
		}
	}

	if !isatty.IsTerminal(os.Stdin.Fd()) {
		SilentMode = true
	}

	if !OptionNorc {
		if !SilentMode {
			fmt.Printf("Nihongo Yet Another GOing Shell %s-%s by %s\n",
				VersionOrStamp(),
				runtime.GOARCH,
				runtime.Version())
			fmt.Println("(c) 2014-2018 NYAOS.ORG <http://www.nyaos.org>")
		}
		if err := LoadScripts(shellEngine, langEngine); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	backupHistory := default_history
	defer func() {
		default_history = backupHistory
	}()

	var stream1 shell.Stream
	if isatty.IsTerminal(os.Stdin.Fd()) {
		constream := NewCmdStreamConsole(
			func() (int, error) {
				fmt.Fprint(readline.Console,
					Format2Prompt(os.Getenv("PROMPT")))
				return 0, nil
			})
		stream1 = constream
		default_history = constream.History
	} else {
		stream1 = NewCmdStreamFile(os.Stdin)
	}

	for {
		_, err := sh.Loop(&MainStream{stream1})
		if err == io.EOF {
			return err
		}
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
