package mains

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/yuin/gopher-lua"
	"github.com/zetamatta/nyagos/texts"
)

const ioTblName = "io"

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
	if io.writer != nil {
		io.writer.Flush()
	}
	if io.closer != nil {
		err := io.closer.Close()
		io.closer = nil
		io.writer = nil
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
	index := L.NewTable()
	L.SetField(index, "lines", L.NewFunction(fileLines))
	L.SetField(index, "close", L.NewFunction(fileClose))
	meta := L.NewTable()
	L.SetField(meta, "__index", index)
	L.SetMetatable(ud, meta)
	return ud
}

func fileClose(L *lua.LState) int {
	if ud, ok := L.Get(1).(*lua.LUserData); ok {
		if this, ok := ud.Value.(io.Closer); ok {
			err := this.Close()
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}
			return 1
		}
	}
	L.Push(lua.LNil)
	L.Push(lua.LString("(file)close: not a file-handle"))
	return 2
}

func newIoLuaWriter(L *lua.LState, w io.Writer, c io.Closer) *lua.LUserData {
	ud := L.NewUserData()
	bw := bufio.NewWriter(w)
	ud.Value = &ioLuaWriter{
		writer: bw,
		closer: c,
	}
	meta := L.NewTable()
	L.SetField(meta, "__gc", L.NewFunction(fileClose))
	index := L.NewTable()
	L.SetField(index, "close", L.NewFunction(fileClose))
	L.SetField(index, "write", L.NewFunction(fileWrite))
	L.SetField(index, "flush", L.NewFunction(fileFlush))
	L.SetField(meta, "__index", index)
	L.SetMetatable(ud, meta)
	return ud
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
		if err == io.EOF && text != "" {
			L.Push(lua.LString(text))
			return 1
		}
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
		ioTbl := L.GetGlobal(ioTblName)
		ud = L.GetField(ioTbl, "stdin").(*lua.LUserData)
	}
	L.Push(L.NewFunction(ioLinesIter))
	L.Push(ud)
	L.Push(lua.LNil)
	return 3
}

func ioWrite(L *lua.LState) int {
	ioTbl := L.GetGlobal(ioTblName)
	if stdout, ok := L.GetField(ioTbl, "stdout").(*lua.LUserData); ok {
		if w, ok := stdout.Value.(*ioLuaWriter); ok {
			for i := 1; i <= L.GetTop(); i++ {
				fmt.Fprint(w.writer, L.Get(i).String())
			}
			return 0
		}
	}
	fmt.Fprintf(os.Stderr, "io.write: %s.stdout is not filehandle\n", ioTblName)
	return 0
}

func _ioOpenWriter(L *lua.LState, fd *os.File, err error) int {
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(newIoLuaWriter(L, fd, fd))
	return 1
}

func ioOpen(L *lua.LState) int {
	fname, ok := L.Get(1).(lua.LString)
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString("io.open: filename is not a string"))
		return 2
	}
	mode, ok := L.Get(2).(lua.LString)
	if !ok {
		mode = "r"
	}
	if mode == "r" {
		fd, err := os.Open(string(fname))
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(newIoLuaReader(L, fd, fd))
		return 1
	}
	if mode == "w" {
		fd, err := os.Create(string(fname))
		return _ioOpenWriter(L, fd, err)
	}
	if mode == "a" {
		fd, err := os.OpenFile(string(fname), os.O_APPEND, 0755)
		return _ioOpenWriter(L, fd, err)
	}
	errmsg := fmt.Sprintf("io.open (nyagos compatible version) does not support mode=\"%s\" yet.", string(mode))
	L.Push(lua.LNil)
	L.Push(lua.LString(errmsg))
	fmt.Fprintln(os.Stderr, errmsg)
	return 2
}

func fileWrite(L *lua.LState) int {
	if ud, ok := L.Get(1).(*lua.LUserData); ok {
		if f, ok := ud.Value.(*ioLuaWriter); ok {
			if f.writer == nil {
				L.Push(lua.LNil)
				L.Push(lua.LString("file:write: handle has already closed"))
				return 2
			}
			for i := 2; i <= L.GetTop(); i++ {
				io.WriteString(f.writer, L.Get(i).String())
			}
			L.Push(ud)
			return 1
		}
	}
	L.Push(lua.LNil)
	L.Push(lua.LString("(file)write: not a file-handle object"))
	return 2
}

func ioPOpen(L *lua.LState) int {
	command, ok := L.Get(1).(lua.LString)
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString("io.popen: command is not a string"))
		return 2
	}
	mode, ok := L.Get(2).(lua.LString)
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString("io.popen: mode is not a string"))
		return 2
	}
	args := texts.SplitLikeShellString(string(command))
	for i, s := range args {
		args[i] = strings.Replace(s, "\"", "", -1)
	}
	xcmd := exec.Command(args[0], args[1:]...)

	if m := string(mode); m == "r" {
		in, err := xcmd.StdoutPipe()
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		if err := xcmd.Start(); err != nil {
			in.Close()
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(newIoLuaReader(L, in, in))
		return 1
	} else if m == "w" {
		out, err := xcmd.StdinPipe()
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		if err := xcmd.Start(); err != nil {
			out.Close()
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(newIoLuaWriter(L, out, out))
		return 1
	} else {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("io.popen(...,\"%s\") is not supported yet", m)))
		return 2
	}
}

func fileLines(L *lua.LState) int {
	L.Push(L.NewFunction(ioLinesIter))
	L.Push(L.Get(1))
	L.Push(lua.LNil)
	return 3
}

func fileFlush(L *lua.LState) int {
	if ud, ok := L.Get(1).(*lua.LUserData); ok {
		if f, ok := ud.Value.(*ioLuaWriter); ok {
			f.writer.Flush()
			L.Push(ud)
			return 1
		}
	}
	L.Push(lua.LNil)
	L.Push(lua.LString("(file):flush: not a file-handle object"))
	return 2
}
