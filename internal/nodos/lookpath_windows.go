package nodos

import (
	"path/filepath"

	"github.com/nyaosorg/go-windows-findfile"
)

func lookPathSkip(f *findfile.FileInfo) bool {
	return f.IsDir() || filepath.Ext(f.Name()) == ""
}
