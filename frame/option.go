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
		U: "\"COMMAND\"\nExecute \"COMMAND\" and continue the command-line.",
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
		U: "\"COMMAND\"\nExecute `COMMAND` and quit.",
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
		U: "\"BASE64edCOMMAND\"\nDecode and execute the command which is encoded with Base64.",
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
		U: "FILE ARG1 ARG2 ...\n" +
			"If FILE's suffix is .lua, execute Lua-code on it.\n" +
			"The script can refer arguments as `arg[]`.\n" +
			"Otherwise, read and execute commands on it.",
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
		U: "\"SCRIPTCODE\"\nExecute SCRIPTCODE with Lua interpretor and quit.",
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
		U: "FILE ARG1 ARG2...\n" +
			"Execute FILE as Lua Script even if FILE's suffix is not .lua .\n" +
			"The script can refer arguments as `arg[]`.\n" +
			"Lines starting with `@` are ignored to embed into batchfile.",
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
		U: "\nshow version only",
		V: func(p *optionArg) (func(context.Context) error, error) {
			OptionNorc = true
			return func(context.Context) error {
				fmt.Printf("%s-%s\n", Version, runtime.GOARCH)
				return io.EOF
			}, nil
		},
	},
	"--disable-virtual-terminal-processing": optionT{
		U: "\nDo not use Windows10's native ESCAPE SEQUENCE.",
		F: func() {
			OptionEnableVirtualTerminalProcessing = false
		},
	},
	"--enable-virtual-terminal-processing": optionT{
		U: "\nEnable Windows10's native ESCAPE SEQUENCE.\nIt should be used with `--no-go-colorable`.",
		F: func() {
			OptionEnableVirtualTerminalProcessing = true
		},
	},
	"--no-go-colorable": optionT{
		U: "\nDo not use the ESCAPE SEQUENCE emulation with go-colorable library.",
		F: func() {
			OptionGoColorable = false
		},
	},
	"--go-colorable": optionT{
		U: "\nUse the ESCAPE SEQUENCE emulation with go-colorable library.",
		F: func() {
			OptionGoColorable = true
		},
	},
	"--norc": optionT{
		U: "\nDo not load the startup-scripts: `~\\.nyagos` , `~\\_nyagos`\nand `(BINDIR)\\nyagos.d\\*`.",
		F: func() {
			OptionNorc = true
		},
	},
	"--look-curdir-first": optionT{
		U: "\nSearch for the executable from the current directory before %PATH%.\n(compatible with CMD.EXE)",
		F: func() {
			shell.LookCurdirOrder = dos.LookCurdirFirst
		},
	},
	"--look-curdir-last": optionT{
		U: "\nSearch for the executable from the current directory after %PATH%.\n(compatible with PowerShell)",
		F: func() {
			shell.LookCurdirOrder = dos.LookCurdirLast
		},
	},
	"--look-curdir-never": optionT{
		U: "\nNever search for the executable from the current directory\nunless %PATH% contains.\n(compatible with UNIX Shells)",
		F: func() {
			shell.LookCurdirOrder = dos.LookCurdirNever
		},
	},
	"--no-use-source": optionT{
		U: "\nforbide batchfile to change environment variables of nyagos",
		F: func() {
			shell.UseSourceRunBatch = false
		},
	},
	"--use-source": optionT{
		U: "\nallow batchfile to change environment variables of nyagos",
		F: func() {
			shell.UseSourceRunBatch = true
		},
	},
}

func Title() {
	fmt.Printf("Nihongo Yet Another GOing Shell %s-%s by %s\n",
		VersionOrStamp(),
		runtime.GOARCH,
		runtime.Version())
	fmt.Println("(c) 2014-2018 NYAOS.ORG <http://www.nyaos.org>")
}

func help(p *optionArg) (func(context.Context) error, error) {
	OptionNorc = true
	return func(context.Context) error {
		Title()
		fmt.Println()
		for key, val := range optionMap {
			fmt.Printf("%s %s\n", key, strings.Replace(val.U, "\n", "\n\t", -1))
		}
		return io.EOF
	}, nil
}

func OptionParse(sh *shell.Shell, e ScriptEngineForOption) (func(context.Context) error, error) {
	args := os.Args[1:]
	optionMap["-h"] = optionT{V: help, U: "\nPrint this usage"}
	optionMap["--help"] = optionT{V: help, U: "\nPrint this usage"}

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
