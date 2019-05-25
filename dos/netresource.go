package dos

import (
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
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

func (nr *NetResource) LocalName() string {
	if nr == nil {
		return ""
	}
	return u2str(nr.localName)
}

func (nr *NetResource) RemoteName() string {
	if nr == nil {
		return `\\`
	}
	return u2str(nr.remoteName)
}

func (nr *NetResource) Comment() string {
	if nr == nil {
		return ""
	}
	return u2str(nr.comment)
}
func (nr *NetResource) Provider() string {
	if nr == nil {
		return ""
	}
	return u2str(nr.provider)
}

func (nr *NetResource) Name() string       { return nr.RemoteName() }
func (nr *NetResource) Size() int64        { return 0 }
func (nr *NetResource) Mode() os.FileMode  { return 0555 }
func (nr *NetResource) ModTime() time.Time { return time.Time{} }
func (nr *NetResource) IsDir() bool        { return true }
func (nr *NetResource) Sys() interface{}   { return nr }

type NetResourceHandle struct {
	handle      uintptr
	netresource *NetResource
}

func (nr *NetResource) open() (*NetResourceHandle, error) {
	var handle uintptr

	rc, _, err := procWNetOpenEnum.Call(
		RESOURCE_GLOBALNET,
		RESOURCETYPE_DISK,
		0,
		uintptr(unsafe.Pointer(nr)),
		uintptr(unsafe.Pointer(&handle)))

	if rc != windows.NO_ERROR {
		return nil, err
	}
	return &NetResourceHandle{
		handle:      handle,
		netresource: nr,
	}, nil
}

func (this *NetResourceHandle) Readdir(_count int) ([]os.FileInfo, error) {
	var buffer [32 * 1024]byte
	count := int32(_count)
	size := len(buffer)
	rc, _, err := procWNetEnumResource.Call(
		uintptr(this.handle),
		uintptr(unsafe.Pointer(&count)),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)))

	if rc == windows.NO_ERROR {
		result := make([]os.FileInfo, 0, count)
		for i := int32(0); i < count; i++ {
			var p NetResource
			p = *(*NetResource)(unsafe.Pointer(&buffer[uintptr(i)*unsafe.Sizeof(p)]))
			result = append(result, &p)
		}
		return result, nil
	} else if rc == ERROR_NO_MORE_ITEMS {
		return nil, io.EOF
	} else {
		return nil, err
	}
}

func (this *NetResourceHandle) Close() error {
	procWNetCloseEnum.Call(this.handle)
	return nil
}

func (this *NetResourceHandle) Read([]byte) (int, error) {
	return 0, errors.New("not support")
}

func (this *NetResourceHandle) Seek(int64, int) (int64, error) {
	return 0, errors.New("not support")
}

func (this *NetResourceHandle) Stat() (os.FileInfo, error) {
	return this.netresource, nil
}

var _ http.File = (*NetResourceHandle)(nil)

func (nr *NetResource) Enum(callback func(*NetResource) bool) error {
	handle, err := nr.open()
	if err != nil {
		return err
	}
	defer handle.Close()

	for {
		records, err := handle.Readdir(-1)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		for _, record := range records {
			if nr1, ok := record.(*NetResource); ok && !callback(nr1) {
				return nil
			}
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
		Scope:       RESOURCE_GLOBALNET,
		Type:        RESOURCETYPE_DISK,
		DisplayType: RESOURCEDISPLAYTYPE_SERVER,
		Usage:       RESOURCEUSAGE_CONTAINER,
		remoteName:  name16,
	}, nil
}
