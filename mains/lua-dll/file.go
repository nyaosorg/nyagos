package lua

import (
	"os"

	"github.com/zetamatta/go-ansicfile"
)

func (this Lua) pushFile(f *os.File, modeFlg int, modeStr string) error {
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
	this.PushStream(fp)
	return nil
}

func (this Lua) PushFileWriter(f *os.File) error {
	return this.pushFile(f, ansicfile.O_APPEND|ansicfile.O_TEXT, "wt")
}

func (this Lua) PushFileReader(f *os.File) error {
	return this.pushFile(f, ansicfile.O_RDONLY|ansicfile.O_TEXT, "rt")
}
