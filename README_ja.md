[![Build status](https://ci.appveyor.com/api/projects/status/bh7866s6oasvchpj?svg=true)](https://ci.appveyor.com/project/zetamatta/nyagos)
[![GoDoc](https://godoc.org/github.com/nyaosorg/nyagos?status.svg)](https://godoc.org/github.com/nyaosorg/nyagos)
[![Go Report Card](https://goreportcard.com/badge/github.com/nyaosorg/nyagos)](https://goreportcard.com/report/github.com/nyaosorg/nyagos)
[![Github latest Releases](https://img.shields.io/github/downloads/nyaosorg/nyagos/latest/total.svg)](https://github.com/nyaosorg/nyagos/releases/latest)

The Nihongo Yet Another GOing Shell
===================================

[&lt;English&gt;](./README.md) / **&lt;Japanese&gt;**

NYAGOS - Nihongo Yet Another GOing Shell は、Bash 風のコマンドライン編集機能と、Windows のファイルパスやバッチファイルとのシームレスな統合を兼ね備えた多機能なコマンドラインシェルです。Lua スクリプトによる柔軟なカスタマイズが可能で、最新の予測入力機能にも対応しています。

![demo-animation](./demo.gif)

### 主な機能

#### UNIX ライクなシェル動作
- **キーバインド**
  - デフォルトでは Bash に近いキーバインドを採用。
  - `%USERPROFILE%\.nyagos` に Lua スクリプトでカスタマイズ可能。
    ```lua
    nyagos.key.c_u = "KILL_WHOLE_LINE"
    ```
  - Lua 関数をキーに割り当て可能：
    ```lua
    nyagos.key.escape = function(this) nyagos.exec("start vim.exe") end
    ```
- **履歴とエイリアス**
  - `Ctrl-P` による履歴検索や `!` を用いたコマンド再実行が可能。
  - DOSKEY 風のエイリアス機能：
    ```lua
    nyagos.alias["g++"] = "g++.exe -std=gnu++17 $*"
    ```
  - Lua でエイリアスを実装可能：
    ```lua
    nyagos.alias["lala"] = function(args) nyagos.exec("ls", "-al", unpack(args)) end
    ```
- **カスタム補完（Bash 風の Tab 補完）**
  - コマンドごとに補完の定義が可能。
    ```lua
    nyagos.complete_for["go"] = function(args)
        if #args == 2 then
            return { "bug", "doc", "fmt", "install", "run", "version",
                     "build", "env", "generate", "list", "test", "vet",
                     "clean", "fix", "get", "mod", "tool" }
        else
            return nil -- ファイル補完
        end
    end
    ```
- **予測補完（PowerShell 7 風）**
  - 過去の履歴をもとに補完候補を提示。
  - `Ctrl-F` または右矢印キーで補完候補を採用可能。

#### Windows 互換機能
- **バッチファイルのシームレスな実行**
  - Windows のバッチファイル（`.bat` や `.cmd`）を CMD.exe 上で実行するのと同じように実行可能。
  - バッチファイル内で変更された環境変数やカレントディレクトリを適切に反映。
- **CMD.EXE ライクな機能**
  - Windows のパス表記（`C:\path\to\file`）に対応。
  - ドライブごとにカレントディレクトリを維持。
  - `copy` や `move` などの DOS 風の組み込みコマンドを搭載。
  - 追加の DLL 不要、レジストリ変更なし。

#### ユーザーエクスペリエンスの向上
- **カラー対応のコマンドライン**
- **Unicode サポート**
  - Windows の Unicode API を完全サポート。
  - クリップボードの Unicode 文字の貼り付け・編集が可能。
  - 特殊な Unicode リテラル `%U+XXXX%` やプロンプトマクロ `$Uxxxx` に対応。
- **組み込み `ls` コマンド**
  - カラー出力に対応（`-o` オプション）。
  - ハードリンク、シンボリックリンク、ジャンクションのターゲットパスを表示。
- **[SKK]（簡易かな漢字変換プログラム）をサポート - [セットアップガイド][SKKSetUpJa]**

### 対応プラットフォーム
- Windows 7, 8.1, 10, 11, Windows Server 2008 以降
- Linux（実験的対応）

[Video by @emisjerry](https://www.youtube.com/watch?v=WsfIrBWwAh0)

[SKK]: https://ja.wikipedia.org/wiki/SKK
[SKKSetUpJa]: doc/10-SetupSKK_ja.md

### ライセンス
NYAGOS は New BSD License のもとで使用・複製・改変が可能です。

インストール
------------

### バイナリをダウンロードする

安定版のバイナリは以下のページからダウンロードできます：

* https://github.com/nyaosorg/nyagos/releases

### Scoop でインストールする

```cmd
C:> scoop install nyagos
```

### Chocolatey でインストールする

```cmd
C:> choco install nyagos
```

### ソースコードからビルドする（スナップショット版）

最新版のスナップショットを試したい場合は、`go install` でインストールできます：

```cmd
C:> go install github.com/nyaosorg/nyagos@latest
```

> ⚠️ この方法でインストールされるのは開発中のスナップショット版であり、安定版ではありません。

目次
----

### リリースノート

[4.4.x](docs/release_note_ja.md)
/ [4.3.x](docs/history-4.3_ja.md)
/ [4.2.x](docs/history-4.2_ja.md)
/ [4.1.x](docs/history-4.1_ja.md)
/ [4.0.x](docs/history-4.0_ja.md)

### ドキュメント

1. [インストール](docs/01-Install_ja.md)
2. [起動オプション](docs/02-Options_ja.md)
3. [編集機能](docs/03-Readline_ja.md)
4. [内蔵コマンド](docs/04-Commands_ja.md)
5. [起動処理](docs/05-Startup_ja.md)
6. [コマンドライン置換](docs/06-Substitution_ja.md)
7. [Lua拡張](docs/07-LuaFunctions_ja.md)
8. [アンインストール](docs/08-Uninstall_ja.md)
9. [ビルド方法](docs/09-Build_ja.md)
10. [SKKの設定](docs/10-SetupSKK_ja.md) (v4.4.14以降)

ライセンス
----------

修正BSDライセンスに基いて、使用・コピー・改変が可能です。

謝辞
----

[nocd5](https://github.com/nocd5)
/ [mattn](https://github.com/mattn)
/ [hattya](https://github.com/hattya)
/ [shiena](https://github.com/shiena)
/ [atotto](https://github.com/atotto)
/ [ironsand](https://github.com/ironsand)
/ [kardianos](https://github.com/kardianos)
/ [malys](https://github.com/malys)
/ [pine613](https://github.com/pine613)
/ [NSP-0123456](https://github.com/NSP-0123456)
/ [hokorobi](https://github.com/hokorobi)
/ [amuramatsu](https://github.com/amuramatsu)
/ [spiegel-im-spiegel](https://github.com/spiegel-im-spiegel)
/ [rururutan](https://github.com/rururutan/)
/ [hogewest](https://github.com/hogewest)
/ [cagechi](https://github.com/cagechi)
/ [Matsuyanagi](https://github.com/Matsuyanagi)
/ [Shougo](https://github.com/Shougo)
/ [orthographic-pedant](https://github.com/orthographic-pedant)
/ HABATA Katsuyuki
/ [hisomura](https://github.com/hisomura)
/ [tsuyoshicho](https://github.com/tsuyoshicho)
/ [rane-hs](https://github.com/rane-hs)
/ [hami-jp](https://github.com/hami-jp)
/ [3bch](https://github.com/3bch)
/ [AoiMoe](https://github.com/aoimoe)
/ [DeaR](https://github.com/DeaR)
/ [gracix](https://github.com/gracix)
/ [orz--](https://github.com/orz--)
/ [zkangaroo](https://github.com/zkangaroo)
/ [maskedw](https://github.com/maskedw)
/ [tyochiai](https://github.com/tyochiai)
/ [masamitsu-murase](https://github.com/masamitsu-murase)
/ [hazychill](https://github.com/hazychill)
/ [erw7](https://github.com/erw7)
/ [tignear](https://github.com/tignear)
/ [crile](https://github.com/crile)
/ [fushihara](https://github.com/fushihara)
/ [ChiyosukeF](https://twitter.com/ChiyosukeF)
/ [beepcap](https://twitter.com/beepcap)
/ [tostos5963](https://github.com/tostos5963)
/ [sambatriste](https://github.com/sambatriste)
/ [terepanda](https://github.com/terepanda)
/ [Takmg](https://github.com/Takmg)
/ nu8 <!-- (https://github.com/nu8) -->
/ [tomato3713](https://github.com/tomato3713)
/ tGqmJHoJKqgK <!-- (https://github.com/tGqmJHoJKqgK) -->
/ [juggler999](https://github.com/juggler999)
/ [zztkm](https://github.com/zztkm)
/ [8exBCYJi5ATL](https://github.com/8exBCYJi5ATL)
/ [ousttrue](https://github.com/ousttrue)
/ [kgasawa](https://github.com/kgasawa)
/ [HAYASHI-Masayuki](https://github.com/HAYASHI-Masayuki)
/ [naoyaikeda](https://github.com/naoyaikeda)
/ [emisjerry](https://github.com/emisjerry)

開発者
------

* [hymkor - HAYAMA Kaoru](https://github.com/hymkor) (a.k.a zetamatta)
<!-- vim:set fenc=utf8 -->
