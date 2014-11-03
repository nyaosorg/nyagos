package dos

import "fmt"
import "unsafe"
import "syscall"

var getFileAttributes = kernel32.NewProc("GetFileAttributesW")

const FILE_ATTRIBUTE_REPARSE_POINT = 0x00000400
const FILE_ATTRIBUTE_HIDDEN = 0x00000002

type FileAttr struct {
	attr uintptr
}

func NewFileAttr(path string) (*FileAttr, error) {
	cpath, err := syscall.UTF16FromString(path)
	if err != nil {
		return &FileAttr{0}, err
	} else if cpath == nil {
		return &FileAttr{0}, fmt.Errorf("sysCall.UTF16FromString(\"%s\") failed", path)
	} else {
		rc, _, _ := getFileAttributes.Call(uintptr(unsafe.Pointer(&cpath[0])))
		return &FileAttr{rc}, nil
	}
}

func (this *FileAttr) IsReparse() bool {
	return (this.attr & FILE_ATTRIBUTE_REPARSE_POINT) != 0
}

func (this *FileAttr) IsHidden() bool {
	return (this.attr & FILE_ATTRIBUTE_HIDDEN) != 0
}
