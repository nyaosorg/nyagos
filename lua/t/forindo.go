package main

import (
	".."
	"fmt"
)

func main() {
	L, err := lua.New()
	if err != nil {
		println(err)
		return
	}
	defer L.Close()
	L.OpenLibs()

	L.GetGlobal("_G")
	if !L.IsTable(-1) {
		println("_G is not a table")
		return
	}
	err = L.ForInDo(-1, func(LL lua.Lua) error {
		switch LL.GetType(-2) {
		case lua.LUA_TNUMBER:
			if val, err := LL.ToInteger(-2); err == nil {
				fmt.Printf("[%d]\n", val)
			} else {
				return err
			}
		case lua.LUA_TSTRING:
			if val, err := LL.ToString(-2); err == nil {
				fmt.Printf("[%s]\n", val)
			} else {
				return err
			}
		default:
			fmt.Printf("Unknown Type Key\n")
		}
		return nil
	})
	if err != nil {
		println(err)
	}
}
