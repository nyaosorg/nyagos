package commands

import "os"
import "os/exec"
import "strings"
import "fmt"
import "bytes"
import "path/filepath"
import "bufio"

import alias "../alias/table"
import "../conio"
import "../history"
import "../interpreter"
import "../ls"

import "github.com/shiena/ansicolor"

func cmd_exit(cmd *exec.Cmd) interpreter.NextT {
	return interpreter.SHUTDOWN
}

func getHome() string {
	home := os.Getenv("HOME")
	if home != "" {
		return home
	}
	homeDrive := os.Getenv("HOMEDRIVE")
	if homeDrive != "" {
		homePath := os.Getenv("HOMEPATH")
		if homePath != "" {
			return homeDrive + homePath
		}
	}
	return ""
}

func cmd_pwd(cmd *exec.Cmd) interpreter.NextT {
	wd, _ := os.Getwd()
	fmt.Fprintln(cmd.Stdout, wd)
	return interpreter.CONTINUE
}

func cmd_cd(cmd *exec.Cmd) interpreter.NextT {
	if len(cmd.Args) >= 2 {
		os.Chdir(cmd.Args[1])
		return interpreter.CONTINUE
	}
	home := getHome()
	if home != "" {
		os.Chdir(home)
		return interpreter.CONTINUE
	}
	return cmd_pwd(cmd)
}

func cmd_ls(cmd *exec.Cmd) interpreter.NextT {
	err := ls.Main(cmd.Args[1:], ansicolor.NewAnsiColorWriter(cmd.Stdout))
	if err != nil {
		fmt.Fprintln(cmd.Stderr, err.Error())
	}
	return interpreter.CONTINUE
}

func cmd_set(cmd *exec.Cmd) interpreter.NextT {
	if len(cmd.Args) <= 1 {
		for _, val := range os.Environ() {
			fmt.Fprintln(cmd.Stdout, val)
		}
		return interpreter.CONTINUE
	}
	for _, arg := range cmd.Args[1:] {
		eqlPos := strings.Index(arg, "=")
		if eqlPos < 0 {
			fmt.Fprintf(cmd.Stdout, "%s=%s\n", arg, os.Getenv(arg))
		} else {
			os.Setenv(arg[:eqlPos], arg[eqlPos+1:])
		}
	}
	return interpreter.CONTINUE
}

func cmd_source(cmd *exec.Cmd) interpreter.NextT {
	envTxtPath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("nyagos-%d.tmp", os.Getpid()))

	var buffer bytes.Buffer

	buffer.WriteString(cmd.Args[1])
	buffer.WriteString(" & set > ")
	buffer.WriteString(envTxtPath)

	args := make([]string, 3)
	args[0] = os.Getenv("COMSPEC")
	args[1] = "/c"
	args[2] = buffer.String()
	var cmd2 exec.Cmd
	cmd2.Path = args[0]
	cmd2.Args = args
	cmd2.Env = nil
	cmd2.Dir = ""
	err := cmd2.Run()
	if err != nil {
		fmt.Fprintf(cmd.Stderr, "%s\n", err.Error())
	} else {
		fp, err := os.Open(envTxtPath)
		if err != nil {
			fmt.Fprintf(cmd.Stderr, "%s\n", err.Error())
			return interpreter.CONTINUE
		}
		defer os.Remove(envTxtPath)
		defer fp.Close()

		scr := bufio.NewScanner(fp)
		for scr.Scan() {
			line := scr.Text()
			eqlPos := strings.Index(line, "=")
			if eqlPos > 0 {
				os.Setenv(line[:eqlPos], line[eqlPos+1:])
			}
		}
	}
	return interpreter.CONTINUE
}

func cmd_echo(cmd *exec.Cmd) interpreter.NextT {
	fmt.Fprintln(cmd.Stdout, strings.Join(cmd.Args[1:], " "))
	return interpreter.CONTINUE
}

func cmd_cls(cmd *exec.Cmd) interpreter.NextT {
	conio.Cls()
	return interpreter.CONTINUE
}

func cmd_alias(cmd *exec.Cmd) interpreter.NextT {
	if len(cmd.Args) <= 1 {
		for key, val := range alias.Table {
			fmt.Fprintf(cmd.Stdout, "%s=%s\n", key, val)
		}
		return interpreter.CONTINUE
	}
	for _, args := range cmd.Args[1:] {
		if eqlPos := strings.IndexRune(args, '='); eqlPos >= 0 {
			key := args[0:eqlPos]
			val := args[eqlPos+1:]
			if len(val) > 0 {
				alias.Table[strings.ToLower(key)] = val
			} else {
				delete(alias.Table, strings.ToLower(key))
			}
		} else {
			key := strings.ToLower(args)
			val := alias.Table[key]

			fmt.Printf("%s=%s\n", key, val)
		}
	}
	return interpreter.CONTINUE
}

var buildInCmd = map[string]func(cmd *exec.Cmd) interpreter.NextT{
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
	"source":  cmd_source,
}

func Exec(cmd *exec.Cmd, IsBackground bool) (interpreter.NextT, error) {
	name := strings.ToLower(cmd.Args[0])
	if len(name) == 2 && strings.HasSuffix(name, ":") {
		os.Chdir(name + ".")
		return interpreter.CONTINUE, nil
	}
	function, ok := buildInCmd[name]
	if ok {
		newArgs := make([]string, 0)
		for _, arg1 := range cmd.Args {
			matches, _ := filepath.Glob(arg1)
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
			go function(cmd)
			return interpreter.CONTINUE, nil
		} else {
			return function(cmd), nil
		}
	} else {
		return interpreter.THROUGH, nil
	}
}
