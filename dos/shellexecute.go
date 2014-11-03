package dos

import "syscall"
import "unsafe"

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

func ShellExecute(action string, path string, param string, directory string) error {
	actionW, actionErr := syscall.UTF16FromString(action)
	if actionErr != nil {
		return actionErr
	}
	pathW, pathErr := syscall.UTF16FromString(path)
	if pathErr != nil {
		return pathErr
	}
	paramW, paramErr := syscall.UTF16FromString(param)
	if paramErr != nil {
		return paramErr
	}
	directoryW, directoryErr := syscall.UTF16FromString(directory)
	if directoryErr != nil {
		return directoryErr
	}
	status, _, _ := shellExecute.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(&actionW[0])),
		uintptr(unsafe.Pointer(&pathW[0])),
		uintptr(unsafe.Pointer(&paramW[0])),
		uintptr(unsafe.Pointer(&directoryW[0])),
		SW_SHOWNORMAL)
	if status < 32 {
		return syscall.GetLastError()
	}
	return nil
}
