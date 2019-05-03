package dos

import (
	"fmt"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ShellExecuteInfo struct {
	size          uint32
	mask          uint32
	hwnd          uintptr
	verb          *uint16
	file          *uint16
	parameter     *uint16
	directory     *uint16
	show          int
	instApp       uintptr
	idList        uintptr
	class         *uint16
	keyClass      uintptr
	hotkey        uint32
	iconOrMonitor uintptr
	hProcess      windows.Handle
}

var shell32 = windows.NewLazySystemDLL("shell32.dll")
var procShellExecute = shell32.NewProc("ShellExecuteExW")
var procGetProcessId = kernel32.NewProc("GetProcessId")

const (
	// EDIT is the action "edit" for ShellExecute
	EDIT = "edit"
	// EXPLORE is the action "explore" for ShellExecute
	EXPLORE = "explore"
	// OPEN is the action "open" for ShellExecute
	OPEN = "open"
	// PRINT is the action "print" for ShellExecute
	PRINT = "print"
	// PROPERTIES is the action "properties" for ShellExecute
	PROPERTIES = "properties"
	// RUNAS is the action "runas" for ShellExecute
	RUNAS = "runas"
)

const (
	_SEE_MASK_NOCLOSEPROCESS = 0x40
	_SEE_MASK_UNICODE        = 0x4000
)

// ShellExecute calls ShellExecute-API: edit,explore,open and so on.
func shellExecute(action, path, param, directory string) (pid uintptr, err error) {
	var p _ShellExecuteInfo

	p.size = uint32(unsafe.Sizeof(p))

	p.mask = _SEE_MASK_UNICODE | _SEE_MASK_NOCLOSEPROCESS

	p.verb, err = windows.UTF16PtrFromString(action)
	if err != nil {
		return 0, err
	}
	p.file, err = windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}
	p.parameter, err = windows.UTF16PtrFromString(param)
	if err != nil {
		return 0, err
	}
	p.directory, err = windows.UTF16PtrFromString(directory)
	if err != nil {
		return 0, err
	}

	p.show = 1
	status, _, err := procShellExecute.Call(uintptr(unsafe.Pointer(&p)))

	if p.hProcess != 0 {
		pid, _, _ = procGetProcessId.Call(uintptr(p.hProcess))
		if err := windows.CloseHandle(p.hProcess); err != nil {
			println("windows.Closehandle()=", err.Error())
		}
	}

	if status == 0 {
		// ShellExecute and ShellExecuteExA's error is lower than 32
		// But, ShellExecuteExW's error is FALSE.

		if err != nil {
			return pid, err
		} else if err = windows.GetLastError(); err != nil {
			return pid, err
		} else {
			return pid, fmt.Errorf("Error(%d) in ShellExecuteExW()", status)
		}
	}
	return pid, nil
}

const haveToEvalSymlinkError = windows.Errno(4294967294)

func ShellExecute(action string, path string, param string, directory string) (uintptr, error) {
	pid, err := shellExecute(action, path, param, directory)
	if err == haveToEvalSymlinkError {
		path, err = filepath.EvalSymlinks(path)
		if err == nil {
			pid, err = shellExecute(action, path, param, directory)
		}
	}
	return pid, err
}
