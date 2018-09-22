package dos

import (
	"fmt"
	"syscall"
	"unsafe"
)

var advapi32 = syscall.NewLazyDLL("advapi32.dll")
var procOpenProcessToken = advapi32.NewProc("OpenProcessToken")
var procGetTokenInformation = advapi32.NewProc("GetTokenInformation")
var procGetCurrentProcess = kernel32.NewProc("GetCurrentProcess")

const ( // from winnt.h
	// cTokenElevationType = 18
	cTokenElevation = 20

	cTokenQuery = 8
)

type tokenElevationT struct {
	TokenIsElevated uint32
}

// IsElevated returns true if the current process runs as Administrator
func IsElevated() (bool, error) {
	var hToken uintptr

	currentProcess, _, _ := procGetCurrentProcess.Call()

	rc, _, err := procOpenProcessToken.Call(uintptr(currentProcess),
		uintptr(cTokenQuery),
		uintptr(unsafe.Pointer(&hToken)))
	if rc == 0 {
		return false, fmt.Errorf("OpenProcessToken: %s", err.Error())
	}

	var tokenElevation tokenElevationT
	dwSize := unsafe.Sizeof(tokenElevation)

	rc, _, err = procGetTokenInformation.Call(uintptr(hToken),
		uintptr(cTokenElevation),
		uintptr(unsafe.Pointer(&tokenElevation)),
		uintptr(dwSize),
		uintptr(unsafe.Pointer(&dwSize)))
	if rc == 0 {
		return false, fmt.Errorf("GetTokenInformation: %s", err.Error())
	}
	return tokenElevation.TokenIsElevated != 0, nil
}
