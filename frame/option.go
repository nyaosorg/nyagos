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

type optionArg struct {
	args []string
	sh   *shell.Shell
	e    ScriptEngineForOption
}

type optionT struct {
	F func()
	V func(*optionArg) (func(context.Context) error, error)
	U string
}

var optionMap = map[string]optionT{
	"-k": optionT{
		V: func(p *optionArg) (func(context.Context) error, error) {
			if len(p.args) <= 0 {
				return nil, errors.New("-k: requires parameters")
			}
			return func(ctx context.Context) error {
				p.sh.Interpret(ctx, p.args[0])
				return nil
			}, nil
		},
	},
	"-c": optionT{
		V: func(p *optionArg) (func(context.Context) error, error) {
			if len(p.args) <= 0 {
				return nil, errors.New("-c: requires parameters")
			}
			return func(ctx context.Context) error {
				p.sh.Interpret(ctx, p.args[0])
				return io.EOF
			}, nil
		},
	},
	"-b": optionT{
		V: func(p *optionArg) (func(context.Context) error, error) {
			if len(p.args) <= 0 {
				return nil, errors.New("-b: requires parameters")
			}
			data, err := base64.StdEncoding.DecodeString(p.args[0])
			if err != nil {
				return nil, err
			}
			text := string(data)
			return func(ctx context.Context) error {
				p.sh.Interpret(ctx, text)
				return io.EOF
			}, nil
		},
	},
	"-f": optionT{
		V: func(p *optionArg) (func(context.Context) error, error) {
			if len(p.args) <= 0 {
				return nil, errors.New("-f: requires parameters")
			}
			if strings.HasSuffix(strings.ToLower(p.args[0]), ".lua") {
				// lua script
				return func(ctx context.Context) error {
					p.e.SetArg(p.args)
					_, err := p.e.RunFile(ctx, p.args[0])
					if err != nil {
						return err
					} else {
						return io.EOF
					}
				}, nil
			} else {
				return func(ctx context.Context) error {
					// command script
					if err := p.sh.Source(ctx, p.args[0]); err != nil {
						return err
					}
					return io.EOF
				}, nil
			}
		},
	},
	"-e": optionT{
		V: func(p *optionArg) (func(context.Context) error, error) {
			if len(p.args) <= 0 {
				return nil, errors.New("-e: requires parameters")
			}
			return func(ctx context.Context) error {
				p.e.SetArg(p.args)
				err := p.e.RunString(ctx, p.args[0])
				if err != nil {
					return err
				} else {
					return io.EOF
				}
			}, nil
		},
	},
	"--lua-file": optionT{
		V: func(p *optionArg) (func(context.Context) error, error) {
			if len(p.args) <= 0 {
				return nil, errors.New("--lua-file: requires parameters")
			}
			return func(ctx context.Context) error {
				p.e.SetArg(p.args)
				_, err := p.e.RunFile(ctx, p.args[0])
				if err != nil {
					return err
				} else {
					return io.EOF
				}
			}, nil
		},
	},
	"--show-version-only": optionT{
		V: func(p *optionArg) (func(context.Context) error, error) {
			OptionNorc = true
			return func(context.Context) error {
				fmt.Printf("%s-%s\n", Version, runtime.GOARCH)
				return io.EOF
			}, nil
		},
	},
	"--disable-virtual-terminal-processing": optionT{
		F: func() {
			OptionEnableVirtualTerminalProcessing = false
		},
	},
	"--enable-virtual-terminal-processing": optionT{
		F: func() {
			OptionEnableVirtualTerminalProcessing = true
		},
	},
	"--no-go-colorable": optionT{
		F: func() {
			OptionGoColorable = false
		},
	},
	"--go-colorable": optionT{
		F: func() {
			OptionGoColorable = true
		},
	},
	"--norc": optionT{
		F: func() {
			OptionNorc = true
		},
	},
	"--look-curdir-first": optionT{
		F: func() {
			shell.LookCurdirOrder = dos.LookCurdirFirst
		},
	},
	"--look-curdir-last": optionT{
		F: func() {
			shell.LookCurdirOrder = dos.LookCurdirLast
		},
	},
	"--look-curdir-never": optionT{
		F: func() {
			shell.LookCurdirOrder = dos.LookCurdirNever
		},
	},
}

func OptionParse(sh *shell.Shell, e ScriptEngineForOption) (func(context.Context) error, error) {
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		if f, ok := optionMap[args[i]]; ok {
			if f.F != nil {
				f.F()
			}
			if f.V != nil {
				return f.V(&optionArg{
					args: args[i+1:],
					sh:   sh,
					e:    e,
				})
			}
		} else {
			fmt.Fprintf(os.Stderr, "%s: unknown parameter\n", args[i])
		}
	}
	return nil, nil
}

var SilentMode = false
