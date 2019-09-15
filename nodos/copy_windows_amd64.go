package nodos

import (
	"unsafe"
)

func progressPrintCallBack(total, transfer, c, d, e, f, g, h, this uintptr) uintptr {
	progressPrint(uint64(total), uint64(transfer), (*progressCopy)(unsafe.Pointer(this)))
	return 0
}
