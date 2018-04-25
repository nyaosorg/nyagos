package lua

// on 386
// sizeof(uint64) = 8
// sizeof(uintptr)= 4

func (value Integer) Expand(list []uintptr) []uintptr {
	return append(list, uintptr(value), uintptr(value>>32))
}
