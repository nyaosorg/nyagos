package frame

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/shell"
)

var OptionNorc = false
var OptionGoColorable = true
var OptionEnableVirtualTerminalProcessing = false

type ScriptEngineForOption interface {
	SetArg([]string)
	RunFile(context.Context, string) ([]byte, error)
	RunString(context.Context, string) error
}

func OptionParse(sh *shell.Shell, e ScriptEngineForOption) (func(context.Context) error, error) {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg1 := args[i]
		if arg1 == "-k" {
			i++
			if i >= len(args) {
				return nil, errors.New("-k: requires parameters")
			}
			return func(ctx context.Context) error {
				sh.Interpret(ctx, args[i])
				return nil
			}, nil
		} else if arg1 == "-c" {
			i++
			if i >= len(args) {
				return nil, errors.New("-c: requires parameters")
			}
			return func(ctx context.Context) error {
				sh.Interpret(ctx, args[i])
				return io.EOF
			}, nil
		} else if arg1 == "-b" {
			i++
			if i >= len(args) {
				return nil, errors.New("-b: requires parameters")
			}
			data, err := base64.StdEncoding.DecodeString(args[i])
			if err != nil {
				return nil, err
			}
			text := string(data)
			return func(ctx context.Context) error {
				sh.Interpret(ctx, text)
				return io.EOF
			}, nil
		} else if arg1 == "-f" {
			i++
			if i >= len(args) {
				return nil, errors.New("-f: requires parameters")
			}
			if strings.HasSuffix(strings.ToLower(args[i]), ".lua") {
				// lua script
				return func(ctx context.Context) error {
					e.SetArg(args[i:])
					_, err := e.RunFile(ctx, args[i])
					if err != nil {
						return err
					} else {
						return io.EOF
					}
				}, nil
			} else {
				return func(ctx context.Context) error {
					// command script
					if err := sh.Source(ctx, args[i]); err != nil {
						return err
					}
					return io.EOF
				}, nil
			}
		} else if arg1 == "-e" {
			i++
			if i >= len(args) {
				return nil, errors.New("-e: requires parameters")
			}
			return func(ctx context.Context) error {
				e.SetArg(args[i:])
				err := e.RunString(ctx, args[i])
				if err != nil {
					return err
				} else {
					return io.EOF
				}
			}, nil
		} else if arg1 == "--norc" {
			OptionNorc = true
		} else if arg1 == "--lua-file" {
			i++
			if i >= len(args) {
				return nil, errors.New("--lua-file: requires parameters")
			}
			return func(ctx context.Context) error {
				e.SetArg(args[i:])
				_, err := e.RunFile(ctx, args[i])
				if err != nil {
					return err
				} else {
					return io.EOF
				}
			}, nil
		} else if arg1 == "--show-version-only" {
			return func(context.Context) error {
				fmt.Printf("%s-%s\n", Version, runtime.GOARCH)
				return io.EOF
			}, nil
		} else if arg1 == "--go-colorable" {
			OptionGoColorable = true
		} else if arg1 == "--no-go-colorable" {
			OptionGoColorable = false
		} else if arg1 == "--enable-virtual-terminal-processing" {
			OptionEnableVirtualTerminalProcessing = true
		} else if arg1 == "--disable-virtual-terminal-processing" {
			OptionEnableVirtualTerminalProcessing = false
		} else if arg1 == "--look-curdir-first" {
			shell.LookCurdirOrder = dos.LookCurdirFirst
		} else if arg1 == "--look-curdir-last" {
			shell.LookCurdirOrder = dos.LookCurdirLast
		} else if arg1 == "--look-curdir-never" {
			shell.LookCurdirOrder = dos.LookCurdirNever
		} else {
			fmt.Fprintf(os.Stderr, "%s: unknwon parameter\n", arg1)
		}
	}

	return nil, nil
}

var SilentMode = false
