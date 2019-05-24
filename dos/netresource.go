package dos

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
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
		var buffer [32 * 1024]byte
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

var rxServerPattern = regexp.MustCompile(`^\\\\[^\\/]+$`)

var netlock sync.RWMutex

func EnumFileServer(callback func(*NetResource) bool) error {
	var me func(*NetResource) bool
	var wg sync.WaitGroup
	me = func(nr *NetResource) bool {
		if rxServerPattern.MatchString(nr.RemoteName()) {
			netlock.Lock()
			callback(nr)
			netlock.Unlock()
		} else {
			nr1 := *nr
			wg.Add(1)
			go func() {
				nr1.Enum(me)
				wg.Done()
			}()
		}
		return true
	}
	err := WNetEnum(me)
	wg.Wait()
	return err
}

func NewFileServer(name string) (*NetResource, error) {
	if !strings.HasPrefix(name, `\\`) {
		name = `\\` + name
	}
	name16, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}
	return &NetResource{
		Scope:       2,
		Type:        1,
		DisplayType: 2,
		Usage:       2,
		remoteName:  name16,
	}, nil
}
