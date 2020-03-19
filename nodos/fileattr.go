package nodos

const (
	REPARSE_POINT = _REPARSE_POINT
)

func GetFileAttributes(path string) (uint32, error) {
	return getFileAttributes(path)
}

func SetFileAttributes(path string, attr uint32) error {
	return setFileAttributes(path, attr)
}
