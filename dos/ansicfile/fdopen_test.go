package ansicfile

import "testing"

func TestFdOpen(t *testing.T) {
	fp, err := FdOpen(1, "w")
	if err == nil {
		fp.Putc(byte('1'))
		fp.Putc(byte('\n'))
		fp.Close()
		print("OK: FdOpen(1)\n")
	} else {
		t.Fatalf("NG: %s", err.Error())
	}
}
