// +build !windows

package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
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

func setWritable(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("os.Stat('%s'): %s", path, err)
	}
	mode := stat.Mode() | 0x600
	return os.Chmod(path, mode)
}

func truncate(path string, _ func(path string, err error) bool, _ io.Writer) error {
	return os.RemoveAll(path)
}
