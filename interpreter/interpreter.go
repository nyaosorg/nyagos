package interpreter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/zetamatta/go-findfile"

	"../dos"
)

const FLAG_AMP2NEWCONSOLE = false

var WildCardExpansionAlways = false

var dbg = false

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

type ErrorLevel int

const (
	NOERROR  ErrorLevel = 0
	THROUGH  ErrorLevel = -1
	SHUTDOWN ErrorLevel = -2
)

func (this ErrorLevel) HasValue() bool {
	return this >= NOERROR
}

func (this ErrorLevel) HasError() bool {
	return this > NOERROR
}

func (this ErrorLevel) String() string {
	switch this {
	case THROUGH:
		return "THROUGH"
	case SHUTDOWN:
		return "SHUTDOWN"
	default:
		return fmt.Sprintf("%d", this)
	}
}

type Interpreter struct {
	exec.Cmd
	Stdio        [3]*os.File
	HookCount    int
	Tag          interface{}
	PipeSeq      [2]uint
	IsBackGround bool
	RawArgs      []string

	OnClone func(*Interpreter) error
	Closers []io.Closer
}

func (this *Interpreter) closeAtEnd() {
	if this.Closers != nil {
		for _, c := range this.Closers {
			c.Close()
		}
		this.Closers = nil
	}
}

func (this *Interpreter) Close() {
	this.closeAtEnd()
}

