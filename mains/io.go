package mains

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
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
	L.SetField(index, "read", L.NewFunction(fileRead))
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
				return lerror(L, err.Error())
			}
			return 1
		}
	}
	return lerror(L, "(file)close: not a file-handle")
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
				return lerror(L, fmt.Sprintf("%s: can not open", filename))
			}
		} else {
			return lerror(L, "io.lines: not a string")
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
		return lerror(L, err.Error())
	}
	L.Push(newIoLuaWriter(L, fd, fd))
	return 1
}

func ioOpen(L *lua.LState) int {
	fname, ok := L.Get(1).(lua.LString)
	if !ok {
		return lerror(L, "io.open: filename is not a string")
	}
	mode, ok := L.Get(2).(lua.LString)
	if !ok {
		mode = "r"
	}
	if mode == "r" {
		fd, err := os.Open(string(fname))
		if err != nil {
			return lerror(L, err.Error())
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
	return lerror(L, fmt.Sprintf("io.open (nyagos compatible version) does not support mode=\"%s\" yet.", string(mode)))
}

func fileWrite(L *lua.LState) int {
	if ud, ok := L.Get(1).(*lua.LUserData); ok {
		if f, ok := ud.Value.(*ioLuaWriter); ok {
			if f.writer == nil {
				return lerror(L, "file:write: handle has already closed")
			}
			for i := 2; i <= L.GetTop(); i++ {
				io.WriteString(f.writer, L.Get(i).String())
			}
			L.Push(ud)
			return 1
		}
	}
	return lerror(L, "(file)write: not a file-handle object")
}

func ioPOpen(L *lua.LState) int {
	command, ok := L.Get(1).(lua.LString)
	if !ok {
		return lerror(L, "io.popen: command is not a string")
	}
	mode, ok := L.Get(2).(lua.LString)
	if !ok {
		return lerror(L, "io.popen: mode is not a string")
	}
	args := texts.SplitLikeShellString(string(command))
	for i, s := range args {
		args[i] = strings.Replace(s, "\"", "", -1)
	}
	xcmd := exec.Command(args[0], args[1:]...)

	if m := string(mode); m == "r" {
		in, err := xcmd.StdoutPipe()
		if err != nil {
			return lerror(L, err.Error())
		}
		if err := xcmd.Start(); err != nil {
			in.Close()
			return lerror(L, err.Error())
		}
		L.Push(newIoLuaReader(L, in, in))
		return 1
	} else if m == "w" {
		out, err := xcmd.StdinPipe()
		if err != nil {
			return lerror(L, err.Error())
		}
		if err := xcmd.Start(); err != nil {
			out.Close()
			return lerror(L, err.Error())
		}
		L.Push(newIoLuaWriter(L, out, out))
		return 1
	} else {
		return lerror(L, fmt.Sprintf("io.popen(...,\"%s\") is not supported yet", m))
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

func openIo(L *lua.LState) *lua.LTable {
	ioTable := L.NewTable()
	L.SetField(ioTable, "lines", L.NewFunction(ioLines))
	L.SetField(ioTable, "write", L.NewFunction(ioWrite))
	L.SetField(ioTable, "open", L.NewFunction(ioOpen))
	L.SetField(ioTable, "popen", L.NewFunction(ioPOpen))
	return ioTable
}

func fileRead(L *lua.LState) int {
	if ud, ok := L.Get(1).(*lua.LUserData); ok {
		if f, ok := ud.Value.(*ioLuaReader); ok {
			r := f.reader
			end := L.GetTop()
			result := make([]lua.LValue, 0, end-1)
			for i := 2; i <= end; i++ {
				val := L.Get(i)
				if num, ok := val.(lua.LNumber); ok {
					data := make([]byte, 0, int(num))
					for len(data) < cap(data) {
						b, err := r.ReadByte()
						if err != nil {
							return lerror(L, err.Error())
						}
						if b != '\r' {
							data = append(data, b)
						}
					}
					result = append(result, lua.LString(string(data)))
				} else if s, ok := val.(lua.LString); ok {
					switch s {
					case "*l":
						line, err := r.ReadString('\n')
						if err != nil {
							return lerror(L, err.Error())
						}
						line = strings.TrimSuffix(line, "\n")
						line = strings.TrimSuffix(line, "\r")
						result = append(result, lua.LString(line))
						break
					case "*a":
						all, err := ioutil.ReadAll(r)
						if err != nil {
							return lerror(L, err.Error())
						}
						text := strings.Replace(string(all), "\r\n", "\n", -1)
						result = append(result, lua.LString(text))
						break
					case "*n":
						var n int
						if _, err := fmt.Fscan(r, &n); err != nil {
							return lerror(L, err.Error())
						}
						result = append(result, lua.LNumber(n))
					default:
						return lerror(L, "(file)read: invalid argument")
					}
				} else {
					return lerror(L, "(file)read: invalid argument")
				}
			}
			for _, v := range result {
				L.Push(v)
			}
			return len(result)
		}
	}
	return lerror(L, "(file).read: not a file-handle")
}
