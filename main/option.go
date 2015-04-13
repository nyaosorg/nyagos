package main

import (
	"flag"
	"fmt"
	"os"

	"../interpreter"
	"../lua"
)

func optionParse(L *lua.Lua) bool {
	optionK := flag.String("k", "", "like `cmd /k`")
	optionC := flag.String("c", "", "like `cmd /c`")
	optionF := flag.String("f", "", "run lua script")
	optionE := flag.String("e", "", "run inline-lua-code")

	flag.Parse()

	result := true

	if *optionK != "" {
		interpreter.New().Interpret(*optionK)
	}
	if *optionC != "" {
		interpreter.New().Interpret(*optionC)
		result = false
	}
	if *optionF != "" {
		L.NewTable()
		L.PushString(*optionF)
		L.RawSetI(-2, 0)
		for i, arg1 := range flag.Args() {
			L.PushString(arg1)
			L.RawSetI(-2, lua.Integer(i+1))
		}
		L.SetGlobal("arg")
		err := L.Source(*optionF)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		result = false
	}
	if *optionE != "" {
		err := L.LoadString(*optionE)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			L.Call(0, 0)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		result = false
	}
	return result
}
