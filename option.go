package main

import "fmt"
import "os"

import "./interpreter"
import "./lua"

func OptionParse(L *lua.Lua, getArg func() (string, bool)) bool {
	for {
		arg, ok := getArg()
		if !ok {
			return true
		}
		if arg[0] != '-' {
			continue
		}
		for _, o := range arg[1:] {
			switch o {
			case 'c', 'k':
				if arg1, ok := getArg(); ok {
					interpreter.New().Interpret(arg1)
				}
				if o == 'c' {
					return false
				}
			case 'f':
				if script, scriptOk := getArg(); scriptOk {
					L.NewTable()
					L.PushString(script)
					L.RawSetI(-2, 0)
					for i := 1; true; i++ {
						arg1, ok := getArg()
						if !ok {
							break
						}
						L.PushString(arg1)
						L.RawSetI(-2, i)
					}
					L.SetGlobal("arg")
					err := L.Source(script)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
				}
				return false
			}
		}
	}
	return true
}
