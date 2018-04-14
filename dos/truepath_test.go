package dos

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestTruePath(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
		return
	}

	c := exec.Command("cmd", "/c", "mklink /J sub ..")
	if err := c.Run(); err != nil {
		t.Fatal(err)
		return
	}
	defer os.Remove("sub")

	result := TruePath(`sub`)
	expect := filepath.Dir(wd)
	if expect != result {
		t.Fatalf("Failed: TruePath(`sub`) -> %s (not %s)", result, expect)
		return
	}

	if err := os.Mkdir("sub/hoge",0777); err != nil {
		t.Fatalf("Failed: could not mkdir: %s", err.Error())
		return
	}
	defer os.Remove("sub/hoge")
	result = TruePath(`sub/hoge`)
	println(result)
}
