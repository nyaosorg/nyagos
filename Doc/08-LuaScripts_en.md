Commands (alias written in Lua)
===============================

`lua_e "INLINE-LUA-COMMANDS"` 
----------------------------
(defined in nyagos.d\aliased.lua)

Execute inline-lua-commands like 'lua.exe -e "..."'.

`cd SHORTCUT.LNK`
-----------------
(defined in nyagos.d\cdlnk.lua)

`cd` cat change the current directory pointed with shortcut-file(`*.lnk`).

`open FILE(s)`
--------------
(defined in nyagos.d\open.lua)

Open the file with associated application.

`su`
----
(defined in nyagos.d\su.lua)

Run another nyagos.exe as Administrator.

`clone`
-------
(defined in nyagos.d\su.lua)

Run another nyagos.exe on another console window.

`sudo COMMAND`
--------------
(defined in nyagos.d\su.lua)

Run COMMAND in nyagos.d\su.lua)

`trash FILE(S)`
---------------
(defined in nyagos.d\trash.lua)

It throws files into trashbox of Windows.

Substitution
============

Command Substitution (nyagos.d\backquote.lua)
---------------------------------------------
(defined in nyagos.d\backquote.lua)

    `COMMAND`

is replaced to what COMMAND print to standard output.

Brace Expansion
---------------
(defined in nyagos.d\brace.lua)

    echo a{b,c,d}e

is replaced to

    echo abe ace ade

Inserting Interpreter-name 
--------------------------
(defined in nyagos.d\suffix.lua)

- `FOO.pl  ...` is replaced to `perl   FOO.pl ...`
- `FOO.py  ...` is replaced to `ipy FOO.py ...` , `py FOO.py` or `python FOO.py ...` (inserted the first found interpreter's name)
- `FOO.rb  ...` is replaced to `ruby   FOO.rb ...`
- `FOO.lua ...` is replaced to `lua    FOO.lua ...`
- `FOO.awk ...` is replaced to `awk -f FOO.awk ...`
- `FOO.js  ...` is replaced to `cscript //nologo FOO.js ...`
- `FOO.vbs ...` is replaced to `cscript //nologo FOO.vbs ...`
- `FOO.ps1 ...` is replaced to `powershell -file FOO.ps1 ...`

To append the new associtation between the suffix and interpreter,
write

    suffix.xxx = "INTERPRETERNAME"
    suffix.xxx = { "INTERPRETERNAME","OPTION" ... }
    suffix[".xxx] = "INTERPRETERNAME"
    suffix[".xxx] = { "INTERPRETERNAME","OPTION" ... }
    suffix(".xxx","INTERPRETERNAME")
    suffix(".xxx",{ "INTERPRETERNAME","OPTION" ... })

in `%USERPROFILE%\\.nyagos`
