package interpreter

import "io"
import "os"
import "os/exec"
import "regexp"
import "strconv"

type NextT int

const (
	THROUGH  NextT = 0
	CONTINUE NextT = 1
	SHUTDOWN NextT = 2
)

type Stdio struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type ArgsHookT func(args []string) []string

var argsHook = func(args []string) []string {
	return args
}

func SetArgsHook(argsHook_ ArgsHookT) (rv ArgsHookT) {
	rv, argsHook = argsHook, argsHook_
	return
}

type HookT func(*exec.Cmd, bool, io.Closer) (NextT, error)

var hook = func(*exec.Cmd, bool, io.Closer) (NextT, error) {
	return THROUGH, nil
}

func SetHook(hook_ HookT) (rv HookT) {
	rv, hook = hook, hook_
	return
}

var errorStatusPattern = regexp.MustCompile("^exit status\\s*(\\d+)")
var ErrorLevel int

func Interpret(text string, stdio *Stdio) (NextT, error) {
	statements := Parse(text)
	for _, pipeline := range statements {
		var pipeIn *os.File = nil
		for _, state := range pipeline {
			//fmt.Println(state)
			var cmd exec.Cmd
			if stdio == nil {
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
			} else {
				if stdio.Stderr != nil {
					cmd.Stdin = stdio.Stdin
				} else {
					cmd.Stdin = os.Stdin
				}
				if stdio.Stdout != nil {
					cmd.Stdout = stdio.Stdout
				} else {
					cmd.Stdout = os.Stdout
				}
				if stdio.Stderr != nil {
					cmd.Stderr = stdio.Stderr
				} else {
					cmd.Stderr = os.Stderr
				}
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
					fd, err = os.OpenFile(state.Redirect[1].Path, os.O_APPEND, 0666)
				} else {
					fd, err = os.OpenFile(state.Redirect[1].Path, os.O_CREATE, 0666)
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
					fd, err = os.OpenFile(state.Redirect[2].Path, os.O_APPEND, 0666)
				} else {
					fd, err = os.OpenFile(state.Redirect[2].Path, os.O_CREATE, 0666)
				}
				if err != nil {
					return CONTINUE, err
				}
				defer fd.Close()
				cmd.Stderr = fd
			}
			var err error = nil
			var pipeOut *os.File = nil
			if state.Term == "|" {
				pipeIn, pipeOut, err = os.Pipe()
				if err != nil {
					return CONTINUE, err
				}
				defer pipeIn.Close()
				cmd.Stdout = pipeOut
			}
			var whatToDo NextT

			isBackGround := (state.Term == "|" || state.Term == "&")

			if len(state.Argv) > 0 {
				if argsHook != nil {
					state.Argv = argsHook(state.Argv)
				}
				cmd.Args = state.Argv
				whatToDo, err = hook(&cmd, isBackGround, pipeOut)
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
				pipeOut.Close()
			}
			if err != nil {
				m := errorStatusPattern.FindStringSubmatch(err.Error())
				if m != nil {
					ErrorLevel, err = strconv.Atoi(m[1])
				} else {
					ErrorLevel = 0
				}
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
