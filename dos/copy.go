package dos

import "golang.org/x/sys/windows"

// Copy calls Win32's CopyFile API.
func Copy(src, dst string, isFailIfExists bool) error {
	rc, err := copyFile(src, dst, isFailIfExists)
	if rc == 0 {
		return err
	}
	return nil
}

// Move calls Win32's MoveFileEx API.
func Move(src, dst string) error {
	_src, err := windows.UTF16PtrFromString(src)
	if err != nil {
		return err
	}
	_dst, err := windows.UTF16PtrFromString(dst)
	if err != nil {
		return err
	}
	return windows.MoveFileEx(
		_src,
		_dst,
		windows.MOVEFILE_REPLACE_EXISTING|
			windows.MOVEFILE_COPY_ALLOWED|
			windows.MOVEFILE_WRITE_THROUGH)
}
