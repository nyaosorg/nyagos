[top](../README_ja.md) &gt; [English](./06-Substitution_en.md) / Japanese

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
* `@` 実行された時のディレクトリ

#### 変数

* `nyagos.histchar`: 置換用のマーク (デフォルト:「!」)
* `nyagos.antihistquot`: 抑制用の引用符 (デフォルト:「'"」)

### 環境変数置換

* コマンドや引数先頭の `~` を `%HOME%` あるいは `%USERPROFILE%` に置換します。

### Unicode リテラル

* `%u+XXXX%` (XXXX:16進数) を Unicode 文字に置換します。

### コマンド出力置換 (nyagos.d\backquote.lua)

    `COMMAND`
  もしくは
    $(COMMAND)

を、COMMAND の標準出力の内容に置換します。

### ブレース展開 (nyagos.d\brace.lua)

    echo a{b,c,d}e

を以下のように置換します。

    echo abe ace ade

### インタプリタ名の追加 (nyagos.d\suffix.lua)

- `FOO.pl  ...` は `perl   FOO.pl ...` に置換されます。
- `FOO.py  ...` は `ipy FOO.pl`、`py FOO.py`、`python FOO.py ...` のいずれかに置換されます。(最初に見付かったインタプリタ名が挿入されます)
- `FOO.rb  ...` は `ruby   FOO.rb ...` に置換されます。
- `FOO.lua ...` は `lua    FOO.lua ...` に置換されます。
- `FOO.awk ...` は `awk -f FOO.awk ...` に置換されます。
- `FOO.js  ...` は `cscript //nologo FOO.js ...` に置換されます。
- `FOO.vbs ...` は `cscript //nologo FOO.vbs ...` に置換されます。
- `FOO.ps1 ...` は `powershell -file FOO.ps1 ...` に置換されます。

拡張子とインタプレタの関連付けを追加したい時は、
`%USERPROFILE%\.nyagos` に

    suffix.拡張子 = "INTERPRETERNAME"
    suffix.拡張子 = {"INTERPRETERNAME","OPTION" ... }
    suffix[".拡張子"] = "INTERPRETERNAME"
    suffix[".拡張子"] = {"INTERPRETERNAME","OPTION" ... }
    suffix(".拡張子","INTERPRETERNAME")
    suffix(".拡張子",{"INTERPRETERNAME","OPTION" ... })

という記述を追加します。

<!-- set:fenc=utf8: -->
