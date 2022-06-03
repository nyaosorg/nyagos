package nodos

import (
	"unsafe"
)

func progressPrintCallBack(totalL, totalH, transferL, transferH, c1, c2, d1, d2, e, f, g, h, this uintptr) uintptr {
	progressPrint(uint64(totalL)|(uint64(totalH)<<32),
		uint64(transferL)|(uint64(transferH)<<32),
		(*progressCopy)(unsafe.Pointer(this)))
	return 0
}
