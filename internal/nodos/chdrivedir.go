package nodos

import (
	"errors"
	"os"
	"path/filepath"
)

func saveLastDirOfDriveToEnv(dir string) {
	os.Setenv("="+filepath.VolumeName(dir), dir)
}

func getLastDirOfdrive(d byte) string {
	dir := os.Getenv(string([]byte{'=', d, ':'}))
	if dir == "" {
		dir = string([]byte{d, ':', os.PathSeparator})
	}
	return dir
}

func isOtherDrivesRelative(path string) bool {
	if len(path) < 2 {
		return false
	}

	if (path[0] < 'A' && 'Z' < path[0]) &&
		(path[0] < 'a' && 'z' < path[0]) {
		return false
	}

	if path[1] != ':' {
		return false
	}
	return len(path) == 2 || !os.IsPathSeparator(path[2])
}

// chDriveRune changes drive and returns the previous working directory.
func chDriveByte(n byte) (string, error) {
	lastDir, err := os.Getwd()
	if err == nil {
		saveLastDirOfDriveToEnv(lastDir)
	}
	return lastDir, os.Chdir(getLastDirOfdrive(n))
}

// Chdrive changes drive without changing the working directory there.
// And returns the previous working directory.
func Chdrive(drive string) (string, error) {
	if len(drive) < 1 {
		return "", errors.New("Chdrive: drive is empty string")
	}
	return chDriveByte(drive[0])
}

func Chdir(folder string) (err error) {
	if isOtherDrivesRelative(folder) {
		folder = filepath.Join(getLastDirOfdrive(folder[0]), folder[2:])
	}
	if absFolder, err := filepath.Abs(folder); err == nil {
		folder = absFolder
	}
	if err := os.Chdir(folder); err != nil {
		return err
	}
	saveLastDirOfDriveToEnv(folder)
	return nil
}
