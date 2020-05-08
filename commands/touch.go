package commands

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
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
	yy, err := strconv.Atoi(m[2])
	if err != nil {
		yy = time.Now().Year() % 100
	}
	cc, err := strconv.Atoi(m[1])
	if err != nil {
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
	if !stampIsValid(year, month, mday, hour, min, sec) {
		return nil
	}
	stamp := time.Date(year, time.Month(month), mday, hour, min, sec, 0, time.Local)
	return &stamp
}

func cmdTouch(ctx context.Context, this Param) (int, error) {
	errcnt := 0
	stamp := time.Now()
	for i := 1; i < len(this.Args()); i++ {
		arg1 := this.Arg(i)
		if arg1 == "-t" {
			i++
			if i >= len(this.Args()) {
				fmt.Fprintf(this.Err(), "-t: Too Few Arguments.\n")
				return 255, nil
			}
			stamp1 := readTimeStamp(this.Arg(i))
			if stamp1 == nil {
				fmt.Fprintf(this.Err(), "-t: %s: Invalid time format.\n",
					this.Arg(i))
				return 255, nil
			}
			stamp = *stamp1
		} else if arg1 == "-r" {
			i++
			if i >= len(this.Args()) {
				fmt.Fprintf(this.Err(), "-r: Too Few Arguments.\n")
				return 255, nil
			}
			stat, statErr := os.Stat(this.Arg(i))
			if statErr != nil {
				fmt.Fprintf(this.Err(), "-r: %s: %s\n", this.Arg(i), statErr)
				return 255, nil
			}
			stamp = stat.ModTime()
		} else if arg1[0] == '-' {
			fmt.Fprintf(this.Err(),
				"%s: built-in touch: Not implemented.\n",
				arg1)
		} else {
			fd, err := os.OpenFile(arg1, os.O_APPEND|os.O_CREATE, 0666)
			if err == nil {
				if err = fd.Close(); err != nil {
					fmt.Fprintln(this.Err(), err.Error())
					errcnt++
					continue
				}
				os.Chtimes(arg1, stamp, stamp)
			} else {
				fmt.Fprintln(this.Err(), err.Error())
				errcnt++
			}
		}
	}
	return errcnt, nil
}
