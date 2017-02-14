English / [Japanese](./09-Build_ja.md)

Build
-----

These sortwares are required.
Modules on github.com are able to be downloaded with `make.cmd get`.

* [go 1.7 for windows](http://golang.org)
* [LuaBinaries 5.3.2 - Release 1 for Win32/64](http://luabinaries.sourceforge.net/download.html)
* https://github.com/atotto/clipboard
* https://github.com/dustin/go-humanize
* https://github.com/go-ole/go-ole
* https://github.com/josephspurrier/goversioninfo
* https://github.com/jteeuwen/go-bindata
* https://github.com/mattn/go-colorable
* https://github.com/mattn/go-isatty
* https://github.com/mattn/go-runewidth
* https://github.com/zetamatta/go-ansicfile
* https://github.com/zetamatta/go-findfile
* https://github.com/zetamatta/go-getch
* https://github.com/zetamatta/go-mbcs

On `%GOPATH%` folder,

    git clone https://github.com/zetamatta/nyagos nyagos
    cd nyagos

For 32bit ,

    unzip PATH\TO\lua-5.3.2_Win32_bin.zip lua53.dll

For 64bit ,

    unzip PATH\TO\lua-5.3.2_Win64_dllw4_lib.zip lua53.dll

Finally

    make.cmd get
    make.cmd
    make.cmd install INSTALLDIR

How to use make.cmd is shown with `make.cmd help`
