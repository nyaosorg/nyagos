package exename

//#include <windows.h>
import "C"

import "bytes"
import "unicode/utf16"

func Query() string {
	var pathW [C.MAX_PATH]C.WCHAR
	C.GetModuleFileNameW(nil, &pathW[0], C.MAX_PATH)

	var path16 [C.MAX_PATH]uint16
	for i := 0; pathW[i] != 0; i++ {
		path16[i] = (uint16)(pathW[i])
	}

	pathRune := utf16.Decode(path16[:])
	var buffer bytes.Buffer
	for _, ch := range pathRune {
		if ch == 0 {
			break
		}
		buffer.WriteRune(ch)
	}
	return buffer.String()
}
