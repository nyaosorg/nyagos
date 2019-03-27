package dos

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var procWNetAddConnection2W = mpr.NewProc("WNetAddConnection2W")
var procWNetCancelConnection = mpr.NewProc("WNetCancelConnection2W")

type _NetResource struct {
	Scope       uint32
	Type        uint32
	DisplayType uint32
	Usage       uint32
	LocalName   *uint16
	RemoteName  *uint16
	Comment     uintptr
	Provider    uintptr
}

func WNetAddConnection2(remote, local, user, pass string) (err error) {
	var rs _NetResource

	rs.LocalName, err = windows.UTF16PtrFromString(local)
	if err != nil {
		return
	}
	rs.RemoteName, err = windows.UTF16PtrFromString(remote)
	if err != nil {
		return
	}
	var _user *uint16
	if user == "" {
		_user = nil
	} else {
		_user, err = windows.UTF16PtrFromString(user)
		if err != nil {
			return
		}
	}
	var _pass *uint16
	if pass == "" {
		_pass = nil
	} else {
		_pass, err = windows.UTF16PtrFromString(pass)
		if err != nil {
			return
		}
	}

	rc, _, err := procWNetAddConnection2W.Call(
		uintptr(unsafe.Pointer(&rs)),
		uintptr(unsafe.Pointer(_pass)),
		uintptr(unsafe.Pointer(_user)),
		0)

	if rc != 0 {
		return err
	}
	return nil
}

func WNetCancelConnection2(name string, update bool, force bool) error {
	_name, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	var _update uintptr
	if update {
		_update = CONNECT_UPDATE_PROFILE
	}
	var _force uintptr
	if force {
		_force = 1
	}
	rc, _, err := procWNetCancelConnection.Call(
		uintptr(unsafe.Pointer(_name)),
		_update,
		_force)
	if rc != 0 {
		return err
	}
	return nil
}
