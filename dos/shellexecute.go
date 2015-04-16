package dos

import (
	"fmt"
	"syscall"
	"unsafe"
)

var shell32 = syscall.NewLazyDLL("shell32")
var shellExecute = shell32.NewProc("ShellExecuteW")

const (
	SW_HIDE            = 0
	SW_MAXIMIZE        = 3
	SW_MINIMIZE        = 6
	SW_RESTORE         = 9
	SW_SHOW            = 5
	SW_SHOWDEFAULT     = 1
	SW_SHOWMAXIMIZED   = 3
	SW_SHOWMINIMIZED   = 2
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8
	SW_SHOWNOACTIVATE  = 4
	SW_SHOWNORMAL      = 1
)

const (
	EDIT       = "edit"
	EXPLORE    = "explore"
	OPEN       = "open"
	PRINT      = "print"
	PROPERTIES = "properties"
	RUNAS      = "runas"
)

// Call ShellExecute-API: edit,explore,open and so on.
func ShellExecute(action string, path string, param string, directory string) error {
	actionP, actionErr := syscall.UTF16PtrFromString(action)
	if actionErr != nil {
		return actionErr
	}
	pathP, pathErr := syscall.UTF16PtrFromString(path)
	if pathErr != nil {
		return pathErr
	}
	paramP, paramErr := syscall.UTF16PtrFromString(param)
	if paramErr != nil {
		return paramErr
	}
	directoryP, directoryErr := syscall.UTF16PtrFromString(directory)
	if directoryErr != nil {
		return directoryErr
	}
	status, _, _ := shellExecute.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(actionP)),
		uintptr(unsafe.Pointer(pathP)),
		uintptr(unsafe.Pointer(paramP)),
		uintptr(unsafe.Pointer(directoryP)),
		SW_SHOWNORMAL)

	if status <= 32 {
		if err := syscall.GetLastError(); err != nil {
			return err
		} else {
			return fmt.Errorf("Error(%d) in ShellExecuteW()", status)
		}
	}
	return nil
}
