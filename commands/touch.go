package commands

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/zetamatta/nyagos/commands/timecheck"
	"github.com/zetamatta/nyagos/shell"
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

func readTimeStamp(s string) *time.Time {
	m := timePattern.FindStringSubmatch(s)
	if m == nil {
		return nil
	}
	yy, yy_err := strconv.Atoi(m[2])
	if yy_err != nil {
		yy = time.Now().Year() % 100
	}
	cc, cc_err := strconv.Atoi(m[1])
	if cc_err != nil {
		if yy <= 68 {
			cc = 20
		} else {
			cc = 19
		}
	} else if cc < 19 {
		return nil
	}
	year := yy + cc*100
	month, _ := strconv.Atoi(m[3])
	mday, _ := strconv.Atoi(m[4])
	hour := atoiOr(m[5], 0)
	min := atoiOr(m[6], 0)
	sec := atoiOr(m[7], 0)
	if !timecheck.IsOk(year, month, mday, hour, min, sec) {
		return nil
	}
	stamp := time.Date(year, time.Month(month), mday, hour, min, sec, 0, time.Local)
	return &stamp
}

func cmd_touch(ctx context.Context, this *shell.Cmd) (int, error) {
	errcnt := 0
	stamp := time.Now()
	for i := 1; i < len(this.Args); i++ {
		arg1 := this.Args[i]
		if arg1 == "-t" {
			i++
			if i >= len(this.Args) {
				fmt.Fprintf(this.Stderr, "-t: Too Few Arguments.\n")
				return 255, nil
			}
			stamp_ := readTimeStamp(this.Args[i])
			if stamp_ == nil {
				fmt.Fprintf(this.Stderr, "-t: %s: Invalid time format.\n",
					this.Args[i])
				return 255, nil
			}
			stamp = *stamp_
		} else if arg1 == "-r" {
			i++
			if i >= len(this.Args) {
				fmt.Fprintf(this.Stderr, "-r: Too Few Arguments.\n")
				return 255, nil
			}
			stat, statErr := os.Stat(this.Args[i])
			if statErr != nil {
				fmt.Fprintf(this.Stderr, "-r: %s: %s\n", this.Args[i], statErr)
				return 255, nil
			}
			stamp = stat.ModTime()
		} else if arg1[0] == '-' {
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
	return errcnt, nil
}
