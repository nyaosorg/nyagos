package shell

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"sync"
	"syscall"

	"github.com/zetamatta/go-findfile"

	"github.com/zetamatta/nyagos/defined"
	"github.com/zetamatta/nyagos/dos"
)

var WildCardExpansionAlways = false

type CommandNotFound struct {
	Name string
	Err  error
}

func (this CommandNotFound) Stringer() string {
	return fmt.Sprintf("'%s' is not recognized as an internal or external command,\noperable program or batch file", this.Name)
}

func (this CommandNotFound) Error() string {
	return this.Stringer()
}

type session struct {
	unreadline []string
}

type CloneCloser interface {
	Clone(context.Context) (context.Context, CloneCloser, error)
	Close() error
}

type Shell struct {
	*session
	Stdout       *os.File
	Stderr       *os.File
	Stdin        *os.File
	Console      io.Writer
	tag          CloneCloser
	IsBackGround bool
}

func (sh *Shell) In() io.Reader          { return sh.Stdin }
func (sh *Shell) Out() io.Writer         { return sh.Stdout }
func (sh *Shell) Err() io.Writer         { return sh.Stderr }
func (sh *Shell) Term() io.Writer        { return sh.Console }
func (sh *Shell) Tag() CloneCloser       { return sh.tag }
func (sh *Shell) SetTag(tag CloneCloser) { sh.tag = tag }

type Cmd struct {
	Shell
	args            []string
	rawArgs         []string
	fullPath        string
	UseShellExecute bool
	Closers         []io.Closer
}

func (cmd *Cmd) Arg(n int) string      { return cmd.args[n] }
func (cmd *Cmd) Args() []string        { return cmd.args }
func (cmd *Cmd) SetArgs(s []string)    { cmd.args = s }
func (cmd *Cmd) RawArg(n int) string   { return cmd.rawArgs[n] }
func (cmd *Cmd) RawArgs() []string     { return cmd.rawArgs }
func (cmd *Cmd) SetRawArgs(s []string) { cmd.rawArgs = s }

func (cmd *Cmd) FullPath() string {
	if cmd.args == nil || len(cmd.args) <= 0 {
		return ""
	}
	if cmd.fullPath == "" {
		cmd.fullPath = dos.LookPath(cmd.args[0], "NYAGOSPATH")
	}
	return cmd.fullPath
}

func (cmd *Cmd) Close() {
	if cmd.Closers != nil {
		for _, c := range cmd.Closers {
			c.Close()
		}
		cmd.Closers = nil
	}
}

func (sh *Shell) Close() {}

func New() *Shell {
	return &Shell{
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		session: &session{},
	}
}

func (sh *Shell) Command() *Cmd {
	cmd := &Cmd{
		Shell: Shell{
			Stdin:   sh.Stdin,
			Stdout:  sh.Stdout,
			Stderr:  sh.Stderr,
			Console: sh.Console,
			tag:     sh.tag,
		},
	}
	if sh.session != nil {
		cmd.session = sh.session
	} else {
		cmd.session = &session{}
	}
	return cmd
}

type ArgsHookT func(ctx context.Context, sh *Shell, args []string) ([]string, error)

var argsHook = func(ctx context.Context, sh *Shell, args []string) ([]string, error) {
	return args, nil
}

func SetArgsHook(argsHook_ ArgsHookT) (rv ArgsHookT) {
	rv, argsHook = argsHook, argsHook_
	return
}

type HookT func(context.Context, *Cmd) (int, bool, error)

var hook = func(context.Context, *Cmd) (int, bool, error) {
	return 0, false, nil
}

func SetHook(hook_ HookT) (rv HookT) {
	rv, hook = hook, hook_
	return
}

var OnCommandNotFound = func(ctx context.Context, cmd *Cmd, err error) error {
	err = &CommandNotFound{cmd.args[0], err}
	return err
}

var LastErrorLevel int

