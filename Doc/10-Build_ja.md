[English](./10-Build_en.md) / Japanese

ビルド方法
----------

次のソフトウェアが必要となります。
github.com 上のモジュールは `make.cmd get` でダウンロード可能です。

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

`%GOPATH%` にて

    git clone https://github.com/zetamatta/nyagos nyagos
    cd nyagos

(32bitの場合)

    unzip PATH\TO\lua-5.3.2_Win32_bin.zip lua53.dll

(64bitの場合)

    unzip PATH\TO\lua-5.3.2_Win64_bin.zip lua53.dll

最後に:

    make.cmd get
    make.cmd
    make.cmd install INSTALLDIR

make.cmd の使い方については `make.cmd help` を参照してください。

<!-- vim:set fenc=utf8: -->
