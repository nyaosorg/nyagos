[top](../readme_ja.md) &gt; [English](./history-4.3_en.md) / Japanese

NYAGOS 4.3.3\_5
===============
(2018.12.24)

* (#345) subcomplete.lua が git 補完で動作しない問題を修正 (Thx! @tsuyoshicho)
* (#347) `dir 2>&1`実行後、dup元の標準出力までクローズされていた不具合を修正(Thx! @Matsuyanagi)
* (#348) ls 後マウスのスクロールが効きにくくなる問題に対応 (Thx! @tyochiai)
    * 要 github.com/zetamatta/go-getch tagged:20181223

NYAGOS 4.3.3\_4
===============
(2018.12.13)

* 出力先が端末でない場合、more を type と等価に
* バッチ実行時に作成する踏み台の一時バッチを廃止。`CMD /V:ON /S /C "..."` を使うようにした
* (#340) 最大ヒストリ保存数を指定する `nyagos.histsize` を追加(Thx! @crile)
* (#343) %COMSPEC% が未定義の時、CMD.EXE を用いるようにした(Thx! @orz--)

NYAGOS 4.3.3\_3
===============
(2018.10.23)

* (#310) copy と move の宛先でショートカットをサポート
* (#313 reopened) `git blame FILES | type | gvim - &` で gvim が空バッファで始まってしまう問題を修正
* 壊れたジャンクションに対する rmdir ができなかった問題を修正
* Luaスクリプトや外部プロセスの一部で Ctrl-C が機能しなかった問題を修正
* (#267) `type` や `more` で UTF16 ファイルを表示できるようにした
* (#336) `io.write` が -e や --lua-file オプション中で機能しない不具合を修正
* (#337) バッチが exit -1 で終了するとクラッシュする不具合を修正(Thx! @hogewest)

NYAGOS 4.3.3\_2
===============
(2018.09.22)

* リダイレクトで存在するファイルを上書きする時のエラーメッセージにファイル名を付与した
* noclobber が設定されている時に nul へのリダイレクトを上書きエラーにしてしまう問題を修正
* diskused: エラーが見付かっても容量計算を続けるようにした
* ls: `-1` があると、`-l` オプションが動かない点を修正
* ls: 出力先が端末でない時、1ファイル1行で出力していなかった点を修正
* 別名定義されていない内蔵コマンドを bash のように `\ls` と呼べるようになった
* for のエイリアス定義が壊れていたのを修正
* ファイル名以外の補完の時もパスの区切り文字が補正されてしまう問題を修正

NYAGOS 4.3.3\_1
===============
(2018.08.29)

* #330,#331 オリジナル版のfile:readの非互換な動作を修正 (Thx! @erw7)
* #332 io.open("w") でバッファリングしないようにした (Thx! @spiegel-im-spiegel)
* #333 Fix file:seek() が読み取り時に期待どおり同しなかった点を修正 (Thx! @erw7)
* #333 Fix file:close() の戻り値がおかしかった点を修正 (Thx! @erw7)
* #319 utf8.len() を実装
* Fix: `which` が拡張子なしのファイルも出力していた点を修正
* `pwd` はデフォルトでは論理パスを出力するようにした
* インクリメンタルサーチを開始した時、表示にゴミが残る不具合を修正
* -lfdflags="-s -w" で実行ファイルのサイズを削減した

NYAGOS 4.3.3\_0
===============
(2018.08.14)

* #283 Ctrl-O での補完で、パスでディレクトリを省略するようにした。
* #326 オプション `nyagos.option.tilde_expansion` を追加
* Fix: `nyagos.option.xxxxxx = true` が機能していなかった
* Fix #328 `start https://...` で URL をブラウザで開けなかった
* #327 のために --read-stdin-as-file を実装(標準入力からファイル扱いでコマンドを読み込む)
* シンボリックリンク先にある GUI アプリケーションの実行が失敗する問題を修正
* (パイプラインではない)リダイレクトがバッググラウンドで起動できなかった不具合を修正
* 文字列を空白で分割する Lua 関数 nyagos.fields を追加
* #185 `ps` , `kill` コマンドを追加
* #329 Lua用数値型として int ではなく float64 を使うようにした

NYAGOS 4.3.2\_0
===============
(2018.07.23)

* #319 github.com/BixData/gluabit32 で、Lua関数 `bit32.*` を全てサポート
* #323 io.lines() , nyagos.lines() がリダイレクトされた標準入力から読み込めない問題を修正
* io.write() がリダイレクトされた標準出力に出力できなかった
* `io.*` を NYAGOS の自前バージョンに置き変えた
* #324 Lua の print で --no-go-colorable が効いていなかった不具合を修正 (Thx @tignear)
* #325 Source 文で空白を含むパスをロードできなかった不具合を修正 (Thx @tignear)
* オプション `--lua-first` and `--cmd-first` を追加

NYAGOS 4.3.1\_3
===============
(2018.06.19)

* #316 %PATH% の中の長さゼロのエントリがカレントディレクトリとみなされていた不具合を修正
* #321 キー機能名の `previous_history` と `next_history` が未登録だった不具合を修正
* -h,--help オプションを追加
* バッチファイル組み込みのため、Luaスクリプトの @ で始まる行を無視するようにした
* #322 バッチファイルの引数のエンコーディングをスレッドのコードページから、コンソールのコードページへ変更した。
* Lua変数 `nyagos.option.*` の全てを nyagos.exe のコマンドラインオプションで設定できるようにした。

NYAGOS 4.3.1\_2
===============
(2018.06.12)

* #320: nyagos.rawexec & raweval が引数内のテーブルを展開していなかった非互換性を修正
* --show-version-only を指定すると --norc を自動的に有効化するようにした

NYAGOS 4.3.1\_1
===============
(2018.06.11)

* lua53.dll 向けのソースコードを削除
* #317: `use subcomplete` が有効で、rclone.exe が見付かった時デッドロックしていた
    - https://github.com/yuin/gopher-lua/issues/181 も参照のこと
* #318,#319 下記の Lua 5.3 互換関数を追加
    - bit32.band/bitor/bxor
    - utf8.char/charpattern/codes

NYAGOS 4.3.1\_0
===============
(2018.06.03)

* `--no-go-colorable` と `--enable-virtual-terminal-processing` で、Windows10 ネイティブのエスケープシーケンスをサポート
* #304,#312, カレントディレクトリから実行ファイルを探す時のオプションを追加
    * --look-curdir-first: %PATH% より前に探す(デフォルト:CMD.EXE互換動作)
    * --look-curdir-last : %PATH% より後に探す(PowerShell互換動作)
    * --look-curdir-never: %PATH% だけから実行ファイルを探す(UNIX Shells互換動作)
* nyagos.prompt にプロンプトテンプレートの文字列を直接代入できるようになった。
* #314 rmdir がジャンクションを削除できなかった問題を修正

NYAGOS 4.3.0\_4
===============
(2018.05.12)

- Fix: #309 nyagos.getkey() が使えない不具合を修正 (Thx @nocd5)
- lnk コマンドの宛先が `*.lnk` でなかったり存在しなかった時のエラーメッセージを修正
- 子プロセスのカーソルがオフになってしまう不具合を修正

NYAGOS 4.3.0\_3
===============
(2018.05.09)

- nyagos.setalias, nyagos.getalias の実装が漏れており、`alias { CMD=XXX}` が動かなくなっていた
- エイリアスの戻り値でテーブルが与えられた時、コマンド名として解釈すべき、要素[0]が使われていなかった不具合を修正
- `doc/09-Build_*.md`: github からのソースダウンロード方法についてドキュメント更新

NYAGOS 4.3.0\_2
===============
(2018.05.07)

- #305: ユーザの .nyagos が二回目以降ロードされない不具合を修正(Thx! @erw7)

NYAGOS 4.3.0\_1
===============
(2018.05.05)

- nyagos.d/start.lua が動作していなかった不具合を修正 (エイリアス関数の rawargs パラメータが実装されていなかった)
- alias 関数の戻り値が評価されていなかった不具合を修正
- -e オプションのスクリプト向けに、arg[] に引数が代入されていなかった
- -e,-f オプションで、`getRegInt: could not find shell in Lua instance` が表示される不具合を修正
- バッチファイルが `exit /b` の値を ERRORLEVEL として返せなかった不具合を修正

NYAGOS 4.3.0\_0
===============
(2018.05.03)

- シンボリックリンクの先を参照するオプション `ls -L` を追加

NYAGOS 4.3\_beta2
=================
(2018.05.01)

- C-o を押すと Enter か Escape が押されるまでハングしたように見える不具合を修正
    - (ライブラリを修正: [go-box](https://github.com/zetamatta/go-box/commit/322b2318471f1ad3ce99a3531118b7095cdf3842))
- chcp が動作しない不具合を修正 (同コマンドは画面幅取得のため別名定義していた)

NYAGOS 4.3\_beta
=================
(2018.04.30)

- **lua53.dll のかわりに Gopher-Lua を採用** #300
    - 旧来の lua53.dll 版 nyagos.exe は `cd mains ; go build` でビルド可能
    - Lua無し版 nyagos.exe を `cd ngs ; go build` でビルド可能
- `nyagos.option.cleanup_buffer` を追加(デフォルトは false)。true の場合、一行入力の前にコンソールバッファをクリアする
- `set -o OPTION_NAME` と `set +o OPTION_NAME` を新設(`nyagos.option.OPTION_NAME=` on Lua と等価)
- コンソール出力をバッファリングするようにした ( go-colorable and bufio.Writer )
