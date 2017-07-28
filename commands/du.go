package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zetamatta/nyagos/shell"
)

func du_(path string, output func(string, int64)) (int64, error) {
	fd, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	fileInfos, err := fd.Readdir(-1)
	if err != nil {
		fd.Close()
		return 0, err
	}
	var diskuse int64 = 0
	dirs := make([]string, 0, len(fileInfos))
	for _, fileInfo1 := range fileInfos {
		if fileInfo1.IsDir() {
			dirs = append(dirs, fileInfo1.Name())
		}
		diskuse += fileInfo1.Size()
	}
	if err := fd.Close(); err != nil {
		return diskuse, err
	}
	for _, dir1 := range dirs {
		fullpath := filepath.Join(path, dir1)
		diskuse1, err := du_(fullpath, output)
		if err != nil {
			return diskuse, err
		}
		diskuse += diskuse1
	}
	output(path, diskuse)
	return diskuse, nil
}

func cmd_du(_ context.Context, cmd *shell.Cmd) (int, error) {
	output := func(name string, size int64) {
		fmt.Fprintf(cmd.Stdout, "%d\t%s\n", size/1000, name)
	}
	if len(cmd.Args) <= 2 {
		_, err := du_(".", output)
		if err != nil {
			return 1, err
		}
	}
	for _, arg1 := range cmd.Args[1:] {
		_, err := du_(arg1, output)
		if err != nil {
			return 1, err
		}
	}
	return 0, nil
}
