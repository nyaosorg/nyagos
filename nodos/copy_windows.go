package nodos

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/zetamatta/nyagos/dos"
)

var kernel32 = windows.NewLazySystemDLL("kernel32.dll")

var procCopyFileW = kernel32.NewProc("CopyFileExW")

type progressCopy struct {
	last time.Time
	n    int
	run  bool
}

func keta(n uint64) int {
	if n < 10 {
		return 1
	}
	return keta(n/10) + 1
}

func progressPrint(total, transfer uint64, this *progressCopy) {
	now := time.Now()

	if now.Sub(this.last) >= time.Second {
		this.n, _ = fmt.Printf("%3d%% %*d/%d\r",
			transfer*100/total,
			keta(total),
			transfer,
			total)

		this.last = now
		this.run = true
	}
}

func progressPrint64bit(total, transfer, c, d, e, f, g, h, this uintptr) uintptr {
	progressPrint(uint64(total), uint64(transfer), (*progressCopy)(unsafe.Pointer(this)))
	return 0
}

func progressPrint32bit(totalL, totalH, transferL, transferH, c1, c2, d1, d2, e, f, g, h, this uintptr) uintptr {
	progressPrint(uint64(totalL)|(uint64(totalH)<<32),
		uint64(transferL)|(uint64(transferH)<<32),
		(*progressCopy)(unsafe.Pointer(this)))
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
	var rc uintptr

	if runtime.GOARCH == "386" {
		rc, _, err = procCopyFileW.Call(
			uintptr(unsafe.Pointer(_src)),
			uintptr(unsafe.Pointer(_dst)),
			windows.NewCallbackCDecl(progressPrint32bit),
			uintptr(unsafe.Pointer(&progressCopy1)),
			uintptr(unsafe.Pointer(&cancel)),
			flag)
	} else {
		rc, _, err = procCopyFileW.Call(
			uintptr(unsafe.Pointer(_src)),
			uintptr(unsafe.Pointer(_dst)),
			windows.NewCallbackCDecl(progressPrint64bit),
			uintptr(unsafe.Pointer(&progressCopy1)),
			uintptr(unsafe.Pointer(&cancel)),
			flag)
	}

	if progressCopy1.run {
		fmt.Printf("%*s\r", progressCopy1.n, "")
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
