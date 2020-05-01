// +build !vanilla

package mains

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/yuin/gopher-lua"
)

const ioTblName = "io"

func newXFile(L *lua.LState, fd *XFile, read, write bool) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = fd
	index := L.NewTable()
	if read {
		L.SetField(index, "lines", L.NewFunction(fileLines))
		L.SetField(index, "read", L.NewFunction(fileRead))
	}
	if write {
		L.SetField(index, "write", L.NewFunction(fileWrite))
		L.SetField(index, "flush", L.NewFunction(fileFlush))
	}
	L.SetField(index, "setvbuf", L.NewFunction(fileSetVBuf))
	L.SetField(index, "close", L.NewFunction(fileClose))
	L.SetField(index, "seek", L.NewFunction(fileSeek))
	meta := L.NewTable()
	L.SetField(meta, "__gc", L.NewFunction(fileClose))
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
			L.Push(lua.LTrue)
			return 1
		}
	}
	return lerror(L, "(file)close: not a file-handle")
}

// ioLineIter is the callback function for `io.lines()`
func ioLinesIter(L *lua.LState) int {
	ud, ok := L.Get(1).(*lua.LUserData)
	if !ok {
		L.Push(lua.LNil)
		return 1
	}
	r, ok := ud.Value.(*XFile)
	if !ok || r.Eof() || r.closed {
		L.Push(lua.LNil)
		return 1
	}

	if text, err := r.ReadString('\n'); err == nil {
		L.Push(lua.LString(strings.TrimSuffix(text, "\n")))
	} else {
		if err == io.EOF {
			r.SetEof()
			if len(text) > 0 {
				L.Push(lua.LString(text))
				r.Close()
				return 1
			}
		}
		L.Push(lua.LNil)
		r.Close()
	}
	return 1
}

