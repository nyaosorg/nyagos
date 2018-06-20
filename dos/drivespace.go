package dos

// GetDiskFreeSpace retunrs disk information.
//   rootPathName - string like "C:"
func GetDiskFreeSpace(rootPathName string) (free uint64, total uint64, totalFree uint64, err error) {
	rc, err1 := getDiskFreeSpaceEx(rootPathName, &free, &total, &totalFree)
	if rc == 0 {
		err = err1
	}
	return
}
