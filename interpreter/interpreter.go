package interpreter

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
)

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

type NextT int

const (
	THROUGH  NextT = 0
	CONTINUE NextT = 1
	SHUTDOWN NextT = 2
)

func (this NextT) String() string {
	switch this {
	case THROUGH:
		return "THROUGH"
	case CONTINUE:
		return "CONTINUE"
	case SHUTDOWN:
		return "SHUTDOWN"
	default:
		return "UNKNOWN"
	}
}

type Interpreter struct {
	exec.Cmd
	Stdio     [3]*os.File
	HookCount int
	Closer    []io.Closer
	Tag       interface{}
}

func (this *Interpreter) closeAtEnd() {
	if this.Closer != nil {
		for _, c := range this.Closer {
			c.Close()
		}
		this.Closer = nil
	}
}

func New() *Interpreter {
	this := Interpreter{
		Stdio: [3]*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	this.Stdin = os.Stdin
	this.Stdout = os.Stdout
	this.Stderr = os.Stderr
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

func (this *Interpreter) Clone() *Interpreter {
	rv := new(Interpreter)
	rv.Stdio[0] = this.Stdio[0]
	rv.Stdio[1] = this.Stdio[1]
	rv.Stdio[2] = this.Stdio[2]
	rv.Stdin = this.Stdin
	rv.Stdout = this.Stdout
	rv.Stderr = this.Stderr
	rv.HookCount = this.HookCount
	rv.Tag = this.Tag
	rv.Closer = nil
	return rv
}

type ArgsHookT func(it *Interpreter, args []string) []string

var argsHook = func(it *Interpreter, args []string) []string {
	return args
}

func SetArgsHook(argsHook_ ArgsHookT) (rv ArgsHookT) {
	rv, argsHook = argsHook, argsHook_
	return
}

type HookT func(*Interpreter) (NextT, error)

var hook = func(*Interpreter) (NextT, error) {
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

var errorStatusPattern = regexp.MustCompile("^exit status ([0-9]+)")
var ErrorLevel string

func nvl(a *os.File, b *os.File) *os.File {
	if a != nil {
		return a
	} else {
		return b
	}
}

func (this *Interpreter) Spawnvp() (NextT, error) {
	var whatToDo NextT = CONTINUE
	var err error = nil

	if len(this.Args) > 0 {
		whatToDo, err = hook(this)
		if whatToDo == THROUGH {
			this.Path, err = exec.LookPath(this.Args[0])
			if err == nil {
				err = this.Run()
			} else {
				err = OnCommandNotFound(this, err)
			}
			whatToDo = CONTINUE
		}
	}
	this.Stdio[1].Sync()
	this.Stdio[2].Sync()
	return whatToDo, err
}

type result_t struct {
	NextValue NextT
	Error     error
}

func (this *Interpreter) Interpret(text string) (next NextT, err error) {
	next = CONTINUE
	err = nil

	statements, statementsErr := Parse(text)
	if statementsErr != nil {
		return CONTINUE, statementsErr
	}
	if argsHook != nil {
		for _, pipeline := range statements {
			for _, state := range pipeline {
				state.Argv = argsHook(this, state.Argv)
			}
		}
	}
	for _, pipeline := range statements {
		var result chan result_t = nil
		var pipeOut *os.File = nil
		for i := len(pipeline) - 1; i >= 0; i-- {
			state := pipeline[i]

			cmd := new(Interpreter)
			cmd.Tag = this.Tag
			cmd.HookCount = this.HookCount
			cmd.SetStdin(nvl(this.Stdio[0], os.Stdin))
			cmd.SetStdout(nvl(this.Stdio[1], os.Stdout))
			cmd.SetStderr(nvl(this.Stdio[2], os.Stderr))

			var err error = nil

			if state.Term[0] == '|' {
				cmd.SetStdout(pipeOut)
				if state.Term == "|&" {
					cmd.SetStderr(pipeOut)
				}
				cmd.Closer = append(cmd.Closer, pipeOut)
			}

			if i > 0 && pipeline[i-1].Term[0] == '|' {
				var pipeIn *os.File
				pipeIn, pipeOut, err = os.Pipe()
				if err != nil {
					return CONTINUE, err
				}
				cmd.SetStdin(pipeIn)
				cmd.Closer = append(cmd.Closer, pipeIn)
			} else {
				pipeOut = nil
			}

			for _, red := range state.Redirect {
				err = red.OpenOn(cmd)
				if err != nil {
					return CONTINUE, err
				}
			}
			cmd.Args = state.Argv
			if i == len(pipeline)-1 && state.Term != "&" {
				result = make(chan result_t)
				go func() {
					whatToDo, err := cmd.Spawnvp()
					cmd.closeAtEnd()
					result <- result_t{whatToDo, err}
				}()
			} else {
				go func() {
					cmd.Spawnvp()
					cmd.closeAtEnd()
				}()
			}
		}
		if result != nil {
			resultValue := <-result
			if resultValue.Error != nil {
				m := errorStatusPattern.FindStringSubmatch(
					resultValue.Error.Error())
				if m != nil {
					ErrorLevel = m[1]
					resultValue.Error = nil
				} else {
					ErrorLevel = "-1"
				}
			} else {
				ErrorLevel = "0"
			}
			next = resultValue.NextValue
			err = resultValue.Error
		}
	}
	return
}
