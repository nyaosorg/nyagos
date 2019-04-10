package dos

import (
	"fmt"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

var shell32 = windows.NewLazySystemDLL("shell32.dll")
var procShellExecute = shell32.NewProc("ShellExecuteW")

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

// ShellExecute calls ShellExecute-API: edit,explore,open and so on.
func shellExecute(action string, path string, param string, directory string) error {
	actionP, actionErr := windows.UTF16PtrFromString(action)
	if actionErr != nil {
		return actionErr
	}
	pathP, pathErr := windows.UTF16PtrFromString(path)
	if pathErr != nil {
		return pathErr
	}
	paramP, paramErr := windows.UTF16PtrFromString(param)
	if paramErr != nil {
		return paramErr
	}
	directoryP, directoryErr := windows.UTF16PtrFromString(directory)
	if directoryErr != nil {
		return directoryErr
	}
	status, _, err := procShellExecute.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(actionP)),
		uintptr(unsafe.Pointer(pathP)),
		uintptr(unsafe.Pointer(paramP)),
		uintptr(unsafe.Pointer(directoryP)),
		windows.SW_SHOWNORMAL)

	if status <= 32 {
		if err != nil {
			return err
		} else if err = windows.GetLastError(); err != nil {
			return err
		} else {
			return fmt.Errorf("Error(%d) in ShellExecuteW()", status)
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
