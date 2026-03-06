package shell_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/nyaosorg/nyagos/internal/shell"
)

func TestAlreadyReportedError(t *testing.T) {
	a := shell.AlreadyReportedError{Err: io.EOF}

	if !errors.Is(a, io.EOF) {
		t.Fatal("AleadyReportedError failed to implement Unwrap():case 1")
	}

	b := shell.AlreadyReportedError{}
	if errors.Is(b, io.EOF) {
		t.Fatal("AleadyReportedError failed to implement Unwrap():case 2")
	}
}

func TestInterpret(t *testing.T) {
	ctx := context.Background()
	tempFilePath := filepath.Join(os.TempDir(), "hogehoge")
	var testShell = `cmd /c "echo 12345" > "%TEMP%\hogehoge"`
	if runtime.GOOS != "windows" {
		testShell = fmt.Sprintf(`sh -c "echo 12345%s" > "%s"`, "\r", tempFilePath)
	}
	_, err := shell.New().Interpret(ctx, testShell)
	if err != nil {
		t.Fatalf("Fail: %s", err.Error())
	}

	tempFileData, err := os.ReadFile(tempFilePath)
	if err != nil {
		t.Fatalf("Fail: `%s` not found", tempFilePath)
	}
	defer os.Remove(tempFilePath)
	if data := string(tempFileData); data != "12345\r\n" {
		t.Fatalf("Fail: %s's contents is expected as \"12345\\r\\n\", but %q",
			tempFilePath, string(tempFileData))
	}
}
