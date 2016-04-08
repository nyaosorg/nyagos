package commands

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	. "../interpreter"
)

var timePattern = regexp.MustCompile(
	"^(?:(\\d\\d)?(\\d\\d))?(\\d\\d)(\\d\\d)(?:(\\d\\d)(\\d\\d))(?:\\.(\\d\\d))?$")

func atoiOr(s string, orelse int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return orelse
	}
	return val
}

func cmd_touch(this *Interpreter) (ErrorLevel, error) {
	errcnt := 0
	stamp := time.Now()
	for _, arg1 := range this.Args[1:] {
		if arg1[0] == '-' {
			fmt.Fprintf(this.Stderr,
				"%s: built-in touch: Not implemented.\n",
				arg1)
		} else if m := timePattern.FindStringSubmatch(arg1); m != nil {
			yy := atoiOr(m[2], stamp.Year()%100)
			cc, cc_err := strconv.Atoi(m[1])
			if cc_err != nil {
				if yy <= 68 {
					cc = 20
				} else {
					cc = 19
				}
			}
			year := yy + cc*100
			month, _ := strconv.Atoi(m[3])
			mday, _ := strconv.Atoi(m[4])
			hour := atoiOr(m[5], 0)
			min := atoiOr(m[6], 0)
			sec := atoiOr(m[7], 0)
			stamp = time.Date(year, time.Month(month), mday, hour, min, sec, 0, time.Local)
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
