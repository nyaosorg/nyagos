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
	TokenElevationType = 18
	TokenElevation     = 20

	TOKEN_QUERY = 8
)

type token_elevation_t struct {
	TokenIsElevated uint32
}

func IsElevated() (bool, error) {
	var hToken uintptr

	currentProcess, _, _ := procGetCurrentProcess.Call()

	rc, _, err := procOpenProcessToken.Call(uintptr(currentProcess),
		uintptr(TOKEN_QUERY),
		uintptr(unsafe.Pointer(&hToken)))
	if rc == 0 {
		return false, fmt.Errorf("OpenProcessToken: %s", err.Error())
	}

	var token_elevation token_elevation_t
	var dwSize uintptr = unsafe.Sizeof(token_elevation)

	rc, _, err = procGetTokenInformation.Call(uintptr(hToken),
		uintptr(TokenElevation),
		uintptr(unsafe.Pointer(&token_elevation)),
		uintptr(dwSize),
		uintptr(unsafe.Pointer(&dwSize)))
	if rc == 0 {
		return false, fmt.Errorf("GetTokenInformation: %s", err.Error())
	}
	return token_elevation.TokenIsElevated != 0, nil
}
