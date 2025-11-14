package frame

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/go-ole/go-ole"
	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-windows-consoleicon"

	"github.com/nyaosorg/nyagos/internal/alias"
	"github.com/nyaosorg/nyagos/internal/commands"
	"github.com/nyaosorg/nyagos/internal/completion"
	"github.com/nyaosorg/nyagos/internal/onexit"
	"github.com/nyaosorg/nyagos/internal/shell"
)

func Setup(version string) {
	Version = strings.TrimSpace(version)

	shell.SetHook(func(ctx context.Context, it *shell.Cmd) (int, bool, error) {
		rc, done, err := commands.Exec(ctx, it)
		return rc, done, err
	})
	completion.AppendCommandLister(commands.AllNames)
	completion.AppendCommandLister(alias.AllNames)

	if ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED) == nil {
		onexit.Register(ole.CoUninitialize)
	}

	alias.Init()

	disableColors := colorable.EnableColorsStdout(nil)
	onexit.Register(disableColors)

	if clean, err := consoleicon.SetFromExe(); err == nil {
		onexit.Register(func() { clean(true) })
	}

	signal.Ignore(os.Interrupt, syscall.SIGINT)
}

func PanicHandler() {
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
