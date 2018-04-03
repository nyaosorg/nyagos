package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/zetamatta/go-getch"

	"github.com/zetamatta/nyagos/alias"
	"github.com/zetamatta/nyagos/commands"
	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/defined"
	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/mains"
	"github.com/zetamatta/nyagos/shell"
)

func whenPanic() {
	err := recover()
	if err == nil {
		return
	}
	var dump bytes.Buffer
	w := io.MultiWriter(os.Stderr, &dump)

	fmt.Fprintln(w, "************ Panic Occurred. ***********")
	fmt.Fprintln(w, err)
	w.Write(debug.Stack())
	fmt.Fprintln(w, "*** Please copy these error message ***")
	fmt.Fprintln(w, "*** And hit ENTER key to quit.      ***")

	ioutil.WriteFile("nyagos.dump", dump.Bytes(), 0666)

	var dummy [1]byte
	os.Stdin.Read(dummy[:])
}

var stamp string
var commit string
var version string

func startMain() error {
	defer whenPanic()

	if len(os.Args) >= 2 && os.Args[1] == "--show-version-only" {
		fmt.Printf("%s-%s\n", version, runtime.GOARCH)
		return nil
	}

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
