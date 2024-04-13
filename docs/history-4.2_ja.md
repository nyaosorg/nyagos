[top](../README_ja.md) &gt; [English](./history-4.2_en.md) / Japanese

NYAGOS 4.2.5\_1
===============
(2018.04.14)

- ブロックif で `if [not] errorlevel N` が動かなかった不具合を修正
- リパースポイント先の実行ファイルが見付からなくなっている問題を修正
- `ls -1F` が `/`,`*` や `@` といったインジケーターを出力しない問題を修正
- `ls -F` が「リパースポイントではあるが、ジャンクション、シンボリックリンクでないファイル・ディレクトリ」に @ マークをつけていた不具合を修正
- `_nyagos` で `history` コマンドを使った時のエラーメッセージを変更

NYAGOS 4.2.5\_0
===============
(2018.03.31)

- luaフラグ nyagos.option.usesource を追加。false の時、バッチファイルは NYAGOS の環境変数を変更できなくなる(default:true)

NYAGOS 4.2.5\_beta2
===================
(2018.03.27)

- #296 ユーザ名にマルチバイト文字が入っていると、バッチが正常動作しない不具合を修正
    - 一時バッチファイルのエンコーディングが UTF8 になっていた
    - 一時バッチファイルの改行コードが CRLF ではなく LF になっていた
- #297 /b なしの exit をバッチファイルが実行した時の、一時ファイルが無い旨のエラーがでていた

NYAGOS 4.2.5\_beta
=================
(2018.03.26)

- CMD.EXE と同様に、バッチファイルが変更した環境変数の値を読み取るようにした。
- ソースの幾つかをリファクタリングした。

NYAGOS 4.2.4\_0
===============
(2018.03.09)

* lua: ole: `variable = OLE.property` が `OLE:_get('property')` のかわりに使えるようになった
* lua: ole: `OLE.property = value` が `OLE:_set('property',value)` のかわりに使えるようになった
* `nyagos.d/*.ny` のコマンドファイルも読み込むようにした
* #266: `lua_e "nyagos.option.noclobber = true"` でリダイレクトでのファイル上書きを禁止
* #269: `>| FILENAME` もしくは `>! FILENAME` で、`nyagos.option.noclobber = true` の時も上書きできるようにした
* #270: プロンプト表示時にコンソール入力バッファをクリアするようにした
* #228: $ENV[TAB] という補完をネイティブでサポート
* #275: `!str:$` や `!str?:$` といったヒストリ置換が機能しない不具合を修正
* ! で指定されるヒストリが存在しない時「event not found」エラーを出させるようにした
* #285: パイプラインを使っていない GUIプログラムは CMD.EXE 同様終了を待たないようにした (CreateProcess ではなく ShellExecute を使用する)
* (bytes.Buffer を strings.Builder に置き換えた。Go 1.10 が必要になった)
* 複数のファイルが「open」で一度に開こうとした時、`open: ambiguous shellexecute` とエラーを表示するようにした。
* `nyagos.alias.NAME = nil` で alias を削除できなかった動作を修正

NYAGOS 4.2.3\_4
===============
(2018.03.04)

* `ls -h` のかわりに `ls -?` をヘルプに用意した
* make.cmd のかわりに go build でビルドした時、バージョンを `snapshot-GOARCH` と表示するようにした
* `type DIRECTORY` が実行された時にエラーにするようにした。
* `del 存在しないファイル` を実行した時のエラーをシンプルにした.
* #279 環境変数置換(%VAR:OLD=NEW%)で、英大文字/小文字を区別していた不具合を修正
* #281 `cd \\host-name\share-name ; open` で `C:\Windows\system32` 開く不具合を修正
* #286 Fix: 二重引用符内の空白に続く ~ が %USERPROFILE% と解釈されていた不具合を修正
* #287 ヒストリの最後のエントリの時、↓をタイプしても何もしないようにした

NYAGOS 4.2.3\_3
===============
(2018.01.28)

* `print(nil,true,false)` が何も出力しない不具合を修正
* 検索にヒットしないヒストリ置換で `!notfoundstr`  が `!n` になってしまう不具合を修正
* #271: Ctrl-O が環境変数を含んだパスで効かない不具合を修正 (go-findfile)
* 補完の際、パーネントの後にスペースを追加しないようにした
* #276 source コマンドで実行されるバッチの標準出力が閉じていた不具合を修正 (Thx @tyochiai )

NYAGOS 4.2.3\_2
===============
(2018.01.06)

* #265 `ls` + 空白 + TAB でコマンド名補完が動いていた不具合を修正

NYAGOS 4.2.3\_1
===============
(2017.12.30)

* 改行コード等が単語の区切りとして認識していなかった不具合を修正
* #264 画面バッファの幅が広すぎる時に、画面にゴミが現われる不具合を修正 
    (You have to do `go get -u github.com/mattn/go-colorable`)

