package functions

var Table = map[string]func([]interface{}) []interface{}{
	"access":         CmdAccess,
	"atou":           CmdAtoU,
	"bitand":         CmdBitAnd,
	"bitor":          CmdBitOr,
	"chdir":          CmdChdir,
	"commonprefix":   CmdCommonPrefix,
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
	"resetcharwidth": CmdResetCharWidth,
	"setenv":         CmdSetEnv,
	"setrunewidth":   CmdSetRuneWidth,
	"shellexecute":   CmdShellExecute,
	"stat":           CmdStat,
	"utoa":           CmdUtoA,
	"which":          CmdWhich,
}

var Table2 = map[string]func(*Param) []interface{}{
	"box":            CmdBox,
	"raweval":        CmdRawEval,
	"rawexec":        CmdRawExec,
	"write":          CmdWrite,
	"writerr":        CmdWriteErr,
	"default_prompt": Prompt,
}
