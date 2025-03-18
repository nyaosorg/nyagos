[English](./04-Commands_en.md) / Japanese

## 内蔵コマンド

これらのコマンドはコマンド名とは別にエイリアスを持っています。
たとえば `ls` は `__ls__` というエイリアスを持っています。

#### `bindkey キー名 機能名`

一行入力のキー操作をカスタマイズします。

キー名

        "C_A" "C_B" ... "C_Z" "M_A" "M_B" ... "M_Z"
        "F1" "F2" ..."F24"
        "BACKSPACE" "CTRL" "DEL" "DOWN" "END"
        "ENTER" "ESCAPE" "HOME" "LEFT" "RIGHT" "SHIFT" "UP"
        "C_BREAK" "CAPSLOCK" "PAGEUP", "PAGEDOWN" "PAUSE"

機能名

        "BACKWARD_DELETE_CHAR" "BACKWARD_CHAR" "CLEAR_SCREEN" "DELETE_CHAR"
        "DELETE_OR_ABORT" "ACCEPT_LINE" "KILL_LINE" "UNIX_LINE_DISCARD"
        "FORWARD_CHAR" "BEGINNING_OF_LINE" "PASS" "YANK" "KILL_WHOLE_LINE"
        "END_OF_LINE" "COMPLETE" "PREVIOUS_HISTORY" "NEXT_HISTORY" "INTR"
        "ISEARCH_BACKWARD" "REPAINT_ON_NEWLINE"

#### `cd ドライブ:ディレクトリ`

現在のカレントドライブ、ディレクトリを変更します。
引数を省略すると、CMD.EXE と違い、環境変数 HOME 、あるいは 
USERPROFILE の差す先のディレクトリへ移動します。
CMD.EXE と違い、ドライブも同時に変更します。

* `cd -` : 一つ前にいたディレクトリへ移動します
* `cd -N` : N 回前のディレクトリへ移動します
* `cd -h` , `cd ?` : 過去いたディレクトリを表示します
* `cd --history` : 過去いたディレクトリを全て装飾なしで表示します
* `cd shortcut.lnk` : ショートカットの差すディレクトリへ移動します

ディレクトリ名が : / \ といった文字を含まず、カレントディレクトリに
存在しない場合、環境変数 CDPATH にリストされたディレクトリの
サブディレクトリを検索します。

#### `chmod ooo FILE(s)`

#### `env ENVVAR1=VAL1 ENVVAR2=VAL2 ... COMMAND ARG(s)`

COMMAND が実行されている間だけ、環境変数の値を変更します。

#### `more`

UTF8 と ANSI テキストの双方をサポートします。(自動判別)

#### `exit`

NYAGOS を終了します。

#### foreach

`foreach` *VAR* *VAL1* *VAL2* ...
    STATEMENTS
`end`

#### `history [件数]`

ヒストリ内容を表示します。件数を省略すると、最近の10件が表示されます。

#### if

##### inline-if

`if` *COND* *THEN-STATEMENT*

##### block-if

`if` *COND* [`then`]
   *THEN-BLOCK*
`else`
   *ELSE-BLOCK*
`end`

* `endif` は `end` の別名として使用可能です(nyaos-3000 との互換性のため)
* `then` は省略可能です

*COND* is:

* `not` *COND*
* `/i` *COND*
* *LEFT* `==` *RIGHT*
* `EXIST` *filename*
* `ERRORLEVEL` *n*

* if *COND* is true, execute *THEN-BLOCK* or *THEN-STATEMENT*
* if *COND* is false, execute *ELSE-BLOCK* or nothing.

#### `kill PID`

PID で示されるプロセスを強制終了します

#### `killall NAME...`

キーワードを含むプロセスを強制終了します

#### `ln [-s] SRC DST`

ハードリンク、もしくは、シンボリックリンクを作成します。
`nyagos.d\lns.lua` で定義されるエイリアス lns は UAC 昇格と
`ln -s` を実行します。

#### `lnk FILENAME SHORTCUT [WORKING-DIRECTORY]`

ショートカットを作成します

#### `ls [-オプション] …`

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
* `-h` -l 使用時に、人間が読みやすい形式でサイズを表記します (例:1K 234M 2G)
* `-S` ファイルサイズでソートします。
* `-?` ヘルプを表示します。
* `-L` リンク自体ではなく、リンクの参照先の情報を表示する

ロングフォーマットの例 (`ls -al`)

```
drwx-sh    0 Feb 19 20:16:53 System Volume Information/
drwxa--    0 Sep  3 2016     TDM-GCC-64/
```

フラグの意味

- `d` - ディレクトリ
- `r` - 読み込み可
- `w` - 書き込み可 (not read-only file)
- `x` - 実行可能
- `a` - アーカイブ済み
- `s` - システムファイル
- `h` - 隠しファイル

