package nodos

func CoInitializeEx(res uintptr, opt uintptr) {
	coInitializeEx(res, opt)
}

func CoUninitialize() {
	coUninitialize()
}
