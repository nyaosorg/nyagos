package shell

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"sync"
	"syscall"

	"github.com/zetamatta/go-findfile"

	"../dos"
	. "../ifdbg"
)

const FLAG_AMP2NEWCONSOLE = false

var WildCardExpansionAlways = false

type CommandNotFound struct {
	Name string
	Err  error
}

// from "TDM-GCC-64/x86_64-w64-mingw32/include/winbase.h"
const (
	CREATE_NEW_CONSOLE       = 0x10
	CREATE_NEW_PROCESS_GROUP = 0x200
)

func (this CommandNotFound) Stringer() string {
	return fmt.Sprintf("'%s' is not recognized as an internal or external command,\noperable program or batch file", this.Name)
}

func (this CommandNotFound) Error() string {
	return this.Stringer()
}

type Cmd struct {
	exec.Cmd
	Stdio        [3]*os.File
	HookCount    int
	Tag          interface{}
	PipeSeq      [2]uint
	IsBackGround bool
	RawArgs      []string

	OnClone func(*Cmd) error
	Closers []io.Closer
}

func (this *Cmd) GetRawArgs() []string {
	return this.RawArgs
}

func (this *Cmd) Close() {
	if this.Closers != nil {
		for _, c := range this.Closers {
			c.Close()
		}
		this.Closers = nil
	}
}

