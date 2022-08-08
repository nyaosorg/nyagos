package nodos

// IsGui returns true if fname is Windows GUI Application
func IsGui(fname string) bool {
	if fname == "" {
		return false
	}
	result, err := isGui(fname)
	if err != nil {
		return false
	}
	return result
}
