The Nihongo Yet Another GOing Shell
===================================

[English](./readme.md) / Japanese

NYAGOS は Go と Lua で記述された Windows 用コマンドラインシェルです。

* UNIX風シェル
  * Emacs風キーバインド
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

アンインストール
----------------

UNZIP で展開されたファイルと %APPDATA%\NYAOS.ORG 以下、デスクトップ
アイコンを削除してください。NYAGOS.exe はレジストリを読み書きしません。

ビルド方法
----------

次のソフトウェアが必要となります。

* [go 1.5 for windows](http://golang.org)
* [LuaBinaries 5.3 - Release 1 for Win32/64](http://luabinaries.sourceforge.net/download.html)
* [NYOLE 0.0.0.5 or later](https://github.com/zetamatta/nyole/releases) (任意です。無い場合、幾つかの Lua 拡張が動きませんが、nyagos.exe 自体は動作します)

`%GOPATH%` にて

    git clone https://github.com/zetamatta/nyagos nyagos
    cd nyagos

(32bitの場合)

    unzip PATH\TO\lua-5.3_Win32_bin.zip lua53.dll
    unzip PATH\TO\nyole-0.0.0.5.zip nyole.dll

(64bitの場合)

    unzip PATH\TO\lua-5.3_Win64_bin.zip lua53.dll
    unzip PATH\TO\nyole-0.0.0.5_x64.zip nyole.dll

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
* [NSP-0123456](https://github.com/NSP-0123456)
* [hokorobi](https://github.com/hokorobi)
* [amuramatsu](https://github.com/amuramatsu)
* [spiegel-im-spiegel](https://github.com/spiegel-im-spiegel)
* [rururutan](https://github.com/rururutan/)
* [hogewest](https://github.com/hogewest)
* [cagechi](https://github.com/cagechi)
* [Matsuyanagi](https://github.com/Matsuyanagi)
* [Shougo](https://github.com/Shougo)

開発者
------

* はやまかおる : [zetamatta](https://github.com/zetamatta) 
