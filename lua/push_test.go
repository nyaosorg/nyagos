package lua

import (
	"fmt"
	"testing"
)

func TestUpValueIndex(t *testing.T) {
	print(UpValueIndex(0), "\n")
}

const testintvalue = -7

func TestPushInteger(t *testing.T) {
	L, err := New()
	if err != nil {
		t.Fatal("New() failed for" + err.Error())
		return
	}
	defer L.Close()
	L.PushInteger(testintvalue)
	L.SetGlobal("x")
	L.LoadString("return x")
	err = L.Call(0, 1)
	if err != nil {
		t.Fatal("Call() failed for " + err.Error())
		return
	}
	defer L.Pop(1)
	var str string
	str, err = L.ToString(-1)
	if err != nil {
		t.Fatal("ToString failed for " + err.Error())
		return
	}
	expected := fmt.Sprintf("%d", testintvalue)
	if str != expected {
		t.Fatalf("PushInteger failed '%s' != '%s'", str, expected)
	} else {
		t.Logf("PushInteger succeeded '%s' == '%s'", str, expected)
	}
}
