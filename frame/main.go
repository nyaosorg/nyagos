package frame

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/go-ole/go-ole"
	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-windows-consoleicon"

	"github.com/nyaosorg/nyagos/completion"
	"github.com/nyaosorg/nyagos/history"
	"github.com/nyaosorg/nyagos/internal/alias"
	"github.com/nyaosorg/nyagos/internal/commands"
	"github.com/nyaosorg/nyagos/internal/shell"
)

var DefaultHistory *history.Container

func Start(mainHandler func() error) error {
	defer panicHandler()

	shell.SetHook(func(ctx context.Context, it *shell.Cmd) (int, bool, error) {
		rc, done, err := commands.Exec(ctx, it)
		return rc, done, err
	})
	completion.AppendCommandLister(commands.AllNames)
	completion.AppendCommandLister(alias.AllNames)

	if ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED) == nil {
		defer ole.CoUninitialize()
	}

	alias.Init()

	disableColors := colorable.EnableColorsStdout(nil)
	defer disableColors()

	if clean, err := consoleicon.SetFromExe(); err == nil {
		defer clean(true)
	}

	signal.Ignore(os.Interrupt, syscall.SIGINT)

	return mainHandler()
}

func panicHandler() {
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

	os.WriteFile("nyagos.dump", dump.Bytes(), 0666)

	var dummy [1]byte
	os.Stdin.Read(dummy[:])
}
