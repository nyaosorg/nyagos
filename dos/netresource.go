package dos

import (
	"fmt"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var procWNetOpenEnum = mpr.NewProc("WNetOpenEnumW")
var procWNetEnumResource = mpr.NewProc("WNetEnumResourceW")
var procWNetCloseEnum = mpr.NewProc("WNetCloseEnum")

type NetResource struct {
	Scope       uint32
	Type        uint32
	DisplayType uint32
	Usage       uint32
	localName   *uint16
	remoteName  *uint16
	comment     *uint16
	provider    *uint16
}

func u2str(u *uint16) string {
	if u == nil {
		return ""
	}
	buffer := make([]uint16, 0, 100)
	for *u != 0 {
		buffer = append(buffer, *u)
		u = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(u)) + 2))
	}
	return windows.UTF16ToString(buffer)
}

func (nr *NetResource) LocalName() string  { return u2str(nr.localName) }
func (nr *NetResource) RemoteName() string { return u2str(nr.remoteName) }
func (nr *NetResource) Comment() string    { return u2str(nr.comment) }
func (nr *NetResource) Provider() string   { return u2str(nr.provider) }

func (nr *NetResource) Enum(callback func(*NetResource) bool) error {
	var handle uintptr

	rc, _, err := procWNetOpenEnum.Call(
		RESOURCE_GLOBALNET,
		RESOURCETYPE_DISK,
		0,
		uintptr(unsafe.Pointer(nr)),
		uintptr(unsafe.Pointer(&handle)))
	if rc != windows.NO_ERROR {
		return fmt.Errorf("NetOpenEnum: %s", err)
	}
	defer procWNetCloseEnum.Call(handle)
	for {
		var buffer [16 * 1024]byte
		count := int32(-1)
		size := len(buffer)
		rc, _, err := procWNetEnumResource.Call(
			handle,
			uintptr(unsafe.Pointer(&count)),
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(unsafe.Pointer(&size)))

		if rc == windows.NO_ERROR {
			for i := int32(0); i < count; i++ {
				var p *NetResource
				p = (*NetResource)(unsafe.Pointer(&buffer[uintptr(i)*unsafe.Sizeof(*p)]))
				if !callback(p) {
					return nil
				}
			}
		} else if rc == ERROR_NO_MORE_ITEMS {
			return nil
		} else {
			return fmt.Errorf("NetEnumResource: %s", err)
		}
	}
}

func WNetEnum(callback func(nr *NetResource) bool) error {
	var nr *NetResource
	return nr.Enum(callback)
}

func EachMachine(callback func(*NetResource) bool) error {
	return WNetEnum(func(all *NetResource) bool {
		if strings.EqualFold(all.RemoteName(), "Microsoft Windows Network") {
			all.Enum(func(domain *NetResource) bool {
				domain.Enum(callback)
				return true
			})
			return false
		} else {
			return true
		}
	})
}

func EachMachineNode(name string, callback func(*NetResource) bool) error {
	return EachMachine(func(machine *NetResource) bool {
		if strings.EqualFold(name, machine.RemoteName()) {
			machine.Enum(callback)
		}
		return true
	})
}
