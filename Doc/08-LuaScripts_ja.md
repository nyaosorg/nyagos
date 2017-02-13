[English](./08-LuaScripts_en.md) / Japanese

コマンド(Lua でエイリアス定義)
===============================

`lua_e "INLINE-LUA-COMMANDS"` 
----------------------------
(nyagos.d\aliased.lua にて定義) 

内蔵Lua で引数の Lua コードを実行します。

`cd SHORTCUT.LNK`
-----------------
(nyagos.d\cdlnk.lua にて定義)

cd のパラメータにショートカット(`*.lnk`)を指定できるようになっています

`open FILE(s)`
--------------
(nyagos.d\open.lua にて定義)

Windows で関連付けられたアプリケーションでファイルを開きます。

`su`
----
(nyagos.d\su.lua にて定義)

UAC 昇格された NYAGOS を別ウインドウで開きます。

`clone`
-------
(nyagos.d\su.lua に定義)

NYAGOS を別ウインドウで開きます。

`sudo COMMAND`
--------------
(nyagos.d\su.lua にて定義)

UAC 昇格させて、コマンドを実行します。

`trash FILE(S)`
---------------
(nyagos.d\trash.lua にて定義)

ファイルを Windows のゴミ箱に移動させます。

置換など
========

コマンド出力置換
----------------
(nyagos.d\backquote.lua にて定義)

    `COMMAND`

を、COMMAND の標準出力の内容に置換します。

ブレース展開
------------
(nyagos.d\brace.lua にて定義)

    echo a{b,c,d}e

を以下のように置換します。

    echo abe ace ade

インタプリタ名の追加
--------------------
(nyagos.d\suffix.lua にて定義)

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
