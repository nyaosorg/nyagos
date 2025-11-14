package functions

import (
	"io"

	"github.com/nyaosorg/nyagos/internal/shell"
)

type Env struct {
	Value interface {
		In() io.Reader
		Out() io.Writer
		Err() io.Writer
		GetHistory() shell.History
	}
}

func (e *Env) Table() map[string]func([]any) []any {
	return map[string]func([]any) []any{
		"access":             e.CmdAccess,
		"atou":               e.CmdAtoU,
		"atou_if_needed":     e.CmdAnsiToUtf8IfNeeded,
		"bitand":             e.CmdBitAnd,
		"bitor":              e.CmdBitOr,
		"chdir":              e.CmdChdir,
		"commonprefix":       e.CmdCommonPrefix,
		"complete_for_files": e.CmdCompleteForFiles,
		"dirname":            e.CmdDirName,
		"elevated":           e.CmdElevated,
		"envadd":             e.CmdEnvAdd,
		"envdel":             e.CmdEnvDel,
		"fields":             e.CmdFields,
		"getenv":             e.CmdGetEnv,
		"gethistory":         e.CmdGetHistory,
		"getkey":             e.CmdGetKey,
		"getkeys":            e.CmdGetKeys,
		"getviewwidth":       e.CmdGetViewWidth,
		"getwd":              e.CmdGetwd,
		"glob":               e.CmdGlob,
		"pathjoin":           e.CmdPathJoin,
		"setenv":             e.CmdSetEnv,
		"shellexecute":       e.CmdShellExecute,
		"skk":                e.CmdSkk,
		"stat":               e.CmdStat,
		"utoa":               e.CmdUtoA,
		"which":              e.CmdWhich,
	}
}

func (e *Env) Table2() map[string]func(*Param) []any {
	return map[string]func(*Param) []any{
		"box":            e.CmdBox,
		"raweval":        e.CmdRawEval,
		"rawexec":        e.CmdRawExec,
		"write":          e.CmdWrite,
		"writerr":        e.CmdWriteErr,
		"default_prompt": e.Prompt,
	}
}
