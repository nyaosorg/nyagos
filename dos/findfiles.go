package dos

import (
	"os"
	"syscall"
	"time"
)

type FindFiles struct {
	handle syscall.Handle
	data1  syscall.Win32finddata
}

func (this *FindFiles) Clone() *FindFiles {
	return &FindFiles{this.handle, this.data1}
}

func FindFirst(pattern string) (*FindFiles, error) {
	pattern16, err := syscall.UTF16PtrFromString(pattern)
	if err != nil {
		return nil, err
	}
	this := new(FindFiles)
	this.handle, err = syscall.FindFirstFile(pattern16, &this.data1)
	if err != nil {
		return nil, err
	}
	return this, nil
}

func (this *FindFiles) Name() string {
	return syscall.UTF16ToString(this.data1.FileName[:])
}

func (this *FindFiles) Size() int64 {
	return int64((this.data1.FileSizeHigh << 32) | this.data1.FileSizeLow)
}

func (this *FindFiles) ModTime() time.Time {
	return time.Unix(0, this.data1.LastWriteTime.Nanoseconds())
}

func (this *FindFiles) Mode() os.FileMode {
	m := os.FileMode(0444)
	if this.IsDir() {
		m |= 0111 | os.ModeDir
	}
	if !this.IsReadOnly() {
		m |= 0222
	}
	return m
}

func (this *FindFiles) Sys() interface{} {
	return this.data1
}

func (this *FindFiles) FindNext() error {
	return syscall.FindNextFile(this.handle, &this.data1)
}

func (this *FindFiles) Close() {
	syscall.FindClose(this.handle)
}

func (this *FindFiles) Attribute() uint32 {
	return this.data1.FileAttributes
}

func (this *FindFiles) IsReparsePoint() bool {
	return (this.Attribute() & FILE_ATTRIBUTE_REPARSE_POINT) != 0
}

func (this *FindFiles) IsReadOnly() bool {
	return (this.Attribute() & syscall.FILE_ATTRIBUTE_READONLY) != 0
}

func (this *FindFiles) IsExecutable() bool {
	return IsExecutableSuffix(this.Name())
}

func (this *FindFiles) IsDir() bool {
	return (this.Attribute() & syscall.FILE_ATTRIBUTE_DIRECTORY) != 0
}

func (this *FindFiles) IsHidden() bool {
	return (this.Attribute() & syscall.FILE_ATTRIBUTE_HIDDEN) != 0
}

func (this *FindFiles) IsSystem() bool {
	return (this.Attribute() & syscall.FILE_ATTRIBUTE_SYSTEM) != 0
}

func DirName(path string) string {
	lastroot := -1
	for i, i_end := 0, len(path); i < i_end; i++ {
		switch path[i] {
		case '\\', '/', ':':
			lastroot = i
		}
	}
	if lastroot >= 0 {
		return path[0:(lastroot + 1)]
	} else {
		return ""
	}
}

func ForFiles(pattern string, callback func(*FindFiles) bool) error {
	this, err := FindFirst(pattern)
	if err != nil {
		return err
	}
	defer this.Close()
	for {
		if !callback(this.Clone()) {
			return nil
		}
		if err := this.FindNext(); err != nil {
			return nil
		}
	}
}
