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
	HookCount    int
	IsBackGround bool
	Closer       io.Closer
}

func New() *Interpreter {
	return new(Interpreter)
}

func (this *Interpreter) Clone() *Interpreter {
	rv := new(Interpreter)
	rv.Stdout = this.Stdout
	rv.Stderr = this.Stderr
	rv.Stdin = this.Stdin
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

func (this *Interpreter) Interpret(text string) (NextT, error) {
	statements := Parse(text)
	for _, pipeline := range statements {
		var pipeIn *os.File = nil
		for _, state := range pipeline {
			var cmd Interpreter
			cmd.HookCount = this.HookCount
			if this.Stderr != nil {
				cmd.Stdin = this.Stdin
			} else {
				cmd.Stdin = os.Stdin
			}
			if this.Stdout != nil {
				cmd.Stdout = this.Stdout
			} else {
				cmd.Stdout = os.Stdout
			}
			if this.Stderr != nil {
				cmd.Stderr = this.Stderr
			} else {
				cmd.Stderr = os.Stderr
			}
			if pipeIn != nil {
				cmd.Stdin = pipeIn
				pipeIn = nil
			}
			if state.Redirect[0].Path != "" {
				fd, err := os.Open(state.Redirect[0].Path)
				if err != nil {
					return CONTINUE, err
				}
				defer fd.Close()
				cmd.Stdin = fd
			}
			if state.Redirect[1].Path != "" {
				var fd *os.File
				var err error
				if state.Redirect[1].IsAppend {
					fd, err = os.OpenFile(state.Redirect[1].Path,
						os.O_APPEND, 0666)
				} else {
					fd, err = os.OpenFile(state.Redirect[1].Path,
						os.O_CREATE|os.O_TRUNC, 0666)
				}
				if err != nil {
					return CONTINUE, err
				}
				defer fd.Close()
				cmd.Stdout = fd
			}
			if state.Redirect[2].Path != "" {
				var fd *os.File
				var err error
				if state.Redirect[2].IsAppend {
					fd, err = os.OpenFile(state.Redirect[2].Path,
						os.O_APPEND, 0666)
				} else {
					fd, err = os.OpenFile(state.Redirect[2].Path,
						os.O_CREATE|os.O_TRUNC, 0666)
				}
				if err != nil {
					return CONTINUE, err
				}
				defer fd.Close()
				cmd.Stderr = fd
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
				cmd.Stdout = pipeOut
				if state.Term == "|&" {
					cmd.Stderr = pipeOut
				}
			case "&":
				isBackGround = true
			}
			var whatToDo NextT

			if len(state.Argv) > 0 {
				if argsHook != nil {
					state.Argv = argsHook(state.Argv)
				}
				cmd.Args = state.Argv
				cmd.IsBackGround = isBackGround
				cmd.Closer = pipeOut
				whatToDo, err = hook(&cmd)
				if whatToDo == THROUGH {
					cmd.Path, err = exec.LookPath(state.Argv[0])
					if err == nil {
						if isBackGround {
							err = cmd.Start()
						} else {
							err = cmd.Run()
						}
					}
				} else {
					pipeOut = nil
				}
			}
			if pipeOut != nil {
				// pipeOut.Close()
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
			if whatToDo == SHUTDOWN {
				return SHUTDOWN, err
			}
			if err != nil {
				return CONTINUE, err
			}
		}
	}
	return CONTINUE, nil
}
