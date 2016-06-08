package dos

func DirName(path string) string {
	lastroot := -1
	for i, i_end := 0, len(path); i < i_end; i++ {
		switch path[i] {
		case '\\', '/', ':':
			lastroot = i
		}
	}
	if lastroot >= 0 {
		return path[0:(lastroot + 1)]
	} else {
		return ""
	}
}
