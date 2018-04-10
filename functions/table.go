package functions

var Table = map[string]func([]interface{}) []interface{}{
	"access":         CmdAccess,
	"atou":           CmdAtoU,
	"box":            CmdBox,
	"chdir":          CmdChdir,
	"commonprefix":   CmdCommonPrefix,
	"default_prompt": Prompt,
	"elevated":       CmdElevated,
	"getenv":         CmdGetEnv,
	"gethistory":     CmdGetHistory,
	"getkey":         CmdGetKey,
	"getviewwidth":   CmdGetViewWidth,
	"getwd":          CmdGetwd,
	"glob":           CmdGlob,
	"msgbox":         CmdMsgBox,
	"netdrivetounc":  CmdNetDriveToUNC,
	"pathjoin":       CmdPathJoin,
	"raweval":        CmdRawEval,
	"resetcharwidth": CmdResetCharWidth,
	"setenv":         CmdSetEnv,
	"setrunewidth":   CmdSetRuneWidth,
	"shellexecute":   CmdShellExecute,
	"stat":           CmdStat,
	"utoa":           CmdUtoA,
	"which":          CmdWhich,
}

var Table2 = map[string]func(*Param) []interface{}{
	"rawexec": CmdRawExec,
	"write":   CmdWrite,
	"writerr": CmdWriteErr,
}
