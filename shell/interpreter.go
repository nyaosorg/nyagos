package shell

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/zetamatta/go-findfile"

	"github.com/zetamatta/nyagos/defined"
	"github.com/zetamatta/nyagos/nodos"
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

type History interface {
	Len() int
	DumpAt(n int) string
	IgnorePush(newvalue bool) bool
}

type _NulHistory struct{}

func (nul *_NulHistory) DumpAt(n int) string           { return "" }
func (nul *_NulHistory) Len() int                      { return 0 }
func (nul *_NulHistory) IgnorePush(newvalue bool) bool { return false }

type Shell struct {
	Stream
	History  History
	LineHook func(context.Context, *Cmd) (int, bool, error)
	ArgsHook func(context.Context, *Shell, []string) ([]string, error)
	*session
	Stdio        [3]*os.File
	Console      io.Writer
	tag          CloneCloser
	IsBackGround bool
}

func (sh *Shell) In() io.Reader          { return sh.Stdio[0] }
func (sh *Shell) Out() io.Writer         { return sh.Stdio[1] }
func (sh *Shell) Err() io.Writer         { return sh.Stdio[2] }
func (sh *Shell) Term() io.Writer        { return sh.Console }
func (sh *Shell) Tag() CloneCloser       { return sh.tag }
func (sh *Shell) SetTag(tag CloneCloser) { sh.tag = tag }
func (sh *Shell) GetHistory() History    { return sh.History }

type Cmd struct {
	Shell
	args            []string
	rawArgs         []string
	fullPath        string
	UseShellExecute bool
	Closers         []io.Closer
	env             map[string]string
	OnBackExec      func(int)
	OnBackDone      func(int)
}

func (cmd *Cmd) Arg(n int) string      { return cmd.args[n] }
func (cmd *Cmd) Args() []string        { return cmd.args }
func (cmd *Cmd) SetArgs(s []string)    { cmd.args = s }
func (cmd *Cmd) RawArg(n int) string   { return cmd.rawArgs[n] }
func (cmd *Cmd) RawArgs() []string     { return cmd.rawArgs }
func (cmd *Cmd) SetRawArgs(s []string) { cmd.rawArgs = s }

func (cmd *Cmd) Getenv(key string) string {
	if cmd.env != nil {
		if val, ok := cmd.env[strings.ToUpper(key)]; ok {
			return val
		}
	}
	return os.Getenv(key)
}

func (cmd *Cmd) Setenv(key, val string) {
	if cmd.env == nil {
		cmd.env = make(map[string]string)
	}
	cmd.env[strings.ToUpper(key)] = val
}

func (cmd *Cmd) DumpEnv() []string {
	if cmd.env == nil {
		return nil
	}
	osEnv := os.Environ()
	result := make([]string, 0, len(cmd.env)+len(osEnv))
	for _, equation := range osEnv {
		eqIndex := strings.IndexRune(equation, '=')
		if _, ok := cmd.env[strings.ToUpper(equation[:eqIndex])]; !ok {
			result = append(result, equation)
		}
	}
	for key, val := range cmd.env {
		result = append(result, key+"="+val)
	}
	return result
}

var LookCurdirOrder = nodos.LookCurdirFirst

func (cmd *Cmd) FullPath() string {
	if cmd.args == nil || len(cmd.args) <= 0 {
		return ""
	}
	if cmd.fullPath == "" {
		cmd.fullPath = cmd.lookpath()
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
		Stream:  NulStream,
		History: &_NulHistory{},
		LineHook: func(ctx context.Context, cmd1 *Cmd) (int, bool, error) {
			if hook != nil {
				return hook(ctx, cmd1)
			}
			return 0, false, nil
		},
		ArgsHook: func(ctx context.Context, sh *Shell, args []string) ([]string, error) {
			if argsHook != nil {
				return argsHook(ctx, sh, args)
			}
			return args, nil
		},
		Stdio:   [3]*os.File{os.Stdin, os.Stdout, os.Stderr},
		session: &session{},
	}
}

