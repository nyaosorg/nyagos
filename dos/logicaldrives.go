package dos

var procGetLogicalDrives = kernel32.NewProc("GetLogicalDrives")

func GetLogicalDrives() (uint32, error) {
	rc, _, err := procGetLogicalDrives.Call()
	if rc == 0 {
		return 0, err
	}
	return uint32(rc), nil
}
