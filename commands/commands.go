package commands

import "io"
import "os"
import "os/exec"
import "strings"
import "fmt"

import "../alias"
import "../conio"
import "../history"
import "../interpreter"
import "./ls"
import "../dos"

import "github.com/shiena/ansicolor"

func cmd_exit(cmd *exec.Cmd) (interpreter.NextT, error) {
	return interpreter.SHUTDOWN, nil
}

func cmd_pwd(cmd *exec.Cmd) (interpreter.NextT, error) {
	wd, _ := os.Getwd()
	fmt.Fprintln(cmd.Stdout, dos.ReplaceHomeToTildeSlash(wd))
	return interpreter.CONTINUE, nil
}

var prevDir string

func cmd_cd(cmd *exec.Cmd) (interpreter.NextT, error) {
	if len(cmd.Args) >= 2 {
		prevDir_, err := os.Getwd()
		if err != nil {
			return interpreter.CONTINUE, err
		}
		if cmd.Args[1] == "-" {
			err = dos.Chdir(prevDir)
		} else {
			err = dos.Chdir(cmd.Args[1])
		}
		prevDir = prevDir_
		return interpreter.CONTINUE, err
	}
	home := dos.GetHome()
	if home != "" {
		prevDir, _ = os.Getwd()
		return interpreter.CONTINUE, dos.Chdir(home)
	}
	return cmd_pwd(cmd)
}

func cmd_ls(cmd *exec.Cmd) (interpreter.NextT, error) {
	return interpreter.CONTINUE,
		ls.Main(cmd.Args[1:], ansicolor.NewAnsiColorWriter(cmd.Stdout))
}

func cmd_set(cmd *exec.Cmd) (interpreter.NextT, error) {
	if len(cmd.Args) <= 1 {
		for _, val := range os.Environ() {
			fmt.Fprintln(cmd.Stdout, val)
		}
		return interpreter.CONTINUE, nil
	}
	for _, arg := range cmd.Args[1:] {
		eqlPos := strings.Index(arg, "=")
		if eqlPos < 0 {
			fmt.Fprintf(cmd.Stdout, "%s=%s\n", arg, os.Getenv(arg))
		} else {
			os.Setenv(arg[:eqlPos], arg[eqlPos+1:])
		}
	}
	return interpreter.CONTINUE, nil
}

func cmd_echo(cmd *exec.Cmd) (interpreter.NextT, error) {
	fmt.Fprintln(cmd.Stdout, strings.Join(cmd.Args[1:], " "))
	return interpreter.CONTINUE, nil
}

func cmd_cls(cmd *exec.Cmd) (interpreter.NextT, error) {
	conio.Cls()
	return interpreter.CONTINUE, nil
}

func cmd_alias(cmd *exec.Cmd) (interpreter.NextT, error) {
	if len(cmd.Args) <= 1 {
		for key, val := range alias.Table {
			fmt.Fprintf(cmd.Stdout, "%s=%s\n", key, val.String())
		}
		return interpreter.CONTINUE, nil
	}
	for _, args := range cmd.Args[1:] {
		if eqlPos := strings.IndexRune(args, '='); eqlPos >= 0 {
			key := args[0:eqlPos]
			val := args[eqlPos+1:]
			if len(val) > 0 {
				alias.Table[strings.ToLower(key)] = alias.New(val)
			} else {
				delete(alias.Table, strings.ToLower(key))
			}
		} else {
			key := strings.ToLower(args)
			val := alias.Table[key]

			fmt.Fprintf(cmd.Stdout, "%s=%s\n", key, val.String())
		}
	}
	return interpreter.CONTINUE, nil
}

var buildInCmd = map[string]func(cmd *exec.Cmd) (interpreter.NextT, error){
	".":       cmd_source,
	"alias":   cmd_alias,
	"cd":      cmd_cd,
	"cls":     cmd_cls,
	"echo":    cmd_echo,
	"exit":    cmd_exit,
	"history": history.CmdHistory,
	"ls":      cmd_ls,
	"pwd":     cmd_pwd,
	"set":     cmd_set,
	"rem":     cmd_rem,
	"source":  cmd_source,
}

func Exec(cmd *exec.Cmd, IsBackground bool, closer io.Closer) (interpreter.NextT, error) {
	name := strings.ToLower(cmd.Args[0])
	if len(name) == 2 && strings.HasSuffix(name, ":") {
		err := dos.Chdrive(name)
		return interpreter.CONTINUE, err
	}
	function, ok := buildInCmd[name]
	if ok {
		newArgs := make([]string, 0)
		for _, arg1 := range cmd.Args {
			matches, _ := dos.Glob(arg1)
			if matches == nil {
				newArgs = append(newArgs, arg1)
			} else {
				for _, s := range matches {
					newArgs = append(newArgs, s)
				}
			}
		}
		cmd.Args = newArgs
		if IsBackground {
			go func(cmd *exec.Cmd, closer io.Closer) {
				function(cmd)
				closer.Close()
			}(cmd, closer)
			return interpreter.CONTINUE, nil
		} else {
			next, err := function(cmd)
			closer.Close()
			return next, err
		}
	} else {
		return interpreter.THROUGH, nil
	}
}
