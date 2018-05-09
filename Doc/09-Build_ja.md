[English](./09-Build_en.md) / Japanese

ビルド方法
----------

次のソフトウェアが必要となります。
これらは `make.cmd get` でダウンロード可能です。

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

`%GOPATH%` にて

    git clone https://github.com/zetamatta/nyagos src/github.com/zetamatta/nyagos
    cd src/github.com/zetamatta/nyagos

(安定板の時)

    git checkout master

(最新版の時)

    git checkout develop

    make.cmd get
    make.cmd
    make.cmd install INSTALLDIR

make.cmd の使い方については `make.cmd help` を参照してください。

Lua を使わない最小バージョンをビルドする
----------------------------------------

    cd nyagos/ngs
    go build

lua53.dll を Lua のエンジンとして使用するバージョンをビルドする
---------------------------------------------------------------

4.3 から Lua エンジンは GopherLua に切り変わりましたが、
lua53.dll をまだ使用することもできます。

    cd nyagos/mains
    go build

<!-- vim:set fenc=utf8: -->
