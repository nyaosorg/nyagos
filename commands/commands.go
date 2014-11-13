package commands

import "strings"
import "regexp"

import "../history"
import "../interpreter"
import "../dos"

var buildInCmd = map[string]func(cmd *interpreter.Interpreter) (interpreter.NextT, error){
	".":       cmd_source,
	"alias":   cmd_alias,
	"cd":      cmd_cd,
	"cls":     cmd_cls,
	"copy":    cmd_copy,
	"del":     cmd_del,
	"erase":   cmd_del,
	"exit":    cmd_exit,
	"history": history.CmdHistory,
	"ls":      cmd_ls,
	"md":      cmd_mkdir,
	"mkdir":   cmd_mkdir,
	"move":    cmd_move,
	"pwd":     cmd_pwd,
	"rd":      cmd_rmdir,
	"rmdir":   cmd_rmdir,
	"set":     cmd_set,
	"source":  cmd_source,
}

var unscoNamePattern = regexp.MustCompile("^__(.*)__$")

func Exec(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	name := strings.ToLower(cmd.Args[0])
	if len(name) == 2 && strings.HasSuffix(name, ":") {
		err := dos.Chdrive(name)
		return interpreter.CONTINUE, err
	}
	function, ok := buildInCmd[name]
	if !ok {
		m := unscoNamePattern.FindStringSubmatch(name)
		if m == nil {
			return interpreter.THROUGH, nil
		}
		name = m[1]
		function, ok = buildInCmd[name]
		if !ok {
			return interpreter.THROUGH, nil
		}
	}
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
	if cmd.IsBackGround {
		go func(cmd *interpreter.Interpreter) {
			function(cmd)
			if cmd.Closer != nil {
				cmd.Closer.Close()
				cmd.Closer = nil
			}
		}(cmd)
		return interpreter.CONTINUE, nil
	} else {
		next, err := function(cmd)
		if cmd.Closer != nil {
			cmd.Closer.Close()
			cmd.Closer = nil
		}
		return next, err
	}
}

func Init() {
	interpreter.SetHook(Exec)
}
