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
