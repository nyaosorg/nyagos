package dos

import "syscall"

func GetFileAttributes(path string) (uint32, error) {
	cpath, cpathErr := syscall.UTF16PtrFromString(path)
	if cpathErr != nil {
		return 0, cpathErr
	}
	return syscall.GetFileAttributes(cpath)
}

func SetFileAttributes(path string, attr uint32) error {
	cpath, cpathErr := syscall.UTF16PtrFromString(path)
	if cpathErr != nil {
		return cpathErr
	}
	return syscall.SetFileAttributes(cpath, attr)
}

// Windows original attribute.
type FileAttr struct {
	attr uint32
}

// Get Windows original attributes such as Hidden,Reparse and so on.
func NewFileAttr(path string) (*FileAttr, error) {
	attr, attrErr := GetFileAttributes(path)
	return &FileAttr{attr}, attrErr
}

func (this *FileAttr) IsReparse() bool {
	return (this.attr & FILE_ATTRIBUTE_REPARSE_POINT) != 0
}

func (this *FileAttr) IsHidden() bool {
	return (this.attr & FILE_ATTRIBUTE_HIDDEN) != 0
}
