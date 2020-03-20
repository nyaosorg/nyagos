// +build !windows

package commands

import (
	"context"
	"fmt"
	"io"
	"os"
)

func init() {
	buildInCommand = map[string]func(context.Context, Param) (int, error){
		".":        cmdSource,
		"alias":    cmdAlias,
		"bindkey":  cmdBindkey,
		"box":      cmdBox,
		"cd":       cmdCd,
		"clip":     cmdClip,
		"cls":      cmdCls,
		"cmdexesc": cmdExeSc,
		"chmod":    cmdChmod,
		"copy":     cmdCopy,
		"del":      cmdDel,
		"dirs":     cmdDirs,
		"diskused": cmdDiskUsed,
		"echo":     cmdEcho,
		"env":      cmdEnv,
		"exit":     cmdExit,
		"foreach":  cmdForeach,
		"history":  cmdHistory,
		"if":       cmdIf,
		"ln":       cmdLn,
		"kill":     cmdKill,
		"killall":  cmdKillAll,
		"md":       cmdMkdir,
		"mkdir":    cmdMkdir,
		"more":     cmdMore,
		"move":     cmdMove,
		"popd":     cmdPopd,
		"ps":       cmdPs,
		"pushd":    cmdPushd,
		"pwd":      cmdPwd,
		"rd":       cmdRmdir,
		"rem":      cmdRem,
		"rmdir":    cmdRmdir,
		"select":   cmdShOpenWithDialog,
		"set":      cmdSet,
		"source":   cmdSource,
		"touch":    cmdTouch,
		"type":     cmdType,
		"which":    cmdWhich,
	}
}

func newMbcsReader(r io.Reader) io.Reader {
	return r
}

func setWritable(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("os.Stat('%s'): %s", path, err)
	}
	mode := stat.Mode() | 0x600
	return os.Chmod(path, mode)
}
