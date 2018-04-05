package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/zetamatta/go-getch"

	"github.com/zetamatta/nyagos/alias"
	"github.com/zetamatta/nyagos/commands"
	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/defined"
	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/mains"
	"github.com/zetamatta/nyagos/shell"
)

var stamp string
var commit string
var version string

func startMain() error {
	defer mains.PanicHandler()

	mains.Stamp = stamp
	mains.Commit = commit
	mains.Version = version

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

	return switchMain()
}

func main() {
	if err := startMain(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		if err != io.EOF {
			if defined.DBG {
				var dummy [1]byte
				os.Stdin.Read(dummy[:])
			}
			os.Exit(1)
		}
	}
	if defined.DBG {
		var dummy [1]byte
		os.Stdin.Read(dummy[:])
	}
	os.Exit(0)
}
