package nodos

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/nyaosorg/go-windows-findfile"
)

func lookPathSkip(f *findfile.FileInfo) bool {
	return f.IsDir() || filepath.Ext(f.Name()) == ""
}
