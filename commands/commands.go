package commands

import (
	"context"
	"os/exec"
	"regexp"
	"strings"

	"github.com/zetamatta/go-findfile"

	"../completion"
	"../dos"
	"../history"
)

var BuildInCommand map[string]func(context.Context, *exec.Cmd) (int, error)
var unscoNamePattern = regexp.MustCompile("^__(.*)__$")

func Exec(ctx context.Context, cmd *exec.Cmd) (int, bool, error) {
	name := strings.ToLower(cmd.Args[0])
	if len(name) == 2 && strings.HasSuffix(name, ":") {
		err := dos.Chdrive(name)
		return 0, true, err
	}
	function, ok := BuildInCommand[name]
	if !ok {
		m := unscoNamePattern.FindStringSubmatch(name)
		if m == nil {
			return 0, false, nil
		}
		name = m[1]
		function, ok = BuildInCommand[name]
		if !ok {
			return 0, false, nil
		}
	}
	cmd.Args = findfile.Globs(cmd.Args)
	next, err := function(ctx, cmd)
	return next, true, err
}

func AllNames() []completion.Element {
	names := make([]completion.Element, 0, len(BuildInCommand))
	for name1, _ := range BuildInCommand {
		names = append(names, completion.Element{name1, name1})
	}
	return names
}

func Init() {
	BuildInCommand = map[string]func(context.Context, *exec.Cmd) (int, error){
		".":       cmd_source,
		"alias":   cmd_alias,
		"bindkey": cmd_bindkey,
		"box":     cmd_box,
		"cd":      cmd_cd,
		"clone":   cmd_clone,
		"cls":     cmd_cls,
		"copy":    cmd_copy,
		"del":     cmd_del,
		"dirs":    cmd_dirs,
		"echo":    cmd_echo,
		"env":     cmd_env,
		"erase":   cmd_del,
		"exit":    cmd_exit,
		"history": history.CmdHistory,
		"if":      cmd_if,
		"ln":      cmd_ln,
		"ls":      cmd_ls,
		"md":      cmd_mkdir,
		"mkdir":   cmd_mkdir,
		"move":    cmd_move,
		"open":    cmd_open,
		"popd":    cmd_popd,
		"pushd":   cmd_pushd,
		"pwd":     cmd_pwd,
		"rd":      cmd_rmdir,
		"rem":     cmd_rem,
		"rmdir":   cmd_rmdir,
		"set":     cmd_set,
		"source":  cmd_source,
		"su":      cmd_su,
		"sudo":    cmd_sudo,
		"touch":   cmd_touch,
		"which":   cmd_which,
	}
}
