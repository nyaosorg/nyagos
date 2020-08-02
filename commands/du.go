package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var errCtrlC = errors.New("^C")

func printDu1Line(out io.Writer, name string, size int64) {
	fmt.Fprintf(out, "%7s %s\n", formatByHumanize(size), name)
}

// _du returns the sum (bytes) of path
func _du(path string, output func(string, int64) error, stderr io.Writer, blocksize int64) (int64, error) {
	fd, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	stat, err := fd.Stat()
	if err != nil {
		return 0, err
	}
	if !stat.IsDir() {
		fd.Close()
		return ((stat.Size() + blocksize - 1) / blocksize) * blocksize, nil
	}
	files, err := fd.Readdir(-1)
	if err != nil {
		fd.Close()
		return 0, err
	}
	var diskuse int64
	dirs := make([]string, 0, len(files))
	for _, file1 := range files {
		if file1.IsDir() {
			dirs = append(dirs, file1.Name())
		} else {
			diskuse += ((file1.Size() + blocksize - 1) / blocksize) * blocksize
		}
	}
	if err := fd.Close(); err != nil {
		return diskuse, err
	}
	for _, dir1 := range dirs {
		fullpath := filepath.Join(path, dir1)
		diskuse1, err := _du(fullpath, output, stderr, blocksize)
		if err == nil {
			if err = output(fullpath, diskuse1); err == nil {
				diskuse += diskuse1
				continue
			}
		}
		fmt.Fprintf(stderr, "%s: %s\n", fullpath, err)
	}
	return diskuse, nil
}

func cmdDiskUsed(ctx context.Context, cmd Param) (int, error) {
	output := func(name string, size int64) error {
		printDu1Line(cmd.Out(), name, size)
		if ctx != nil {
			select {
			case <-ctx.Done():
				return errCtrlC
			default:
			}
		}
		return nil
	}
	count := 0
	for _, arg1 := range cmd.Args()[1:] {
		if arg1 == "-s" {
			output = func(_ string, _ int64) error {
				if ctx != nil {
					select {
					case <-ctx.Done():
						return errCtrlC
					default:
					}
				}
				return nil
			}
			continue
		}
		size, err := _du(arg1, output, cmd.Err(), 4096)
		count++
		if err != nil {
			fmt.Fprintf(cmd.Err(), "%s: %s\n", arg1, err)
			continue
		}
		printDu1Line(cmd.Out(), arg1, size)
	}
	if count <= 0 {
		size, err := _du(".", output, cmd.Err(), 4096)
		if err != nil {
			return 1, err
		}
		printDu1Line(cmd.Out(), ".", size)
	}
	return 0, nil
}
