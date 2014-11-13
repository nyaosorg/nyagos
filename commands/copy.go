package commands

import "fmt"
import "os"
import "path/filepath"

import "../dos"
import "../interpreter"

func cmd_copy(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	return interpreter.CONTINUE, cmd_xxxx(cmd, func(src, dst string) error {
		return dos.Copy(src, dst, false)
	})
}

func cmd_move(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	return interpreter.CONTINUE, cmd_xxxx(cmd, func(src, dst string) error {
		return dos.Move(src, dst)
	})
}

func cmd_xxxx(cmd *interpreter.Interpreter, cmds func(src, dst string) error) error {
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
		return cmds(src, dst)
	default:
		for i, n := 1, len(cmd.Args)-1; i < n; i++ {
			src := cmd.Args[i]
			dst := dos.Join(cmd.Args[n], filepath.Base(src))
			fmt.Fprintf(cmd.Stderr, "%s -> %s\n", src, dst)
			err := cmds(src, dst)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
