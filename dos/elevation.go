package dos

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type tokenElevationT struct {
	TokenIsElevated uint32
}

// IsElevated returns true if the current process runs as Administrator
func IsElevated() (bool, error) {
	var hToken windows.Token

	currentProcess, err := windows.GetCurrentProcess()
	if err != nil {
		return false, err
	}

	err = windows.OpenProcessToken(currentProcess,
		windows.TOKEN_QUERY,
		&hToken)
	if err != nil {
		return false, err
	}

	var tokenElevation tokenElevationT
	dwSize := uint32(unsafe.Sizeof(tokenElevation))

	err = windows.GetTokenInformation(hToken,
		windows.TokenElevation,
		(*byte)(unsafe.Pointer(&tokenElevation)),
		dwSize,
		&dwSize)
	return tokenElevation.TokenIsElevated != 0, err
}
