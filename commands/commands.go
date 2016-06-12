package commands

import (
	"regexp"
	"strings"

	"../dos"
	"../history"
	. "../interpreter"
)

var BuildInCommand map[string]func(*Interpreter) (ErrorLevel, error)
var unscoNamePattern = regexp.MustCompile("^__(.*)__$")

func Exec(cmd *Interpreter) (ErrorLevel, error) {
	name := strings.ToLower(cmd.Args[0])
	if len(name) == 2 && strings.HasSuffix(name, ":") {
		err := dos.Chdrive(name)
		return NOERROR, err
	}
	function, ok := BuildInCommand[name]
	if !ok {
		m := unscoNamePattern.FindStringSubmatch(name)
		if m == nil {
			return THROUGH, nil
		}
		name = m[1]
		function, ok = BuildInCommand[name]
		if !ok {
			return THROUGH, nil
		}
	}
	cmd.Args = dos.Globs(cmd.Args)
	next, err := function(cmd)
	return next, err
}

func Init() {
	BuildInCommand = map[string]func(*Interpreter) (ErrorLevel, error){
		".":       cmd_source,
		"alias":   cmd_alias,
		"cd":      cmd_cd,
		"cls":     cmd_cls,
		"copy":    cmd_copy,
		"del":     cmd_del,
		"dirs":    cmd_dirs,
		"echo":    cmd_echo,
		"erase":   cmd_del,
		"exit":    cmd_exit,
		"history": history.CmdHistory,
		"ln":      cmd_ln,
		"ls":      cmd_ls,
		"md":      cmd_mkdir,
		"mkdir":   cmd_mkdir,
		"move":    cmd_move,
		"popd":    cmd_popd,
		"pushd":   cmd_pushd,
		"pwd":     cmd_pwd,
		"rd":      cmd_rmdir,
		"rem":     cmd_rem,
		"rmdir":   cmd_rmdir,
		"set":     cmd_set,
		"source":  cmd_source,
		"touch":   cmd_touch,
		"which":   cmd_which,
	}
	SetHook(Exec)
}
