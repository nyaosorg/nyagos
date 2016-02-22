package dos

import "io"

func ReadAnsiLine(f io.Reader) (string, error) {
	line := make([]byte, 0, 1024)
	var ch [1]byte
	for {
		n, err := f.Read(ch[:])
		if err != nil {
			return "", err
		}
		if n <= 0 || ch[0] == '\n' {
			break
		}
		if ch[0] != '\r' {
			line = append(line, ch[0])
		}
	}
	return AtoU(line)
}
