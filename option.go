package main

import "./alias"
import "./interpreter"

func OptionParse(getArg func() (string, bool)) {
	for {
		arg, ok := getArg()
		if !ok {
			return
		}
		if arg[0] != '-' {
			continue
		}
		for _, o := range arg[1:] {
			switch o {
			case 'c', 'k':
				if fname, fnameOk := getArg(); fnameOk {
					interpreter.Interpret(fname, alias.Hook, nil)
				}
				if o == 'c' {
					return
				}
			}
		}
	}
}
