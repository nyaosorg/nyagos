package commands

import (
	"context"
	"fmt"
	"strconv"

	"golang.org/x/sys/windows"

	"github.com/dustin/go-humanize"

	"github.com/zetamatta/nyagos/dos"
)

func driveType(rootPathName string) (uint32, string) {
	_rootPathName, err := windows.UTF16PtrFromString(rootPathName)
	if err != nil {
		return 0, "UNKNOWN"
	}

	t := windows.GetDriveType(_rootPathName)
	switch t {
	case windows.DRIVE_REMOVABLE:
		return t, "REMOVABLE"
	case windows.DRIVE_FIXED:
		return t, "FIXED"
	case windows.DRIVE_REMOTE:
		return t, "REMOTE"
	case windows.DRIVE_CDROM:
		return t, "CDROM"
	case windows.DRIVE_RAMDISK:
		return t, "RAMDISK"
	default:
		return t, "UNKNOWN"
	}
}

func df(rootPathName string) ([]string, error) {
	label, fs, err := dos.VolumeName(rootPathName)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", rootPathName, err)
	}
	free, total, totalFree, err := dos.GetDiskFreeSpace(rootPathName)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", rootPathName, err)
	}
	driveTypeId, driveTypeStr := driveType(rootPathName)
	var uncPath string
	if driveTypeId == windows.DRIVE_REMOTE {
		uncPath, _ = dos.WNetGetConnectionUTF16a(uint16(rootPathName[0]))
	}

	return []string{
		rootPathName,
		humanize.Comma(int64(free)),
		humanize.Comma(int64(total)),
		humanize.Comma(int64(totalFree)),
		strconv.FormatUint(100*(total-free)/total, 10),
		fs,
		driveTypeStr,
		label,
		uncPath,
	}, nil
}

func cmdDiskFree(_ context.Context, cmd Param) (int, error) {
	bits, err := windows.GetLogicalDrives()
	if err != nil {
		return 0, err
	}
	dfs := [][]string{
		[]string{"", "Available", "Total", "TotalFree", "Use%"},
	}

	for _, arg1 := range cmd.Args()[1:] {
		if df1, err := df(arg1); err != nil {
			fmt.Fprintln(cmd.Err(), err)
		} else {
			dfs = append(dfs, df1)
		}
	}
	if len(dfs) <= 1 {
		for d := 'A'; d <= 'Z'; d++ {
			if (bits & 1) != 0 {
				rootPathName := fmt.Sprintf("%c:\\", d)
				if df1, err := df(rootPathName); err != nil {
					fmt.Fprintln(cmd.Err(), err)
				} else {
					dfs = append(dfs, df1)
				}
			}
			bits >>= 1
		}
	}

	colsiz := []int{}
	for _, df1 := range dfs {
		for i, s := range df1 {
			if i >= len(colsiz) {
				colsiz = append(colsiz, len(s))
			} else if len(s) > colsiz[i] {
				colsiz[i] = len(s)
			}
		}
	}
	for _, df1 := range dfs {
		for i, s := range df1 {
			if i > 0 {
				cmd.Out().Write([]byte{' '})
			}
			if i == len(df1)-1 {
				fmt.Fprintln(cmd.Out(), s)
			} else if i >= 5 {
				fmt.Fprintf(cmd.Out(), "%-*s", colsiz[i], s)
			} else {
				fmt.Fprintf(cmd.Out(), "%*s", colsiz[i], s)
			}
		}
	}

	return 0, nil
}
