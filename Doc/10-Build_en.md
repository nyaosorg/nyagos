Build
-----

These sortwares are required.

* [go 1.6 for windows](http://golang.org)
* [LuaBinaries 5.3.2 - Release 1 for Win32/64](http://luabinaries.sourceforge.net/download.html)
* [NYOLE 0.0.0.5 or later](https://github.com/zetamatta/nyole/releases) (Optionally. Without it, some lua-extensions does not work, but nyagos.exe itself is available.)
- http://github.com/mattn/go-runewidth
- http://github.com/shiena/ansicolor
- http://github.com/atotto/clipboard
- http://github.com/mattn/go-isatty

On `%GOPATH%` folder,

    git clone https://github.com/zetamatta/nyagos nyagos
    cd nyagos

For 32bit ,

    unzip PATH\TO\lua-5.3.2_Win32_bin.zip lua53.dll
    unzip PATH\TO\nyole-0.0.0.5.zip nyole.dll

For 64bit ,

    unzip PATH\TO\lua-5.3.2_Win64_dllw4_lib.zip lua53.dll
    unzip PATH\TO\nyole-0.0.0.5_x64.zip nyole.dll

Finally

    make.cmd get
    make.cmd
    make.cmd install INSTALLDIR

How to use make.cmd is shown with `make.cmd help`
