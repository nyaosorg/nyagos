package dos

//go:generate go run $GOROOT/src/syscall/mksyscall_windows.go -output z$GOFILE $GOFILE

//sys CopyFile(src string,dst string,isFailIfExist bool)(n uint32,err error) = kernel32.CopyFileW
//sys getDiskFreeSpaceEx(rootPathName string,free *uint64,total *uint64,totalFree *uint64)(n uint32,err error) = kernel32.GetDiskFreeSpaceExW
//sys CoInitializeEx(res uintptr,opt uintptr) = ole32.CoInitializeEx
//sys CoUninitialize() = ole32.CoUninitialize
