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
	RawArgs() []string
	Spawn(context.Context, []string, []string) (int, error)
	Loop(s shell.Stream) (int, error)
	ReadCommand(context.Context, shell.Stream) (context.Context, string, error)
}

type ParamImpl struct{ *shell.Cmd }

func (this *ParamImpl) Arg(n int) string  { return this.Cmd.Args[n] }
func (this *ParamImpl) Args() []string    { return this.Cmd.Args }
func (this *ParamImpl) In() io.Reader     { return this.Cmd.Stdin }
func (this *ParamImpl) Out() io.Writer    { return this.Cmd.Stdout }
func (this *ParamImpl) Err() io.Writer    { return this.Cmd.Stderr }
func (this *ParamImpl) RawArgs() []string { return this.Cmd.RawArgs }

func (this *ParamImpl) Spawn(ctx context.Context, args, rawargs []string) (int, error) {
	subCmd, err := this.Clone()
	if err != nil {
		return 0, err
	}
	subCmd.Args = args
	subCmd.RawArgs = rawargs
	return subCmd.SpawnvpContext(ctx)
}

func cmd2param(f func(context.Context, Param) (int, error)) func(context.Context, *shell.Cmd) (int, error) {
	return func(ctx context.Context, cmd *shell.Cmd) (int, error) {
		return f(ctx, &ParamImpl{cmd})
	}
}

func Init() {
	BuildInCommand = map[string]func(context.Context, *shell.Cmd) (int, error){
		".":        cmd2param(cmdSource),
		"alias":    cmd2param(cmdAlias),
		"attrib":   cmd2param(cmdAttrib),
		"bindkey":  cmd2param(cmdBindkey),
		"box":      cmd2param(cmdBox),
		"cd":       cmd2param(cmdCd),
		"clip":     cmd2param(cmdClip),
		"clone":    cmd2param(cmdClone),
		"cls":      cmd2param(cmdCls),
		"chmod":    cmd2param(cmdChmod),
		"copy":     cmd2param(cmdCopy),
		"del":      cmd2param(cmdDel),
		"dirs":     cmd2param(cmdDirs),
		"diskfree": cmd2param(cmdDiskFree),
		"diskused": cmd2param(cmdDiskUsed),
		"echo":     cmd2param(cmdEcho),
		"env":      cmd_env,
		"erase":    cmd2param(cmdDel),
		"exit":     cmd2param(cmdExit),
		"foreach":  cmd_foreach,
		"history":  history.CmdHistory,
		"if":       cmd2param(cmdIf),
		"ln":       cmd2param(cmdLn),
		"lnk":      cmd2param(cmdLnk),
		"ls":       cmd2param(cmdLs),
		"md":       cmd2param(cmdMkdir),
		"mkdir":    cmd2param(cmdMkdir),
		"more":     cmd2param(cmdMore),
		"move":     cmd2param(cmdMove),
		"open":     cmd2param(cmdOpen),
		"popd":     cmd2param(cmdPopd),
		"pushd":    cmd2param(cmdPushd),
		"pwd":      cmd2param(cmdPwd),
		"rd":       cmd2param(cmdRmdir),
		"rem":      cmd2param(cmdRem),
		"rmdir":    cmd2param(cmdRmdir),
		"set":      cmd2param(cmdSet),
		"source":   cmd2param(cmdSource),
		"su":       cmd2param(cmdSu),
		"touch":    cmd2param(cmdTouch),
		"type":     cmd2param(cmdType),
		"which":    cmd2param(cmdWhich),
	}
}
