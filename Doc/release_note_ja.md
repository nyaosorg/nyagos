[English](release_note_en.md) / Japanese

* git.lua: `git add` 向け補完で"\343\201\202\343\201\257\343\201\257"といったファイル名のクォーテーションを解除するようにした
* diskused: サイズ表記を `ls -h` のように
* diskused が Ctrl-C で止まらなかった不具合を修正

NYAGOS 4.4.7\_0
===============
(2020.07.18)

* cd,pushd とその補完で bash のような %CDPATH% をサポートした
* `%APPDATA%\NYAOS_ORG\nyagos.d` のスクリプトも読むようにした
* WindowsTerminal上では、サロゲートペアなUnicodeを&lt;nnnnn&gt;のようにエスケープしないようにした
* バイナリファイルを置くディレクトリを Cmd から bin へ変更した
* catalog/subcomplete.lua
    - 新補完API `nyagos.complete_for` を使うようにした
    - 起動を早くするため、補完するサブコマンド名をファイルにキャッシング
    - キャッシュクリアコマンド `clear_subcomands_cache` を実装
    - `fsutil` と `go` のサブコマンド補完
* catalog/git.lua
    - `subcomplete.lua` を自動でロードするようにした
    - commit-hash も branch-name 同様に補完する
    - `git checkout`で commit-hash,ブランチ名、修正されたファイル名を補完
