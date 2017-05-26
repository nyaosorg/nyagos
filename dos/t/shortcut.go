package main

import (
	"path/filepath"

	"github.com/zetamatta/nyagos/dos"
)

func main() {
	dos.CoInitializeEx(0, dos.COINIT_MULTITHREADED)
	defer dos.CoUninitialize()

	path1, err := filepath.Abs("shortcut.go.lnk")
	if err != nil {
		println(err.Error())
		return
	}
	println("path=" + path1)
	target, dir, err := dos.ReadShortcut(path1)
	if err != nil {
		println(err.Error())
		return
	}
	println("target=" + target)
	println("dir=" + dir)
}
