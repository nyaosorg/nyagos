package lua

import (
	"fmt"
	"testing"
)

func TestToInterface(t *testing.T) {
	L, err := New()
	if err != nil {
		t.Fatal("New() failed for" + err.Error())
		return
	}
	L.OpenLibs()
	defer L.Close()
	L.LoadString(`t = { ["a"]="b",[2]=3 }`)
	err = L.Call(0, 0)
	if err != nil {
		t.Fatal("Call() failed for " + err.Error())
		return
	}
	L.GetGlobal("t")
	inte, err := L.ToInterface(-1)
	L.Pop(1)
	if err != nil {
		t.Fatal("ToInterface(-1) failed for " + err.Error())
		return
	}
	table, ok := inte.(map[interface{}]interface{})
	if !ok {
		t.Fatal("table cast failed")
		return
	}
	for key, val := range table {
		fmt.Printf("%v=%v\n", key, val)
	}
}