NYAGOS 4.2.3\_0
===============
(2017.12.25)

* 起動スクリプトのロードを抑制する --norc オプションを追加
* #132 foreach 文とブロック if 文をサポート
* 拡張子が .lua でない場合でも Lua スクリプトとして実行するオプション --lua-file を追加
* `complete_hook(c)` の パラメータ c に項目を追加
    * `c.field` : `c.text` を空白で分割したもの
    * `c.left` : カーソル前の文字列
* `|`, `&`, `;` の直後でも、コマンド名補完が有効になるようにした
* #245 Lua の print がリダイレクトに対応
* インクリメンタルサーチ中に Ctrl-S で逆方向サーチできるようにした
* #248 バックスラッシュ直後の環境変数展開が機能しない不具合を修正
* lua関数 `nyagos.msgbox(MESSAGE,TITLE)` を追加

NYAGOS 4.2.2\_2
===============
(2017.11.26)

* #255 `start` コマンドでコマンドを %PATH% から探すようにした
* #254 `nyagos -f SCRIPT -xxxx` の -xxxx が SCRIPT のオプションではなく、nyagos のオプションとして扱われていた問題を修正
* コマンドラインフィルターが設定されていない時に Lua のスタックがオーバーフローしてクラッシュする不具合を修正

NYAGOS 4.2.2\_1
===============
(2017.10.11)

* #250 引数なしの `bindkey` でクラッシュする不具合を修正 (Thx @masamitsu-murase)
* #252 Shift/Ctrl キーのタイプで、画面のスクロールがキャンセルされてしまう問題を修正 ( Shift/Ctrl キーのタイプでカーソルOFF/ONの出力を省くようにした ) (Thx @masamitsu-murase)
* #253 `nyagos-4.2.2_0-386` が make.cmd の不具合で 64bitでビルドされていた (Thx @hazychill)

NYAGOS 4.2.2\_0
===============
(2017.10.08)

* 新Lua製コマンド(`abspath`,`chompf`,`wildcard`)を追加
* 漏れていたLua製コマンドのリファレンスを追記: `lua_f` , `kill` , `killall`
* #246 クラッシュ回避のため、Lua の userdata を `share[]` に代入したり、Lua インスタンスの fork 時にコピーしないようにした (Thx @masamitsu-murase)
* #247 Go の Garbage Collector が Lua で参照中のデータを開放してクラッシュする問題を修正した (Thx @masamitsu-murase)
* #248 補完用フックで、補完文字列とは別にリストアップ用の表示テキストを指定できるようになった。(Thx @masamitsu-murase)
* #249 `nyagos.completion_slash` を追加。これが true の時、ファイル名補完はデフォルトでパス区切り文字に / を使う(Thx @masamitsu-murase)
* PowerShell で記述したあたらしいビルドスクリプト(make.cmd) を用意

NYAGOS 4.2.1\_0
===============
(2017.08.31)

* #241 `completion_hook` で戻るリストの順番が反映されていなかった問題を修正 (Thx @masamitsu-murase)
* #242,#243 readline のキーに Alt+Backspace と Alt+"/" を追加 (Thx @masamitsu-murase)
* 内蔵コマンドの sudo を削除
* 内蔵コマンド more を追加(カラー & utf8 サポート)
* 一行入力で `C-q`,`C-v` をサポート(`QUOTED_INSERT`)
* 内蔵コマンド pwd に -P(全てのリンクをたどる) ,-L(環境からPWDを得る) を追加
* パニックが発生した時、nyagos.dump を出力するようにした
* `diskused`: du ライクな新コマンド
* `rmdir` : 進捗を表示する仕様を復活させた
* `diskfree`: df ライクな新コマンド

NYAGOS 4.2.0\_5
===============
(2017.08.16)

* Windows7 でのビルドで、バージョン情報が実行ファイルのプロパティーに記入されない問題があり、修正した。原因は goversioninfo 向けの JSON を作るスクリプトが PowerShell 3.0 の ConvertTo-JSON メソッドを必要としていたが、Windows 7 はサポートしていなかった。
* nyagos.box(LIST)関数が LIST の順番を無視していた

NYAGOS 4.2.0\_4
===============
(2017.07.29)

