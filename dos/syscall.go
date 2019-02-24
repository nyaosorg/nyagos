package dos

//go:generate go run $GOROOT/src/syscall/mksyscall_windows.go -output z$GOFILE $GOFILE

//sys getDiskFreeSpaceEx(rootPathName string,free *uint64,total *uint64,totalFree *uint64)(n uint32,err error) = kernel32.GetDiskFreeSpaceExW
