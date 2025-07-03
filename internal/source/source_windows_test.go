package source_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/nyaosorg/nyagos/internal/source"
)

func TestBatchCall(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "script.bat")

	fd, err := os.Create(scriptPath)
	if err != nil {
		t.Fatalf("os.Create(\"%s\")=%s", scriptPath, err.Error())
	}
	moveDir := os.Getenv("TEMP")
	fmt.Fprintf(fd, "cd \"%s\"\r\n", moveDir)
	io.WriteString(fd, "set BATCHTEST=SUCCESS\r\n")
	io.WriteString(fd, "exit /b 1\r\n")
	fd.Close()
	defer os.Remove(scriptPath)

	t.Setenv("BATCHTEST", "FAILURE")
	rc, err := source.ExecBatch(
		[]string{scriptPath},
		io.Discard,
		false,
		os.Stdin,
		os.Stdout,
		os.Stderr,
		nil)
	if err != nil {
		t.Fatalf("system.Run(\"%s\")=%d,%s", scriptPath, rc, err.Error())
	}
	if rc != 1 {
		t.Fatalf("system.Run(\"%s\")=%d,nil", scriptPath, rc)
	}
	batchTest := os.Getenv("BATCHTEST")
	if batchTest != "SUCCESS" {
		t.Fatalf("BATCHTEST=\"%s\" (expect \"SUCCESS\")", batchTest)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd():%s", err.Error())
	}
	if wd != moveDir {
		t.Fatalf("os.Getwd():\"%s\" != \"%s\"", wd, moveDir)
	}
}
