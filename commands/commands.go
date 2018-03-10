package commands

import (
	"context"
	"io"
	"regexp"
	"strings"

	"github.com/zetamatta/go-findfile"

	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/history"
	"github.com/zetamatta/nyagos/shell"
)

var BuildInCommand map[string]func(context.Context, *shell.Cmd) (int, error)
var unscoNamePattern = regexp.MustCompile("^__(.*)__$")

func Exec(ctx context.Context, cmd *shell.Cmd) (int, bool, error) {
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
	for name1 := range BuildInCommand {
		names = append(names, completion.Element{InsertStr: name1, ListupStr: name1})
	}
	return names
}

type Param interface {
	Arg(int) string
	Args() []string
	In() io.Reader
	Out() io.Writer
	Err() io.Writer
}

type ParamImpl struct{ *shell.Cmd }

func (this *ParamImpl) Arg(n int) string { return this.Cmd.Args[n] }
func (this *ParamImpl) Args() []string   { return this.Cmd.Args }
func (this *ParamImpl) In() io.Reader    { return this.Cmd.Stdin }
func (this *ParamImpl) Out() io.Writer   { return this.Cmd.Stdout }
func (this *ParamImpl) Err() io.Writer   { return this.Cmd.Stderr }

func cmd2param(f func(context.Context, Param) (int, error)) func(context.Context, *shell.Cmd) (int, error) {
	return func(ctx context.Context, cmd *shell.Cmd) (int, error) {
		return f(ctx, &ParamImpl{cmd})
	}
}

func Init() {
	BuildInCommand = map[string]func(context.Context, *shell.Cmd) (int, error){
		".":        cmd2param(cmdSource),
		"alias":    cmd_alias,
		"attrib":   cmd_attrib,
		"bindkey":  cmd_bindkey,
		"box":      cmd_box,
		"cd":       cmd_cd,
		"clip":     cmd_clip,
		"clone":    cmd_clone,
		"cls":      cmd_cls,
		"chmod":    cmd_chmod,
		"copy":     cmd_copy,
		"del":      cmd_del,
		"dirs":     cmd_dirs,
		"diskfree": cmd_df,
		"diskused": cmd_du,
		"echo":     cmd_echo,
		"env":      cmd_env,
		"erase":    cmd_del,
		"exit":     cmd_exit,
		"foreach":  cmd_foreach,
		"history":  history.CmdHistory,
		"if":       cmd_if,
		"ln":       cmd_ln,
		"lnk":      cmd_lnk,
		"ls":       cmd_ls,
		"md":       cmd_mkdir,
		"mkdir":    cmd_mkdir,
		"more":     cmd_more,
		"move":     cmd_move,
		"open":     cmd_open,
		"popd":     cmd_popd,
		"pushd":    cmd_pushd,
		"pwd":      cmd_pwd,
		"rd":       cmd_rmdir,
		"rem":      cmd_rem,
		"rmdir":    cmd_rmdir,
		"set":      cmd_set,
		"source":   cmd2param(cmdSource),
		"su":       cmd_su,
		"touch":    cmd_touch,
		"type":     cmd_type,
		"which":    cmd_which,
	}
}
