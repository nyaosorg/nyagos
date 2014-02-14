package interpreter

import "os"
import "os/exec"
import "../parser"

// import "fmt"

func Interpret(text string) (int, error) {
	statements := parser.Parse(text)
	for _, pipeline := range statements {
		var pipeIn *os.File = nil
		for _, state := range pipeline {
			//fmt.Println(state)
			path, err := exec.LookPath(state.Argv[0])
			if err != nil {
				return -1, err
			}
			cmd := new(exec.Cmd)
			cmd.Path = path
			cmd.Args = state.Argv
			cmd.Env = nil
			cmd.Dir = ""
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			if pipeIn != nil {
				cmd.Stdin = pipeIn
				pipeIn = nil
			}
			if state.Redirect[0].Path != "" {
				fd, err := os.Open(state.Redirect[0].Path)
				if err != nil {
					return -1, err
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
					return -1, err
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
					return -1, err
				}
				defer fd.Close()
				cmd.Stderr = fd
			}
			var pipeOut *os.File = nil
			if state.Term == "|" {
				pipeIn, pipeOut, err = os.Pipe()
				if err != nil {
					return -1, err
				}
				defer pipeIn.Close()
				cmd.Stdout = pipeOut
			}
			if state.Term == "|" || state.Term == "&" {
				err = cmd.Start()
			} else {
				err = cmd.Run()
			}
			if pipeOut != nil {
				pipeOut.Close()
			}
			if err != nil {
				return -1, err
			}
		}
	}
	return 0, nil
}