func makeCmdline(args, rawargs []string) string {
	var buffer strings.Builder
	for i, s := range args {
		if i > 0 {
			buffer.WriteRune(' ')
		}
		if (len(rawargs) > i && len(rawargs[i]) > 0 && rawargs[i][0] == '"') || strings.ContainsAny(s, " &|<>\t\"") {
			fmt.Fprintf(&buffer, `"%s"`, strings.Replace(s, `"`, `\"`, -1))
		} else {
			buffer.WriteString(s)
		}
	}
	return buffer.String()
}

var UseSourceRunBatch = true

func (cmd *Cmd) spawnvpSilent(ctx context.Context) (int, error) {
	// command is empty.
	if len(cmd.args) <= 0 {
		return 0, nil
	}
	if defined.DBG {
		print("spawnvpSilent('", cmd.args[0], "')\n")
	}

	// aliases and lua-commands
	if errorlevel, done, err := hook(ctx, cmd); done || err != nil {
		return errorlevel, err
	}

	// command not found hook
	fullpath := cmd.FullPath()
	if fullpath == "" {
		return 255, OnCommandNotFound(ctx, cmd, os.ErrNotExist)
	}
	cmd.args[0] = fullpath

	if defined.DBG {
		print("exec.LookPath(", cmd.args[0], ")==", fullpath, "\n")
	}
	if WildCardExpansionAlways {
		cmd.args = findfile.Globs(cmd.args)
	}
	if cmd.UseShellExecute {
		// GUI Application
		cmdline := makeCmdline(cmd.args[1:], cmd.rawArgs[1:])
		err := dos.ShellExecute("open", fullpath, cmdline, "")
		return 0, err
	}
	if UseSourceRunBatch {
		lowerName := strings.ToLower(cmd.args[0])
		if strings.HasSuffix(lowerName, ".cmd") || strings.HasSuffix(lowerName, ".bat") {
			// Batch files
			return RawSource(cmd.RawArgs(), nil, false, cmd.Stdin, cmd.Stdout, cmd.Stderr)
		}
	}
	xcmd := exec.Command(cmd.args[0], cmd.args[1:]...)
	xcmd.Stdin = cmd.Stdin
	xcmd.Stdout = cmd.Stdout
	xcmd.Stderr = cmd.Stderr

	if xcmd.SysProcAttr == nil {
		xcmd.SysProcAttr = new(syscall.SysProcAttr)
	}
	cmdline := makeCmdline(xcmd.Args, cmd.rawArgs)
	if defined.DBG {
		println(cmdline)
	}
	xcmd.SysProcAttr.CmdLine = cmdline
	err := xcmd.Run()
	errorlevel, errorlevelOk := dos.GetErrorLevel(xcmd)
	if errorlevelOk {
		return errorlevel, err
	} else {
		return 255, err
	}
}

type AlreadyReportedError struct {
	Err error
}

func (_ AlreadyReportedError) Error() string {
	return ""
}

func IsAlreadyReported(err error) bool {
	_, ok := err.(AlreadyReportedError)
	return ok
}

func (cmd *Cmd) Spawnvp(ctx context.Context) (int, error) {
	errorlevel, err := cmd.spawnvpSilent(ctx)
	if err != nil && err != io.EOF && !IsAlreadyReported(err) {
		if defined.DBG {
			val := reflect.ValueOf(err)
			fmt.Fprintf(cmd.Stderr, "error-type=%s\n", val.Type())
		}
		fmt.Fprintln(cmd.Stderr, err.Error())
		err = AlreadyReportedError{err}
	}
	return errorlevel, err
}

func (sh *Shell) Spawnlp(ctx context.Context, args, rawargs []string) (int, error) {
	cmd := sh.Command()
	defer cmd.Close()
	cmd.SetArgs(args)
	cmd.SetRawArgs(rawargs)
	return cmd.Spawnvp(ctx)
}

