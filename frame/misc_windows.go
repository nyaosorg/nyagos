package frame

import (
	"github.com/zetamatta/nyagos/dos"
)

func coInitialize() {
	dos.CoInitializeEx(0, dos.COINIT_MULTITHREADED)
}

func coUnInitialize() {
	dos.CoUninitialize()
}
