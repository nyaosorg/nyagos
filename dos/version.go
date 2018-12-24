package dos

import (
	"syscall"
	"unsafe"
)

const (
	verBuildNumber      = 0x0000004
	verMajorVersion     = 0x0000002
	verMinorVersion     = 0x0000001
	verPlatformId       = 0x0000008
	verProductType      = 0x0000080
	verServicePackMajor = 0x0000020
	verServicePackMinor = 0x0000010
	verSuiteName        = 0x0000040

	verEqual        = 1
	verGreater      = 2
	verGreaterEqual = 3
	verLess         = 4
	verLessEqual    = 5

	errorOldWinVersion syscall.Errno = 1150
)

type OSVersionInfoEx struct {
	osVersionInfoSize uint32
	MajorVersion      uint32
	MinorVersion      uint32
	buildNumber       uint32
	platformId        uint32
	csdVersion        [128]uint16
	servicePackMajor  uint16
	servicePackMinor  uint16
	SuiteMask         uint16
	productType       byte
	reserve           byte
}

var (
	Windows10    = &OSVersionInfoEx{MajorVersion: 10}
	Windows81    = &OSVersionInfoEx{MajorVersion: 6, MinorVersion: 3}
	Windows8     = &OSVersionInfoEx{MajorVersion: 6, MinorVersion: 2}
	Windows7     = &OSVersionInfoEx{MajorVersion: 6, MinorVersion: 1}
	WindowsVista = &OSVersionInfoEx{MajorVersion: 6}
)

var procVerSetConditionMask = kernel32.NewProc("VerSetConditionMask")

var procVerifyVersionInfo = kernel32.NewProc("VerifyVersionInfoW")

func (vi OSVersionInfoEx) Verify() bool {
	var m1, m2 uintptr

	var mask uintptr = verMajorVersion
	m1, m2, _ = procVerSetConditionMask.Call(m1, m2, verMajorVersion, verGreaterEqual)
	if vi.MinorVersion > 0 {
		m1, m2, _ = procVerSetConditionMask.Call(m1, m2, verMinorVersion, verGreaterEqual)
		mask |= verMinorVersion
	}

	vi.osVersionInfoSize = uint32(unsafe.Sizeof(vi))
	r, _, _ := procVerifyVersionInfo.Call(
		uintptr(unsafe.Pointer(&vi)),
		mask,
		m1, m2)
	return r != 0
}