func ioLines(L *lua.LState) int {
	var ud *lua.LUserData
	if L.GetTop() >= 1 {
		if filename, ok := L.Get(1).(lua.LString); ok {
			// io.lines("filename")
			//   requires close()
			if fd, err := os.Open(string(filename)); err == nil {
				ud = L.NewUserData()
				ud.Value = &XFile{File: fd}
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
		if w, ok := stdout.Value.(io.Writer); ok {
			for i := 1; i <= L.GetTop(); i++ {
				fmt.Fprint(w, L.Get(i).String())
			}
			return 0
		}
	}
	fmt.Fprintf(os.Stderr, "io.write: %s.stdout is not filehandle\n", ioTblName)
	return 0
}

func ioOpen(L *lua.LState) int {
	fname, ok := L.Get(1).(lua.LString)
	if !ok {
		return lerror(L, "io.open: filename is not a string")
	}
	mode := os.O_RDONLY
	read := true
	write := false
	if m, ok := L.Get(2).(lua.LString); ok {
		switch m {
		case "r", "rb":
			mode = os.O_RDONLY
			read = true
			write = false
		case "w", "wb":
			mode = os.O_WRONLY | os.O_CREATE
			read = false
			write = true
		case "a", "ab":
			mode = os.O_WRONLY | os.O_APPEND | os.O_CREATE
			read = false
			write = true
		case "r+", "rb+":
			mode = os.O_RDWR
			read = true
			write = true
		case "w+", "wb+":
			mode = os.O_RDWR | os.O_TRUNC | os.O_CREATE
			read = true
			write = true
		case "a+", "ab+":
			mode = os.O_APPEND | os.O_RDWR | os.O_CREATE
			read = true
			write = true
		default:
			return lerror(L, fmt.Sprintf("io.open (nyagos compatible version) does not support mode=\"%s\" yet.", string(mode)))
		}
	}
	fd, err := os.OpenFile(string(fname), mode, 0666)
	if err != nil {
		return lerror(L, err.Error())
	}
	L.Push(newXFile(L, &XFile{File: fd}, read, write))
	return 1
}

func fileWrite(L *lua.LState) int {
	if ud, ok := L.Get(1).(*lua.LUserData); ok {
		if w, ok := ud.Value.(io.Writer); ok {
			for i, end := 2, L.GetTop(); i <= end; i++ {
				fmt.Fprint(w, L.Get(i).String())
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
	xcmd := newCommand(string(command))
	// Append one space to enclose with double quotation by exec.Command
	xcmd.Stderr = os.Stderr

	if m := string(mode); m == "r" {
		xcmd.Stdin = os.Stdin
		in, out, err := os.Pipe()
		if err != nil {
			return lerror(L, err.Error())
		}
		xcmd.Stdout = out
		if err := xcmd.Start(); err != nil {
			in.Close()
			out.Close()
			return lerror(L, err.Error())
		}
		L.Push(newXFile(L, &XFile{File: in}, true, false))
		go func() {
			xcmd.Wait()
			out.Close()
		}()
		return 1
	} else if m == "w" {
		xcmd.Stdout = os.Stdout
		in, out, err := os.Pipe()
		if err != nil {
			return lerror(L, err.Error())
		}
		xcmd.Stdin = in
		if err := xcmd.Start(); err != nil {
			in.Close()
			out.Close()
			return lerror(L, err.Error())
		}
		L.Push(newXFile(L, &XFile{File: out}, false, true))
		go func() {
			xcmd.Wait()
			in.Close()
		}()
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
		type Syncer interface{ Sync() error }
		if fd, ok := ud.Value.(Syncer); ok {
			fd.Sync()
			L.Push(ud)
			return 1
		}
	}
	L.Push(lua.LNil)
	L.Push(lua.LString("(file):flush: not a file-handle object"))
	return 2
}

func ioType(L *lua.LState) int {
	if ud, ok := L.Get(1).(*lua.LUserData); ok {
		if x, ok := ud.Value.(*XFile); ok {
			if x.closed {
				L.Push(lua.LString("closed file"))
			} else {
				L.Push(lua.LString("file"))
			}
			return 1
		}
		if _, ok := ud.Value.(*os.File); ok {
			L.Push(lua.LString("file"))
			return 1
		}
	}
	L.Push(lua.LNil)
	return 1
}

func openIo(L *lua.LState) *lua.LTable {
	ioTable := L.NewTable()
	L.SetField(ioTable, "lines", L.NewFunction(ioLines))
	L.SetField(ioTable, "write", L.NewFunction(ioWrite))
	L.SetField(ioTable, "open", L.NewFunction(ioOpen))
	L.SetField(ioTable, "close", L.NewFunction(fileClose))
	L.SetField(ioTable, "popen", L.NewFunction(ioPOpen))
	L.SetField(ioTable, "type", L.NewFunction(ioType))
	L.SetField(ioTable, "stdin",
		newXFile(L, &XFile{File: os.Stdin, dontClose: true}, true, false))
	L.SetField(ioTable, "stdout",
		newXFile(L, &XFile{File: os.Stdout, dontClose: true}, false, true))
	L.SetField(ioTable, "stderr",
		newXFile(L, &XFile{File: os.Stderr, dontClose: true}, false, true))
	return ioTable
}

type Eofer interface {
	SetEof()
	Eof() bool
}

func fileRead(L *lua.LState) int {
	var err error
	ud, ok := L.Get(1).(*lua.LUserData)
	if !ok {
		L.ArgError(1, "not a file-handle")
		return 0
	}
	r, ok := ud.Value.(*XFile)
	if !ok {
		L.ArgError(1, "not a xfile-handle")
		return 0
	}

	end := L.GetTop()
	if end == 1 {
		L.Push(lua.LString("*l"))
		end++
	}
	result := make([]lua.LValue, 0, end-1)
	for i := 2; i <= end; i++ {
		if r.Eof() {
			break
		}
		val := L.Get(i)
		if num, ok := val.(lua.LNumber); ok {
			if num == 0 {
				_, err = r.ReadByte()
				if err == io.EOF {
					r.SetEof()
					result = append(result, lua.LNil)
					goto normalreturn
				}
				r.UnreadByte()
			}
			data := make([]byte, 0, int(num))
			for len(data) < cap(data) {
				b, err := r.ReadByte()
				if err == io.EOF {
					r.SetEof()
					if len(data) == 0 {
						result = append(result, lua.LNil)
					} else {
						result = append(result, lua.LString(string(data)))
					}
					goto normalreturn
				}
				if err != nil {
					goto errreturn
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
					if err == io.EOF {
						r.SetEof()
						if line == "" {
							result = append(result, lua.LNil)
							goto normalreturn
						}
					} else {
						goto errreturn
					}
				}
				line = strings.TrimSuffix(line, "\n")
				line = strings.TrimSuffix(line, "\r")
				result = append(result, lua.LString(line))
				break
			case "*a":
				var all []byte
				all, err = ioutil.ReadAll(r.reader())
				if err != nil {
					if err == io.EOF {
						r.SetEof()
						if len(all) <= 0 {
							result = append(result, lua.LString(""))
							goto normalreturn
						}
					} else {
						goto errreturn
					}
				}
				text := strings.Replace(string(all), "\r\n", "\n", -1)
				result = append(result, lua.LString(text))
				break
			case "*n":
				var n int
				if _, err = fmt.Fscan(r, &n); err != nil {
					if err == io.EOF ||
						(err != nil && err.Error() == "expected integer") {
						result = append(result, lua.LNil)
						goto normalreturn
					}
					goto errreturn
				}
				result = append(result, lua.LNumber(n))
			default:
				L.ArgError(i, "invalid format")
			}
		} else {
			L.ArgError(i, "invalid argument")
		}
	}
normalreturn:
	for _, v := range result {
		L.Push(v)
	}
	return len(result)
errreturn:
	L.RaiseError(err.Error())
	return 2
}

func fileSeek(L *lua.LState) int {
	ud, ok := L.Get(1).(*lua.LUserData)
	if !ok {
		return lerror(L, "(file)seek: not file-handle")
	}
	seeker, ok := ud.Value.(io.Seeker)
	if !ok {
		return lerror(L, "(file)seek: not seekable file handle")
	}
	whence := 1
	offset := int64(0)
	if L.GetTop() >= 2 {
		_whence, ok := L.Get(2).(lua.LString)
		if !ok {
			return lerror(L, "(file)seek: invalid whence string")
		}
		switch strings.ToLower(string(_whence)) {
		case "set":
			whence = io.SeekStart
		case "cur":
			whence = io.SeekCurrent
		case "end":
			whence = io.SeekEnd
		default:
			return lerror(L, "(file)seek: invalid whence string")
		}
		if L.GetTop() >= 3 {
			_offset, ok := L.Get(3).(lua.LNumber)
			if !ok {
				return lerror(L, "(file)seek: invalid offset number")
			}
			offset = int64(_offset)
		}
	}
	result, err := seeker.Seek(offset, whence)
	if err != nil {
		return lerror(L, err.Error())
	}
	L.Push(lua.LNumber(result))
	return 1
}

func fileSetVBuf(L *lua.LState) int {
	const msg = "file:setvbuf is not implemented yet"
	println(msg)
	return lerror(L, msg)
}
