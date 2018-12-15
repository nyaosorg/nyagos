// +build !windows

package commands

import (
	"context"
	"errors"
	"io"
)

func init() {
	buildInCommand = map[string]func(context.Context, Param) (int, error){
		"alias":    cmdAlias,
		"bindkey":  cmdBindkey,
		"box":      cmdBox,
		"cd":       cmdCd,
		"clip":     cmdClip,
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
		"more":     cmdMore,
		"popd":     cmdPopd,
		"ps":       cmdPs,
		"pushd":    cmdPushd,
		"pwd":      cmdPwd,
		"rem":      cmdRem,
		"set":      cmdSet,
		"touch":    cmdTouch,
		"type":     cmdType,
		"which":    cmdWhich,
	}
}

func newMbcsReader(r io.Reader) io.Reader {
	return r
}

func readShortCut(dir string) (string, error) {
	return "", errors.New("not support shortcut")
}
