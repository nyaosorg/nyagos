package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/dustin/go-humanize"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/shell"
)

func df(rootPathName string, w io.Writer) (err error) {
	fmt.Fprint(w, rootPathName)
	free, total, totalFree, err1 := dos.GetDiskFreeSpace(rootPathName)
	if err1 != nil {
		err = fmt.Errorf("%s: %s", rootPathName, err1)
	} else {
		fmt.Fprintf(w, " %20s %20s %20s",
			humanize.Comma(int64(free)),
			humanize.Comma(int64(total)),
			humanize.Comma(int64(totalFree)))
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
			fmt.Fprint(w, " [REMOVABLE]")
		case dos.DRIVE_FIXED:
			fmt.Fprint(w, " [FIXED]")
		case dos.DRIVE_REMOTE:
			fmt.Fprint(w, " [REMOTE]")
		case dos.DRIVE_CDROM:
			fmt.Fprint(w, " [CDROM]")
		case dos.DRIVE_RAMDISK:
			fmt.Fprint(w, " [RAMDISK]")
		}
	}
	fmt.Fprintln(w)
	return
}

func cmd_df(_ context.Context, cmd *shell.Cmd) (int, error) {
	bits, err := dos.GetLogicalDrives()
	if err != nil {
		return 0, err
	}
	fmt.Fprintf(cmd.Stdout, "   %20s %20s %20s\n",
		"FreeBytesAvailable",
		"TotalNumberOfBytes",
		"TotalNumberOfFreeBytes")

	count := 0
	for _, arg1 := range cmd.Args[1:] {
		if err := df(arg1, cmd.Stdout); err != nil {
			return 0, err
		}
		count++
	}
	if count <= 0 {
		for d := 'A'; d <= 'Z'; d++ {
			if (bits & 1) != 0 {
				rootPathName := fmt.Sprintf("%c:", d)
				if err := df(rootPathName, cmd.Stdout); err != nil {
					fmt.Fprintln(cmd.Stderr, err)
				}
			}
			bits >>= 1
		}
	}
	return 0, nil
}
