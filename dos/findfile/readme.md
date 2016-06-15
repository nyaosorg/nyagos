go-findfile
===========

Windows native FindFile-API's Wrapper.(`package findfile`)

func Walk
---------

        func Walk(pattern string, callback func(*FileInfo) bool) error
            func (this *FileInfo) Name() string
            func (this *FileInfo) Size() int64
            func (this *FileInfo) ModTime() time.Time
            func (this *FileInfo) Mode() os.FileMode
            func (this *FileInfo) Attribute() uint32
            func (this *FileInfo) IsReparsePoint() bool
            func (this *FileInfo) IsReadOnly() bool
            func (this *FileInfo) IsDir() bool
            func (this *FileInfo) IsHidden() bool
            func (this *FileInfo) IsSystem() bool

- `Walk` calls `callback` for each file which matches `pattern`. (not recursive)
- `findfile.FileInfo` is compatible to `os.FileInfo`.

func Glob,Globs
---------------

        func Glob(pattern string) ([]string, error)
        func Globs(patterns []string) []string

`Glob` and `Globs` expand filename matching with wildcard.

<!-- vim:set et: -->
