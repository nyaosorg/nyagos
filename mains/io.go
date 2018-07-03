package mains

import (
	"bufio"
	"io"
	"os"

	"github.com/yuin/gopher-lua"
)

type fileHandleT struct {
	scanner *bufio.Scanner
	closer  io.Closer
}

func ioLinesIter(L *lua.LState) int {
	ud, ok := L.Get(1).(*lua.LUserData)
	if !ok {
		L.Push(lua.LNil)
		return 1
	}
	fh, ok := ud.Value.(*fileHandleT)
	if !ok {
		L.Push(lua.LNil)
		return 1
	}
	if fh.scanner.Scan() {
		L.Push(lua.LString(fh.scanner.Text()))
		return 1
	} else {
		L.Push(lua.LNil)
		if fh.closer != nil {
			fh.closer.Close()
		}
		return 1
	}
}

func ioLines(L *lua.LState) int {
	ud := L.NewUserData()
	_, sh := getRegInt(L)
	if sh != nil {
		ud.Value = &fileHandleT{
			scanner: bufio.NewScanner(sh.In()),
			closer:  nil,
		}
	} else {
		ud.Value = &fileHandleT{
			scanner: bufio.NewScanner(os.Stdin),
			closer:  nil,
		}
	}
	L.Push(L.NewFunction(ioLinesIter))
	L.Push(ud)
	L.Push(lua.LNil)
	return 3
}
