package dos

//go:generate go run $GOROOT/src/syscall/mksyscall_windows.go -output zsyscall.go syscall.go

//sys copyFile(src string,dst string,isFailIfExist bool)(n uint32,err error) = kernel32.CopyFileW
//sys moveFileEx(src string,dst string,flag uintptr)(n uint32,err error) = kernel32.MoveFileExW
