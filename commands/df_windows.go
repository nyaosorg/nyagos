package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/dustin/go-humanize"

	"github.com/zetamatta/nyagos/dos"
)

func df(rootPathName string, w io.Writer) (err error) {
	io.WriteString(w, rootPathName)
	free, total, totalFree, err1 := dos.GetDiskFreeSpace(rootPathName)
	if err1 != nil {
		fmt.Fprintf(w, " %20s %20s %20s     ", "", "", "")
		err = fmt.Errorf("%s: %s", rootPathName, err1)
	} else {
		fmt.Fprintf(w, " %20s %20s %20s %3d%%",
			humanize.Comma(int64(free)),
			humanize.Comma(int64(total)),
			humanize.Comma(int64(totalFree)),
			100*(total-free)/total)
	}
	t, err1 := dos.GetDriveType(rootPathName)
	if err1 != nil {
		if err != nil {
			err = fmt.Errorf("%s,%s", err, err1)
		} else {
			err = fmt.Errorf("%s: %s", rootPathName, err1)
		}
	} else {
		switch t {
		case dos.DRIVE_REMOVABLE:
			io.WriteString(w, " [REMOVABLE]")
		case dos.DRIVE_FIXED:
			io.WriteString(w, " [FIXED]")
		case dos.DRIVE_REMOTE:
			io.WriteString(w, " [REMOTE]")
		case dos.DRIVE_CDROM:
			io.WriteString(w, " [CDROM]")
		case dos.DRIVE_RAMDISK:
			io.WriteString(w, " [RAMDISK]")
		}
	}
	fmt.Fprintln(w)
	return
}

func cmdDiskFree(_ context.Context, cmd Param) (int, error) {
	bits, err := dos.GetLogicalDrives()
	if err != nil {
		return 0, err
	}
	fmt.Fprintf(cmd.Out(), "   %20s %20s %20s Use%%\n",
		"Available",
		"TotalNumber",
		"TotalNumberOfFree")

	count := 0
	for _, arg1 := range cmd.Args()[1:] {
		if err := df(arg1, cmd.Out()); err != nil {
			return 0, err
		}
		count++
	}
	if count <= 0 {
		for d := 'A'; d <= 'Z'; d++ {
			if (bits & 1) != 0 {
				rootPathName := fmt.Sprintf("%c:", d)
				if err := df(rootPathName, cmd.Out()); err != nil {
					fmt.Fprintln(cmd.Err(), err)
				}
			}
			bits >>= 1
		}
	}
	return 0, nil
}
