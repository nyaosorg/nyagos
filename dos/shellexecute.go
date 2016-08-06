package dos

import (
	"fmt"
	"syscall"
	"unsafe"
)

var shell32 = syscall.NewLazyDLL("shell32")
var shellExecute = shell32.NewProc("ShellExecuteW")

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
	status, _, err := shellExecute.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(actionP)),
		uintptr(unsafe.Pointer(pathP)),
		uintptr(unsafe.Pointer(paramP)),
		uintptr(unsafe.Pointer(directoryP)),
		SW_SHOWNORMAL)

	if status <= 32 {
		if err != nil {
			return err
		} else if err = syscall.GetLastError(); err != nil {
			return err
		} else {
			return fmt.Errorf("Error(%d) in ShellExecuteW()", status)
		}
	}
	return nil
}
