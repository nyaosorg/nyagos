package onexit

var functions = [](func()){}

func Register(f func()) {
	functions = append(functions, f)
}

func Done() {
	for i := len(functions) - 1; i >= 0; i-- {
		functions[i]()
	}
}
