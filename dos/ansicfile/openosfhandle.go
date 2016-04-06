package ansicfile

import "errors"

var open_osfhandle = msvcrt.NewProc("_open_osfhandle")

func OpenOsFHandle(handle uintptr, flags int) (uintptr, error) {
	fd, _, err := open_osfhandle.Call(handle, uintptr(flags))
	if int(fd) == -1 {
		if err != nil {
			return 0, err
		} else {
			return 0, errors.New("ansicfile.OpenOsFHandle failed.")
		}
	} else {
		return fd, nil
	}
}