func (sh *Shell) Interpret(ctx context.Context, text string) (errorlevel int, finalerr error) {
	if defined.DBG {
		print("Interpret('", text, "')\n")
	}
	if sh == nil {
		return 255, errors.New("Fatal Error: Interpret: instance is nil")
	}
	errorlevel = 0
	finalerr = nil

	statements, statementsErr := Parse(text)
	if statementsErr != nil {
		if defined.DBG {
			print("Parse Error:", statementsErr.Error(), "\n")
		}
		return 0, statementsErr
	}
	if argsHook != nil {
		if defined.DBG {
			print("call argsHook\n")
		}
		for _, pipeline := range statements {
			for _, state := range pipeline {
				var err error
				state.Args, err = argsHook(ctx, sh, state.Args)
				if err != nil {
					return 255, err
				}
			}
		}
		if defined.DBG {
			print("done argsHook\n")
		}
	}
	for _, pipeline := range statements {
		for i, state := range pipeline {
			if state.Term == "|" && (i+1 >= len(pipeline) || len(pipeline[i+1].Args) <= 0) {
				return 255, errors.New("The syntax of the command is incorrect.")
			}
		}
	}

	for _, pipeline := range statements {

		var pipeIn *os.File = nil
		isBackGround := sh.IsBackGround
		for _, state := range pipeline {
			if state.Term == "&" {
				isBackGround = true
				break
			}
		}
		var wg sync.WaitGroup
		shutdown_immediately := false
		for i, state := range pipeline {
			if defined.DBG {
				print(i, ": pipeline loop(", state.Args[0], ")\n")
			}
			cmd := sh.Command()
			cmd.IsBackGround = isBackGround

			if pipeIn != nil {
				cmd.Stdin = pipeIn
				cmd.Closers = append(cmd.Closers, pipeIn)
				pipeIn = nil
			}

			var err error
			if state.Term[0] == '|' {
				var pipeOut *os.File
				pipeIn, pipeOut, err = os.Pipe()
				cmd.Stdout = pipeOut
				if state.Term == "|&" {
					cmd.Stderr = pipeOut
				}
				cmd.Closers = append(cmd.Closers, pipeOut)
			}

			for _, red := range state.Redirect {
				var fd *os.File
				fd, err = red.OpenOn(cmd)
				if err != nil {
					return 0, err
				}
				defer fd.Close()
			}

			cmd.args = state.Args
			cmd.rawArgs = state.RawArgs
			if i > 0 {
				cmd.IsBackGround = true
			}
			if len(pipeline) == 1 && dos.IsGui(cmd.FullPath()) {
				cmd.UseShellExecute = true
			}
			if i == len(pipeline)-1 && state.Term != "&" {
				// foreground execution.
				errorlevel, finalerr = cmd.Spawnvp(ctx)
				LastErrorLevel = errorlevel
				cmd.Close()
			} else {
				// background
				if !isBackGround {
					wg.Add(1)
				}
				newctx := ctx
				if tag := cmd.Tag(); tag != nil {
					var newtag CloneCloser
					if newctx, newtag, err = tag.Clone(ctx); err != nil {
						fmt.Fprintln(os.Stderr, err.Error())
						return -1, err
					} else {
						cmd.SetTag(newtag)
					}
				}
				go func(ctx1 context.Context, cmd1 *Cmd) {
					if !isBackGround {
						defer wg.Done()
					}
					cmd1.Spawnvp(ctx1)
					if tag := cmd1.Tag(); tag != nil {
						if err := tag.Close(); err != nil {
							fmt.Fprintln(os.Stderr, err.Error())
						}
					}
					cmd1.Close()
				}(newctx, cmd)
			}
		}
		if !isBackGround {
			wg.Wait()
			if shutdown_immediately {
				return errorlevel, nil
			}
			if len(pipeline) > 0 {
				switch pipeline[len(pipeline)-1].Term {
				case "&&":
					if errorlevel != 0 {
						return errorlevel, nil
					}
				case "||":
					if errorlevel == 0 {
						return errorlevel, nil
					}
				}
			}
		}
	}
	return
}
