package nodos

import (
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/zetamatta/nyagos/dos"
)

var kernel32 = windows.NewLazySystemDLL("kernel32.dll")

var procCopyFileW = kernel32.NewProc("CopyFileExW")

type progressCopy struct {
	last time.Time
	run  bool
}

func keta(n uintptr) int {
	if n < 10 {
		return 1
	}
	return keta(n/10) + 1
}

func progressPrint(total, transfer, c, d, e, f, g, h, _this uintptr) uintptr {
	this := (*progressCopy)(unsafe.Pointer(_this))
	now := time.Now()

	if now.Sub(this.last) >= time.Second {
		fmt.Printf("%3d%% %*d/%d\r",
			transfer*100/total,
			keta(total),
			transfer,
			total)

		this.last = now
		this.run = true
	}
	return 0
}

// Copy calls Win32's CopyFile API.
func copyFile(src, dst string, isFailIfExists bool) error {
	_src, err := windows.UTF16PtrFromString(src)
	if err != nil {
		return err
	}
	_dst, err := windows.UTF16PtrFromString(dst)
	if err != nil {
		return err
	}
	var flag uintptr
	if isFailIfExists {
		flag |= 1
	}
	var progressCopy1 progressCopy
	progressCopy1.last = time.Now()

	var cancel uintptr

	rc, _, err := procCopyFileW.Call(
		uintptr(unsafe.Pointer(_src)),
		uintptr(unsafe.Pointer(_dst)),
		windows.NewCallbackCDecl(progressPrint),
		uintptr(unsafe.Pointer(&progressCopy1)),
		uintptr(unsafe.Pointer(&cancel)),
		flag)

	if progressCopy1.run {
		fmt.Println()
	}
	if rc == 0 {
		return err
	}
	return nil
}

// Move calls Win32's MoveFileEx API.
func moveFile(src, dst string) error {
	_src, err := windows.UTF16PtrFromString(src)
	if err != nil {
		return err
	}
	_dst, err := windows.UTF16PtrFromString(dst)
	if err != nil {
		return err
	}
	return windows.MoveFileEx(
		_src,
		_dst,
		windows.MOVEFILE_REPLACE_EXISTING|
			windows.MOVEFILE_COPY_ALLOWED|
			windows.MOVEFILE_WRITE_THROUGH)
}

func readShortcut(path string) (string, string, error) {
	return dos.ReadShortcut(path)
}
