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

	"../dos"
	. "../ifdbg"
)

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
	Stdout       *os.File
	Stderr       *os.File
	Stdin        *os.File
	Args         []string
	HookCount    int
	Tag          interface{}
	PipeSeq      [2]uint
	IsBackGround bool
	RawArgs      []string

	OnFork  func(*Cmd) error
	OffFork func(*Cmd) error
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
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	this.PipeSeq[0] = pipeSeq
	this.PipeSeq[1] = 0
	return &this
}

func (this *Cmd) Clone() (*Cmd, error) {
	rv := new(Cmd)
	rv.Args = this.Args
	rv.RawArgs = this.RawArgs
	rv.Stdin = this.Stdin
	rv.Stdout = this.Stdout
	rv.Stderr = this.Stderr
	rv.HookCount = this.HookCount
	rv.Tag = this.Tag
	rv.PipeSeq = this.PipeSeq
	rv.Closers = nil
	rv.OnFork = this.OnFork
	rv.OffFork = this.OffFork
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

func makeCmdline(args, rawargs []string) string {
	buffer := make([]byte, 0, 1024)
	for i, s := range args {
		if i > 0 {
			buffer = append(buffer, ' ')
		}
		if (len(rawargs) > i && len(rawargs[i]) > 0 && rawargs[i][0] == '"') || strings.ContainsAny(s, " &|<>\t\"") {
			buffer = append(buffer, '"')
			qs := strings.Replace(s, `"`, `\"`, -1)
			buffer = append(buffer, qs...)
			buffer = append(buffer, '"')
		} else {
			buffer = append(buffer, s...)
		}
	}
	return string(buffer)
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
	path1 := dos.LookPath(this.Args[0], "NYAGOSPATH")
	if path1 == "" {
		return 255, OnCommandNotFound(this, os.ErrNotExist)
	}
	this.Args[0] = path1

	if DBG {
		print("exec.LookPath(", this.Args[0], ")==", path1, "\n")
	}

	if WildCardExpansionAlways {
		this.Args = findfile.Globs(this.Args)
	}

	cmd1 := exec.Command(this.Args[0], this.Args[1:]...)
	cmd1.Stdin = this.Stdin
	cmd1.Stdout = this.Stdout
	cmd1.Stderr = this.Stderr

	if cmd1.SysProcAttr == nil {
		cmd1.SysProcAttr = new(syscall.SysProcAttr)
	}
	cmdline := makeCmdline(cmd1.Args, this.RawArgs)
	if DBG {
		println(cmdline)
	}
	cmd1.SysProcAttr.CmdLine = cmdline
	err = cmd1.Run()

	errorlevel, errorlevelOk := dos.GetErrorLevel(cmd1)
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

type gotoEol struct{}

var GotoEol = gotoEol{}

func (this *Cmd) InterpretContext(ctx_ context.Context, text string) (errorlevel int, finalerr error) {
	if DBG {
		print("Interpret('", text, "')\n")
	}
	if this == nil {
		return 255, errors.New("Fatal Error: Interpret: instance is nil")
	}
	errorlevel = 0
	finalerr = nil

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
				var err error
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
			cmd, err := this.Clone()
			if err != nil {
				return 255, err
			}
			cmd.PipeSeq[0] = pipeSeq
			cmd.PipeSeq[1] = uint(1 + i)
			cmd.IsBackGround = isBackGround

			ctx := context.WithValue(ctx_, GotoEol, func() {
				shutdown_immediately = true
				gotoeol, ok := ctx_.Value(GotoEol).(func())
				if ok {
					gotoeol()
				}
			})

			if pipeIn != nil {
				cmd.Stdin = pipeIn
				cmd.Closers = append(cmd.Closers, pipeIn)
				pipeIn = nil
			}

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

			cmd.Args = state.Args
			cmd.RawArgs = state.RawArgs
			if i > 0 {
				cmd.IsBackGround = true
			}
			if i == len(pipeline)-1 && state.Term != "&" {
				// foreground execution.
				errorlevel, finalerr = cmd.SpawnvpContext(ctx)
				LastErrorLevel = errorlevel
				cmd.Close()
			} else {
				// background
				if !isBackGround {
					wg.Add(1)
				}
				if cmd.OnFork != nil {
					if err := cmd.OnFork(cmd); err != nil {
						fmt.Fprintln(cmd.Stderr, err.Error())
						return -1, err
					}
				}
				go func(cmd1 *Cmd) {
					if !isBackGround {
						defer wg.Done()
					}
					cmd1.SpawnvpContext(ctx)
					if cmd1.OffFork != nil {
						if err := cmd1.OffFork(cmd1); err != nil {
							fmt.Fprintln(cmd1.Stderr, err.Error())
							goto exit
						}
					}
				exit:
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
