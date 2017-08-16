package dos

import (
	"fmt"
	"testing"
)

func TestGetLogicalDrives(t *testing.T) {
	bits, err := GetLogicalDrives()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 'z'-'a'; i++ {
		if (bits & 1) != 0 {
			fmt.Printf("%c:\n", 'A'+i)
		}
		bits >>= 1
	}

}