func New() *Interpreter {
	this := Interpreter{
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

func (this *Interpreter) SetStdin(f *os.File) {
	this.Stdio[0] = f
	this.Stdin = f
}
func (this *Interpreter) SetStdout(f *os.File) {
	this.Stdio[1] = f
	this.Stdout = f
}
func (this *Interpreter) SetStderr(f *os.File) {
	this.Stdio[2] = f
	this.Stderr = f
}

func (this *Interpreter) Clone() (*Interpreter, error) {
	rv := new(Interpreter)
	rv.Stdio[0] = this.Stdio[0]
	rv.Stdio[1] = this.Stdio[1]
	rv.Stdio[2] = this.Stdio[2]
	rv.Stdin = this.Stdin
	rv.Stdout = this.Stdout
	rv.Stderr = this.Stderr
	rv.HookCount = this.HookCount
	rv.Tag = this.Tag
	rv.PipeSeq = rv.PipeSeq
	rv.Closers = nil
	rv.OnClone = this.OnClone
	if this.OnClone != nil {
		if err := this.OnClone(rv); err != nil {
			return nil, err
		}
	}
	return rv, nil
}

type ArgsHookT func(it *Interpreter, args []string) ([]string, error)

var argsHook = func(it *Interpreter, args []string) ([]string, error) {
	return args, nil
}

func SetArgsHook(argsHook_ ArgsHookT) (rv ArgsHookT) {
	rv, argsHook = argsHook, argsHook_
	return
}

type HookT func(*Interpreter) (ErrorLevel, error)

var hook = func(*Interpreter) (ErrorLevel, error) {
	return THROUGH, nil
}

func SetHook(hook_ HookT) (rv HookT) {
	rv, hook = hook, hook_
	return
}

var OnCommandNotFound = func(this *Interpreter, err error) error {
	err = &CommandNotFound{this.Args[0], err}
	return err
}

var ErrorLevelStr string

func nvl(a *os.File, b *os.File) *os.File {
	if a != nil {
		return a
	} else {
		return b
	}
}

func (this *Interpreter) spawnvp_noerrmsg() (ErrorLevel, error) {
	// command is empty.
	if len(this.Args) <= 0 {
		return NOERROR, nil
	}
	if dbg {
		print("spawnvp_noerrmsg('", this.Args[0], "')\n")
	}

	// aliases and lua-commands
	if errorlevel, err := hook(this); errorlevel != THROUGH || err != nil {
		return errorlevel, err
	}

	// command not found hook
	var err error
	this.Path, err = exec.LookPath(this.Args[0])
	if err != nil {
		return ErrorLevel(255), OnCommandNotFound(this, err)
	}

	if WildCardExpansionAlways {
		this.Args = findfile.Globs(this.Args)
	}

	// executable-file
	if FLAG_AMP2NEWCONSOLE {
		if this.SysProcAttr != nil && (this.SysProcAttr.CreationFlags&CREATE_NEW_CONSOLE) != 0 {
			err = this.Start()
			return ErrorLevel(0), err
		}
	}
	err = this.Run()

	errorlevel, errorlevelOk := dos.GetErrorLevel(&this.Cmd)
	if errorlevelOk {
		return ErrorLevel(errorlevel), err
	} else {
		return ErrorLevel(255), err
	}
}

func (this *Interpreter) Spawnvp() (ErrorLevel, error) {
	errorlevel, err := this.spawnvp_noerrmsg()
	if err != nil {
		fmt.Fprintln(this.Stderr, err.Error())
	}
	return errorlevel, err
}

type result_t struct {
	NextValue ErrorLevel
	Error     error
}

var pipeSeq uint = 0

func (this *Interpreter) Interpret(text string) (errorlevel ErrorLevel, err error) {
	if dbg {
		print("Interpret('", text, "')\n")
	}
	if this == nil {
		return ErrorLevel(255), errors.New("Fatal Error: Interpret: instance is nil")
	}
	errorlevel = NOERROR
	err = nil

	statements, statementsErr := Parse(text)
	if statementsErr != nil {
		if dbg {
			print("Parse Error:", statementsErr.Error(), "\n")
		}
		return NOERROR, statementsErr
	}
	if argsHook != nil {
		if dbg {
			print("call argsHook\n")
		}
		for _, pipeline := range statements {
			for _, state := range pipeline {
				state.Args, err = argsHook(this, state.Args)
				if err != nil {
					return ErrorLevel(255), err
				}
			}
		}
		if dbg {
			print("done argsHook\n")
		}
	}
	for _, pipeline := range statements {
		for i, state := range pipeline {
			if state.Term == "|" && (i+1 >= len(pipeline) || len(pipeline[i+1].Args) <= 0) {
				return ErrorLevel(255), errors.New("The syntax of the command is incorrect.")
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
		for i, state := range pipeline {
			if dbg {
				print(i, ": pipeline loop(", state.Args[0], ")\n")
			}
			cmd := new(Interpreter)
			cmd.PipeSeq[0] = pipeSeq
			cmd.PipeSeq[1] = uint(1 + i)
			cmd.IsBackGround = isBackGround
			cmd.Tag = this.Tag
			cmd.HookCount = this.HookCount
			cmd.SetStdin(nvl(this.Stdio[0], os.Stdin))
			cmd.SetStdout(nvl(this.Stdio[1], os.Stdout))
			cmd.SetStderr(nvl(this.Stdio[2], os.Stderr))
			cmd.OnClone = this.OnClone
			if this.OnClone != nil {
				if err := this.OnClone(cmd); err != nil {
					return ErrorLevel(255), err
				}
			}

			var err error = nil

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
					return NOERROR, err
				}
				defer fd.Close()
			}

			cmd.Args = state.Args
			cmd.RawArgs = state.RawArgs
			if i > 0 {
				cmd.IsBackGround = true
			}
			if i == len(pipeline)-1 && state.Term != "&" {
				errorlevel, err = cmd.Spawnvp()
				cmd.closeAtEnd()
				ErrorLevelStr = errorlevel.String()
				cmd.Close()
			} else {
				if !isBackGround {
					wg.Add(1)
				}
				go func(cmd1 *Interpreter) {
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
					cmd1.Spawnvp()
					cmd1.closeAtEnd()
					cmd1.Close()
				}(cmd)
			}
		}
		if !isBackGround {
			wg.Wait()
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
