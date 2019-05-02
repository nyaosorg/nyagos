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
	hProcess      uintptr
}

var shell32 = windows.NewLazySystemDLL("shell32.dll")
var procShellExecute = shell32.NewProc("ShellExecuteExW")

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
	_SEE_MASK_UNICODE = 0x4000
)

// ShellExecute calls ShellExecute-API: edit,explore,open and so on.
func shellExecute(action string, path string, param string, directory string) (err error) {
	var p _ShellExecuteInfo

	p.size = uint32(unsafe.Sizeof(p))

	p.mask = _SEE_MASK_UNICODE

	p.verb, err = windows.UTF16PtrFromString(action)
	if err != nil {
		return err
	}
	p.file, err = windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	p.parameter, err = windows.UTF16PtrFromString(param)
	if err != nil {
		return err
	}
	p.directory, err = windows.UTF16PtrFromString(directory)
	if err != nil {
		return err
	}

	p.show = 1

	status, _, err := procShellExecute.Call(uintptr(unsafe.Pointer(&p)))

	if status == 0 {
		// ShellExecute and ShellExecuteExA's error is lower than 32
		// But, ShellExecuteExW's error is FALSE.

		if err != nil {
			return err
		} else if err = windows.GetLastError(); err != nil {
			return err
		} else {
			return fmt.Errorf("Error(%d) in ShellExecuteExW()", status)
		}
	}
	return nil
}

const haveToEvalSymlinkError = windows.Errno(4294967294)

func ShellExecute(action string, path string, param string, directory string) error {
	err := shellExecute(action, path, param, directory)
	if err == haveToEvalSymlinkError {
		path, err = filepath.EvalSymlinks(path)
		if err == nil {
			err = shellExecute(action, path, param, directory)
		}
	}
	return err
}
