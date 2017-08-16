package dos

import (
	"testing"
)

func TestGetDriveType(t *testing.T) {
	rc, err := GetDriveType("C:\\")
	if err != nil {
		t.Fatal(err)
	}
	switch rc {
	case DRIVE_REMOVABLE:
		println("REMOVABLE")
	case DRIVE_FIXED:
		println("FIXED")
	case DRIVE_REMOTE:
		println("REMOTE")
	case DRIVE_CDROM:
		println("CDROM")
	case DRIVE_RAMDISK:
		println("RAMDISK")
	default:
		println("???")
	}
}
