// +build vanilla

package mains

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"

	"github.com/zetamatta/nyagos/frame"
	"github.com/zetamatta/nyagos/functions"
	"github.com/zetamatta/nyagos/shell"
)

type scriptEngineForOptionImpl struct{}

func (*scriptEngineForOptionImpl) SetArg(args []string) {}

func (*scriptEngineForOptionImpl) RunFile(ctx context.Context, fname string) ([]byte, error) {
	println("Script is not supported.")
	return nil, nil
}

func (*scriptEngineForOptionImpl) RunString(ctx context.Context, code string) error {
	println("Script is not supported.")
	return nil
}

// Main is the main routine on the build without Lua
func Main() error {
	disableColors := colorable.EnableColorsStdout(nil)
	defer disableColors()

	sh := shell.New()
	defer sh.Close()
	sh.Console = colorable.NewColorableStdout()
	ctx := context.Background()

	langEngine := func(fname string) ([]byte, error) {
		return nil, nil
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
		}
		return err
	}

	script, err := frame.OptionParse(ctx, sh, &scriptEngineForOptionImpl{})
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
			fmt.Println("(c) 2014-2020 NYAOS.ORG <http://www.nyaos.org>")
		}
		if err := frame.LoadScripts(shellEngine, langEngine); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
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
	if isatty.IsTerminal(os.Stdin.Fd()) {
		constream := frame.NewCmdStreamConsole(
			func() (int, error) {
				functions.Prompt(
					&functions.Param{
						Args: []interface{}{frame.Format2Prompt(os.Getenv("PROMPT"))},
						Out:  os.Stdout,
						Err:  os.Stderr,
						In:   os.Stdin,
						Term: colorable.NewColorableStdout(),
					},
				)
				return 0, nil
			})
		stream1 = constream
		frame.DefaultHistory = constream.History
		sh.History = constream.History
	} else {
		stream1 = shell.NewCmdStreamFile(os.Stdin)
	}
	sh.ForEver(ctx, stream1)
	return nil
}
