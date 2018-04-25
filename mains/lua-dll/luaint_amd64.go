package lua

// on amd64
// sizeof(uint64) = 8
// sizeof(uintptr)= 8

func (value Integer) Expand(list []uintptr) []uintptr {
	return append(list, uintptr(value))
}
