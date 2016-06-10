package ansicfile

import (
	"bufio"
	"testing"
)

func TestOpen(t *testing.T) {
	fp, err := Open("あ", "w")
	if err == nil {
		w := bufio.NewWriter(fp)
		w.WriteString("12345678910")
		w.Flush()
		fp.Close()
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
