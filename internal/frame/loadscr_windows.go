package frame

import (
	"os"
	"path/filepath"
)

var _appDataPath string

func appDataDir() string {
	if _appDataPath == "" {
		_appDataPath = filepath.Join(os.Getenv("APPDATA"), "NYAOS_ORG")
		os.Mkdir(_appDataPath, 0777)
	}
	return _appDataPath
}
