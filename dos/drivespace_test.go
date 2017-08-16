package dos

import (
	"testing"
)

func TestGetDiskFreeSpace(t *testing.T) {
	free, total, totalfree, err := GetDiskFreeSpace("C:")
	if err != nil {
		t.Fatal(err)
	}
	println("free=", free, "total=", total, "totalfree=", totalfree)
}
