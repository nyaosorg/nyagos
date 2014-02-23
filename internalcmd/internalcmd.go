package internalcmd

import "io"
import "os"
import "os/exec"
import "strings"

import "../interpreter"

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

var buildInCmd = map[string]func(cmd *exec.Cmd) interpreter.WhatToDoAfterCmd{
	"cd":   cmd_cd,
	"exit": cmd_exit,
}

func Exec(cmd *exec.Cmd, IsBackground bool) (interpreter.WhatToDoAfterCmd, error) {
	name := strings.ToLower(cmd.Args[0])
	if len(name) == 2 && strings.HasSuffix(name,":") {
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
