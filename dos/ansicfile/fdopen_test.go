package ansicfile

import "testing"

func TestFdOpen(t *testing.T) {
	fp, err := FdOpen(1, "w")
	if err == nil {
		fp.Write([]byte{'1', '\n'})
		fp.Close()
		print("OK: FdOpen(1)\n")
	} else {
		t.Fatalf("NG: %s", err.Error())
	}
}
