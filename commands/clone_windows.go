package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	var pid int
	pid, err = dos.ShellExecute(action, me, "", wd)
	if err != nil {
		pid, err = dos.ShellExecute(action, "CMD.EXE", "/c \""+me+"\"", wd)
		if err != nil {
			return 1, err // return original error
		}
	}
	if pid > 0 {
		fmt.Fprintf(out, "[%d]\n", pid)
		if process, err := os.FindProcess(pid); err == nil {
			go func() {
				process.Wait()
				fmt.Fprintf(os.Stderr, "[%d]+ Done\n", pid)
			}()
		}
	}
	return 0, nil
}

func cmdClone(ctx context.Context, cmd Param) (int, error) {
	return _clone("open", cmd.Err())
}

func cmdSu(ctx context.Context, cmd Param) (int, error) {
	me, err := os.Executable()
	if err != nil {
		return 1, err
	}
	if me2, err2 := filepath.EvalSymlinks(me); err2 == nil {
		me = me2
	}
	wd, err := os.Getwd()
	if err != nil {
		return 2, err
	}
	if strings.HasSuffix(wd, `\`) {
		wd += "."
	}

	var buffer strings.Builder

	if netdrives, err := dos.GetNetDrives(); err == nil {
		for _, n := range netdrives {
			fmt.Fprintf(&buffer, ` --netuse "%c:=%s"`, n.Letter, n.Remote)
		}
	}
	fmt.Fprintf(&buffer, ` --chdir "%s"`, wd)

	pid, err := dos.ShellExecute("runas", me, buffer.String(), "")
	if err != nil {
		return 3, err
	}
	if pid > 0 {
		fmt.Fprintf(cmd.Err(), "[%d]\n", pid)
	}
	return 0, nil
}
