package frame

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zetamatta/nyagos/commands"
	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/shell"
	"github.com/zetamatta/nyagos/texts"
)

// OptionNorc is true, then rcfiles are not executed.
var OptionNorc = false

// OptionGoColorable is true,
// then escape sequences are interpreted by go-colorable library.
var OptionGoColorable = true

// OptionEnableVirtualTerminalProcessing is true,
// then Windows10's ENABLE_VIRTUAL_TERMINAL_PROCESSING is enabled.
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
	ctx  context.Context // ctx is the Context object at parsing
}

type optionT struct {
	F func()
	V func(*optionArg) (func(context.Context) error, error)
	U string
}

var optionMap = map[string]optionT{
	"--lua-first": {
		U: "\"LUACODE\"\nExecute \"LUACODE\" before processing any rcfiles and continue shell",
		V: func(p *optionArg) (func(context.Context) error, error) {
			if len(p.args) <= 0 {
				return nil, errors.New("--lua-first: requires parameters")
			}
			if err := p.e.RunString(p.ctx, p.args[0]); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
			return nil, nil
		},
	},
	"--cmd-first": {
		U: "\"COMMAND\"\nExecute \"COMMAND\" before processing any rcfiles and continue shell",
		V: func(p *optionArg) (func(context.Context) error, error) {
			if len(p.args) <= 0 {
				return nil, errors.New("--cmd-first: requires parameters")
			}
			p.sh.Interpret(p.ctx, p.args[0])
			return nil, nil
		},
	},
	"-k": {
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
	"-c": {
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
	"-b": {
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
	"-f": {
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
					}
					return io.EOF
				}, nil
			}
			return func(ctx context.Context) error {
				// command script
				if err := p.sh.Source(ctx, p.args[0]); err != nil {
					return err
				}
				return io.EOF
			}, nil
		},
	},
	"-e": {
		U: "\"SCRIPTCODE\"\nExecute SCRIPTCODE with Lua interpreter and quit.",
		V: func(p *optionArg) (func(context.Context) error, error) {
			if len(p.args) <= 0 {
				return nil, errors.New("-e: requires parameters")
			}
			return func(ctx context.Context) error {
				p.e.SetArg(p.args)
				err := p.e.RunString(ctx, p.args[0])
				if err != nil {
					return err
				}
				return io.EOF
			}, nil
		},
	},
	"--lua-file": {
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
				}
				return io.EOF
			}, nil
		},
	},
	"--show-version-only": {
		U: "\nshow version only",
		V: func(p *optionArg) (func(context.Context) error, error) {
			OptionNorc = true
			return func(context.Context) error {
				fmt.Printf("%s-%s\n", Version, runtime.GOARCH)
				return io.EOF
			}, nil
		},
	},
	"--disable-virtual-terminal-processing": {
		U: "\nDo not use Windows10's native ESCAPE SEQUENCE.",
		F: func() {
			OptionEnableVirtualTerminalProcessing = false
		},
	},
	"--enable-virtual-terminal-processing": {
		U: "\nEnable Windows10's native ESCAPE SEQUENCE.\nIt should be used with `--no-go-colorable`.",
		F: func() {
			OptionEnableVirtualTerminalProcessing = true
		},
	},
	"--no-go-colorable": {
		U: "\nDo not use the ESCAPE SEQUENCE emulation with go-colorable library.",
		F: func() {
			OptionGoColorable = false
		},
	},
	"--go-colorable": {
		U: "\nUse the ESCAPE SEQUENCE emulation with go-colorable library.",
		F: func() {
			OptionGoColorable = true
		},
	},
	"--norc": {
		U: "\nDo not load the startup-scripts: `~\\.nyagos` , `~\\_nyagos`\nand `(BINDIR)\\nyagos.d\\*`.",
		F: func() {
			OptionNorc = true
		},
	},
	"--look-curdir-first": {
		U: "\nSearch for the executable from the current directory before %PATH%.\n(compatible with CMD.EXE)",
		F: func() {
			shell.LookCurdirOrder = dos.LookCurdirFirst
		},
	},
	"--look-curdir-last": {
		U: "\nSearch for the executable from the current directory after %PATH%.\n(compatible with PowerShell)",
		F: func() {
			shell.LookCurdirOrder = dos.LookCurdirLast
		},
	},
	"--look-curdir-never": {
		U: "\nNever search for the executable from the current directory\nunless %PATH% contains.\n(compatible with UNIX Shells)",
		F: func() {
			shell.LookCurdirOrder = dos.LookCurdirNever
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
		for _, key := range texts.SortedKeys(optionMap) {
			val := optionMap[key]
			fmt.Printf("%s %s\n", key, strings.Replace(val.U, "\n", "\n\t", -1))
		}

		fmt.Println("\nThese script are called on startup")
		if me, err := os.Executable(); err == nil {
			binDir := filepath.Dir(me)
			nyagosD := filepath.Join(binDir, "nyagos.d")
			fmt.Printf("  %s\\*.lua\n", nyagosD)
			file1 := filepath.Join(binDir, ".nyagos")
			fmt.Printf("  %s (Lua)\n", file1)
			file1 = filepath.Join(binDir, "_nyagos")
			fmt.Printf("  %s (Command-lines)\n", file1)
		}

		home := strings.TrimSpace(os.Getenv("HOME"))
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		file1 := filepath.Join(home, ".nyagos")
		fmt.Printf("  %s (Lua)\n", file1)
		file1 = filepath.Join(home, "_nyagos")
		fmt.Printf("  %s (Command-lines)\n", file1)

		return io.EOF
	}, nil
}

func isDefault(value bool) string {
	if value {
		return " [default]"
	}
	return ""
}

func OptionParse(_ctx context.Context, sh *shell.Shell, e ScriptEngineForOption) (func(context.Context) error, error) {
	args := os.Args[1:]
	optionMap["-h"] = optionT{V: help, U: "\nPrint this usage"}
	optionMap["--help"] = optionT{V: help, U: "\nPrint this usage"}

	for key, val := range commands.BoolOptions {
		_key := strings.Replace(key, "_", "-", -1)
		_val := val
		optionMap["--"+_key] = optionT{
			F: func() {
				*_val.V = true
			},
			U: fmt.Sprintf("(lua: `nyagos.option.%s=true`)%s\n%s",
				key,
				isDefault(*val.V),
				_val.Usage),
		}
		optionMap["--no-"+_key] = optionT{
			F: func() {
				*_val.V = false
			},
			U: fmt.Sprintf("(lua: `nyagos.option.%s=false`)%s\n%s",
				key,
				isDefault(!*val.V),
				_val.NoUsage),
		}
	}

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
					ctx:  _ctx,
				})
			}
		} else {
			fmt.Fprintf(os.Stderr, "%s: unknown parameter\n", args[i])
		}
	}
	return nil, nil
}

var SilentMode = false
