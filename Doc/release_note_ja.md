* 簡易OLEインターフェイスを実装した。NYOLE.DLL は不要になった。
* デフォルトのプロンプト表示関数を `nyagos.default_prompt` と定義し、第二引数で端末タイトルを変更できるようにした
* Fix: nyagos.lines() が改行を削除していなかった(io.lines非互換だった)
* Fix: Lua のデフォルトファイルハンドル(標準入出力)がバイナリモードでオープンされていた(#146)
* nyagos.d/catalog/peco.lua: 表示順を反転させて、速度を改善した。

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
