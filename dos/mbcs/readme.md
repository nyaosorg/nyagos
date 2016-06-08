mbcs
====

`mbcs` is the library for the programming language Go for Windows,
to convert characters between the current codepage and UTF8

	var ansi []byte
	var ansi_err error

	ansi, ansi_err = mbcs.UtoA("UTF8文字列")
	if ansi_err != nil {
		fmt.Fprintln(os.Stderr, ansi_err)
		return
	}

	var utf8 string
	var utf8_err error

	utf8, utf8_err = mbcs.AtoU(ansi)
	if utf8_err != nil {
		fmt.Fprintln(os.Stderr, utf8_err)
		return
	}
	fmt.Printf("Ok: %s\n", utf8)

<!-- vim:set fenc=utf8: -->
