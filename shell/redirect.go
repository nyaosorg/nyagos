package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// NoClobber is the switch to forbide to overwrite the exist file.
var NoClobber = false

type _Redirecter struct {
	path     string
	isAppend bool
	no       int
	dupFrom  int
	force    bool
}

func newRedirecter(no int) *_Redirecter {
	return &_Redirecter{
		path:     "",
		isAppend: false,
		no:       no,
		dupFrom:  -1}
}

func (r *_Redirecter) FileNo() int {
	return r.no
}

func (r *_Redirecter) DupFrom(fileno int) {
	r.dupFrom = fileno
}

func (r *_Redirecter) SetPath(path string) {
	r.path = path
}

func (r *_Redirecter) SetAppend() {
	r.isAppend = true
}

var deviceName = map[string]struct{}{
	"AUX":    {},
	"CON":    {},
	"NUL":    {},
	"PRN":    {},
	"CLOCK$": {},
	"COM1":   {},
	"COM2":   {},
	"COM3":   {},
	"COM4":   {},
	"COM5":   {},
	"COM6":   {},
	"COM7":   {},
	"COM8":   {},
	"COM9":   {},
	"LPT1":   {},
	"LPT2":   {},
	"LPT3":   {},
	"LPT4":   {},
	"LPT5":   {},
	"LPT6":   {},
	"LPT7":   {},
	"LPT8":   {},
	"LPT9":   {},
}

func (r *_Redirecter) open() (*os.File, error) {
	if r.path == "" {
		return nil, errors.New("_Redirecter.open(): path=\"\"")
	}
	if r.no == 0 {
		return os.Open(r.path)
	} else if r.isAppend {
		return os.OpenFile(r.path, os.O_APPEND|os.O_CREATE, 0666)
	} else {
		if NoClobber && !r.force {
			_, err := os.Stat(r.path)
			if err == nil {
				name := strings.ToUpper(filepath.Base(r.path))
				if pos := strings.IndexRune(name, '.'); pos >= 0 {
					name = name[:pos]
				}
				if _, ok := deviceName[name]; !ok {
					return nil, fmt.Errorf("%s: cannot overwrite existing file", r.path)
				}
			}
		}
		return os.Create(r.path)
	}
}

type dontCloseHandle struct{}

func (this dontCloseHandle) Close() error {
	return nil
}

func (r *_Redirecter) OpenOn(cmd *Cmd) (closer io.Closer, err error) {
	var fd *os.File

	switch r.dupFrom {
	case 0:
		fd = cmd.Stdin
		closer = &dontCloseHandle{}
	case 1:
		fd = cmd.Stdout
		closer = &dontCloseHandle{}
	case 2:
		fd = cmd.Stderr
		closer = &dontCloseHandle{}
	default:
		fd, err = r.open()
		if err != nil {
			return nil, err
		}
		closer = fd
	}
	switch r.FileNo() {
	case 0:
		cmd.Stdin = fd
	case 1:
		cmd.Stdout = fd
	case 2:
		cmd.Stderr = fd
	default:
		panic("Assertion failed: _Redirecter.OpenAs: r.no not in (0,1,2)")
	}
	return
}
