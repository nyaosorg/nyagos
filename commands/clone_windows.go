package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"

	"github.com/nyaosorg/go-windows-netresource"
	"github.com/nyaosorg/go-windows-su"
	"github.com/nyaosorg/go-windows-subst"

	"github.com/nyaosorg/nyagos/internal/shell"
)

var isWindowsTerminal = os.Getenv("WT_SESSION") != "" && os.Getenv("WT_PROFILE_ID") != ""

var wtExePath = filepath.Join(os.Getenv("LOCALAPPDATA"), `Microsoft\WindowsApps\wt.exe`)

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
	if _me, err := filepath.EvalSymlinks(me); err == nil {
		me = _me
	}
	if isWindowsTerminal {
		pid, err = su.ShellExecute(action, wtExePath,
			fmt.Sprintf(`--window 0 new-tab -- "%s" --chdir "%s" "`, me, wd), wd)
	} else {
		pid, err = su.ShellExecute(action, me, "", wd)
	}
	if err != nil {
		pid, err = su.ShellExecute(action, "CMD.EXE", "/c \""+me+"\"", wd)
		if err != nil {
			return 1, err // return original error
		}
	}
	if pid > 0 {
		fmt.Fprintf(out, "[%d]\n", pid)
		if process, err := os.FindProcess(pid); err == nil {
			go func() {
				process.Wait()
				shell.Message("[%d]+ Done\n", pid)
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

	if netdrives, err := netresource.GetNetDrives(); err == nil {
		for _, n := range netdrives {
			fmt.Fprintf(&buffer, ` --netuse "%c:=%s"`, n.Letter, n.Remote)
		}
	}
	if drives, err := netresource.GetDrives(); err == nil {
		for _, d := range drives {
			if d.Type == windows.DRIVE_FIXED {
				mountPt := string([]byte{byte(d.Letter), ':'})
				target, err := subst.QueryRaw(mountPt)
				if err == nil && target[:4] == `\??\` {
					fmt.Fprintf(&buffer, ` --subst "%c:=%s"`, d.Letter, target[4:])
				}

			}
		}
	}
	fmt.Fprintf(&buffer, ` --chdir "%s"`, wd)

	var pid int

	if isWindowsTerminal {
		pid, err = su.ShellExecute("runas", wtExePath,
			fmt.Sprintf(`--window 0 new-tab -- "%s" %s`, me, buffer.String()), "")
	} else {
		pid, err = su.ShellExecute("runas", me, buffer.String(), "")
	}
	if err != nil {
		return 3, err
	}
	if pid > 0 {
		fmt.Fprintf(cmd.Err(), "[%d]\n", pid)
	}
	return 0, nil
}
