Release notes
=============
( [English](release_note_en.md) / **Japanese** )

- readline の次回プロンプトに挿入する初期テキストを設定する `nyagos.setnextline(STR)` を実装 (#458, #466, thanks to @emisjerry)

- `nyagos.d/catalog/complete-jj.lua` (#473, #474, Thanks to @tsuyoshicho)
    - `jj` のサブコマンド補完をv0.35 ベースに更新
    - 改行コードが LF になっていたので、CRLF となるよう、生成スクリプト make-complete-jj.lua を修正

4.4.18\_1
----------
Nov 23, 2025

### 不具合修正

- Lua関数 `io.open(..., "w")` でファイルをオープンした時、ファイルが truncate されていなかった不具合を修正 (#471, #472 Thanks to @emisjerry)

### 内部的な変更

- [github.com/nyaosorg/go-box] を v2.2.1 から v3.0.0 にバージョンアップした (#453)
- `nyagos.getkeys` 内で使用していたが、Deprecated となっていた [readline.GetKey] を [go-ttyadapter]/tty8 に置き換えた (#454)
- nilに対するlen()は0だから、nil チェックを省略すべきという staticcheck の warning を解消した (#455)
- go-readline-ny での冗長な文字幅データのキャッシュ機能廃止にともなって Deprecated となった関数を使っていた、メンテ・調査用の非公開Lua関数 `nyagos.resetcharwidth`, `nyagos.setrunewidth` を廃止した (#457)

[readline.GetKey]: https://pkg.go.dev/github.com/nyaosorg/go-readline-ny@v1.11.0#GetKey
[go-ttyadapter]: https://github.com/nyaosorg/go-ttyadapter
[github.com/nyaosorg/go-box]: https://github.com/nyaosorg/go-box

4.4.18\_0
---------
Oct 25, 2025

- `nyagos.d/` 直下の Lua スクリプトを実行ファイルに組み込んだ
    - 他環境への展開が容易になり、各位の .nyagos と nyagos.exe のコピーだけで利用可能になりました。
    - `nyagos.d/` 直下のスクリプトは自動で読まなくなりました。かわりに `%APPDATA%/NYAOS_ORG/nyagos.d` を使ってください
    - `nyagos.d/catalog` 以下のスクリプトは従来どおりの扱いです。これらを利用する場合は、あわせてコピーが必要です。
- Windows/arm64 向けバイナリをビルドができるようにした(ビルドのみで動作は未検証)
- 環境変数名の補完において候補リストが表示されなかった不具合を修正
- SKK かな漢字変換ライブラリ: go-readline-skk を v0.6.0 へ更新:
    - 変換結果にスラッシュを含む単語も変換・単語登録できるようにした
    - emacslisp で書かれた変換結果について `(concat)`, `(pwd)`, `(substring)`, `(skk-current-date)` 程度は評価できるようにした (`(lambda)` はまだ)
- [go-readline-ny] を v1.11.0 へ更新
    - [#452] への対応のため以下を実施 (Thanks @emisjerry ) 
        - キーを表す識別子・シンボル文字列を追加 
            | シンボル       | キー組み合わせ  |
            |----------------|-----------------|
            |`"C_PAGEDOWN"`  |`Ctrl`+`PageDown`|
            |`"C_PAGEUP"`    |`Ctrl`+`PageUp`  |
            |`"C_HOME"`      |`Ctrl`+`Home`    |
            |`"C_END"`       |`Ctrl`+`End`     |
        - 初期キー設定を追加  
            |キー組み合わせ| 機能 |
            |--------------|------|
            |`Ctrl`+`Home` | 先頭からカーソル位置までを削除(`Ctrl`+`U`と等価) |
            |`Ctrl`+`End`  | カーソル位置から末尾までの削除(`Ctrl`+`K`と等価) |
- [mattn/go-runewidth] を v0.0.16 から v0.0.19 へ更新

[go-readline-ny]: https://github.com/nyaosorg/go-readline-ny
[#452]: https://github.com/nyaosorg/nyagos/issues/452
[mattn/go-runewidth]: https://github.com/mattn/go-runewidth

4.4.17\_2
---------
Jul 27, 2025

- コマンド処理中に複数回 Ctrl-C が押下された際、二回目以降が次のコマンド実行時に処理されてしまう問題を修正  
  ( [#449]: Thanks to [@fushihara] )

[#449]: https://github.com/nyaosorg/nyagos/issues/449
[@fushihara](https://github.com/fushihara)

4.4.17\_1
---------
Jul 3, 2025

- `%USERPROFILE%` が `C:\Users\foo` の時、`C:\Users\fool` へ chdir すると、プロンプトの `$P` が `~l` となってしまう不具合を修正
- 環境変数名補完の際、変数名よりも前の位置に日本語など１文字が２バイト以上の文字が存在するとクラッシュする不具合を修正。環境変数名の開始を示すパーセント文字の位置を指すインデックスの単位がバイト単位だったのを文字個数位置と取り違えていた

4.4.17\_0
---------
May 7, 2025


- NYAGOS のコピーバッファは、これまで OS のクリップボードと連動していたが、デフォルトでは連動しないように変更。クリップボードと連動させるには、nyagos.option.clipboard = true を設定する
- シンタックスハイライト
    - 対象とする環境変数は英字から始まる英数字のみとした
    - 対象とするオプションは先頭だけでなく、どこでも `-` を含められるようにした
- 入力予測機能の色をカスタマイズできるようにした。  
    - 例: `nyagos.option.predict_color = "\027[0;31;1m"` (赤に設定)
    - 変更を反映させるには、.nyagos に記載後、nyagos の再起動が必要
- v4.4.15 より非推奨としていた設定ファイル `_nyagos` を廃止
- `set -o` で現在の設定内容を Lua 構文として表示するようになった
- 実行ファイルにアイコンを複数含め、ショートカット作成時に選択できるようにした
- ドキュメントの間違いを訂正
    - `nyagos.option.cleanup_buffer` は削除されていたのに記載が残っていた

4.4.16\_0
----------
Oct 14, 2024

### 廃止・非推奨

- `nyagos.d/catalog/neco.lua` を削除
- Lua関数: `nyagos.msgbox` を削除
- grep.exe が存在しない時、勝手に findstr.exe を grep の別名にする動作を削除

### 新機能

#### 一行入力

- PowerShell 7 風の入力予測機能の実装
    - 一文字以上入力すると、その内容で始まる履歴の最新エントリをインラインで表示(青の斜体)
    - `→` もしくは `Ctrl-F` で今の予測表示内容を採用する
    - デフォルトでオン。起動オプション `--no-predict` もしくは `nyagos.option.prediction=false` でオフ
- Ctrl-P/N: 履歴を切り替えるときに変更したエントリを保存し、(Enterが入力されるまでは)再度切り替えたときに復元するようにした

#### Lua 拡張

- キー入力の最初のコードの Unicode しか返さなくなっていた nyagos.getkey のかわりに、入力キーを`"\027[A"` といった[キーシーケンス]で返す nyagos.getkeys() を実装した
- `nyagos.key[KEY]=function(this)...end` の中で使える機能を拡充
    - KEY として従来の`"BACKSPACE"`, `"UP"` などの名前の他、`"\007"`, `"\027[A"` などの[キーシーケンス]も使えるようにした
    - `this:eval("キーシーケンス")` で[キーシーケンス]に設定された機能を呼び出せるようにした
    - 更新内容を画面に反映するメソッド`this:repaint()`を追加
    - 更新系のメソッドを呼び出した際に`this.pos`と`this.text` を自動的に更新するようにした
- `make` コマンドのエントリ補完を追加: `require "makefile-complete"` にて有効化
- `jj` コマンドのサブコマンド名補完を追加: `require "complete-jj"` にて有効化
- `gmnlisp` コマンドの引用機能を追加: `require "gmnlisp"` にて有効化
    - `@(Lispコマンド)` を gmnlisp で処理した結果に置換
    - コマンド先頭が `(` で始まっていた場合、gmnlisp.exe で実行
- UNIX風シングルクォーテーション機能を追加: `require "sq2dq"` にて有効化
    - `'..".."..'` を `"..\"..\".."` へ置換
    - `"..'..'.."` はそのまま (二重引用符内の一重引用符は変換しない)

[キーシーケンス]: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#input-sequences

### ドキュメント

- readme.md → README.md でのファイル名変更でドキュメント間のリンクが切れてしまっていた点の修正 ( Thx @HAYASHI-Masayuki )

4.4.15\_1
---------
May 2, 2024

当バージョンのバイナリは Go 1.20.14 でビルドしました。  
サポート対象は Windows 7, 8.1, 10, 11, WindowsServer 2008以降, Linux となります。

- (#442) 4.4.15\_0 で `nyagos.d/catalog` 以下へ移動させて自動ロード対象外とした `nyagos.d/aliasandset.lua` を元の場所へ戻し、Luaコード: `set "ENV=VALUE" `, `alias "NAME=DEFINE"` を再び、そのまま使えるようにした (Thx @naoyaikeda )
- Windows の日付の設定に曜日が含まれているとき %DATE% の結果がおかしくなる不具合を修正 (`2024/04/19 金` と出て欲しいのに `2024/04/19 1919` と出てしまう)
- 一行入力のキーハンドル関数の中で `this:replacefrom(0,...)` を呼び出すと、ランタイムエラー: `runtime error: slice bounds out of range [-1:]` でクラッシュする不具合を修正。かわりに {nil,エラーメッセージ} を返すようにした。

```lua
-- On 4.4.15, typing C-I causes crash of nyagos.
nyagos.key["C-I"] = function(this)
    assert(this:replacefrom(0,"XXXXX"))
end
```

4.4.15\_0
---------
Apr 7, 2024

当バージョンのバイナリは Go 1.20.14 でビルドしました。  
サポート対象は Windows 7, 8.1, 10, 11, WindowsServer 2008以降, Linux となります。

### 廃止・非推奨

- Lua 文法ではない設定ファイル `_nyagos` は非推奨とし添付をやめた。使用されている場合は実行するが、警告を表示するようにした
- `use` は非推奨とした
    - `use "mmmm.lua"` のかわりに `require "mmmm"` を使用してください
    - これに伴い、Luaの標準検索パス: `package.path` に `nyagos.d/catalog` を追加しました
- `nyagos.d/catalog/ezoe.lua` を削除
- エイリアス: `chompf`, `wordpad`, `abspath` を削除
- ツール Lua 関数: `set("ENV=VALUE")`, `alias("NAME=DEFINE")` を定義する nyagos.d/aliasandset.lua を nyagos.d/catalog/ 以下へ移動
- Lua関数: `addpath` を削除

### 新機能

- `%NO_COLOR%` が定義されていたら、プロンプト・コマンドライン・ls の着色を無効化
- サブコマンド補完(`require "subcomplete"`)関連
    - [#436] curl のオプション補完をサポート ( Thanks to [@tsuyoshicho] )

### 不具合修正

- SKK関連
    - `UTta`,`UTTa` が`打った`ではなく`打っtあ`,`▽う*t*t` になってしまう不具合を修正
    - 手入力した逆三角形が変換マーカーと認識される問題を修正した
    - 次のローマ字かな変換を追加
        - `z,`→`‥`, `z-`→`～`, `z.`→`…`, `z/`→`・`, `z[`→`『`, `z]`→`』`,
            `z1`→`○`, `z2`→`▽`, `z3`→`△`, `z4`→`□`, `z5`→`◇`,
            `z6`→`☆`, `z7`→`◎`, `z8`→`〔`, `z9`→`〕`, `z0`→`∞`,
            `z^`→`※`, `z\\`→`￥`, `z@`→`〃`, `z;`→`゛`, `z:`→`゜`,
            `z!`→`●`, `z"`→`▼`, `z#`→`▲`, `z$`→`■ `, `z%`→`◆`,
            `z&`→`★`, `z'`→`♪`, `z(`→`【`, `z)`→`】`, `z=`→`≒`,
            `z~`→`≠`, `z|`→`〒`, ``z` ``→`“`, `z+`→`±`, `z*`→`×`,
            `z<`→`≦`, `z>`→`≧`, `z?`→`÷`, `z_`→`―`,
        - `bya`→`びゃ` or `ビャ` ... `byo`→`びょ` or `ビョ`
        - `pya`→`ぴゃ` or `ピャ` ... `pyo`→`ぴょ` or `ピョ`
        - `tha`→`てぁ` or `テァ` ... `tho`→`てょ` or `テョ`
    - 変換中の q で、入力済みの平仮名・片仮名を相互変換する機能を実装
- サブコマンド補完(`require "subcomplete"`)関連
    - scoop.cmd であるべき scoop の実行ファイル名が scoop.exe になっていてサブコマンド名補完できない問題を修正した
- Linux でプロンプトが表示されない問題を修正  
  (デフォルトの .nyagos で nyagos.env.prompt を設定していたが、Linux では環境変数の英大文字・小文字を区別するので、nyagos.env.PROMPT でなければいけなかった)
- `foreach` で繰り返すコマンド入力のループ を Ctrl-C で中断できない不具合を修正  
  (本問題は [4.4.6\_2] のコミット[8bf0a2acb25b152d3d40c188de09858c2ef572ae] で混入した模様。[#383]と[4.4.6\_0] も参照のこと)

[8bf0a2acb25b152d3d40c188de09858c2ef572ae]: https://github.com/nyaosorg/nyagos/commit/8bf0a2acb25b152d3d40c188de09858c2ef572ae
[4.4.6\_2]: https://github.com/nyaosorg/nyagos/releases/tag/4.4.6_2
[4.4.6\_0]: https://github.com/nyaosorg/nyagos/releases/tag/4.4.6_0
[#383]: https://github.com/nyaosorg/nyagos/issues/383

### ドキュメント

- docs/07-LuaFunctions_\*.md: `nyagos.shellexecute` について記述 (忘れてた!)

[#436]: https://github.com/nyaosorg/nyagos/pull/436
[@tsuyoshicho]: https://github.com/tsuyoshicho

4.4.14\_0
---------
Oct 6, 2023

- 2023年下期版
- 当バージョンのバイナリは Go 1.20.9 でビルド
- サポート対象は Windows 7, 8.1, 10, 11, WindowsServer 2008以降, Linux となります。

### 新機能

- nyagos.d/suffix.lua: 環境変数 NYAGOSEXPANDWILDCARD にリストされているコマンドのパラメータはワイルドカードを自動展開するようにした （例：`nyagos.env.NYAGOSEXPANDWILDCARD="gorename;gofmt"` ）
- [#432] 新オプション `glob_slash` を追加（起動オプション：`--glob-slash`, lua関数： `nyagos.option.glob_slash=true`、コマンドライン：`set -o glob_slash`）。設定されている時、ワイルドカード展開で `/` を使う
- [SKK] \(Simple Kana Kanji conversion program\) サポート - [設定方法][SKKSetUpJa]
- 適切なUTF8文字列でない時は ANSI文字列とみなして UTF8変換を試みる関数 `nyagos.atou_if_needed` を追加

### 不具合修正

- [#432] `set -o glob` 時、二重引用符内の`*`,`?` がワイルドカードとして展開されていた(本来されるべきではない)
- Linux版で逆クォートがエラーになって機能しない不具合を修正 (Lua関数 atou が常に "not supopported" を返していたので、引数と同じ値を戻すようにした)
- [#433] 文字化けを避けるために、逆クォートでは `nyagos.atou_if_needed` を使って、UTF8 を更に UTF8 化させないようにした
- [v4.4.13\_3] で、`more`, `nyagos.getkey`, `nyagos.getviewwidth` が Windows 7, 8.1 や WindowsServer 2008 で動かない問題を修正
- [#434] Lua で `nyagos.which('cp')` が機能しない問題を修正 (Thanks to [@ousttrue])
- 全角空白(U+3000)の背景色が赤くなっていなかった不具合を修正

### 破壊的変更

- `nyagos.default_prompt` や `nyagos.prompt` は直接ターミナルへプロンプトを出力するのではなく、プロンプト文字列を戻り値として返すようにした。([go-readline-ny.Editor] の deprecated フィールドの Prompt ではなく、PromptWriter を使用するための修正)

[SKK]: https://ja.wikipedia.org/wiki/SKK
[SKKSetUpJa]: https://github.com/nyaosorg/nyagos/blob/master/docs/10-SetupSKK_ja.md
[#432]: https://github.com/nyaosorg/nyagos/issues/432
[#433]: https://github.com/nyaosorg/nyagos/discussions/433
[#434]: https://github.com/nyaosorg/nyagos/pull/434
[@ousttrue]: https://github.com/ousttrue
[v4.4.13\_3]: https://github.com/nyaosorg/nyagos/releases/tag/4.4.13_3
[go-readline-ny.Editor]: https://pkg.go.dev/github.com/nyaosorg/go-readline-ny#Editor

4.4.13\_3
---------
Apr 30, 2023

- (#431) バッチファイル実行で変更した環境変数や、more/typeの出力など非UTF8からUTF8へ変更する時、4096バイトを越えるような行の変換で失敗する不具合を修正 (Thx. @8exBCYJi5ATL)
- (#431 とは別件で) 行のサイズが大きすぎると、more が行を出力しないことがある不具合を修正

4.4.13\_2
---------
Apr 25, 2023

2023年上期版

- (#428) `rmdir /s` でシンボリックリンクやジャンクションの削除に失敗する問題を修正
- \" で引用領域の色を反転させないようにした
- (#429) カレントディレクトリが `C:` の時、`cd c:` が失敗する不具合を修正
- ANSIとUTF8間の変換に go-windows-mbcs v0.4 とgolang.org/x/text/transform を使うようにした

4.4.13\_1
----------
Oct 15, 2022

- (#425) nyagos.d/suffix.lua で設定した拡張子を環境変数 PATHEXT ではなく、NYAGOSPATHEXT に登録するようにした。コマンド名補完は PATHEXT に加え、NYAGOSPATHEXT も参照するようにした  (Thx. @tsuyoshicho)
- (#426) nyagos.argsfilter でコマンドラインが変換された時、空パラメータ("") が消えてしまう問題を修正 (Thx. @juggler999)
- (#426) 外部コマンドに対するワイルドカード展開が有効になっている時、空パラメータ("")が消えてしまう問題を修正 (Thx. @juggler999)
- (#427) '""' が BEEP音プラス '(数字)' に置換されてしまう不具合を修正 (Thx.@hogewest)

4.4.13\_0
---------
Sep 24, 2022

- セキュリティー警告対応のため、gopkg.in/yaml.v3 に依存するモジュールを直接利用しないようにした( https://github.com/nyaosorg/nyagos/security/dependabot/1 )
- (#420) macOS でのビルドをサポート (Thanks to @zztkm)
- (#421) バッチファイルが環境変数を削除した時、それが反映されていない問題を修正(Thanks to @tsuyoshicho)
- (#422) プロンプトに$hを使った時、編集開始位置が右にずれる不具合を修正(Thanks to @Matsuyanagi)
- 端末の背景色が白の時、入力文字が全く見えない不具合を修正(端末デフォルト文字色を使用)
- コマンド名補完において、タイムアウトが効きづらい問題を改善
- サブパッケージを internal フォルダー以下へ移動
- キーの英大文字・小文字を区別しない辞書に Generics を使用するようにした。
- Windows以外で Makefile がエラーになる問題を修正
- (#424) fuzzy finder 拡張機能の統合 (Thx @tsuyoshicho)

4.4.12\_0
---------
Apr 29, 2022

- カラー化コマンドラインの修正
    - オプション文字列の色を黄土色へ変更し、範囲を修正
    - 全角空白の背景色を赤へ変更
    - -0...-9 にオプション用のカラーがついていなかった
    - 端末の背景透過が効くように、背景を黒(ESC[40m)からデフォルト色(ESC[49m)にした(go-readline-ny v0.7.0)
    - WezTerm ではサロゲートペア表示を有効にした
- Linux版のみの不具合修正
    - set -o noclobber 設定時のリダイレクト出力がゼロバイトになってしまう不具合を修正
    - ヒストリがファイルに保存されず、次回起動の際に復元されない不具合を修正
- start コマンドのパラメータは %PATH% 上の任意のファイル・ディレクトリに補完する
- cmd.exe と同様に `rmdir /s` は readonly のフォルダーも削除できるようにした
- `rmdir FOLDER /s` とオプションをフォルダーの後に書けるようにした
- (#418) コマンドラインが ^ で終わっていた場合、Enter入力後も行入力を継続するようにした

4.4.11\_0
---------
Dec 10, 2021

- コマンドラインをカラー化した

4.4.10\_3
---------
Aug 30, 2021

- (#412) Windows10の、WindowsTerminalでない端末で、罫線キャラクターの幅が不正確になっていた問題に対応
- パッケージに新しいアイコンファイルを添付

4.4.10\_2
---------
Jul 23, 2021

- コードページ437 で、%DATE% の置換結果が CMD.EXE 非互換になっている不具合を修正
- go-readline-ny v0.4.13: Windows Terminal で Mathematical Bold Capital (U+1D400 - U+1D7FF) の編集をサポート
- -oオプションがついていも、リダイレクトされた ls の出力から ESC[0m を除くようにした
- (#411) 英語部と日本語部が入れ変わっていたドキュメントを修正 (Thx! @tomato3713)
- テストコードの整理、自動化

4.4.10\_1
---------
Jul 2, 2021

- `./dll` というフォルダが存在して、CDPATH上に DLL がある時、入力されたパス `dll` が補完で DLL に置き変わっていた動作を修正(英大文字・小文字が変わってしまうのが問題)
- 空白を含むディレクトリでの clone コマンドで、カレントディレクトリが維持されない不具合を修正

4.4.10\_0
--------
Jun 25, 2021

- nyagos.d/aliases.lua: abspath をワイルドカード対応にした
- Luaで、`OLEOBJECT:_release()` ではなく、 `OLEOBJECT._release()` が使われた時に適切なエラーが起こされるようにした(glua-oleパッケージを更新)
- CMD.EXE と同様、echo で二重引用符を出力させるようにした
- `"..\"..\".."` という文字列が字句解析の結果、`"..\"..".."` になってしまう不具合を修正
- (#410) 端末の閉じるボタンが押された時にただちに終了させるため SIGTERM を無視しないようにした (Thx @nocd5)
- WindowsTerminal 1.8 以降で、cloneコマンドは同じウインドウの別タブで起動するようにした
- アプリ実行エイリアス経由で wt.exe が呼べず、WindowsTerminal下で clone コマンドが動かなくなっていた問題を修正
- suでネットワークドライブが引き継がれてない不具合を修正
- ビルドするのに PowerShell(make.cmd) ではなく Makefile(GNU Make)を使うようにした
- Linux でもビルドできるようにした

4.4.9\_7
--------
May 22, 2021

- (#409) `set -o glob` や `nyagos.option.glob=true` による外部コマンド向けワイルドカード展開が効かなくなっていた不具合を修正 (Thx @juggler999)

4.4.9\_6
--------
May 7, 2021

- (#406) nyagos.argsfilter で生引数がコンバートされず、suffixコマンドが期待どおり機能しなくなっていた不具合を修正 (Thx @tGqmJHoJKqgK)

4.4.9\_5
--------
May 3, 2021

- go-readline-ny v0.4.10: Yes/Noの回答のYが次のコマンドラインに入力される不具合を修正
- go-readline-ny v0.4.11: Emoji Moifier Sequence (skin tone) をサポート
- Windows 8.1でCPU負荷が高い時のカラーlsの速度を改善した
- ( "io/ioutil" を使わないようにした )
- open で開いたプロセスが閉じる時のメッセージがコマンドラインに重ならないようにした
- go-readline-ny v0.4.12: VisualStudioCodeのターミナルでは絵文字編集はオフにするようにした
- (#403) CMD.EXE のような -S,-C,-K オプションをサポート
- (#403) コマンドラインの不規則な二重引用符が外部コマンドに渡される時に削除される問題を修正
- (#405) fuzzyfinder catalog module を追加 (Thx @tsuyoshicho)

4.4.9\_4
--------
Mar 6, 2021

- (#400) サブコマンド補完向けのコマンドの存在チェック追加(Thx @tsuyoshicho )
- (#401) choco/chocolaty 向けサブコマンド名補完追加(Thx @tsuyoshicho )
- WindowsTerminal で ls や Ctrl-Oの選択時のレイアウトが崩れる問題を修正
- go-readline-ny v0.4.4: 任意の一字 + Ctrl-B + 合字が入力された時、表示が乱れる問題を修正
- go-readline-ny v0.4.5: 合字の異体字サポート
- go-readline-ny v0.4.6: 異体字の後の囲み記号の編集をサポート(&#x0023;&#xFE0F;&#x20E3;)
- (#402) "echo !xxx" でシェルがいきなり終了してしまう問題を修正 (Thx @masamitsu-murase)
- go-readline-ny v0.4.7: REGIONAL INDICATOR (U+1F1E6..U+1F1FF) でカーソル位置が狂わないようにした
- go-readline-ny v0.4.8: WAVING WHITE FLAG and its variations (U+1F3F3 U+FE0F?)
- go-readline-ny v0.4.9: RAINBOW FLAG (U+1F3F3 U+200D U+1F308)

4.4.9\_3
--------
Feb 20, 2021

- WindowsTerminal利用下の一行入力で Unicode の異体字をサポート
- (#397) scoopコマンドのサブコマンド補完を追加 (`use "subcomplete.lua"`) (Thx @tomato3713)
- 補完時に最も短い候補に英大文字/小文字をあわせるようにした
- (#398) io.popen の第二引数のデフォルトが機能していなかった (Thx @ironsand)
- (#399) utf8 offset の改良 (Thx @masamitsu-murase)
- ALT-/ キーバインドのサポート (Thx @masamitsu-murase) https://github.com/zetamatta/go-readline-ny/pull/1
- WindowsTerminal 1.5 で絵文字や丸数字が入力できなくなっていた問題を修正

4.4.9\_2
--------
Jan 8, 2021

- (#342) Ctrl-C 押下時に子プロセスを kill しないようにした

4.4.9\_1
--------
Dec 21, 2020

- パス引数なしの`make install` が失敗する不具合を修正
- (#396) Ctrl-W で左へのスクロールが必要な時に panic する不具合を修正
- コンソール入力の more/clip/type がエコーバックしない時がある不具合を修正
- (#342) クラッシュさせないよう Ctrl-C 割り込みハンドリングを改善

4.4.9\_0
--------
Dec 5, 2020

- (#390,#394) Unicode の合字をサポート
- 異字体コード1-16があるとカーソル位置がおかしくなる不具合を修正  
  (異字体コード自体は未対応なので &lt;FE0F&gt; などと表示する)
- su と clone で WindowsTerminal をサポート
- 編集中はバックグランドプロセスの開始・終了メッセージを出させないようにした
- C-r: インクリメンタルサーチでは英大文字・小文字を区別しないようにした
- || や && の後でコマンド名補完が効かなかった不具合を修正
- C-y: ペースト時に最後の CRLF を除くようにした
- Fix: (#393) ウィンドウアクティブ後の最初のキーが２つ入力されます (Thanks to @tostos5963)
- アンチウィルスが誤判断をするので upx.exe を使用しないようにした

4.4.8\_0
--------
Oct 3, 2020

- git.lua: `git add` 向け補完
    - "\343\201\202\343\201\257\343\201\257"といったファイル名のクォーテーションを解除するようにした
    - untrackなディレクトリの下のファイルも補完対象とした
- diskused: サイズ表記を `ls -h` のように
- diskused が Ctrl-C で止まらなかった不具合を修正
- %ENV:~10,5% のような環境変数抽出を実装
- (#308) UNCパスで表現されていないネットワーク上の GUI 実行ファイルを起動しようとすると `The operation was canceled by the user` というエラーになる問題を修正
- nyagos がネットワーク上にある時、clone コマンドでエラーダイアログが出る問題を修正
- (#389) su: SUBST コマンドのドライブマウントを維持するようにした
- (#390) U+2000～U+2FFF の Unicode が入力できない不具合を修正
- (#390) サロゲートペアな文字が入力できない不具合を修正
- box.lua: Ctrl+O→ESCAPE でユーザが入力した単語が消える不具合を修正
- (#391) subcommand.lua: ghコマンド向けサブコマンド補完を追加 (Thanks to @tsuyoshicho)

4.4.7\_0
--------
Jul 18, 2020

- cd,pushd とその補完で bash のような %CDPATH% をサポートした
- `%APPDATA%\NYAOS_ORG\nyagos.d` のスクリプトも読むようにした
- WindowsTerminal上では、サロゲートペアなUnicodeを&lt;nnnnn&gt;のようにエスケープしないようにした
- バイナリファイルを置くディレクトリを Cmd から bin へ変更した
- catalog/subcomplete.lua
    - 新補完API `nyagos.complete_for` を使うようにした
    - 起動を早くするため、補完するサブコマンド名をファイルにキャッシング
    - キャッシュクリアコマンド `clear_subcomands_cache` を実装
    - `fsutil` と `go` のサブコマンド補完
- catalog/git.lua
    - `subcomplete.lua` を自動でロードするようにした
    - commit-hash も branch-name 同様に補完する
    - `git checkout`で commit-hash,ブランチ名、修正されたファイル名を補完
- (#386) `ls -h` のサイズ出力を単位付きで表示するよう修正 (Thx! [@Matsuyanagi](https://github.com/Matsuyanagi))
- Fix: `nyagos.exec{ ALIAS-COMMAND-USING $@ }` がパニックを引き起す不具合を修正
- 補完可能なファイルのテーブルを返す関数 `nyagos.complete_for_files` を追加

4.4.6\_2
--------
Jun 9, 2020

- Fix: Ctrl-C で Ctrl-D のように終了していた (`4.4.6_0` で #383 修正時に発生)

4.4.6\_1
--------
May 31, 2020

- (#385) 最後にいたフォルダーが削除されたドライブの任意のフォルダーへ移動できなかった不具合を修正
- cd のディレクトリヒストリがパスの大文字小文字を区別していなかった問題を修正
- ドライブ移動(`x:`) でディレクトリヒストリにディレクトリをスタックしていなかった問題を修正
- `nyagos.rawexec{...}`の最後の要素が無視されていた不具合を修正

4.4.6\_0
--------
May 8, 2020

- %DATE% と %TIME% を実装した。
- nyagos.envdel は削除したディレクトリを戻り値として返すようになった。
- `dos/net*.go` などを github.com/zetamatta/go-windows-netresource へ移行
- (#379) nyagos.preexechook & postexechook を追加
- (#383) 端末がクラッシュした時、無限ループに突入してしまう不具合を修正
- `start` の後のタブキーは `which` のようにコマンド名補完をするようにした
- `cd x:\y\z` が失敗した時、`x:\` (ルートディレクトリ)に移動する不具合を修正した

4.4.5\_4
--------
Mar 13, 2020

- github.com/BixData/gluabit32 が消えて C-xC-r C-xC-h , C-xC-g が動かなくなった問題を修正
- (#319) 自前版 bit32.band , bor , bxor を再び追加
- (#378) nyagos.d/catalog/subcomplete.lua: こちらのサブコマンド補完でも拡張子なし・英大文字・小文字は区別しないでコマンドを照合する動作を標準にした
- (#377) scoop でインストールされた git で `git gui` を実行すると、エスケープシーケンスが効かなくなる問題に対応
- パッケージを作成する時だけ、upx で実行ファイルを圧縮し、毎回のビルドでは使わないようにした

4.4.5\_3
--------
Mar 8, 2020

- UNCパスのキャッシュを `~/appdata/local/nyaos.org/computers.txt` ではなく `~/appdata/local/nyaos_org/computers.txt` にセーブするようにした ( 他の機能は `nyaos_org` フォルダーを使っているため )
- サブコマンド補完(`complete_for`)では拡張子は無視してコマンドのマッチングを行うようにした
- UPX.EXE で、実行ファイルを圧縮するようにした
- github.com/BixData/gluabit32 が 404 になって、Lua関数 `bit32.*` が利用できなくなった。
- Windows10 のネイティブANSIエスケープシーケンスも mattn/go-colorable 経由で利用するようにした。
- `echo $(gawk "BEGIN{ print \"\x22\x22\" }")` で二重引用符が出ない不具合を修正

4.4.5\_2
--------
Oct 26, 2019

- (#375) `~randomstring` でクラッシュする不具合を修正
- (#374) 未来のタイムスタンプのファイルの`ls -l`で西暦がでなかった不具合を修正

4.4.5\_1
--------
Oct 20, 2019

- 内蔵boxコマンドが複数アイテム選択に対応していなかった不具合を修正した
- プロセスを開始終了させる時、[PID]表示する際にカーソルを移動させないようにした
- Ctrl-O: 最後の \ の後に引用符を不可しないようにした。(NG: `"Program Files\"` -> OK:`"Program Files\`)
- nyagos.stat/access で ~ や %ENV% を解釈できるようにした

4.4.5\_0
--------
Sep 1, 2019

- Lua関数: `nyagos.dirname()` を実装
- C-o で複数ファイル選択をサポート(Space,BackSpace,Shift-H/J/K/L,Ctrl-Left/Right/Down/Up)
- Alt-Y(引用符つきペースト)で、改行前後に引用符を置くようにした
- C-o で表示される選択肢がディレクトリの時、末尾に \ (Linux では /) をつけるようにした。
- `nyagos.envadd("ENVNAME","DIR")` と `nyagos.envdel("ENVNAME","PATTERN")` を実装
- `nyagos.pathjoin()` で `%ENVNAME%` と `~\`,`~/` を展開するようにした

4.4.4\_3
--------
Jun 14, 2019

- (#371) ファイル名に.を含む実行ファイルを参照できなかった
- diskfree でネットワークドライブに割り当てられた UNC パスを表示

4.4.4\_2
--------
Jun 14, 2019

- バックグラウンドでキャッシュを更新することで `\\host-name` の補完を高速化

4.4.4\_1
--------
May 30, 2019

- Linux 版バイナリがビルドできなかった問題を修正

4.4.4\_0 令和版
--------
May 27, 2019

- (#233) `\\host-name\share-name` を補完できるようになった
- (#238) copyコマンドで進捗表示をするようにした
- `環境変数名=値　コマンド名　パラメータ…` をサポート
- バッチファイル用の一時ファイル名が重複する問題を修正
- (#277) set /a 式を実装
- (#291) バックグラウンド実行のプロセスのIDを表示するようにした
- (#361) GUIアプリの標準出力がリダイレクトできなかった問題を修正
- Linux用の `.` と `source` を実装(/bin/sh を想定)
- 一行入力で、ユーザが待っている時にカーソルの点滅がオフになっていなかった不具合を修正
- `mklink /J マウントポイント 相対パス` で作るジャンクションが壊れていた(絶対パス化が抜けていた)
- 起動オプション `--chdir "DIR"` and `--netuse "X:=\\host-name\share-name"` を追加
- `su`を実行する際にCMD.EXEを使わないようにした(アイコンをNYAGOSのにするため)
- 100個を越える補完候補がある時、確認するようにした
- ps: nyagos.exe 自身の行に `[self]` と表示するようにした
- (#272) `!(ヒストリ番号)@` をそのコマンドが実行された時のディレクトリに置換するようにした
- (#130) ヒアドキュメントをサポート
- Alt-O でショートカットのパス(例:SHORTCUT.lnk) をリンク先のファイル名に置換するようにした
- (#368) Lua関数 io.close() が未定義だった。
- (#332)(#369) io.open() のモード r+/w+/a+ を実装した。

4.4.3\_0
--------
Apr 27, 2019

- (#116) readline: Ctrl-Z,Ctrl-`_` による操作取り消しを実装
- (#194) コンソールウインドウの左上のアイコンを更新するようにした
- CMD.EXE 内蔵 date,time を使うためのエイリアスを追加
- `cd 相対パス` の後のドライブ毎のカレントディレクトリが狂う不具合を修正  
  ( `cd C:\x\y\z ; cd .. ; cd \\localhost\c$ ; c: ; pwd` -> `C:\x` (not `C:\x\y`) )

4.4.2\_3
--------
Apr 13, 2019

- Ctrl-RIGHT,ALT-F(次の単語へ), Ctrl-LEFT,ALT-B(前の単語へ)を実装
- インクリメンタルサーチ開始時にトップへ移動する時のバックスペースの数が間違っていた不具合を修正
- (#364) `ESC[0A` というエスケープシーケンスが使われていた不具合を修正

4.4.2\_1
--------
Apr 5, 2019

- diskfree: 行末の空白を削除
- `~"\Program Files"` の最初の引用符が消えて、Files が引数に含まれない不具合を修正

4.4.2\_0
--------
Apr 2, 2019

- OLEオブジェクトからLuaオブジェクトへの変換が日付型などでパニックを起こす不具合を修正
- Luaの数値が実数として OLE に渡されるべきだったのに、整数として渡されていた。
- Lua: 関数: `nyagos.to_ole_integer(n)` (数値を OLE 向けの整数に変換)を追加(trash.lua用)
- Lua: OLEObject に列挙用オブジェクトを得るメソッド `_iter()` を追加
- Lua: OLEObject を開放するメソッド `OLEObject:_release()` を追加
- trash.lua が COM の解放漏れを起こしていた問題を修正
- Lua: `create_object`生成された IUnkown インスタンスが解放されていなかった不具合を修正
- 「~ユーザ名」の展開を実装
- バッチファイル以外の実行ファイルの exit status が表示されなくなっていた不具合を修正
- %COMSPEC% が未定義の時に CMD.EXE を用いるエイリアス(ren,mklink,dir,...)が動かなくなっていた不具合を修正
- 全角空白(%U+3000%)がパラメータの区切り文字と認識されていた点を修正
- (#359) -c,-k オプションで CMD.EXE のように複数の引数をとれるようにした
- 「存在しないディレクトリ\何か」を補完しようとすると「The system cannot find the path specified.」と表示される不具合を修正 (Thx! [tsuyoshicho](https://github.com/tsuyoshicho))
- (#360) 幅ゼロやサロゲートペアな Unicode は`<NNNNN>` と表示するようにした (Thx! [tsuyoshicho](https://github.com/tsuyoshicho))
- サロゲートペアな Unicode をそのまま出力するオプション --output-surrogate-pair を追加
- suコマンドで、ネットワークドライブが失なわれないようにした
- (#197) ソースがディレクトリで -s がない時、`ln` はジャンクションを作成するようにした
- 内蔵の mklink コマンドを実装し、`CMD.exe /c mklink` のエイリアス `mklink` を削除
- ゼロバイトの Lua ファイルを削除(cdlnk.lua, open.lua, su.lua, swapstdfunc.lua )
- (#262) `diskfree` でボリュームラベルとファイルシステムを表示するようにした
- UNCパスがカレントディレクトリでもバッチファイルを実行できるようにした。
- UNCパスがカレントディレクトリの時、ren,assoc,dir,for が動作しない不具合を修正
- (#363) nyagos.alias.COMMAND="string" 中では逆クォート置換が機能しない問題を修正 (Thx! [tostos5963](https://github.com/tostos5963) & [sambatriste](https://github.com/sambatriste) )
- (#259) アプリケーションをダイアログで選んでファイルを開くコマンド `select` を実装
- `diskfree` の出力フォーマットを修正

4.4.1\_1
--------
Feb 15, 2019

- `print(nyagos.complete_for["COMMAND"])`が機能するようにした
- (#356) `type` が LF を含まない最終行を表示しない不具合を修正 (Thx! @spiegel-im-spiegel)
    - 要 [zetamatta/go-texts](https://github.com/zetamatta/go-texts) v1.0.1～
- ビルドに `Go Modules` を使うようにした
- `killall`,`taskkill` コマンド向け補完
- `kill` & `killall`: 自分自身のプロセスを停止できなくした。
- (#261) 補完や1フォルダのlsは10秒でタイムアウトするようにした
- Lua で OLE オブジェクトのセッター(`__newindex`)が効かなかった不具合を修正
- (#357) 仏語キーボードで AltGrシフトが効かない問題を修正 (Thx! @crile)
- (#358) `foo.exe`と`foo.cmd`があった時、`foo`で`foo.exe`ではなく`foo.cmd` が呼び出される不具合を修正

4.4.1\_0
--------
Feb 2, 2019

- `which`,`set`,`cd`,`pushd`,`rmdir`,`env` コマンド向け補完 (Thx! [ChiyosukeF](https://twitter.com/ChiyosukeF))
- (#353) OpenSSHでパスワード入力中に Ctrl-C で中断すると、画面表示がおかしくなる問題を修正 (コマンド実行後にコンソールモードを復旧するようにした) (Thx! [beepcap](https://twitter.com/beepcap))
- (#350) `-l` なしの `ls -F` で os.Readlink を呼ぶのをやめた
- `nyagos.complete_for["COMMANDNAME"] = function(args) ... end` 形式の補完
- (#345) subcomplete.lua で git/svn/hg が効かない問題を修正(Thx! @tsuyoshicho)
- リダイレクトが含まれている時、Lua関数 io.popen が機能しない不具合を修正(Thx! @tsuyoshicho)
- (#354) box.lua のヒストリ補完が C-X h で起動していなかった不具合を修正 (Thx! @fushihara)
- nyagos.d/catalog/subcomplete.lua で `hub` コマンドの補完をサポート (Thx! @tsuyoshicho)

4.4.0\_1
--------
Jan 19, 2019

- "--go-colorable" と "--enable-virtual-terminal-processing" を廃止
- `killall` コマンドを実装
- Linux用の copy と move を実装
- (#351) `END` と `F11` キーが動作もキー割り当てもできなかった不具合を修正

4.4.0\_0
-----------
Jan 12, 2019

- バッチファイルを呼ぶ時に、`/V:ON` を CMD.EXE に使わないようにした

4.4.0\_beta
-----------
Jan 2, 2019

- Linux サポート(実験レベル)
- ドライブ毎のカレントディレクトリが子プロセスに継承されなかった問題を修正
- ライブラリ "zetamatta/go-getch" のかわりに "mattn/go-tty" を使うようにした
- msvcrt.dll を直接syscall経由で使わないようにした。
- Linux でも NUL を /dev/null 相当へ
- Lua変数 nyagos.goos を追加
- (#341) Windows10で全角文字の前に文字を挿入すると、不要な空白が入る不具合を修正
    - それに伴い、Windows10 では virtual terminal processing を常に有効に
    - `git.exe push`が無効にしても再び有効にする
- (#339) ワイルドカード `.??*` が `..` にマッチする問題を修正
    - 要 github.com/zetamatta/go-findfile tagged 20181230-2
