package lua

import (
	"testing"
)

func TestRawGet(t *testing.T) {
	L, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer L.Close()
	L.NewTable()
	L.PushInteger(1)
	L.PushString("ahaha")
	L.RawSet(-3)

	L.PushInteger(1)
	L.RawGet(-2)
	value, err := L.ToString(-1)
	if err != nil {
		t.Fatal(err.Error())
	}
	if value != "ahaha" {
		t.Fatalf("value differs (%s)\n", value)
	}
}
