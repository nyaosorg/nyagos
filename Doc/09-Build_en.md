English / [Japanese](./09-Build_ja.md)

Build
-----

These sortwares are required.
They are able to be downloaded with `make.cmd get`.

* [go 1.10 for windows](http://golang.org)
* https://github.com/atotto/clipboard
* https://github.com/dustin/go-humanize
* https://github.com/go-ole/go-ole
* https://github.com/hillu/go-pefile
* https://github.com/josephspurrier/goversioninfo
* https://github.com/mattn/go-colorable
* https://github.com/mattn/go-isatty
* https://github.com/mattn/go-runewidth
* https://github.com/mattn/msgbox
* https://github.com/yuin/gopher-lua
* https://github.com/zetamatta/go-ansicfile
* https://github.com/zetamatta/go-box
* https://github.com/zetamatta/go-findfile
* https://github.com/zetamatta/go-getch
* https://github.com/zetamatta/go-mbcs
* https://golang.org/x/sys/windows

    go get github.com/zetamatta/nyagos
    cd "%GOPATH%\src\github.com\zetamatta\nyagos"

(for stable version)

    git checkout master

(for latest version)

    git checkout develop

    make.cmd get
    make.cmd
    make.cmd install INSTALLDIR

How to use make.cmd is shown with `make.cmd help`

Build the minimum version (no-lua)
----------------------------------

    cd nyagos/ngs
    go build

Build the version using lua53.dll as lua-engine
-----------------------------------------------

Since 4.3, GopherLua is used as the lua engine,
but we can still use lua53.dll.

    cd nyagos/mains
    go build

<!-- vim:set fenc=utf8: -->
