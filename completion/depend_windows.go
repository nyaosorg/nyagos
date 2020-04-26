package completion

import (
	"path/filepath"

	"github.com/zetamatta/nyagos/nodos"
)

func isExecutable(path string) bool {
	return nodos.IsExecutableSuffix(filepath.Ext(path))
}
