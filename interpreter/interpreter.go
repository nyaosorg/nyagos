package interpreter

import (
	"io"
	"os"
	"os/exec"
	"regexp"
)

type NextT int

const (
	THROUGH  NextT = 0
	CONTINUE NextT = 1
	SHUTDOWN NextT = 2
)

type Interpreter struct {
	exec.Cmd
	Stdio        [3]*os.File
	HookCount    int
	IsBackGround bool
	Closer       io.Closer
}

func New() *Interpreter {
	this := Interpreter{
		Stdio: [3]*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	this.Stdin = os.Stdin
	this.Stdout = os.Stdout
	this.Stderr = os.Stderr
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
	// Dont Copy 'Closer' and 'IsBackGround'
	return rv
}

type ArgsHookT func(args []string) []string

var argsHook = func(args []string) []string {
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

	if argsHook != nil {
		this.Args = argsHook(this.Args)
	}
	if len(this.Args) > 0 {
		whatToDo, err = hook(this)
		if whatToDo == THROUGH {
			this.Path, err = exec.LookPath(this.Args[0])
			if err == nil {
				if this.IsBackGround {
					go func() {
						this.Run()
						if this.Closer != nil {
							this.Closer.Close()
							this.Closer = nil
						}
					}()
				} else {
					err = this.Run()
				}
			}
		}
	}
	if err != nil {
		m := errorStatusPattern.FindStringSubmatch(err.Error())
		if m != nil {
			ErrorLevel = m[1]
			err = nil
		} else {
			ErrorLevel = "-1"
		}
	} else {
		ErrorLevel = "0"
	}
	return whatToDo, err
}

func (this *Interpreter) Interpret(text string) (NextT, error) {
	statements, statementsErr := Parse(text)
	if statementsErr != nil {
		return CONTINUE, statementsErr
	}
	for _, pipeline := range statements {
		var pipeIn *os.File = nil
		for _, state := range pipeline {
			cmd := new(Interpreter)
			cmd.HookCount = this.HookCount
			cmd.SetStdin(nvl(this.Stdio[0], os.Stdin))
			cmd.SetStdout(nvl(this.Stdio[1], os.Stdout))
			cmd.SetStderr(nvl(this.Stdio[2], os.Stderr))
			if pipeIn != nil {
				cmd.SetStdin(pipeIn)
				pipeIn = nil
			}
			var err error = nil
			var pipeOut *os.File = nil
			isBackGround := false

			switch state.Term {
			case "|", "|&":
				isBackGround = true
				pipeIn, pipeOut, err = os.Pipe()
				if err != nil {
					return CONTINUE, err
				}
				// defer pipeIn.Close()
				cmd.SetStdout(pipeOut)
				if state.Term == "|&" {
					cmd.SetStderr(pipeOut)
				}
			case "&":
				isBackGround = true
			}

			for _, red := range state.Redirect {
				err = red.OpenOn(cmd)
				if err != nil {
					return CONTINUE, err
				}
			}
			cmd.Args = state.Argv
			cmd.IsBackGround = isBackGround
			cmd.Closer = pipeOut
			whatToDo, err := cmd.Spawnvp()
			if whatToDo == SHUTDOWN || err != nil {
				return whatToDo, err
			}
		}
	}
	return CONTINUE, nil
}
