package commands

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/sys/windows"

	"github.com/dustin/go-humanize"

	"github.com/zetamatta/nyagos/dos"
)

func driveType(rootPathName string) string {
	_rootPathName, err := windows.UTF16PtrFromString(rootPathName)
	if err != nil {
		return "UNKNOWN"
	}

	t := windows.GetDriveType(_rootPathName)
	switch t {
	case windows.DRIVE_REMOVABLE:
		return "REMOVABLE"
	case windows.DRIVE_FIXED:
		return "FIXED"
	case windows.DRIVE_REMOTE:
		return "REMOTE"
	case windows.DRIVE_CDROM:
		return "CDROM"
	case windows.DRIVE_RAMDISK:
		return "RAMDISK"
	default:
		return "UNKNOWN"
	}
}

func df(rootPathName string, w io.Writer) error {
	label, fs, err := dos.VolumeName(rootPathName)
	if err != nil {
		return fmt.Errorf("%s: %s", rootPathName, err)
	}
	free, total, totalFree, err := dos.GetDiskFreeSpace(rootPathName)
	if err != nil {
		return fmt.Errorf("%s: %s", rootPathName, err)
	}
	fmt.Fprintf(w, "%s %16s %16s %16s %3d%% \"%s\" (%s/%s)\n",
		rootPathName,
		humanize.Comma(int64(free)),
		humanize.Comma(int64(total)),
		humanize.Comma(int64(totalFree)),
		100*(total-free)/total,
		label,
		fs,
		driveType(rootPathName))
	return nil
}

func cmdDiskFree(_ context.Context, cmd Param) (int, error) {
	bits, err := windows.GetLogicalDrives()
	if err != nil {
		return 0, err
	}
	fmt.Fprintf(cmd.Out(), "   %16s %16s %16s Use%%\n",
		"Available",
		"TotalNumber",
		"TotalNumberOfFree")

	count := 0
	for _, arg1 := range cmd.Args()[1:] {
		if err := df(arg1, cmd.Out()); err != nil {
			fmt.Fprintln(cmd.Err(), err)
		}
		count++
	}
	if count <= 0 {
		for d := 'A'; d <= 'Z'; d++ {
			if (bits & 1) != 0 {
				rootPathName := fmt.Sprintf("%c:\\", d)
				if err := df(rootPathName, cmd.Out()); err != nil {
					fmt.Fprintln(cmd.Err(), err)
				}
			}
			bits >>= 1
		}
	}
	return 0, nil
}