func (sh *Shell) Command() *Cmd {
	cmd := &Cmd{
		Shell: Shell{
			Stream:   sh.Stream,
			History:  sh.History,
			LineHook: sh.LineHook,
			ArgsHook: sh.ArgsHook,
			Stdio:    sh.Stdio,
			Console:  sh.Console,
			tag:      sh.tag,
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
			buffer.WriteByte(' ')
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
	for {
		// command is empty.
		if len(cmd.args) <= 0 {
			return 0, nil
		}
		eq := strings.IndexRune(cmd.args[0], '=')
		if eq <= 0 {
			break
		}
		envName := cmd.args[0][:eq]
		envNewValue := cmd.args[0][eq+1:]
		cmd.Setenv(envName, envNewValue)
		cmd.args = cmd.args[1:]
		cmd.rawArgs = cmd.rawArgs[1:]
	}

	if defined.DBG {
		print("spawnvpSilent('", cmd.args[0], "')\n")
	}

	// aliases and lua-commands
	if errorlevel, done, err := cmd.LineHook(ctx, cmd); done || err != nil {
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
	return cmd.startProcess(ctx)
}

func startAndWaitProcess(ctx context.Context, name string, args []string, procAttr *os.ProcAttr, onExec, onDone func(int)) (int, error) {
	if ctx != nil {
		select {
		case <-ctx.Done():
			return 252, ctx.Err()
		default:
		}
	}

	process, err := os.StartProcess(name, args, procAttr)
	if err != nil {
		return 255, err
	}

	if onExec != nil {
		onExec(process.Pid)
	}

	if ctx != nil {
		done := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
				os.Stderr.WriteString("^C\n")
				process.Kill()
			case <-done:
			}
		}()
		defer func() {
			close(done)
		}()
	}
	processState, err := process.Wait()

	if onDone != nil {
		onDone(process.Pid)
	}

	return processState.ExitCode(), err
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
			fmt.Fprintf(cmd.Err(), "error-type=%s\n", val.Type())
		}
		fmt.Fprintln(cmd.Err(), err.Error())
		err = AlreadyReportedError{err}
	}
	return errorlevel, err
}

func (sh *Shell) Spawnlpe(ctx context.Context, args, rawargs []string, env map[string]string) (int, error) {
	cmd := sh.Command()
	defer cmd.Close()
	cmd.SetArgs(args)
	cmd.SetRawArgs(rawargs)
	cmd.env = env
	return cmd.Spawnvp(ctx)
}

type TmpCloser struct {
	Closer func()
}

func (this *TmpCloser) Close() error {
	this.Closer()
	return nil
}

func (sh *Shell) Spawnlp(ctx context.Context, args, rawargs []string) (int, error) {

	return sh.Spawnlpe(ctx, args, rawargs, nil)
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

	statements, statementsErr := parse(sh.Stream, text)
	if statementsErr != nil {
		if defined.DBG {
			print("Parse Error:", statementsErr.Error(), "\n")
		}
		return 0, statementsErr
	}
	if sh.ArgsHook != nil {
		for _, pipeline := range statements {
			for _, state := range pipeline {
				var err error
				state.Args, err = sh.ArgsHook(ctx, sh, state.Args)
				if err != nil {
					return 255, err
				}
			}
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
				cmd.Stdio[0] = pipeIn
				cmd.Closers = append(cmd.Closers, pipeIn)
				pipeIn = nil
			}

			var err error
			if state.Term[0] == '|' {
				var pipeOut *os.File
				pipeIn, pipeOut, err = os.Pipe()
				cmd.Stdio[1] = pipeOut
				if state.Term == "|&" {
					cmd.Stdio[2] = pipeOut
				}
				cmd.Closers = append(cmd.Closers, pipeOut)
			}

			for _, f := range state.Redirect {
				c, err := f(cmd.Stdio[:])
				if err != nil {
					return 0, err
				}
				cmd.Closers = append(cmd.Closers, &TmpCloser{Closer: c})
			}

			cmd.args = state.Args
			cmd.rawArgs = state.RawArgs
			if i > 0 {
				cmd.IsBackGround = true
			}
			if len(pipeline) == 1 && isGui(cmd.FullPath()) {
				if len(state.Redirect) > 0 {
					// Use CreateProcess even if it is GUI application
					// bacause process by ShellExecute can not redirect. #361
					state.Term = "&"
				} else {
					cmd.UseShellExecute = true
					cmd.OnBackExec = func(pid int) {
						fmt.Fprintf(os.Stderr, "[%d]\n", pid)
					}
					cmd.OnBackDone = func(pid int) {
						fmt.Fprintf(os.Stderr, "[%d]+ Done\n", pid)
					}
				}
			}
			if i == len(pipeline)-1 && state.Term == "&" {
				cmd.OnBackExec = func(pid int) {
					fmt.Fprintf(os.Stderr, "[%d]\n", pid)
				}
				cmd.OnBackDone = func(pid int) {
					fmt.Fprintf(os.Stderr, "[%d]+ Done\n", pid)
				}
			}
			if i == len(pipeline)-1 && state.Term != "&" {
				// foreground execution.
				errorlevel, finalerr = cmd.Spawnvp(ctx)
				LastErrorLevel = errorlevel
				cmd.Close()
			} else {
				// background
				var newctx context.Context
				if isBackGround {
					// let Context not terminate background-work (#313's 2nd)
					// for the problem gvim starts with empty buffer
					// executing `git blame FILE | type | gvim - &`.
					newctx = context.Background()
				} else {
					wg.Add(1)
					newctx = ctx
				}
				if tag := cmd.Tag(); tag != nil {
					var newtag CloneCloser
					if newctx, newtag, err = tag.Clone(newctx); err != nil {
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
