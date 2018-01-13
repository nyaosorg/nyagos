package shell

import (
	"errors"
	"os"
)

var NoClobber = false

type Redirecter struct {
	path     string
	isAppend bool
	no       int
	dupFrom  int
	force    bool
}

func NewRedirecter(no int) *Redirecter {
	return &Redirecter{
		path:     "",
		isAppend: false,
		no:       no,
		dupFrom:  -1}
}

func (this *Redirecter) FileNo() int {
	return this.no
}

func (this *Redirecter) DupFrom(fileno int) {
	this.dupFrom = fileno
}

func (this *Redirecter) SetPath(path string) {
	this.path = path
}

func (this *Redirecter) SetAppend() {
	this.isAppend = true
}

func (this *Redirecter) open() (*os.File, error) {
	if this.path == "" {
		return nil, errors.New("Redirecter.open(): path=\"\"")
	}
	if this.no == 0 {
		return os.Open(this.path)
	} else if this.isAppend {
		f, err := os.OpenFile(this.path, os.O_APPEND, 0666)
		if err != nil && os.IsNotExist(err) {
			f, err = os.Create(this.path)
		}
		return f, err
	} else {
		if NoClobber && !this.force {
			_, err := os.Stat(this.path)
			if err == nil {
				return nil, os.ErrExist
			}
		}
		return os.Create(this.path)
	}
}

func (this *Redirecter) OpenOn(cmd *Cmd) (*os.File, error) {
	var fd *os.File
	var err error

	switch this.dupFrom {
	case 0:
		fd = cmd.Stdin
	case 1:
		fd = cmd.Stdout
	case 2:
		fd = cmd.Stderr
	default:
		fd, err = this.open()
		if err != nil {
			return nil, err
		}
	}
	switch this.FileNo() {
	case 0:
		cmd.Stdin = fd
	case 1:
		cmd.Stdout = fd
	case 2:
		cmd.Stderr = fd
	default:
		panic("Assertion failed: Redirecter.OpenAs: this.no not in (0,1,2)")
	}
	return fd, nil
}