func New() *Cmd {
	this := Cmd{
		Stdio: [3]*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	this.Stdin = os.Stdin
	this.Stdout = os.Stdout
	this.Stderr = os.Stderr
	this.PipeSeq[0] = pipeSeq
	this.PipeSeq[1] = 0
	this.Tag = nil
	return &this
}

func (this *Cmd) SetStdin(f *os.File) {
	this.Stdio[0] = f
	this.Stdin = f
}
func (this *Cmd) SetStdout(f *os.File) {
	this.Stdio[1] = f
	this.Stdout = f
}
func (this *Cmd) SetStderr(f *os.File) {
	this.Stdio[2] = f
	this.Stderr = f
}

func (this *Cmd) Clone() (*Cmd, error) {
	rv := new(Cmd)
	rv.Args = this.Args
	rv.RawArgs = this.RawArgs
	rv.Stdio[0] = this.Stdio[0]
	rv.Stdio[1] = this.Stdio[1]
	rv.Stdio[2] = this.Stdio[2]
	rv.Stdin = this.Stdin
	rv.Stdout = this.Stdout
	rv.Stderr = this.Stderr
	rv.HookCount = this.HookCount
	rv.Tag = this.Tag
	rv.PipeSeq = this.PipeSeq
	rv.Closers = nil
	rv.OnClone = this.OnClone
	if this.OnClone != nil {
		if err := this.OnClone(rv); err != nil {
			return nil, err
		}
	}
	return rv, nil
}

type ArgsHookT func(it *Cmd, args []string) ([]string, error)

var argsHook = func(it *Cmd, args []string) ([]string, error) {
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

var OnCommandNotFound = func(this *Cmd, err error) error {
	err = &CommandNotFound{this.Args[0], err}
	return err
}

var LastErrorLevel int

func nvl(a *os.File, b *os.File) *os.File {
	if a != nil {
		return a
	} else {
		return b
	}
}

func (this *Cmd) spawnvp_noerrmsg(ctx context.Context) (int, error) {
	// command is empty.
	if len(this.Args) <= 0 {
		return 0, nil
	}
	if DBG {
		print("spawnvp_noerrmsg('", this.Args[0], "')\n")
	}

	// aliases and lua-commands
	if errorlevel, done, err := hook(ctx, this); done || err != nil {
		return errorlevel, err
	}

	// command not found hook
	var err error
	this.Path = dos.LookPath(this.Args[0], "NYAGOSPATH")
	if this.Path == "" {
		return 255, OnCommandNotFound(this, os.ErrNotExist)
	}
	this.Args[0] = this.Path
	if DBG {
		print("exec.LookPath(", this.Args[0], ")==", this.Path, "\n")
	}

	if WildCardExpansionAlways {
		this.Args = findfile.Globs(this.Args)
	}

	// executable-file
	if FLAG_AMP2NEWCONSOLE {
		if this.SysProcAttr != nil && (this.SysProcAttr.CreationFlags&CREATE_NEW_CONSOLE) != 0 {
			err = this.Start()
			return 0, err
		}
	}
	err = this.Run()

	errorlevel, errorlevelOk := dos.GetErrorLevel(&this.Cmd)
	if errorlevelOk {
		return errorlevel, err
	} else {
		return 255, err
	}
}

type AlreadyReportedError struct {
	Err error
}

func (this AlreadyReportedError) Error() string {
	return ""
}

func IsAlreadyReported(err error) bool {
	_, ok := err.(AlreadyReportedError)
	return ok
}

func (this *Cmd) Spawnvp() (int, error) {
	return this.SpawnvpContext(context.Background())
}

func (this *Cmd) SpawnvpContext(ctx context.Context) (int, error) {
	errorlevel, err := this.spawnvp_noerrmsg(ctx)
	if err != nil && err != io.EOF && !IsAlreadyReported(err) {
		if DBG {
			val := reflect.ValueOf(err)
			fmt.Fprintf(this.Stderr, "error-type=%s\n", val.Type())
		}
		fmt.Fprintln(this.Stderr, err.Error())
		err = AlreadyReportedError{err}
	}
	return errorlevel, err
}

var pipeSeq uint = 0

func (this *Cmd) Interpret(text string) (int, error) {
	return this.InterpretContext(context.Background(), text)
}

func (this *Cmd) InterpretContext(ctx_ context.Context, text string) (errorlevel int, err error) {
	if DBG {
		print("Interpret('", text, "')\n")
	}
	if this == nil {
		return 255, errors.New("Fatal Error: Interpret: instance is nil")
	}
	errorlevel = 0
	err = nil

	statements, statementsErr := Parse(text)
	if statementsErr != nil {
		if DBG {
			print("Parse Error:", statementsErr.Error(), "\n")
		}
		return 0, statementsErr
	}
	if argsHook != nil {
		if DBG {
			print("call argsHook\n")
		}
		for _, pipeline := range statements {
			for _, state := range pipeline {
				state.Args, err = argsHook(this, state.Args)
				if err != nil {
					return 255, err
				}
			}
		}
		if DBG {
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
		pipeSeq++
		isBackGround := this.IsBackGround
		for _, state := range pipeline {
			if state.Term == "&" {
				isBackGround = true
				break
			}
		}
		var wg sync.WaitGroup
		shutdown_immediately := false
		for i, state := range pipeline {
			if DBG {
				print(i, ": pipeline loop(", state.Args[0], ")\n")
			}
			cmd := new(Cmd)
			cmd.PipeSeq[0] = pipeSeq
			cmd.PipeSeq[1] = uint(1 + i)
			cmd.IsBackGround = isBackGround
			cmd.Tag = this.Tag
			cmd.HookCount = this.HookCount
			cmd.SetStdin(nvl(this.Stdio[0], os.Stdin))
			cmd.SetStdout(nvl(this.Stdio[1], os.Stdout))
			cmd.SetStderr(nvl(this.Stdio[2], os.Stderr))
			cmd.OnClone = this.OnClone

			ctx := context.WithValue(ctx_, "rawargs", state.RawArgs)
			ctx = context.WithValue(ctx, "exec",
				func(cmdline string) (int, error) {
					return cmd.InterpretContext(ctx, cmdline)
				})
			ctx = context.WithValue(ctx, "gotoeol", func() {
				shutdown_immediately = true
				gotoeol, ok := ctx_.Value("gotoeol").(func())
				if ok {
					gotoeol()
				}
			})
			ctx = context.WithValue(ctx, "errorlevel", LastErrorLevel)
			if this.OnClone != nil {
				if err := this.OnClone(cmd); err != nil {
					return 255, err
				}
			}

			if pipeIn != nil {
				cmd.SetStdin(pipeIn)
				cmd.Closers = append(cmd.Closers, pipeIn)
				pipeIn = nil
			}

			if state.Term[0] == '|' {
				var pipeOut *os.File
				pipeIn, pipeOut, err = os.Pipe()
				cmd.SetStdout(pipeOut)
				if state.Term == "|&" {
					cmd.SetStderr(pipeOut)
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

			cmd.Args = state.Args
			cmd.RawArgs = state.RawArgs
			if i > 0 {
				cmd.IsBackGround = true
			}
			if i == len(pipeline)-1 && state.Term != "&" {
				// foreground execution.
				errorlevel, err = cmd.SpawnvpContext(ctx)
				LastErrorLevel = errorlevel
				cmd.Close()
			} else {
				// background
				if !isBackGround {
					wg.Add(1)
				}
				go func(cmd1 *Cmd) {
					if isBackGround {
						if FLAG_AMP2NEWCONSOLE {
							if len(pipeline) == 1 {
								cmd1.SysProcAttr = &syscall.SysProcAttr{
									CreationFlags: CREATE_NEW_CONSOLE |
										CREATE_NEW_PROCESS_GROUP,
								}
							}
						}
					} else {
						defer wg.Done()
					}
					cmd1.SpawnvpContext(ctx)
					cmd1.Close()
				}(cmd)
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
