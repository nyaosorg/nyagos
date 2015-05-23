# NYAGOS - Nihongo Yet Another GOing Shell

NYAGOS は、go言語による Windows 用コマンドラインシェルです。
Nihongo とありますが、Unicode ベースなので、特に特定の自然言語に
特化しているわけではありません。

NYAGOS は、Windows の文化を尊重しつつ、UNIX に慣れた人が、
Windows であまりストレスを感じないような環境を構築するために
開発されているシェルです。bash など特定のシェルの互換を目指す
ものではありません。

特徴

* 現在のコードページに関わらず、JISに未登録の Unicode も扱えます。
   * クリップボードにある Unicode 文字のペースト・編集可能
   * 直接入力できずとも %U+XXXX% というリテラルで変換可能
   * プロンプトにも $Uxxxx というマクロが使用可能
* 高機能内蔵 ls
   * カラー表示(-o オプション)
   * ジャンクションを区別可能(-F オプションで「@」がファイル名末尾につく)
* UNIXライクなシェル機能
   * ヒストリ機能(Ctrl-P や ! 文字による置換)
   * エイリアス
   * ファイル名・コマンド名補完
   * 引数を囲むのにシングルクォーテーションが利用可能
* Lua言語を使ったカスタマイズ機能
   * Lua で内蔵コマンドを作成可能
   * 入力文字列を Lua で加工できる
   * コードページ文字列⇔UTF8変換関数や、eval関数などの支援関数も用意

## インストール

ファイル:`nyagos.exe`,`nyagos.lua`,`lua53.dll`、ディレクトリ`nyagos.d`を
`%PATH%` の差すディレクトリに置いてください。
(同一のディレクトリに置いてください)

カスタマイズ用ファイル `.nyagos` は、`%USERPROFILE%` か `%HOME%`
の差すディレクトリに置いて、必要に応じて修正してください。

## 起動オプション

### `-h`

オプションのヘルプを表示します。

### `-c "コマンド"`

コマンドを実行して、ただちに終了します。

### `-k "コマンド"`

コマンドを実行してから、通常起動します。

### `-f スクリプトファイル名 引数1 …`

Luaインタプリタでスクリプトファイルを実行後、終了します。
引数は配列 arg[] という形で参照できます。

### `-e "スクリプトコード"`

Luaインタプリタでスクリプトコードを実行後、終了します。

## 編集機能

UNIX系シェルに近いキーバインドで、コマンドラインを編集可能です。

