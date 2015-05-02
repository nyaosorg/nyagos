package dos

import "syscall"

type FindFiles struct {
	handle syscall.Handle
	data1  syscall.Win32finddata
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

func (this *FindFiles) FindNext() error {
	return syscall.FindNextFile(this.handle, &this.data1)
}

func (this *FindFiles) Close() {
	syscall.FindClose(this.handle)
}

func (this *FindFiles) Attribute() uint32 {
	return this.data1.FileAttributes
}

func (this *FindFiles) IsDir() bool {
	return (this.Attribute() & syscall.FILE_ATTRIBUTE_DIRECTORY) != 0
}

func (this *FindFiles) IsHidden() bool {
	return (this.Attribute() & syscall.FILE_ATTRIBUTE_HIDDEN) != 0
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
