package mains

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yuin/gopher-lua"
)

type ioLuaReader struct {
	reader *bufio.Reader
	closer io.Closer
}

func (io *ioLuaReader) Close() error {
	if io.closer != nil {
		err := io.closer.Close()
		io.closer = nil
		return err
	}
	return nil
}

type ioLuaWriter struct {
	writer *bufio.Writer
	closer io.Closer
}

func (io *ioLuaWriter) Close() error {
	if io.closer != nil {
		err := io.closer.Close()
		io.closer = nil
		return err
	}
	return nil
}

func newIoLuaReader(L *lua.LState, r io.Reader, c io.Closer) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &ioLuaReader{
		reader: bufio.NewReader(r),
		closer: c,
	}
	return ud
}

func newIoLuaWriter(L *lua.LState, w io.Writer, c io.Closer) (*lua.LUserData, *bufio.Writer) {
	ud := L.NewUserData()
	bw := bufio.NewWriter(w)
	ud.Value = &ioLuaWriter{
		writer: bw,
		closer: c,
	}
	return ud, bw
}

func ioLinesIter(L *lua.LState) int {
	ud, ok := L.Get(1).(*lua.LUserData)
	if !ok {
		L.Push(lua.LNil)
		return 1
	}
	r, ok := ud.Value.(*ioLuaReader)
	if !ok {
		L.Push(lua.LNil)
		return 1
	}
	if text, err := r.reader.ReadString('\n'); err == nil {
		L.Push(lua.LString(strings.TrimSuffix(text, "\n")))
	} else {
		L.Push(lua.LNil)
		if r.closer != nil {
			r.closer.Close()
			r.closer = nil
		}
	}
	return 1
}

func ioLines(L *lua.LState) int {
	var ud *lua.LUserData
	if L.GetTop() >= 1 {
		if filename, ok := L.Get(1).(lua.LString); ok {
			if fd, err := os.Open(string(filename)); err == nil {
				ud = newIoLuaReader(L, fd, fd)
			} else {
				L.Push(lua.LNil)
				L.Push(lua.LString(fmt.Sprintf("%s: can not open", filename)))
				return 2
			}
		} else {
			L.Push(lua.LNil)
			L.Push(lua.LString("io.lines: not a string"))
			return 2
		}
	} else {
		nyagosTbl := L.GetGlobal("nyagos")
		ud = L.GetField(nyagosTbl, "stdin").(*lua.LUserData)
	}
	L.Push(L.NewFunction(ioLinesIter))
	L.Push(ud)
	L.Push(lua.LNil)
	return 3
}

func ioWrite(L *lua.LState) int {
	nyagosTbl := L.GetGlobal("nyagos")
	if stdout, ok := L.GetField(nyagosTbl, "stdout").(*lua.LUserData); ok {
		if w, ok := stdout.Value.(*ioLuaWriter); ok {
			for i := 1; i <= L.GetTop(); i++ {
				fmt.Fprint(w.writer, L.Get(i).String())
			}
			return 0
		}
	}
	fmt.Fprintln(os.Stderr, "nyagos.stdout is not filehandle")
	return 0
}
