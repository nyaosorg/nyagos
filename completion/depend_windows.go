package completion

import (
	"path/filepath"

	"github.com/zetamatta/nyagos/dos"
)

func join(dir, name string) string {
	return dos.Join(dir, name)
}

func isExecutable(path string) bool {
	return dos.IsExecutableSuffix(filepath.Ext(path))
}
