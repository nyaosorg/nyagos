package commands

import (
	"fmt"
	"os"
	"time"

	. "../interpreter"
)

func cmd_touch(this *Interpreter) (ErrorLevel, error) {
	errcnt := 0
	stamp := time.Now()
	for _, arg1 := range this.Args[1:] {
		if arg1[0] == '-' {
			fmt.Fprintf(this.Stderr,
				"%s: built-in touch: Not implemented.\n",
				arg1)
		} else {
			fd, err := os.OpenFile(arg1, os.O_APPEND, 0666)
			if err != nil && os.IsNotExist(err) {
				fd, err = os.Create(arg1)
			}
			if err == nil {
				fd.Close()
				os.Chtimes(arg1, stamp, stamp)
			} else {
				fmt.Fprintln(this.Stderr, err.Error())
				errcnt++
			}
		}
	}
	return ErrorLevel(errcnt), nil
}
