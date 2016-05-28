[![Build status](https://ci.appveyor.com/api/projects/status/bh7866s6oasvchpj?svg=true)](https://ci.appveyor.com/project/zetamatta/nyagos)

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
  * ハードリンク・シンボリックリンク・ジャンクションのリンク先を表示
* Lua によるカスタマイズ
  * Lua で内蔵コマンドを組込み
  * コマンドラインフィルター
  * コードページ文字列とUTF8とのコンバート関数
  * COM サポート

目次
----

### リリースノートと履歴

- [最新](Doc/release_note_ja.md)
- [4.0まで](Doc/history_4.0_ja.md)
- [4.1で変わったこと](Doc/since_4.1_ja.md)

### ドキュメント

1. [インストール](Doc/01-Install_ja.md)
2. [起動オプション](Doc/02-Options_ja.md)
3. [編集機能](Doc/03-Readline_ja.md)
4. [内蔵コマンド](Doc/04-Commands_ja.md)
5. [起動処理](Doc/05-Startup_ja.md)
6. [コマンドライン置換](Doc/06-Substitution_ja.md)
7. [Lua拡張](Doc/07-LuaFunctions_ja.md)
8. [付属のLuaスクリプト](Doc/08-LuaScripts_ja.md)
9. [アンインストール](Doc/09-Uninstall_ja.md)
10. [ビルド方法](Doc/10-Build_ja.md)

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
* [orthographic-pedant](https://github.com/orthographic-pedant)
* HABATA Katsuyuki
* [hisomura](https://github.com/hisomura)
* [tsuyoshicho](https://github.com/tsuyoshicho)
* [rane-hs](https://github.com/rane-hs)
* [hami-jp](https://github.com/hami-jp)
* [3bch](https://github.com/3bch)

開発者
------

* はやまかおる : [zetamatta](https://github.com/zetamatta) 

<!-- vim:set fenc=utf8 -->