* (#386) `ls -h` のサイズ出力を単位付きで表示するよう修正 (Thx! [@Matsuyanagi](https://github.com/Matsuyanagi))
* Fix: `nyagos.exec{ ALIAS-COMMAND-USING $@ }` がパニックを引き起す不具合を修正
* 補完可能なファイルのテーブルを返す関数 `nyagos.complete_for_files` を追加

NYAGOS 4.4.6\_2
===============
(2020.06.09)

* Fix: Ctrl-C で Ctrl-D のように終了していた (`4.4.6_0` で #383 修正時に発生)

NYAGOS 4.4.6\_1
===============
(2020.05.31)

* (#385) 最後にいたフォルダーが削除されたドライブの任意のフォルダーへ移動できなかった不具合を修正
* cd のディレクトリヒストリがパスの大文字小文字を区別していなかった問題を修正
* ドライブ移動(`x:`) でディレクトリヒストリにディレクトリをスタックしていなかった問題を修正
* `nyagos.rawexec{...}`の最後の要素が無視されていた不具合を修正

NYAGOS 4.4.6\_0
===============
(2020.05.08)

* %DATE% と %TIME% を実装した。
* nyagos.envdel は削除したディレクトリを戻り値として返すようになった。
* `dos/net*.go` などを github.com/zetamatta/go-windows-netresource へ移行
* (#379) nyagos.preexechook & postexechook を追加
* (#383) 端末がクラッシュした時、無限ループに突入してしまう不具合を修正
* `start` の後のタブキーは `which` のようにコマンド名補完をするようにした
* `cd x:\y\z` が失敗した時、`x:\` (ルートディレクトリ)に移動する不具合を修正した

NYAGOS 4.4.5\_4
===============
(2020.03.13)

* github.com/BixData/gluabit32 が消えて C-xC-r C-xC-h , C-xC-g が動かなくなった問題を修正
* (#319) 自前版 bit32.band , bor , bxor を再び追加
* (#378) nyagos.d/catalog/subcomplete.lua: こちらのサブコマンド補完でも拡張子なし・英大文字・小文字は区別しないでコマンドを照合する動作を標準にした
* (#377) scoop でインストールされた git で `git gui` を実行すると、エスケープシーケンスが効かなくなる問題に対応
* パッケージを作成する時だけ、upx で実行ファイルを圧縮し、毎回のビルドでは使わないようにした

NYAGOS 4.4.5\_3
===============
(2020.03.08)

* UNCパスのキャッシュを `~/appdata/local/nyaos.org/computers.txt` ではなく `~/appdata/local/nyaos_org/computers.txt` にセーブするようにした ( 他の機能は `nyaos_org` フォルダーを使っているため )
* サブコマンド補完(`complete_for`)では拡張子は無視してコマンドのマッチングを行うようにした
* UPX.EXE で、実行ファイルを圧縮するようにした
* github.com/BixData/gluabit32 が 404 になって、Lua関数 `bit32.*` が利用できなくなった。
* Windows10 のネイティブANSIエスケープシーケンスも mattn/go-colorable 経由で利用するようにした。
* `echo $(gawk "BEGIN{ print \"\x22\x22\" }")` で二重引用符が出ない不具合を修正

NYAGOS 4.4.5\_2
===============
(2019.10.26)

* (#375) `~randomstring` でクラッシュする不具合を修正
* (#374) 未来のタイムスタンプのファイルの`ls -l`で西暦がでなかった不具合を修正

NYAGOS 4.4.5\_1
===============
(2019.10.20)

* 内蔵boxコマンドが複数アイテム選択に対応していなかった不具合を修正した
* プロセスを開始終了させる時、[PID]表示する際にカーソルを移動させないようにした
* Ctrl-O: 最後の \ の後に引用符を不可しないようにした。(NG: `"Program Files\"` -> OK:`"Program Files\`)
* nyagos.stat/access で ~ や %ENV% を解釈できるようにした

NYAGOS 4.4.5\_0
===============
(2019.09.01)

* Lua関数: `nyagos.dirname()` を実装
* C-o で複数ファイル選択をサポート(Space,BackSpace,Shift-H/J/K/L,Ctrl-Left/Right/Down/Up)
* Alt-Y(引用符つきペースト)で、改行前後に引用符を置くようにした
* C-o で表示される選択肢がディレクトリの時、末尾に \ (Linux では /) をつけるようにした。
* `nyagos.envadd("ENVNAME","DIR")` と `nyagos.envdel("ENVNAME","PATTERN")` を実装
* `nyagos.pathjoin()` で `%ENVNAME%` と `~\`,`~/` を展開するようにした

NYAGOS 4.4.4\_3
===============
(2019.06.14)

* (#371) ファイル名に.を含む実行ファイルを参照できなかった
* diskfree でネットワークドライブに割り当てられた UNC パスを表示

NYAGOS 4.4.4\_2
===============
(2019.06.14)

* バックグラウンドでキャッシュを更新することで `\\host-name` の補完を高速化

NYAGOS 4.4.4\_1
===============
(2019.05.30)

* Linux 版バイナリがビルドできなかった問題を修正

NYAGOS 4.4.4\_0 令和版
======================
(2019.05.27)

* (#233) `\\host-name\share-name` を補完できるようになった
* (#238) copyコマンドで進捗表示をするようにした
* `環境変数名=値　コマンド名　パラメータ…` をサポート
* バッチファイル用の一時ファイル名が重複する問題を修正
* (#277) set /a 式を実装
* (#291) バックグラウンド実行のプロセスのIDを表示するようにした
* (#361) GUIアプリの標準出力がリダイレクトできなかった問題を修正
* Linux用の `.` と `source` を実装(/bin/sh を想定)
* 一行入力で、ユーザが待っている時にカーソルの点滅がオフになっていなかった不具合を修正
* `mklink /J マウントポイント 相対パス` で作るジャンクションが壊れていた(絶対パス化が抜けていた)
* 起動オプション `--chdir "DIR"` and `--netuse "X:=\\host-name\share-name"` を追加
* `su`を実行する際にCMD.EXEを使わないようにした(アイコンをNYAGOSのにするため)
* 100個を越える補完候補がある時、確認するようにした
* ps: nyagos.exe 自身の行に `[self]` と表示するようにした
* (#272) `!(ヒストリ番号)@` をそのコマンドが実行された時のディレクトリに置換するようにした
* (#130) ヒアドキュメントをサポート
* Alt-O でショートカットのパス(例:SHORTCUT.lnk) をリンク先のファイル名に置換するようにした
* (#368) Lua関数 io.close() が未定義だった。
* (#332)(#369) io.open() のモード r+/w+/a+ を実装した。

NYAGOS 4.4.3\_0
===============
(2019.04.27)

* (#116) readline: Ctrl-Z,Ctrl-`_` による操作取り消しを実装
* (#194) コンソールウインドウの左上のアイコンを更新するようにした
* CMD.EXE 内蔵 date,time を使うためのエイリアスを追加
* `cd 相対パス` の後のドライブ毎のカレントディレクトリが狂う不具合を修正  
  ( `cd C:\x\y\z ; cd .. ; cd \\localhost\c$ ; c: ; pwd` -> `C:\x` (not `C:\x\y`) )

NYAGOS 4.4.2\_3
===============
(2019.04.13)

* Ctrl-RIGHT,ALT-F(次の単語へ), Ctrl-LEFT,ALT-B(前の単語へ)を実装
* インクリメンタルサーチ開始時にトップへ移動する時のバックスペースの数が間違っていた不具合を修正
* (#364) `ESC[0A` というエスケープシーケンスが使われていた不具合を修正

NYAGOS 4.4.2\_1
===============
(2019.04.05)

* diskfree: 行末の空白を削除
* `~"\Program Files"` の最初の引用符が消えて、Files が引数に含まれない不具合を修正

NYAGOS 4.4.2\_0
===============
(2019.04.02)

* OLEオブジェクトからLuaオブジェクトへの変換が日付型などでパニックを起こす不具合を修正
* Luaの数値が実数として OLE に渡されるべきだったのに、整数として渡されていた。
* Lua: 関数: `nyagos.to_ole_integer(n)` (数値を OLE 向けの整数に変換)を追加(trash.lua用)
* Lua: OLEObject に列挙用オブジェクトを得るメソッド `_iter()` を追加
* Lua: OLEObject を開放するメソッド `OLEObject:_release()` を追加
* trash.lua が COM の解放漏れを起こしていた問題を修正
* Lua: `create_object`生成された IUnkown インスタンスが解放されていなかった不具合を修正
* 「~ユーザ名」の展開を実装
* バッチファイル以外の実行ファイルの exit status が表示されなくなっていた不具合を修正
* %COMSPEC% が未定義の時に CMD.EXE を用いるエイリアス(ren,mklink,dir,...)が動かなくなっていた不具合を修正
* 全角空白(%U+3000%)がパラメータの区切り文字と認識されていた点を修正
* (#359) -c,-k オプションで CMD.EXE のように複数の引数をとれるようにした
* 「存在しないディレクトリ\何か」を補完しようとすると「The system cannot find the path specified.」と表示される不具合を修正 (Thx! [tsuyoshicho](https://github.com/tsuyoshicho))
* (#360) 幅ゼロやサロゲートペアな Unicode は`<NNNNN>` と表示するようにした (Thx! [tsuyoshicho](https://github.com/tsuyoshicho))
* サロゲートペアな Unicode をそのまま出力するオプション --output-surrogate-pair を追加
* suコマンドで、ネットワークドライブが失なわれないようにした
* (#197) ソースがディレクトリで -s がない時、`ln` はジャンクションを作成するようにした
* 内蔵の mklink コマンドを実装し、`CMD.exe /c mklink` のエイリアス `mklink` を削除
* ゼロバイトの Lua ファイルを削除(cdlnk.lua, open.lua, su.lua, swapstdfunc.lua )
* (#262) `diskfree` でボリュームラベルとファイルシステムを表示するようにした
* UNCパスがカレントディレクトリでもバッチファイルを実行できるようにした。
* UNCパスがカレントディレクトリの時、ren,assoc,dir,for が動作しない不具合を修正
* (#363) nyagos.alias.COMMAND="string" 中では逆クォート置換が機能しない問題を修正 (Thx! [tostos5963](https://github.com/tostos5963) & [sambatriste](https://github.com/sambatriste) )
* (#259) アプリケーションをダイアログで選んでファイルを開くコマンド `select` を実装
* `diskfree` の出力フォーマットを修正

NYAGOS 4.4.1\_1
===============
(2019.02.15)

* `print(nyagos.complete_for["COMMAND"])`が機能するようにした
* (#356) `type` が LF を含まない最終行を表示しない不具合を修正 (Thx! @spiegel-im-spiegel)
    * 要 [zetamatta/go-texts](https://github.com/zetamatta/go-texts) v1.0.1～
* ビルドに `Go Modules` を使うようにした
* `killall`,`taskkill` コマンド向け補完
* `kill` & `killall`: 自分自身のプロセスを停止できなくした。
* (#261) 補完や1フォルダのlsは10秒でタイムアウトするようにした
* Lua で OLE オブジェクトのセッター(`__newindex`)が効かなかった不具合を修正
* (#357) 仏語キーボードで AltGrシフトが効かない問題を修正 (Thx! @crile)
* (#358) `foo.exe`と`foo.cmd`があった時、`foo`で`foo.exe`ではなく`foo.cmd` が呼び出される不具合を修正

NYAGOS 4.4.1\_0
===============
(2019.02.02)

* `which`,`set`,`cd`,`pushd`,`rmdir`,`env` コマンド向け補完 (Thx! [ChiyosukeF](https://twitter.com/ChiyosukeF))
* (#353) OpenSSHでパスワード入力中に Ctrl-C で中断すると、画面表示がおかしくなる問題を修正 (コマンド実行後にコンソールモードを復旧するようにした) (Thx! [beepcap](https://twitter.com/beepcap))
* (#350) `-l` なしの `ls -F` で os.Readlink を呼ぶのをやめた
* `nyagos.complete_for["COMMANDNAME"] = function(args) ... end` 形式の補完
* (#345) subcomplete.lua で git/svn/hg が効かない問題を修正(Thx! @tsuyoshicho)
* リダイレクトが含まれている時、Lua関数 io.popen が機能しない不具合を修正(Thx! @tsuyoshicho)
* (#354) box.lua のヒストリ補完が C-X h で起動していなかった不具合を修正 (Thx! @fushihara)
* nyagos.d/catalog/subcomplete.lua で `hub` コマンドの補完をサポート (Thx! @tsuyoshicho)

NYAGOS 4.4.0\_1
===============
(2019.01.19)

* "--go-colorable" と "--enable-virtual-terminal-processing" を廃止
* `killall` コマンドを実装
* Linux用の copy と move を実装
* (#351) `END` と `F11` キーが動作もキー割り当てもできなかった不具合を修正

NYAGOS 4.4.0\_0
===============
(2019.01.12)

* バッチファイルを呼ぶ時に、`/V:ON` を CMD.EXE に使わないようにした

NYAGOS 4.4.0\_beta
===================
(2019.01.02)

* Linux サポート(実験レベル)
* ドライブ毎のカレントディレクトリが子プロセスに継承されなかった問題を修正
* ライブラリ "zetamatta/go-getch" のかわりに "mattn/go-tty" を使うようにした
* msvcrt.dll を直接syscall経由で使わないようにした。
* Linux でも NUL を /dev/null 相当へ
* Lua変数 nyagos.goos を追加
* (#341) Windows10で全角文字の前に文字を挿入すると、不要な空白が入る不具合を修正
    * それに伴い、Windows10 では virtual terminal processing を常に有効に
    * `git.exe push`が無効にしても再び有効にする
* (#339) ワイルドカード `.??*` が `..` にマッチする問題を修正
    * 要 github.com/zetamatta/go-findfile tagged 20181230-2

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

NYAGOS 4.1.9\_3
===============
(4017.05.13)

* Fix #214: .nyagos でのバッチ実行時に `main/lua_cmd.go: cmdExec: not found interpreter object` と表示される

NYAGOS 4.1.9\_2
===============
(2017.04.03)

* Fix #191: `-c` オプションが `option parse error` を表示していた。
* 昇格していたら true を返す Lua 関数 `nyagos.elevated()`
* デフォルトのタイトルバーは昇格時に `(admin)` と表示

NYAGOS 4.1.9\_1
===============
(2017.03.28)

* Fix: 4.1.9\_0 の一行入力でカーソルが時々見えなくなる問題を修正
* 新go-colorableの機能で、 でコマンドプロンプトのタイトルを変更するエスケープシーケンス `\033]0;タイトル\007` が使えるようになった。

NYAGOS 4.1.9\_0
===============
(2017.03.27)

* Fix: `open http(s)://...` が機能しなかった不具合を修正
* `cd file:///...` をサポート
* ALT-y: クリップボード文字列が空白を含んでいる時、二重引用符で囲んでペースト
* ファイル名補完の一覧表示で、フルパスのうちのディレクトリ部分を省くようにした
* `history`コマンドで「!」マークで使用する ID を表示していなかった
* %NYAGOSPATH% にあるコマンドも補完されるようにした。
* 補完で環境変数を展開しないようにした。
* 補完で `~/` や `~\` を二重引用符で囲まないようにした。
* `;` や `=` の前の文字列は補完では無視するようにした(setコマンド用)
* ファイル名の大文字・小文字の補正をしないことによる、プロンプトのカレントディレクトリの取得速度の改善
* 二重引用符無しでも `cd C:\Program Files` が機能するようにした
* cd /D を機能するようにした(#182 CMD.EXE との互換性のため /D オプションは無視される)
* `history` で時間順にソートするようにした
* `open regedit` を機能させるため、`open` でのファイル存在チェックを省く
* `clone`,`su`,`sudo`: ネットワークフォルダーで失敗させないよう、シンボリックリンクの宛先パスで ShellExecute を行うようにした (#122)
* set の動作を CMD.EXE 互換とした(`set FOO=A B` が `set FOO="A B"` と同じ)
* #184 `_nyagos` 内で逆クォートが効かなかった不具合を修正
* `_nyagos`: `bindkey KEYNAME FUNCNAME` を実装
* CMD.EXE と同様の `%環境変数名:被置換文字列=置換文字列%` をサポート
* インクリメンタルサーチで ESCAPE キーを検索モード終了に割り当てた。
* カーソル選択型補完(選択用の内蔵コマンド box を新設)
    * Ctrl-O          : カーソルで選択したファイル名を挿入する (by box.lua)
    * Ctrl-XR , Alt-R : カーソルで選択したヒストリを挿入する (by box.lua)
    * Ctrl-XG , Alt-G : カーソルで選択したGit Revisionを挿入する(by box.lua)
    * Ctrl-XH , Alt-H : カーソルで選択した過去に移動したディレクトリを挿入する(by box.lua)
* `lua_e "nyagos.key = function(this) end"` というキーアサインをサポート

NYAGOS 4.1.8\_0
===============
(2017.02.15)

* COMMAND.COMバッチ風の新カスタマイズファイルとして `_nyagos` を用意
* Fix #173 `ls` や内蔵コマンドを Ctrl-C で止められるようになった
* ls -h のファイルサイズを 1K,2M 等ではなく、カンマ区切りの数値とした
* nyagos.lines(FILENAME,"n") を実装した(ただし、実数ではなく整数)
* nyagos.exe の中だけで機能する %PATH% 的な環境変数 %NYAGOSPATH% を追加
* vim のような SET VAR+=VALUE , VAR^=VALUE をサポート
* Fix #176 `gawk "BEGIN{ print substr(""%01"",2) }"` がエラーになっていた
* アイコンを付けるのに、windres.exe ではなく github.com/josephspurrier/goversioninfo を使うようにした
* command.com と同程度の `if` をサポート(`==`,`not`,`errorlevel`,`/I`)
* alias に新マクロを追加 `$~1` `$~2` ... `$~*` (前後の二重引用符を削除する)
* カレントディレクトリ,時刻,PID もヒストリに記録するようにした (#112)
* ls -l: タイムスタンプのフォーマットを 'Jan 2 15:04:05' or 'Jan 2 2006'へ変更
* lua53.dll が無い時、スタックトレースではなくエラーを表示するようにした
* '#' 以降をコメントとみなすようにした
* open,clone,su,sudo を Lua から Go に書き直した

NYAGOS 4.1.7\_0
===============
(2016.11.29)

* nyagos.lua を廃止した。その役割は nyagos.exe 自身が担うようにした。
* `~/.nyagos` を`%APPDATA%\NYAOS_ORG/dotnyagos.luac` にキャッシング
* `nyagos.d/*` を nyagos.exe 自体にバンドルするようにした
* Fix #167 相対パスにシンボリックリンクされた実行ファイルが動かなかった
* Fix `ls -l` でリンクされた実行ファイルに @ とリンク先が表示されていなかった
* Fix su.lua: clone/su で文字化けしたパスが表示されていた
* Fix #168 `ls 相対パスのシンボリックリンク` がエラーになっていた
* Fix `ls -lh` の時のファイルサイズの表示幅がおかしくなっていた
* `ls -oFh` をデフォルトの ls のエイリアスにした
* `history` で標準出力が端末ではない時、全行を出力するようにした
* `open` で複数のファイルが指定された時にプロンプトを表示するようにした
* `use "cho"` → [cho](https://github.com/mattn/cho) 向け拡張
        * C-r: ヒストリ
        * C-o: ファイル名
        * M-h: ディレクトリヒストリ
        * M-g: Git のリビジョン名
* Fix: {a,b,c} といったブレース展開が、引用符の中でも機能していた不具合を修正

NYAGOS 4.1.6\_1
===============
(2016.09.07)

* Fix: パッケージの ZIP ファイルに lua53.dll が含まれていなかった。

NYAGOS 4.1.6\_0
===============
(2016.09.07)

* スペースとバックスペースで行っていた行末削除に "\x1B[0K" を使うようにした
* m回のバックスペースに "\x1B[mC" を使うようにした。
* Fix #159: 端末幅を変更した時にプロンプトから再表示していたのを廃止
* Fix #164: `cd --history` でカレントディレクトリがホームに移動していた
* stat 取得の成否にかかわらず、`[\\/:]\.{0,2}$` にマッチする宛先パスをディレクトリとみなすようにした。

NYAGOS 4.1.5\_1
===============
(2016.07.31)

* Fix #157++: 端末サイズ変更後、追記でズレる不具合を修正
* 4.0.x の不適切なデフォルト ~/.nyagos 向けに、prompter という名前の上位値がクロージャ(nyagos.prompt)で使われていたらエラーにするようにした (#155,#158)

NYAGOS 4.1.5\_0
===============
(2016.07.31)

* カレントディレクトリのヒストリがゼロの時に peco がハングしないように、`cd --history` の先頭にカレントディレクトリを出力するようにした。
* Luaで `nyagos.option.glob = true` とすると、外部コマンドでもワイルドカード展開するようにした。(#150)
* source の互換性改善を試みた
* nyagos.lines(FILENAME,X) の X='a','l','L',数値のサポート(#147)
* Fix #156: %U+0000% でパニックが発生する
* Fix #152: 「ls -ld Downloads\」の結果が「Downloads\/」となる
* Fix #157: 端末サイズ変更時の、一行入力の表示幅を再設定するようにした
* 内蔵パッケージを別レポジトリヘ外出し

NYAGOS 4.1.4\_1
===============
(2016.06.12)

* `&&` や `||` が ` ;`と等価になっていた不具合を修正(#151)
* @DeaR さん提供の autocd.lua & autols.lua を nyagos.d/catalog に追加(#149)

NYAGOS 4.1.4\_0
===============
(2016.05.29)

* 簡易OLEインターフェイスを実装した。NYOLE.DLL は不要になった。
* デフォルトのプロンプト表示関数を `nyagos.default_prompt` と定義し、第二引数で端末タイトルを変更できるようにした
* Fix: nyagos.lines() が改行を削除していなかった
* Fix: Lua のデフォルトファイルハンドル(標準入出力)がバイナリモードでオープンされていた(#146)
* nyagos.d/catalog/peco.lua: C-r: 表示順を反転させて、速度を改善した。

NYAGOS 4.1.3\_1
===============
(2016.05.08)

* Fix: ヒストリがファイルに保存されない #138
* Fix: nyagos.history を削除すると、exit で終了するまで警告が出続ける
* Fix: nyagos.d/catalog/peco.lua: nyagos.history が存在しないと、peco がハングする

NYAGOS 4.1.3\_0
===============
(2016.05.05)

* Add: `nyagos.open(PATH,MODE)` UTF8版`io.open`
* Add: `nyagos.loadfile(PATH)` UTF8版`loadfile`
* Add: `nyagos.lines(PATH)` UTF8版`io.lines`(注意:戻り値はバイト列、ファイル名だけがUTF8指定になった)
* 内蔵`echo`の改行コードとして LF ではなく CRLF を使うようにした (#124)
* Lua のデフォルト入出力を NYAGOS のリダイレクトに追随させるようにした
* touch コマンドに -r と -t オプションを実装した
* touch コマンドで簡易日時フォーマットチェックを入れた
* `make install` でログを残して、3秒後にインストール窓を閉じるようにした(#107)
* `nyagos < TEXTFILE` が利用可能になった (#125)
* {conio,dos}/const.go を再作成するのに lua.exe,findstr.exe は不要になった
* 標準エイリアス suffix が機能していなかった
* カレントドライブがネットワークドライブでも、`su` は新しい管理者モード nyagos を同じ UNC-Path でディレクトリで起動させられるようにした。
* `nyagos -c 'CMD'` で CMD は `nyagos.lua` の後に実行するようにした。
* `nyagos -[cfe] "..."や `nyagos < TEXTFILE` では著作権表示を出さないようにした
* Fix: `make install DIR` が次回の `make install` 向けに DIR をセーブしていなかった。
* Fix: nyagos.exe が日本語フォルダーに置いてある時、nyagos.lua をロードできていなかった。
* Fix: nyagos.d/catalog/subcomplete.lua が 4.1 以降で動かなくなっていた (#135)
* エスケープシーケンスエミュレータをgithub.com/mattn/go-colorable に変更 (#137)
* Fix: `ls -ltr * `で時系列でソートされていなかった (#136)
* nyagos -f で拡張子が .lua で無い時、シェルコマンドが格納されたファイルと解釈するようにした

(2016.05.17 追記)
-----------------
* ANSI文字列とUTF8文字列の混乱を避けるため、print でエスケープシーケンス入りの UTF8 文字列出力を廃止した。print は lua53.dll 内蔵のもののままとなった( #129 )

NYAGOS 4.1.2\_0
===============
(2016.03.29)

* スクリプトのカタログシステムを作った
    - スクリプト `catalog.d\*.lua` を `nyagos.d\catalog\.` へ移動
    - カタログのスクリプトを .nyagos より `use "NAME"` で利用できるようにした
        - `use "dollar"` → `$PATH`形式で環境変数を展開
        - `use "peco"` → [peco](https://github.com/peco/peco) 向け拡張
            * C-r: ヒストリ
            * C-o: ファイル名
            * M-h: ディレクトリヒストリ
            * M-g: Git のリビジョン名
* ls
    - 壊れたシンボリックリンクがあっても ls は中断しないようにした。
    - `ls -d` をサポート
* .nyagos を nyagos.exe と同じディレクトリに置けるようにした。
* cd のヒストリ全てを `cd --history` で出せるようにした
* 組込みの簡易`touch`コマンドを実装
* ファイルが存在しない時に、>> が失敗する不具合を修正
* Lua関数の第一パラメータテーブルのメンバに rawargs を追加
  (ユーザ入力文字列から引用符が削除されていない文字列を格納したテーブル)
* bindkeyのコールバック関数の引数テーブルに `replacefrom` メソッドを追加

NYAGOS 4.1.1\_2
===============
(2016.02.17)

* Lua の loadfile 等を呼ぶ際に UTF8 を ANSI へコンバートしていなかった不具合を修正 (#110,Thx Mr.HABATA)

NYAGOS 4.1.1\_1
===============
(2016.02.16)

* プロンプトが長すぎる時、強制的に改行するようにした (#104)
* ls でワイルドカードがマッチしない時のメッセージを修正 (#108)
* %ProgramFiles(x86)%のような環境変数が展開できてなかった点を修正(#109,Thx @hattya)

NYAGOS 4.1.1\_0
===============
(2016.01.15)

* キー入力で UTF16 のサロゲートペアをサポート
* mkdirに必要に応じて親ディレクトリを作成する /p オプションを追加

NYAGOS 4.1.0\_0
===============
(2016.01.03)

* 内蔵コマンド ln を追加
* Lua コマンド lns を追加 (UACを表示後、`ln -s` を実行する)
* `ls -l` でシンボリックリンクの宛先を表示
* あるファイルでcopy/move 時に失敗した時、以降のファイルを続けるか問合せるようにした。
* 新変数: `nyagos.histchar`: ヒストリ置換文字(デフォルト「`!`」)
    - ヒストリ置換を完全に無効にする場合、`nyagos.histchar = nil`
* 新変数: `nyagos.antihistquot`: ヒストリ置換を抑制する引用符(デフォルト「`'"`」)
    - 【注意】`"!!"` は「デフォルト」では置換されなくなりました
    - 4.0互換にするには `nyagos.antihistquot = [[']]` とする
* 新変数: `nyagos.quotation`: 補完でのデリミタ文字(デフォルト「`"'`」)。
    - `nyagos.quotation` の最初の文字がデフォルトの引用符となる。
    - 二番目以降の文字は、ユーザが補完前に使用していた場合に採用される
    - `nyagos.quotation=[["']]`の場合
        - `C:\Prog[TAB]` → `"C:\Program Files\ ` (`"` が挿入される)
        - `'C:\Prog[TAB]` → `'C:\Program Files\ ` (`'` が維持される)
        - `"C:\Prog[TAB]` → `"C:\Program Files\ ` (`"` が維持される)

NYAGOS 4.1-beta
================
(2015.12.13)

* クラッシュ回避のため、全てのLua のコールバック関数はそれぞれの Lua
  インスタンスを持つようにした。
* コールバック関数と .nyagos 間で値を共有するため、テーブル share[] を作った
* `*.wsf` を cscript に関連付けた
* `nyagos[]` への不適切な代入を警告するようにした。

<!-- vim:set fenc=utf8: -->
