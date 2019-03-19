package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/zetamatta/nyagos/dos"
)

func _getwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}

func _clone(action string, out io.Writer) (int, error) {
	wd := _getwd()
	var err error
	var me string
	me, err = os.Executable()
	if err != nil {
		return 1, err
	}
	err = dos.ShellExecute(action, me, "", wd)
	if err != nil {
		err = dos.ShellExecute(action, "CMD.EXE", "/c \""+me+"\"", wd)
		if err != nil {
			return 1, err // return original error
		}
	}
	return 0, nil
}

func cmdClone(ctx context.Context, cmd Param) (int, error) {
	return _clone("open", cmd.Err())
}

func cmdSu(ctx context.Context, cmd Param) (int, error) {
	netdrives, err := dos.GetNetDrives()
	if err == nil {
		var buffer strings.Builder
		buffer.WriteString("/S /C \"")
		for _, n := range netdrives {
			fmt.Fprintf(&buffer, "net use %c: \"%s\" 1>nul 2>nul & ", n.Letter, n.Remote)
		}
		me, err := os.Executable()
		if err != nil {
			return 1, err
		}
		wd, err := os.Getwd()
		if err != nil {
			return 2, err
		}
		fmt.Fprintf(&buffer, " cd /d \"%s\" & \"%s\" \"", wd, me)
		err = dos.ShellExecute("runas", "CMD.EXE", buffer.String(), "")
		if err != nil {
			return 3, err
		}
		return 0, nil
	}
	return _clone("runas", cmd.Err())
}