#### `ps`

プロセスのリストを表示します。

#### `pwd`

現在のカレントドライブ + ディレクトリを表示します。

* `pwd -N` : N 回 cd で移動する前のディレクトリを表示します。
* `pwd -L` : 環境から PWD を得る (default)
* `pwd -P` : 全てのシンボリックリンクをたどる

#### `set 変数名=値`

環境変数に値を設定します。値に空白等を含む場合、CMD.EXE と同様に
「`set "変数名=値"`」とします。= 以降を省略すると、現在の変数の内容を
表示します。

以下の変数は特別な意味を持ちます。

* `PROMPT` … プロンプトの文字列を設定します。`$P` 等のマクロ文字はCMD.EXE と同じです。shiena 様開発のモジュールによりエスケープシーケンスが使えます。
* `set ENV^=値` ... `set ENV=値;%ENV%` と等価ですが、重複した値は削除します
* `set ENV+=値` ... `set ENV=%ENV%;値` と等価ですが、重複した値は削除します

#### `set -o OPTION-NAME`, `set +o OPTION-NAME`

`-o` は OPTION を設定し、`+o` は解除します。

- `-o glob` 外部コマンドに対するワイルドカード展開を有効にします。
- `-o noclobber` リダイレクトによる既存ファイルの上書きを禁止します。
- `-o usesource` バッチファイルで NYAGOS の環境変数が変更できるようになります
- `+o usesource` バッチファイルから環境変数の変更を読みとるには source コマンドを使う必要があります。
- `-o cleaup_buffer` 一行入力の前に入力バッファをクリアします。

#### `set -a "EQUATION"`, `set /a "EQUATION"`

CMD.EXE と同じ。式を評価する

#### `select FILENAME(s)...`

アプリケーションを選択するダイアログ付きでファイルを開きます

#### `touch [-t [CC[YY]MMDDhhmm[.ss]]] [-r 参照ファイル] ファイル名…`

ファイルが存在すれば更新日時を更新し、存在しなければ新規作成します。

#### `which [-a] COMMAND-NAME`

コマンド名に対して、どのファイルが実行されるか表示します

* `-a` - %PATH% 上の全ての実行ファイルを表示します。

#### `copy SOURCE-FILENAME DESTINATE-FILENAME`
#### `copy SOURCE-FILENAME(S)... DESINATE-DIRECTORY`
#### `copy SOURCE-FILENAME(S)... SHORTCUT(*.lnk)`
#### `move OLD-FILENAME NEW-FILENAME`
#### `move SOURCE-FILENAME(S)... DESITINATE-DIRECTORY`
#### `move SOURCE-FILENAME(S)... SHORTCUT(*.lnk)`
#### `del FILE(S)...`
#### `erase FILE(S)...`
#### `mkdir [/p] NEWDIR(S)...`
#### `rmdir [/s] DIR(S)...`
#### `pushd`
#### `popd`
#### `dirs`
#### `diskfree`
#### `diskused`

これらの内蔵版は、上書きや削除の際に常にプロンプトで実行可否を問い合わせます。

#### `source バッチファイル名`

バッチファイルを CMD.EXE で実行して、CMD.EXE が変更した環境変数と
カレントディレクトリを NYAGOS.EXE に取り込みます。

- コマンド名として「`source`」の代わりに「`.`」(ドット)一文字も使うことができます
- `source` は一時ファイル `%TEMP%\nyagos-(PID).tmp` を作成します。
    - 更新されたカレントディレクトリと環境変数の内容が書き出されます
- `-d` オプションで、`source` が作成する一時ファイルが削除されなくなります
- `-v` オプションで、各一時ファイルが標準エラー出力に出力されます

#### `open FILE(s)`

Windows で関連付けられたアプリケーションでファイルを開きます。

#### `clone`

NYAGOS を別ウインドウで開きます。

#### `su`

UAC 昇格された NYAGOS を別ウインドウで開きます。

#### `su COMMAND ARGS(s)...`

UAC 昇格させて、コマンドを実行します。

### Lua で実装されたコマンド

#### `lua_e "INLINE-LUA-COMMANDS"` (nyagos.d\aliases.lua) 

内蔵Lua で引数の Lua コードを実行します。

#### `lua_f "LUA-SCRIPT-FILENAME" ARG(s)...` (nyagos.d\aliases.lua)

内蔵Lua で Lua スクリプトを実行します。

#### `trash FILE(S)` (nyagos.d\trash.lua)

ファイルを Windows のゴミ箱に移動させます。

#### `wildcard COMMAND ARG(s)...` (nyagos.d\aliases.lua)

ARG(s) に含まれるワイルドカードを展開して、COMMAND を実行します。

<!-- set:fenc=utf8: -->
