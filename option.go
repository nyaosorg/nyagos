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
				if fname, fnameOk := getArg(); fnameOk {
					interpreter.Interpret(fname, nil)
				}
				if o == 'c' {
					return false
				}
			}
		}
	}
	return true
}
