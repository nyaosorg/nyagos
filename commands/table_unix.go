// +build !windows

package commands

import "context"

func init() {
	buildInCommand = map[string]func(context.Context, Param) (int, error){
		"alias":    cmdAlias,
		"bindkey":  cmdBindkey,
		"box":      cmdBox,
		"cd":       cmdCd,
		"cls":      cmdCls,
		"chmod":    cmdChmod,
		"dirs":     cmdDirs,
		"diskused": cmdDiskUsed,
		"echo":     cmdEcho,
		"env":      cmdEnv,
		"exit":     cmdExit,
		"foreach":  cmdForeach,
		"history":  cmdHistory,
		"if":       cmdIf,
		"kill":     cmdKill,
		"popd":     cmdPopd,
		"ps":       cmdPs,
		"pushd":    cmdPushd,
		"pwd":      cmdPwd,
		"rem":      cmdRem,
		"set":      cmdSet,
		"touch":    cmdTouch,
		"which":    cmdWhich,
	}
}
