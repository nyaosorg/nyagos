ビルド方法
----------

次のソフトウェアが必要となります。

* [go 1.6 for windows](http://golang.org)
* [LuaBinaries 5.3.2 - Release 1 for Win32/64](http://luabinaries.sourceforge.net/download.html)
* [NYOLE 0.0.0.5 or later](https://github.com/zetamatta/nyole/releases) (任意です。無い場合、幾つかの Lua 拡張が動きませんが、nyagos.exe 自体は動作します)
* http://github.com/mattn/go-runewidth
* http://github.com/shiena/ansicolor
* http://github.com/atotto/clipboard
* http://github.com/mattn/go-isatty

`%GOPATH%` にて

    git clone https://github.com/zetamatta/nyagos nyagos
    cd nyagos

(32bitの場合)

    unzip PATH\TO\lua-5.3.2_Win32_bin.zip lua53.dll
    unzip PATH\TO\nyole-0.0.0.5.zip nyole.dll

(64bitの場合)

    unzip PATH\TO\lua-5.3.2_Win64_bin.zip lua53.dll
    unzip PATH\TO\nyole-0.0.0.5_x64.zip nyole.dll

最後に:

    make.cmd get
    make.cmd
    make.cmd install INSTALLDIR

make.cmd の使い方については `make.cmd help` を参照してください。

<!-- vim:set fenc=utf8: -->
