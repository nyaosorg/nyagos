package completion

import (
	"path/filepath"

	"github.com/nyaosorg/nyagos/internal/nodos"
)

func isExecutable(path string) bool {
	return nodos.IsExecutableSuffix(filepath.Ext(path))
}
