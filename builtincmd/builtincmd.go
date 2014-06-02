package builtincmd

import "io"
import "os"
import "os/exec"
import "strings"
import "fmt"
import "bytes"
import "path/filepath"
import "bufio"

import "../interpreter"
import "../ls"
import "github.com/shiena/ansicolor"

func cmd_exit(cmd *exec.Cmd) interpreter.WhatToDoAfterCmd {
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

func cmd_pwd(cmd *exec.Cmd) interpreter.WhatToDoAfterCmd {
	wd, _ := os.Getwd()
	io.WriteString(cmd.Stdout, wd)
	io.WriteString(cmd.Stdout, "\n")
	return interpreter.CONTINUE
}

func cmd_cd(cmd *exec.Cmd) interpreter.WhatToDoAfterCmd {
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

func cmd_ls(cmd *exec.Cmd) interpreter.WhatToDoAfterCmd {
	err := ls.Main(cmd.Args[1:], ansicolor.NewAnsiColorWriter(cmd.Stdout))
	if err != nil {
		io.WriteString(cmd.Stderr, err.Error())
		io.WriteString(cmd.Stderr, "\n")
	}
	return interpreter.CONTINUE
}

func cmd_set(cmd *exec.Cmd) interpreter.WhatToDoAfterCmd {
	if len(cmd.Args) <= 1 {
		for _, val := range os.Environ() {
			fmt.Fprintf(cmd.Stdout, "%s\n", val)
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

func cmd_source(cmd *exec.Cmd) interpreter.WhatToDoAfterCmd {
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

func cmd_echo(cmd *exec.Cmd) interpreter.WhatToDoAfterCmd {
	io.WriteString(cmd.Stdout, strings.Join(cmd.Args[1:], " "))
	io.WriteString(cmd.Stdout, "\n")
	return interpreter.CONTINUE
}

var buildInCmd = map[string]func(cmd *exec.Cmd) interpreter.WhatToDoAfterCmd{
	"cd":     cmd_cd,
	"exit":   cmd_exit,
	"ls":     cmd_ls,
	"set":    cmd_set,
	".":      cmd_source,
	"source": cmd_source,
	"pwd":    cmd_pwd,
	"echo":   cmd_echo,
}

func Exec(cmd *exec.Cmd, IsBackground bool) (interpreter.WhatToDoAfterCmd, error) {
	name := strings.ToLower(cmd.Args[0])
	if len(name) == 2 && strings.HasSuffix(name, ":") {
		os.Chdir(name + ".")
		return interpreter.CONTINUE, nil
	}
	function, ok := buildInCmd[name]
	if ok {
		return function(cmd), nil
	} else {
		return interpreter.THROUGH, nil
	}
}