* BackSpace , Ctrl-H : カーソル左の一文字を削除
* Enter , Ctrl-M     : 入力終結
* Del                : カーソル上の一文字を削除
* Home , Ctrl-A      : カーソルを先頭へ移動
* ← , Ctrl-B        : カーソルを一文字左へ移動
* Ctrl-D             : 0文字の時は NYAGOS を終了、さもなければ Del と同じ
* End , Ctrl-E       : カーソルを末尾へ移動
* → , Ctrl-F        : カーソルを一文字右へ移動
* Ctrl-K             : カーソル以降の文字を全て削除し、クリップボードへコピー
* Ctrl-L             : 画面をクリアして、入力した内容を再表示
* Ctrl-U             : カーソルまでの文字を全て削除し、クリップボードへコピー
* Ctrl-Y             : クリップボードの内容を貼り付ける
* Esc , Ctrl-[       : 入力内容を全て削除する
* ↑ , Ctrl-P        : ヒストリ：一つ前の入力内容を展開する
* ↓ , Ctrl-N        : ヒストリ：一つ後の入力内容を展開する
* TAB , Ctrl-I       : ファイル名・コマンド名補完
* Ctrl-C             : 入力内容を破棄
* Ctrl-R             : インクリメンタルサーチ

## 内蔵コマンド

これらのコマンドはコマンド名とは別にエイリアスを持っています。
たとえば `ls` は `__ls__` というエイリアスを持っています。

### `alias エイリアス名=定義`

エイリアスを設定します。置換内容には以下のマクロが使えます。

* `$n` (n:数字) n番目の引数となります
* `$*` 全ての引数に置換されます。

置換内容が空の時はエイリアスを削除します。
= が無い場合は、その名前のエイリアスの内容を表示します。
引数が無い場合は、全エイリアスを一覧します。

### `cd ドライブ:ディレクトリ`

現在のカレントドライブ、ディレクトリを変更します。
引数を省略すると、CMD.EXE と違い、環境変数 HOME 、あるいは 
USERPROFILE の差す先のディレクトリへ移動します。
CMD.EXE と違い、ドライブも同時に変更します。

* `cd -` : 一つ前にいたディレクトリへ移動します
* `cd -N` : N 回前のディレクトリへ移動します
* `cd -h` , `cd ?` : 過去いたディレクトリを表示します

### `exit`

NYAGOS を終了します。

### `history [件数]`

ヒストリ内容を表示します。件数を省略すると、最近の10件が表示されます。

### `ls [-オプション] …`

ディレクトリの一覧を表示します。
サポートしているオプションは以下の通りです。

* `-l` ロングフォーマットで一覧を表示します。
* `-F` ディレクトリ名の末尾に /  を、実行ファイル名の末尾に * を表示します。
* `-o` カラー化します
* `-a` 隠しファイルや「.」で始まるファイル名を含め、全て表示します。
* `-R` サブディレクトリ以下も表示します。
* `-1` ファイル名だけを表示します。
* `-t` 最終変更日時でソートします。
* `-r` ソート順を逆転します。
* `-h` 説明を表示します。

### `pwd`

現在のカレントドライブ + ディレクトリを表示します。

* `pwd -N` : N 回 cd で移動する前のディレクトリを表示します。

### `set 変数名=値`

環境変数に値を設定します。値に空白等を含む場合、CMD.EXE と同様に
「`set "変数名=値"`」とします。= 以降を省略すると、現在の変数の内容を
表示します。

以下の変数は特別な意味を持ちます。

* `PROMPT` … プロンプトの文字列を設定します。`$P` 等のマクロ文字はCMD.EXE と同じです。shiena 様開発のモジュールによりエスケープシーケンスが使えます。

### `copy SOURCE-FILENAME DESTINATE-FILENAME`
### `copy SOURCE-FILENAME(S)... DESINATE-DIRECTORY`
### `move OLD-FILENAME NEW-FILENAME`
### `move SOURCE-FILENAME(S)... DESITINATE-DIRECTORY`
### `del FILE(S)...`
### `erase FILE(S)...`
### `mkdir NEWDIR(S)...`
### `rmdir [/s] DIR(S)...`
### `pushd`
### `popd`
### `dirs`

これらの内蔵版は、上書きや削除の際に常にプロンプトで実行可否を問い合わせます。

### `source バッチファイル名`

バッチファイルを CMD.EXE で実行して、CMD.EXE が変更した環境変数と
カレントディレクトリを NYAGOS.EXE に取り込みます。

コマンド名として「`source`」の代わりに「`.`」(ドット)一文字も使う
ことができます。

## 起動処理

1. 起動時に nyagos.exe と同じフォルダの nyagos.lua を読み込みます。nyagos.lua はLua で記述されており、ここから更にホームディレクトリ(%HOME% or %USERPROFILE%)の .nyagos の Lua コードを読み込みます(nyagos拡張は後述)。ユーザカスタマイズは、この .nyagos を編集して行うことができます。
2. 過去のヒストリ内容を `%APPDATA%\NYAOS_ORG\nyagos.history` から読み出します。NYAGOS 終了時には、このファイルに再び最後のヒストリ内容が書き出されます。

## コマンドライン置換

### ヒストリ置換

* `!!`  一つ前の入力文字列へ
* `!n`  最初から n 番目に入力文字列へ
* `!-n` n 個前に入力した文字列へ
* `!STR` STR で始まる入力文字列へ
* `!?STR?` STR を含む入力文字列へ

以下のような語尾をつけることができます。

* `:0` コマンド名を引用する。
* `:m` m 番目の引数だけを引用する。
* `^`  最初の引数だけを抜き出す。
* `$`  最後の引数だけを抜き出す。
* `*`  全ての引数を引用する。

### 環境変数置換

* コマンドや引数先頭の `~` を `%HOME%` あるいは `%USERPROFILE%` に置換します。

### Unicode リテラル

* `%u+XXXX%` (XXXX:16進数) を Unicode 文字に置換します。

## Lua拡張

nyagos では、EXE の本体の機能はコンパクトとし、便利機能は 
なるべく Lua で機能を拡張できるよう設計を進めています。
現在は以下のような関数が使用できます。

### `nyagos.setalias("エイリアス名","置換コード")`

エイリアスを設定します。nyagos.lua 内で、これを簡略した

* `alias "エイリアス名=置換コード"`
* `alias{ エイリアス名="置換コード" , エイリアス名="置換コード" … }`

が定義されています(Lua は引数が一つの場合は括弧を省略できます)。
置換コードでは「alias」コマンドと同様、`$1` や `$*` などのマクロが
使用可能です。

### `nyagos.setalias("エイリアス名",function(args)～end)`

Lua 関数をエイリアスコマンドとして呼び出せるようにします。
args には全引数を格納したテーブルが入ります。

### `nyagos.getalias("エイリアス名")`

現在 "エイリアス名" に設定されている文字列もしくは Lua 関数が返します。

### `nyagos.setenv("環境変数名","変数内容")`

環境変数を設定します。nyagos.lua 内で、これを簡略した

* `set "変数名=定義内容"`
* `set "変数名+=追加定義"`
* `set{ 変数名="定義内容" , 変数名="定義内容"}`

が定義されています(Lua は引数が一つの場合は括弧を省略できます)。
`set` は `nyagos.setenv` の処理に %変数名% の展開機能なども
組込まれています。

### `status,err = nyagos.exec("シェルコマンド")`

シェルコマンドを実行します。エラーが発生した時、
status に nil が入り、err にエラーメッセージが入ります。
標準 nyagos.lua では、これの別名として

* `x'シェルコマンド'`

を定義しています。

### `nyagos.eval("シェルコマンド")`

nyagos.exec と同じですが、標準出力を取り込んで、戻り値として返します。
実行に失敗した場合などは nil が戻ります。

### `nyagos.write(テキスト)`

テキストを標準出力に出力しますが、リダイレクトされている場合は
文字コードはUTF8 になります。内蔵 Lua の print は 
nyagos.write(テキスト..'\n') に差し替えられています。

### `nyagos.writerr(テキスト)`

テキストを標準エラー出力に出力しますが、リダイレクトされている場合は
文字コードはUTF8 になります。

### `nyagos.getwd()`

現在のカレントディレクトリを返します。

### `nyagos.utoa(UTF8文字列)`

UTF8文字列を、現在のコードページの文字列に変換します。

### `nyagos.atou(ANSI文字列)`

現在のコードページの文字列を、UTF8 へ変換します。

### `nyagos.glob(ワイルドカード文字列1,ワイルドカード文字列2,...)`

ワイルドカードを展開し、それらを格納したテーブルを返します。

### `path = nyagos.pathjoin('パス1','パス2'...)`

パスの要素を連結して、一つのパスにします。

### `nyagos.bindkey("キー名","機能名")`

一行入力のキーに機能を割り当てます。

キー名として以下が使えます。

        "C_A" "C_B" ... "C_Z" "M_A" "M_B" ... "M_Z"
        "F1" "F2" ... "F24"
        "BACKSPACE" "CTRL" "DEL" "DOWN" "END"
        "ENTER" "ESCAPE" "HOME" "LEFT" "RIGHT" "SHIFT" "UP",
        "C_BREAK" "CAPSLOCK" "PAGEUP", "PAGEDOWN" "PAUSE"

機能名として以下が使えます。

        "BACKWARD_DELETE_CHAR" "BACKWARD_CHAR" "CLEAR_SCREEN" "DELETE_CHAR"
        "DELETE_OR_ABORT" "ACCEPT_LINE" "KILL_LINE" "UNIX_LINE_DISCARD"
        "FORWARD_CHAR" "BEGINNING_OF_LINE" "PASS" "YANK" "KILL_WHOLE_LINE"
        "END_OF_LINE" "COMPLETE" "PREVIOUS_HISTORY" "NEXT_HISTORY" "INTR"
        "ISEARCH_BACKWARD"

成功すると true を、失敗すると nil とエラーメッセージを返します。
大文字・小文字は区別せず、\_ のかわりに - を使うことができます。

### `nyagos.bindkey("キー名",function(this) ... end)`

キーが押下された時、関数を呼び出します。引数 this は次のような
メンバーを持ったテーブルです。

* `this.pos` … バイト数で数えたカーソル位置(先頭は 1 になります)
* `this.text` … utf8 で表現された現在の入力テキスト
* `this:call("FUNCNAME")` ... `this.call("BACKWARD_DELETE_CHAR")` のように機能を呼び出す
* `this:insert("TEXT")` ... TEXT をカーソル位置に挿入します
* `this:firstword()` ... コマンドラインの先頭の単語(コマンド名)を返します
* `this:lastword()` ... コマンドラインの最後の単語とその位置を返します
* `this:boxprint({...})` ... テーブルの要素を補完候補リスト風に表示します

また、戻り値は次のように使われます。

* 文字列の時: カーソル位置に挿入されます。
* true の時: Enter が押下されたのと同様に入力を終結します
* false の時: Ctrl-C が押下されたのと同様に内容を破棄して入力を終結します。
* nil の時: 無視されます。

### `nyagos.filter`

通常ユーザが呼び出すことはありません。
当関数を定義すると、ユーザが入力したコマンドラインの内容を引数として
NYAGOS.EXE から呼び出されます。これを加工して戻り値とすると、
NYAGOS.EXE はコマンドラインを、その文字列と置き換えます。

標準の nyagos.lua では nyagos.filter には、逆クォート機能を実現する関数が
定義されています。処理内容としては nyagos.eval でコマンドの出力を取り込み、
nyagos.atou で UTF8 に変換して、NYAGOS.EXE に返しています。

### `nyagos.argsfilter`

nyagos.argsfilter は nyagos.filter と似ていますが、コマンドライン
を字句解析した後の、引数配列(args)を加工できる点が違います。

標準の nyagos.lua では nyagos.argsfilter を使って、
suffix というコマンドを作成しています。

    コマンド
        suffix 拡張子 インタプリタ名 引数1 引数2 …
    Lua:関数
        suffix("拡張子",{"インタプリタ名","引数1"…})

これはコマンドに特定の拡張子がついた時に、インタプリタ名を
先頭に挿入するものです。

### `length = nyagos.prompt(template)`

通常ユーザが直接呼び出すことはありません。
引数のプロンプトのテンプレート(=%PROMPT%)を展開して、プロンプト文字列を
生成して表示、文字の桁数を戻り値を返す関数が格納されています。
ユーザはこれを横取りして独自のプロンプト表示を改造することができます。

    local prompt_ = nyagos.prompt
    nyagos.prompt = function(template)
        nyagos.echo("xxxxx")
        return prompt_(template)
    end

### `nyagos.gethistory(N)`

N 番目のヒストリ内容を返します。N が負の時は現在から(-N)個過去の
ヒストリを返します。引数が無い場合は、ヒストリの総数を返します。

### `nyagos.access(PATH,MODE)`

PATH で示されるファイルがアクセス可能かどうかを boolean 値で返します。
C言語の access 関数と同じです。

### `nyagos.completion_hook(c)`

補完のフックです。関数を代入してください。
引数 c は下記のような要素を持つテーブルです。

    c.list[1] .. c.list[#c.list] - コマンド名・ファイル名の補完候補
    c.word - 補完元の単語(二重引用符を含まない)
    c.rawword - 補完元の単語(二重引用符を含む場合がある)
    c.pos - 補完元の単語の始まる位置(0起点)
    c.text - コマンドラインの全文字列

`nyagos.completion_hook` は更新した候補リストのテーブルか nil を
戻り値としてください。nil は、更新しない c.list と等価です。

### `nyagos.getkey()`

入力されたキーの、Unicode、スキャンコード、シフト状態を返します。

### `nyagos.exe`

nyagos.exe のフルパスが格納されています。

## その他

NYAGOS は https://github.com/zetamatta/nyagos にて公開しています。
ソースは修正BSDライセンスにて配布・改変が可能です。

NYAGOS のビルドには

* [go1.4.2 windows/386](http://golang.org)
* [LuaBinaries(5.3 for Win32)](http://sourceforge.net/projects/luabinaries/files/5.3/Tools%20Executables/lua-5.3_Win32_bin.zip)

が必要となります。言語標準以外では、以下のモジュールを
利用させていただいております。

- http://github.com/mattn/go-runewidth
- http://github.com/shiena/ansicolor
- http://github.com/atotto/clipboard

以上