* `.nyagos` にエラーがあった時のエラー行番号が表示されない問題を修正
* 前回実行時とEXEファイルのアーキテクチャ(amd64 or 386)が変わった時、`.nyagos` のキャッシュがエラーになる不具合を修正
* Fix: `ls | more` で `ESC[0K` が表示されていた
* (内部) go-colorable の `ESC[%dC` と `ESC[%dD` の挙動変更に追随 ( https://github.com/mattn/go-colorable/commit/3fa8c76f , 感謝 > @tyochiai )
* デフォルトと `_nyagos` で `suffix "lua=nyagos"` は間違っていた。「`.exe -f`」を追記した。
* `nyagos.d` ディレクトリのスクリプトが、lua.exe など nyagos.exe 以外で実行された場合、エラーにするようにした。
* `nyagos.d` ディレクトリで `suffix` とタイプすると、無限に nyagos.exe プロセスが起動する問題 #237 を修正するために、ユーザがタイプしたコマンド名に拡張子が含まれていない場合は、インタプリタ名の挿入をしないようにした。
* Fix #240: 空のディレクトリで C-o を押下すると「`bad argument # 1 to 'find' (string expected, got nil)`」と表示されていた

NYAGOS 4.2.0\_3
===============
(2017.07.13)

* Fix: `box` Enter & Ctrl-C でパニックが発生する不具合を修正
* Fix: `lua_e "nyagos.box({})"` でパニックが発生する不具合を修正
* Fix: `box` でスクロールの際、カーソルが消える不具合を修正(go-boxライブラリの不具合修正)
* `box` コマンドでのチラツキを軽減した(go-boxライブラリの修正)
* Fix: #235 実行ファイルと同じフォルダーの .nyagos が読み込まれていなかった
* 補完で、! マークがある時、"" で囲むようにした。
* Fix: `suffix ps1` が `?:-1: attempt to concatenate a table value` となる不具合を修正

NYAGOS 4.2.0\_2
===============
(2017.06.13)

* `lnk . ~`が失敗する不具合を修正
* ネットワークフォルダーにシンボリックリンクされていて、UAC昇格が必要な実行ファイルを呼び出せない問題を修正 (ShellExecute に物理パスを渡すようにした)
* 一行入力のインクリメンタルサーチ中、BACKSPACE で行が更新されなかった不具合を修正
* Lua で空文字のグローバル変数があると、pipeを使った時に落ちる不具合を修正(#232)

NYAGOS 4.2.0\_1
========================
(2017.06.06)

* Fix: `_nyagos` のサンプルをパッケージに同梱するのを忘れていた (#230)
* `chmod` を実装した (#199)
* nyagos.d/catalog/dollar.lua: $TEMP\xxxx 形式のファイル名補完をサポート(#228)
* nyagos.d/catalog/ezoe.lua: 復活

NYAGOS 4.2.0\_0
===============
(2017.05.29)

* `share[]`直下の以外のLua変数が共有されない制限を改善した(#210,#208)
    * Luaの新インスタンス作成をバックグラウンドスレッド開始時に限った(さもなければインスタンス共有する)
    * Luaの新インスタンス作成時に`share[]`以外のグローバル変数もメインスレッドのインスタンスからフルコピーするようにした。
    * ~/.nyagos をロードした Lua インスタンスでプロンプトを表示するようにした

新機能
------
* `nyagos.completion_hidden`: true の時、隠しファイルも補完候補に入れる
* 内蔵コマンド `env` の追加
* ヒストリ参照テーブル `nyagos.history[..]` , `#nyagos.history` 用意
* 内蔵コマンドとして `type` を追加
* UTF8 / MBCS 両方を読み込める内蔵コマンド `clip` を実装 (#202)
* READONLY属性のファイルも消す `del /f`オプション追加 (#198)
* `attrib` コマンドを内蔵コマンドとして実装 (#199)
* ショートカット作成コマンド lnk を内蔵(`lnk FILENAME SHORTCUT WORKING-DIRECTORY`)
* `$( )` 形式のコマンド出力引用形式をサポート
* `ls -l`: ショートカットのリンク先・作業ディレクトリを表示するようにした
* Lua関数 `nyagos.box()` を追加

Trivial change
--------------
* 起動オプションに「-b (base64化されたコマンド文字列)」を追加した
* Lua(nyagos.d/cdlnk.lua)製の `cd/push ショートカット.lnk` を Go で書き直した
* `nyagos.alias.grep = "findstr.exe"`

Bugfix
------
* `%USERPROFILE:\=/%` で `\` が一度しか置換されていない不具合を修正
* デフォルトの`_nyagos`で`ll`がカラーでない`ls`に別名定義されていた点を修正
* C-o が空白と ~ を含むファイル名を補完できなかった不具合を修正
* Ctrl-O のファイル名選択がパニックを起こす不具合を修正(#204)
* FINDコマンドなどのために、ユーザが明示した二重引用符は決して削除しないようにした(#218,#222)
* UAC昇格が必要なコマンドを呼ぶとエラーになる問題を修正し、UAC昇格ダイアログを出すようにした (#227)
* `FOO.123.EXE` が `FOO` とタイプした時でも実行されてしまう不具合を修正 #229
