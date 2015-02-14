The Nihongo Yet Another GOing Shell
===================================

[English](./readme.md) / Japanese

NYAGOS は Go と Lua で記述された Windows 用コマンドラインシェルです。

* UNIX風シェル
  * ヒストリ (Ctrl-P や ! マークによる)
  * エイリアス
  * ファイル名・コマンド名補完
* Unicodeサポート
  * Unicode文字をコピペ・編集可能
  * Unicodeリテラル %U+XXXX%
  * プロンプト向けマクロ $Uxxxx
* 内蔵ls
  * カラーサポート(-oオプション)
  * ジャンクション・シンボリックリンクを @ 表示(-Fオプション)
* Lua によるカスタマイズ
  * Lua で内蔵コマンドを組込み
  * コマンドラインフィルター
  * コードページ文字列とUTF8とのコンバート関数

インストール
------------

バイナリファイルは https://github.com/zetamatta/nyagos/releases よりダウンロード可能です。

    mkdir PATH\TO\INSTALLDIR
    cd PATH\TO\INSTALLDIR
    unzip PATH\TO\DOWNLOADDIR\nyagos-****.zip
    makeicon.cmd

バッチファイル `makeicon.cmd` はデスクトップにアイコンを作成します。

* [英語マニュアル](Doc/nyagos_en.md)
* [日本語マニュアル](Doc/nyagos_ja.md)

ビルド方法
----------

次のソフトウェアが必要となります。

* [go1.4.1 windows/386](http://golang.org)
* [Lua 5.3](http://www.lua.org)
* [tdm-gcc](http://tdm-gcc.tdragon.net)

`%GOPATH%` にて

    git clone https://github.com/zetamatta/nyagos nyagos
    cd nyagos

lua53.dll が既にある場合:

    copy PATH\TO\lua53.dll lua\.

さもなければ:

    tar zxvf PATH/TO/lua-5.3.0.tar.gz
    cd lua-5.3.0\src
    mingw32-make.exe mingw
    copy lua53.dll ..\..\..
    cd ..\..\..

最後に:

    make.cmd get
    make.cmd
    make.cmd install INSTALLDIR

make.cmd の使い方については `make.cmd help` を参照してください。

ライセンス
----------

修正BSDライセンスに基いて、使用・コピー・改変が可能です。

謝辞
----

* [nocd5](https://github.com/nocd5)
* [mattn](https://github.com/mattn)
* [hattya](https://github.com/hattya)
* [shiena](https://github.com/shiena)
* [atotto](https://github.com/atotto)
* [ironsand](https://github.com/ironsand)
* [kardianos](https://github.com/kardianos)
* [malys](https://github.com/malys)
* [pine613](https://github.com/pine613)

開発者
------

* はやまかおる : [zetamatta](https://github.com/zetamatta) 
