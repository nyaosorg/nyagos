package commands

import (
	"context"

	"golang.org/x/sys/windows"

	"github.com/nyaosorg/nyagos/internal/nodos"
)

func init() {
	data := map[string]func(context.Context, Param) (int, error){
		".":        cmdSource,
		"alias":    cmdAlias,
		"attrib":   cmdAttrib,
		"bindkey":  cmdBindkey,
		"box":      cmdBox,
		"cd":       cmdCd,
		"clip":     cmdClip,
		"clone":    cmdClone,
		"cls":      cmdCls,
		"cmdexesc": cmdExeSc,
		"chmod":    cmdChmod,
		"copy":     cmdCopy,
		"del":      cmdDel,
		"dirs":     cmdDirs,
		"diskfree": cmdDiskFree,
		"diskused": cmdDiskUsed,
		"echo":     cmdEcho,
		"env":      cmdEnv,
		"erase":    cmdDel,
		"exit":     cmdExit,
		"foreach":  cmdForeach,
		"history":  cmdHistory,
		"if":       cmdIf,
		"ln":       cmdLn,
		"lnk":      cmdLnk,
		"mklink":   cmdMklink,
		"kill":     cmdKill,
		"killall":  cmdKillAll,
		"ls":       cmdLs,
		"md":       cmdMkdir,
		"mkdir":    cmdMkdir,
		"more":     cmdMore,
		"move":     cmdMove,
		"open":     cmdOpen,
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
		"su":       cmdSu,
		"touch":    cmdTouch,
		"type":     cmdType,
		"which":    cmdWhich,
	}
	for key, val := range data {
		buildInCommand.Store(key, val)
	}
}

func setWritable(path string) error {
	perm, err := nodos.GetFileAttributes(path)
	if err != nil {
		return err
	}
	return nodos.SetFileAttributes(path, perm&^windows.FILE_ATTRIBUTE_READONLY)
}
