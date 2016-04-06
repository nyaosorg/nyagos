package lua

import (
	"os"

	"../dos/ansicfile"
)

func (L Lua) pushFile(f *os.File, modeFlg int, modeStr string) error {
	// *os.File to file-descripter
	fd, fd_err := ansicfile.OpenOsFHandle(f.Fd(), modeFlg)
	if fd_err != nil {
		return fd_err
	}
	// file-descripter to FILE*
	fp, fp_err := ansicfile.FdOpen(fd, modeStr)
	if fp_err != nil {
		return fp_err
	}
	L.PushStream(fp)
	return nil
}

func (L Lua) PushFileWriter(f *os.File) error {
	return L.pushFile(f, os.O_APPEND, "w")

}

func (L Lua) PushFileReader(f *os.File) error {
	return L.pushFile(f, os.O_RDONLY, "r")
}
