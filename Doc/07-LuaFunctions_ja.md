## Lua拡張

nyagos では、EXE の本体の機能はコンパクトとし、便利機能は 
なるべく Lua で機能を拡張できるよう設計を進めています。
現在は以下のような関数が使用できます。

### `nyagos.alias.エイリアス名 = "置換コード"`

エイリアスを設定します。
置換コードでは、`$1` や `$*` などのマクロが使用可能です。
`$1`… には引用符が除去されていますが、`$*` には引用符が
残されています。

### `nyagos.alias.エイリアス名 = function(args)～end`

Lua 関数をエイリアスコマンドとして呼び出せるようにします。
args には全引数を格納したテーブルが入ります。

    {
        [1]=第一引数,
        [2]=第二引数,
        [3]=第三引数,
            :
        ["rawargs"]={
            [1]=第一引数(引用符を除去していない),
            [2]=第二引数(引用符を除去していない),
            [3]=第三引数(引用符を除去していない),
                :
        }
    }


エラーがあった時、関数は %ERRORLEVEL% に格納すべき「整数値」と
エラーメッセージの二値を返さなくてはいけません。
(return なしの場合は「return 0,nil」と同じです)

戻り値が文字列や、文字列テーブルの場合、その文字列(テーブル)が
新コマンドラインとして実行されます。

エイリアスは Lua の別のインスタンスで実行されるため、.nyagos で
定義された変数は、共有テーブルshare[] を除いて、参照できません。
share[] はユーザが自由に使用可能ですが、全てのインスタンスで、
ただちに同期されるのは share[] 直下のメンバーのみです。

### `nyagos.env.環境変数名`

環境変数にリンクしています。参照・変更が可能です。

### `errorlevel,errormessage = nyagos.exec("シェルコマンド")`

シェルコマンドを実行します。エラーが発生した時、
戻り値は %ERRORLEVEL% に格納すべき整数値とエラーメッセージが入ります。
エラーが無い時は (0,nil) が戻ります。

### `errorlevel,errormessage = nyagos.rawexec("外部コマンド名","引数1","引数2"…)`

外部コマンドを実行します。
戻り値は %ERRORLEVEL% に格納すべき整数値とエラーメッセージが入ります。
エラーが無い時は (0,nil) が戻ります。
(os.execute との違いは引数が UTF8 と解釈される点です)

### `nyagos.eval("シェルコマンド")`

nyagos.exec と同じですが、標準出力を取り込んで、戻り値として返します。
実行に失敗した場合などは nil が戻ります。

### `OUTPUT,ERR = nyagos.raweval("外部コマンド名","引数1","引数2"…)`

外部コマンドを実行して、標準出力の内容を戻り値として返します。
実行に失敗した場合は nil とエラーが戻ります。

### `nyagos.write(テキスト)`

テキストを標準出力に出力しますが、リダイレクトされている場合は
文字コードはUTF8 になります。内蔵 Lua の print は 
nyagos.write(テキスト..'\n') に差し替えられています。

### `nyagos.writerr(テキスト)`

テキストを標準エラー出力に出力しますが、リダイレクトされている場合は
文字コードはUTF8 になります。

### `nyagos.getwd()`

現在のカレントディレクトリを返します。

### `nyagos.chdir('DIRECTORY')`

カレントディレクトリを変更します。

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
* `this:replacefrom(POS,"TEXT")` ... POSからカーソルまでを TEXT と差替えます

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

    share.backup_prompt = nyagos.prompt
    nyagos.prompt = function(template)
        print("xxxxx")
        return share.backup_prompt(template)
    end

### `nyagos.gethistory(N)`

N 番目のヒストリ内容を返します。N が負の時は現在から(-N)個過去の
ヒストリを返します。引数が無い場合は、ヒストリの総数を返します。

### `nyagos.access(PATH,MODE)`

PATH で示されるファイルがアクセス可能かどうかを boolean 値で返します。
C言語の access 関数と同じです。

### `nyagos.completion_hook = function(c) ... end`

補完のフックです。関数を代入してください。
引数 c は下記のような要素を持つテーブルです。

    c.list[1] .. c.list[#c.list] - コマンド名・ファイル名の補完候補
    c.word - 補完元の単語(二重引用符を含まない)
    c.rawword - 補完元の単語(二重引用符を含む場合がある)
    c.pos - 補完元の単語の始まる位置(0起点)
    c.text - コマンドラインの全文字列

`nyagos.completion_hook` は更新した候補リストのテーブルか nil を
戻り値としてください。nil は、更新しない c.list と等価です。

### `nyagos.on_command_not_found = function(args) ... end`

定義されていると、コマンドが見付からなかった時に呼び出されます。
コマンド名とパラメータが args[0] ～ args[#args] にセットされます。
関数が nil か false を返した場合は nyagos.exe は通常のエラーを
表示します。

関数は別の Lua インスタンスで実行されるため、.nyagos で定義された変数への
アクセスはエイリアス同様の制限があります。

### `WIDTH,HEIGHT = nyagos.getviewwidth()`

ターミナルの横幅と高さを返します。

### `STAT = nyagos.stat(FILENAME)`

ファイルの情報を返します。
ファイルが存在する時、テーブル STAT は下記のようなメンバーを持ちます。

    STAT.name
    STAT.isdir (ディレクトリなら true, さもなければ false)
    STAT.size  (バイト数)
    STAT.mtime.year
    STAT.mtime.month
    STAT.mtime.day
    STAT.mtime.hour
    STAT.mtime.minute
    STAT.mtime.second

ファイルがない時、STAT は nil です。

### `nyagos.getkey()`

入力されたキーの、Unicode、スキャンコード、シフト状態を返します。

### `nyagos.open(PATH,MODE)`

PATH が utf8 と解釈される以外は io.open と等価です。

### `nyagos.loadfile(PATH)`

PATH が UTF8 と解釈される以外は、通常の loadfile と等価です。

### `nyagos.lines(PATH)`

PATH が UTF8 と解釈される以外は、通常の io.lines と等価です。

```
for text in nyagos.lines(PATH) do ... end
```

`text` は UTF8 変換などはなく、io.lines 同様、ただのバイト列です


### `nyagos.goversion`

ビルドに使用した Go のバージョン文字列が格納されます。
(例：「go1.6」)

### `nyagos.goarch`

実行ファイルが想定している CPU アーキテクチャを示す文字列が格納されます。
(例：「386」「amd64」)

### `nyagos.exe`

nyagos.exe のフルパスが格納されています。

<!-- set:fenc=utf8: -->
