package dos

import (
	"fmt"
	"testing"
)

func TestGetLogicalDrives(t *testing.T) {
	map1 := GetLogicalDrives()
	print("GetLogicalDrives\n")
	for _, val := range map1 {
		fmt.Printf("%c:\n", val)
	}
}

func TestGetDiskFreeSpaceEx(t *testing.T) {
	free, all, free2, err := GetDiskFreeSpaceEx("C:")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("C: Free %d All %d, Free2 %d\n", free, all, free2)
}
