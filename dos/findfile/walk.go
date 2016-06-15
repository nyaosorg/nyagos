package findfile

import (
	"os"
	"syscall"
	"time"
)

type FileInfo struct {
	syscall.Win32finddata
	handle syscall.Handle
}

func (this *FileInfo) clone() *FileInfo {
	return &FileInfo{this.Win32finddata, this.handle}
}

func findFirst(pattern string) (*FileInfo, error) {
	pattern16, err := syscall.UTF16PtrFromString(pattern)
	if err != nil {
		return nil, err
	}
	this := new(FileInfo)
	this.handle, err = syscall.FindFirstFile(pattern16, &this.Win32finddata)
	if err != nil {
		return nil, err
	}
	return this, nil
}

func (this *FileInfo) Name() string {
	return syscall.UTF16ToString(this.FileName[:])
}

func (this *FileInfo) Size() int64 {
	return int64((this.FileSizeHigh << 32) | this.FileSizeLow)
}

func (this *FileInfo) ModTime() time.Time {
	return time.Unix(0, this.LastWriteTime.Nanoseconds())
}

func (this *FileInfo) Mode() os.FileMode {
	m := os.FileMode(0444)
	if this.IsDir() {
		m |= 0111 | os.ModeDir
	}
	if !this.IsReadOnly() {
		m |= 0222
	}
	return m
}

func (this *FileInfo) Sys() interface{} {
	return &this.Win32finddata
}

func (this *FileInfo) findNext() error {
	return syscall.FindNextFile(this.handle, &this.Win32finddata)
}

func (this *FileInfo) close() {
	syscall.FindClose(this.handle)
}

func (this *FileInfo) Attribute() uint32 {
	return this.FileAttributes
}

func (this *FileInfo) IsReparsePoint() bool {
	return (this.Attribute() & FILE_ATTRIBUTE_REPARSE_POINT) != 0
}

func (this *FileInfo) IsReadOnly() bool {
	return (this.Attribute() & syscall.FILE_ATTRIBUTE_READONLY) != 0
}

func (this *FileInfo) IsDir() bool {
	return (this.Attribute() & syscall.FILE_ATTRIBUTE_DIRECTORY) != 0
}

func (this *FileInfo) IsHidden() bool {
	return (this.Attribute() & syscall.FILE_ATTRIBUTE_HIDDEN) != 0
}

func (this *FileInfo) IsSystem() bool {
	return (this.Attribute() & syscall.FILE_ATTRIBUTE_SYSTEM) != 0
}

func Walk(pattern string, callback func(*FileInfo) bool) error {
	this, err := findFirst(pattern)
	if err != nil {
		return err
	}
	defer this.close()
	for {
		if !callback(this.clone()) {
			return nil
		}
		if err := this.findNext(); err != nil {
			return nil
		}
	}
}
