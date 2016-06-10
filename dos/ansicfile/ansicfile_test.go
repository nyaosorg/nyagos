package ansicfile

import (
	"bufio"
	"testing"
)

const sample_data = "12345678910"

func TestOpen(t *testing.T) {
	fp, err := Open("あ", "w")
	if err == nil {
		w := bufio.NewWriter(fp)
		for i := 0; i < 5; i++ {
			w.WriteString(sample_data)
			w.WriteString("\n")
		}
		w.Flush()
		fp.Close()
	} else {
		t.Fatalf("NG: Write Open(\"あ\") Failed by %s", err.Error())
	}
	fp, err = Open("あ", "r")
	if err == nil {
		r := bufio.NewScanner(fp)
		count := 0
		for r.Scan() {
			print("read: ", r.Text(), "\n")
			if r.Text() != sample_data {
				t.Fatalf("NG: Read data and Write data not equal('%s' != '%s') ",
					r.Text(), sample_data)
			}
			count++
		}
		if count != 5 {
			t.Fatalf("NG: data lack\n")
		}
		fp.Close()
	} else {
		t.Fatalf("NG: Read Open(\"あ\") Failed by %s", err.Error())
	}

	fp, err = Open("*", "w")
	if err == nil {
		t.Fatalf("NG: Open(\"*\") should failed")
	} else {
		print("OK: Open(\"*\") failed by " + err.Error() + "\n")
	}
}

// vim:set fenc=utf8:
