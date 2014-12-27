package dos

import "syscall"

var kernel32 = syscall.NewLazyDLL("kernel32")
