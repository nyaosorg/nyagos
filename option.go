package main

import "./interpreter"

func OptionParse(getArg func() (string, bool)) bool {
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
			}
		}
	}
	return true
}
