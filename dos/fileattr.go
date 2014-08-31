package dos

//#include <windows.h>
import "C"
import "fmt"
import "syscall"

type FileAttr struct {
	attr uint
}

func NewFileAttr(path string) (*FileAttr, error) {
	cpath, err := syscall.UTF16FromString(path)
	if err != nil {
		return &FileAttr{0}, err
	} else if cpath == nil {
		return &FileAttr{0}, fmt.Errorf("sysCall.UTF16FromString() failed")
	} else {
		return &FileAttr{uint(C.GetFileAttributesW((*C.WCHAR)(&cpath[0])))}, nil
	}
}

func (this *FileAttr) IsReparse() bool {
	return (this.attr & C.FILE_ATTRIBUTE_REPARSE_POINT) != 0
}

func (this *FileAttr) IsHidden() bool {
	return (this.attr & C.FILE_ATTRIBUTE_HIDDEN) != 0
}
