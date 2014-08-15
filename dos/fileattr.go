package dos

//#include <windows.h>
import "C"
import "syscall"

type FileAttr struct {
	attr uint
}

func NewFileAttr(path string) *FileAttr {
	cpath, _ := syscall.UTF16FromString(path)
	return &FileAttr{uint(C.GetFileAttributesW((*C.WCHAR)(&cpath[0])))}
}

func (this *FileAttr) IsReparse() bool {
	return (this.attr & C.FILE_ATTRIBUTE_REPARSE_POINT) != 0
}

func (this *FileAttr) IsHidden() bool {
	return (this.attr & C.FILE_ATTRIBUTE_HIDDEN) != 0
}
