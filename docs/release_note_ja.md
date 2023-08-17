[top](../readme_ja.md) &gt; [English](release_note_en.md) / Japanese

当バージョンのバイナリは Go 1.20.7 でビルド。  
サポート対象は Windows 7, 8.1, 10, 11, WindowsServer 2012R以降, Linux となります。

* nyagos.d/suffix.lua: 環境変数 NYAGOSEXPANDWILDCARD にリストされているコマンドのパラメータはワイルドカードを自動展開するようにした。
* (#432) `set -o glob` 時、二重引用符内の`*`,`?` がワイルドカードとして展開されていた(本来されるべきではない)
* (#432) 新オプション `glob_slash` を追加。設定されている時、ワイルドカード展開で `/` を使う
* Linux版で逆クォートがエラーになって機能しない不具合を修正 (Lua関数 atou が常に "not supopported" を返していたので、引数と同じ値を戻すようにした)
* [SKK] \(Simple Kana Kanji conversion program\) サポート :[設定方法][SKKSetUp]
* 適切なUTF8文字列でない時は ANSI文字列とみなして UTF8変換を試みる関数 `nyagos.atou_if_needed` を追加
* (#433) 文字化けを避けるために、逆クォートでは `nyagos.atou_if_needed` を使って、UTF8 を更に UTF8 化させないようにした
* `more`, `nyagos.getkey`, `nyagos.getviewwidth` が Windows 7, 8.1 や WindowsServer 2012 で動かない可能性があった問題を修正。それらは Windows10,11 の新端末に依存する "golang.org/x/term" を使用していました。(本件は v4.4.13\_3 のみに含まれていた)

[SKK]: https://ja.wikipedia.org/wiki/SKK
[SKKSetUp]: https://github.com/nyaosorg/nyagos/blob/master/docs/10-SetupSKK_ja.md

NYAGOS 4.4.13\_3
================
(2023.04.30)

* (#431) バッチファイル実行で変更した環境変数や、more/typeの出力など非UTF8からUTF8へ変更する時、4096バイトを越えるような行の変換で失敗する不具合を修正 (Thx. @8exBCYJi5ATL)
* (#431 とは別件で) 行のサイズが大きすぎると、more が行を出力しないことがある不具合を修正

NYAGOS 4.4.13\_2
================
(2023.04.25)

* (#428) `rmdir /s` でシンボリックリンクやジャンクションの削除に失敗する問題を修正
* \" で引用領域の色を反転させないようにした
* (#429) カレントディレクトリが `C:` の時、`cd c:` が失敗する不具合を修正
* ANSIとUTF8間の変換に go-windows-mbcs v0.4 とgolang.org/x/text/transform を使うようにした

NYAGOS 4.4.13\_1
================
(2022.10.15)

* (#425) nyagos.d/suffix.lua で設定した拡張子を環境変数 PATHEXT ではなく、NYAGOSPATHEXT に登録するようにした。コマンド名補完は PATHEXT に加え、NYAGOSPATHEXT も参照するようにした  (Thx. @tsuyoshicho)
* (#426) nyagos.argsfilter でコマンドラインが変換された時、空パラメータ("") が消えてしまう問題を修正 (Thx. @juggler999)
* (#426) 外部コマンドに対するワイルドカード展開が有効になっている時、空パラメータ("")が消えてしまう問題を修正 (Thx. @juggler999)
* (#427) '""' が BEEP音プラス '(数字)' に置換されてしまう不具合を修正 (Thx.@hogewest)

NYAGOS 4.4.13\_0
================
(2022.09.24)

* セキュリティー警告対応のため、gopkg.in/yaml.v3 に依存するモジュールを直接利用しないようにした( https://github.com/nyaosorg/nyagos/security/dependabot/1 )
* (#420) macOS でのビルドをサポート (Thanks to @zztkm)
* (#421) バッチファイルが環境変数を削除した時、それが反映されていない問題を修正(Thanks to @tsuyoshicho)
* (#422) プロンプトに$hを使った時、編集開始位置が右にずれる不具合を修正(Thanks to @Matsuyanagi)
* 端末の背景色が白の時、入力文字が全く見えない不具合を修正(端末デフォルト文字色を使用)
* コマンド名補完において、タイムアウトが効きづらい問題を改善
* サブパッケージを internal フォルダー以下へ移動
* キーの英大文字・小文字を区別しない辞書に Generics を使用するようにした。
* Windows以外で Makefile がエラーになる問題を修正
* (#424) fuzzy finder 拡張機能の統合 (Thx @tsuyoshicho)

NYAGOS 4.4.12\_0
================
(2022.04.29)

* カラー化コマンドラインの修正
    * オプション文字列の色を黄土色へ変更し、範囲を修正
    * 全角空白の背景色を赤へ変更
    * -0...-9 にオプション用のカラーがついていなかった
    * 端末の背景透過が効くように、背景を黒(ESC[40m)からデフォルト色(ESC[49m)にした(go-readline-ny v0.7.0)
    * WezTerm ではサロゲートペア表示を有効にした
* Linux版のみの不具合修正
    * set -o noclobber 設定時のリダイレクト出力がゼロバイトになってしまう不具合を修正
    * ヒストリがファイルに保存されず、次回起動の際に復元されない不具合を修正
* start コマンドのパラメータは %PATH% 上の任意のファイル・ディレクトリに補完する
* cmd.exe と同様に `rmdir /s` は readonly のフォルダーも削除できるようにした
* `rmdir FOLDER /s` とオプションをフォルダーの後に書けるようにした
* (#418) コマンドラインが ^ で終わっていた場合、Enter入力後も行入力を継続するようにした

NYAGOS 4.4.11\_0
================
(2021.12.10)

* コマンドラインをカラー化した

NYAGOS 4.4.10\_3
================
(2021.08.30)

* (#412) Windows10の、WindowsTerminalでない端末で、罫線キャラクターの幅が不正確になっていた問題に対応
* パッケージに新しいアイコンファイルを添付

NYAGOS 4.4.10\_2
================
(2021.07.23)

* コードページ437 で、%DATE% の置換結果が CMD.EXE 非互換になっている不具合を修正
* go-readline-ny v0.4.13: Windows Terminal で Mathematical Bold Capital (U+1D400 - U+1D7FF) の編集をサポート
* -oオプションがついていも、リダイレクトされた ls の出力から ESC[0m を除くようにした
* (#411) 英語部と日本語部が入れ変わっていたドキュメントを修正 (Thx! @tomato3713)
* テストコードの整理、自動化

NYAGOS 4.4.10\_1
================
(2021.07.02)

* `./dll` というフォルダが存在して、CDPATH上に DLL がある時、入力されたパス `dll` が補完で DLL に置き変わっていた動作を修正(英大文字・小文字が変わってしまうのが問題)
* 空白を含むディレクトリでの clone コマンドで、カレントディレクトリが維持されない不具合を修正

NYAGOS 4.4.10\_0
================
(2021.06.25)

* nyagos.d/aliases.lua: abspath をワイルドカード対応にした
* Luaで、`OLEOBJECT:_release()` ではなく、 `OLEOBJECT._release()` が使われた時に適切なエラーが起こされるようにした(glua-oleパッケージを更新)
* CMD.EXE と同様、echo で二重引用符を出力させるようにした
* `"..\"..\".."` という文字列が字句解析の結果、`"..\"..".."` になってしまう不具合を修正
* (#410) 端末の閉じるボタンが押された時にただちに終了させるため SIGTERM を無視しないようにした (Thx @nocd5)
* WindowsTerminal 1.8 以降で、cloneコマンドは同じウインドウの別タブで起動するようにした
* アプリ実行エイリアス経由で wt.exe が呼べず、WindowsTerminal下で clone コマンドが動かなくなっていた問題を修正
* suでネットワークドライブが引き継がれてない不具合を修正
* ビルドするのに PowerShell(make.cmd) ではなく Makefile(GNU Make)を使うようにした
* Linux でもビルドできるようにした

NYAGOS 4.4.9\_7
===============
(2021.05.22)

* (#409) `set -o glob` や `nyagos.option.glob=true` による外部コマンド向けワイルドカード展開が効かなくなっていた不具合を修正 (Thx @juggler999)

NYAGOS 4.4.9\_6
===============
(2021.05.07)

* (#406) nyagos.argsfilter で生引数がコンバートされず、suffixコマンドが期待どおり機能しなくなっていた不具合を修正 (Thx @tGqmJHoJKqgK)

NYAGOS 4.4.9\_5
===============
(2021.05.03)

* go-readline-ny v0.4.10: Yes/Noの回答のYが次のコマンドラインに入力される不具合を修正
* go-readline-ny v0.4.11: Emoji Moifier Sequence (skin tone) をサポート
* Windows 8.1でCPU負荷が高い時のカラーlsの速度を改善した
* ( "io/ioutil" を使わないようにした )
* open で開いたプロセスが閉じる時のメッセージがコマンドラインに重ならないようにした
* go-readline-ny v0.4.12: VisualStudioCodeのターミナルでは絵文字編集はオフにするようにした
* (#403) CMD.EXE のような -S,-C,-K オプションをサポート
* (#403) コマンドラインの不規則な二重引用符が外部コマンドに渡される時に削除される問題を修正
* (#405) fuzzyfinder catalog module を追加 (Thx @tsuyoshicho)

NYAGOS 4.4.9\_4
===============
(2021.03.06)

* (#400) サブコマンド補完向けのコマンドの存在チェック追加(Thx @tsuyoshicho )
* (#401) choco/chocolaty 向けサブコマンド名補完追加(Thx @tsuyoshicho )
* WindowsTerminal で ls や Ctrl-Oの選択時のレイアウトが崩れる問題を修正
* go-readline-ny v0.4.4: 任意の一字 + Ctrl-B + 合字が入力された時、表示が乱れる問題を修正
* go-readline-ny v0.4.5: 合字の異体字サポート
* go-readline-ny v0.4.6: 異体字の後の囲み記号の編集をサポート(&#x0023;&#xFE0F;&#x20E3;)
* (#402) "echo !xxx" でシェルがいきなり終了してしまう問題を修正 (Thx @masamitsu-murase)
* go-readline-ny v0.4.7: REGIONAL INDICATOR (U+1F1E6..U+1F1FF) でカーソル位置が狂わないようにした
* go-readline-ny v0.4.8: WAVING WHITE FLAG and its variations (U+1F3F3 U+FE0F?)
* go-readline-ny v0.4.9: RAINBOW FLAG (U+1F3F3 U+200D U+1F308)

NYAGOS 4.4.9\_3
===============
(2021.02.20)

* WindowsTerminal利用下の一行入力で Unicode の異体字をサポート
* (#397) scoopコマンドのサブコマンド補完を追加 (`use "subcomplete.lua"`) (Thx @tomato3713)
* 補完時に最も短い候補に英大文字/小文字をあわせるようにした
* (#398) io.popen の第二引数のデフォルトが機能していなかった (Thx @ironsand)
* (#399) utf8 offset の改良 (Thx @masamitsu-murase)
* ALT-/ キーバインドのサポート (Thx @masamitsu-murase) https://github.com/zetamatta/go-readline-ny/pull/1
* WindowsTerminal 1.5 で絵文字や丸数字が入力できなくなっていた問題を修正

NYAGOS 4.4.9\_2
===============
(2021.01.08)

* (#342) Ctrl-C 押下時に子プロセスを kill しないようにした

NYAGOS 4.4.9\_1
===============
(2020.12.21)

* パス引数なしの`make install` が失敗する不具合を修正
* (#396) Ctrl-W で左へのスクロールが必要な時に panic する不具合を修正
* コンソール入力の more/clip/type がエコーバックしない時がある不具合を修正
* (#342) クラッシュさせないよう Ctrl-C 割り込みハンドリングを改善

NYAGOS 4.4.9\_0
===============
(2020.12.05)

* (#390,#394) Unicode の合字をサポート
* 異字体コード1-16があるとカーソル位置がおかしくなる不具合を修正  
  (異字体コード自体は未対応なので &lt;FE0F&gt; などと表示する)
* su と clone で WindowsTerminal をサポート
* 編集中はバックグランドプロセスの開始・終了メッセージを出させないようにした
* C-r: インクリメンタルサーチでは英大文字・小文字を区別しないようにした
* || や && の後でコマンド名補完が効かなかった不具合を修正
* C-y: ペースト時に最後の CRLF を除くようにした
* Fix: (#393) ウィンドウアクティブ後の最初のキーが２つ入力されます (Thanks to @tostos5963)
* アンチウィルスが誤判断をするので upx.exe を使用しないようにした

NYAGOS 4.4.8\_0
===============
(2020.10.03)

* git.lua: `git add` 向け補完
    * "\343\201\202\343\201\257\343\201\257"といったファイル名のクォーテーションを解除するようにした
    * untrackなディレクトリの下のファイルも補完対象とした
* diskused: サイズ表記を `ls -h` のように
* diskused が Ctrl-C で止まらなかった不具合を修正
* %ENV:~10,5% のような環境変数抽出を実装
* (#308) UNCパスで表現されていないネットワーク上の GUI 実行ファイルを起動しようとすると `The operation was canceled by the user` というエラーになる問題を修正
* nyagos がネットワーク上にある時、clone コマンドでエラーダイアログが出る問題を修正
* (#389) su: SUBST コマンドのドライブマウントを維持するようにした
* (#390) U+2000～U+2FFF の Unicode が入力できない不具合を修正
* (#390) サロゲートペアな文字が入力できない不具合を修正
* box.lua: Ctrl+O→ESCAPE でユーザが入力した単語が消える不具合を修正
* (#391) subcommand.lua: ghコマンド向けサブコマンド補完を追加 (Thanks to @tsuyoshicho)

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
