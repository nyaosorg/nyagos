package ansicfile

import "testing"

func TestOpen(t *testing.T) {
	fp, err := Open("あ", "w")
	if err == nil {
		Putc(byte('1'), fp)
		Close(fp)
	} else {
		t.Fatalf("NG: Open(\"あ\") Failed by %s", err.Error())
	}
	fp, err = Open("*", "w")
	if err == nil {
		t.Fatalf("NG: Open(\"*\") should failed")
	} else {
		print("OK: Open(\"*\") failed by " + err.Error() + "\n")
	}
}

// vim:set fenc=utf8:
