package commands

import "fmt"
import "os"
import "path/filepath"

import "../dos"
import "../interpreter"

func cmd_copy(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	switch len(cmd.Args) {
	case 0, 1, 2:
		fmt.Fprintf(cmd.Stderr,
			"Usage: %s SOURCE-FILENAME DESITINATE-FILENAME\n"+
				"       %s FILENAMES... DESINATE-DIRECTORY\n",
			cmd.Args[0], cmd.Args[0])
	case 3:
		src := cmd.Args[1]
		dst := cmd.Args[2]
		fi, err := os.Stat(dst)
		if err == nil && fi != nil && fi.Mode().IsDir() {
			dst = dos.Join(dst, filepath.Base(src))
		}
		fmt.Fprintf(cmd.Stderr, "%s -> %s\n", src, dst)
		err = dos.Copy(src, dst, false)
		if err != nil {
			return interpreter.CONTINUE, err
		}
	default:
		for i, n := 1, len(cmd.Args)-1; i < n; i++ {
			src := cmd.Args[i]
			dst := dos.Join(cmd.Args[n], filepath.Base(src))
			fmt.Fprintf(cmd.Stderr, "%s -> %s\n", src, dst)
			err := dos.Copy(src, dst, false)
			if err != nil {
				return interpreter.CONTINUE, err
			}
		}
	}
	return interpreter.CONTINUE, nil
}
