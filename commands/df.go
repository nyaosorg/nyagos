package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/dustin/go-humanize"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/shell"
)

func df(rootPathName string, w io.Writer) error {
	free, total, totalFree, err := dos.GetDiskFreeSpace(rootPathName)
	if err != nil {
		return fmt.Errorf("%s: %s", rootPathName, err)
	}
	fmt.Fprintf(w, "%s %20s %20s %20s\n",
		rootPathName,
		humanize.Comma(int64(free)),
		humanize.Comma(int64(total)),
		humanize.Comma(int64(totalFree)))
	return nil
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
