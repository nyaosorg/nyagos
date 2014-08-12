package fileattr

//#include <windows.h>
import "C"
import "unsafe"

type FileAttr struct {
	attr uint
}

func New(path string) *FileAttr {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return &FileAttr{uint(C.GetFileAttributes((*C.CHAR)(cpath)))}
}

func (this *FileAttr) IsReparse() bool {
	return (this.attr & C.FILE_ATTRIBUTE_REPARSE_POINT) != 0
}

func (this *FileAttr) IsHidden() bool {
	return (this.attr & C.FILE_ATTRIBUTE_HIDDEN) != 0
}
