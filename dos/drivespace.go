package dos

func GetDiskFreeSpace(rootPathName string) (free uint64, total uint64, totalFree uint64, err error) {
	rc, err1 := getDiskFreeSpaceEx(rootPathName, &free, &total, &totalFree)
	if rc == 0 {
		err = err1
	}
	return
}
