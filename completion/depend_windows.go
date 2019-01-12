package completion

import (
	"path/filepath"

	"github.com/zetamatta/nyagos/dos"
)

func isExecutable(path string) bool {
	return dos.IsExecutableSuffix(filepath.Ext(path))
}
