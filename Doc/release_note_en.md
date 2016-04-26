* Add: `nyagos.open(PATH,MODE)` which `PATH` is in UTF8 and compatible with `io.open`.
* Add: `nyagos.loadfile(PATH)` which `PATH` is in UTF8 and compatible with `loadfile`.
* Add: `nyagos.lines(PATH)` which `PATH` is in UTF8 and compatible with `io.lines`. (Be careful that it returns bytearray-string not always UTF8!)
* Built-in `echo` uses CRLF not LF as the end of line.(#124)
* Lua's default file handles follow nyagos's redirect and pipeline
* Implemented touch's -r and -t option
* `touch` do tiny validation for timestamp format.
* `make install` makes log and closes installing window after 3sec(#107)
* `nyagos < TEXTFILE` is available.(#125)
* lua.exe & findstr.exe is no longer needed to make {conio,dos}/const.go.
* Fix: alias `suffix` dit not work.
* When the current working drive is a network drive, `su` starts new nyagos.exe as administrator on the same directory with UNC-Path.
* On `nyagos -c "CMD"`, CMD is executed after executing `nyagos.lua`.
* `nyagos -[cfe] "..."` & `nyagos < TEXTFILE` do not display copyrights.

NYAGOS 4.1.2\_0
===============
* Made scripts-catalog system
    - Moved `catalog.d\*.lua` to `nyagos.d\catalog\.`
    - We can import cataloged functions with `use "NAME"` in `.nyagos`
        - `use "dollar"` -> Expand the environment variable like `$PATH`
        - `use "peco"` -> powered by [peco](https://github.com/peco/peco)
            * C-r: History
            * C-o: Filename completion
            * M-h: Directory history
            * M-g: Git-revision
* ls 
    - not stop listing even if broken symbolic file exists.
    - Support: `ls -d`
* .nyagos can be put on the same directory with nyagos.exe
* Add: `cd --history`: print all the directory stayed with no decorations.
* Implemented built-in command tiny `touch`
* Fix: `>> bar` fails when `bar` does not exist(#121)
* Add the field `rawargs` to lua-command's parameter table,
  which contains parameters not removed quotations from user-typed ones.
* Add the method `replacefrom` to bindkey-function's parameter table.


NYAGOS 4.1.1\_2
===============
* Fix the miss to convert filename UTF8 to ANSI for loadfile() of Lua (#110,Thx Mr.HABATA)

NYAGOS 4.1.1\_1
===============
* Force to insert a line feed when prompt is too wide (#104)
* Fix the error message when no files matches with a given wildcard (#108)
* Fix the environment variable like %ProgramFiles(x86)% were not expanded. (#109 Thx @hattya)

NYAGOS 4.1.1\_0
===============

* Support UTF-16 surrogate pair on getkey
* `mkdir` suppports /p: make parent directories as needed.

NYAGOS 4.1.0\_0
===============

* Add build-in `ln`.
* Add lua-command `lns` which shows UAC and do `ln -s`
* `ls -l` shows the destination of the symbolic-link.
* Query continue or not when copy/move failed on one of parameters.
* New variable: `nyagos.histchar`: a character for history-substitution (default: `!`)
    - To disable history-substitution, do `nyagos.histchar = nil`
* New variable: `nyagos.antihistquot`: characters to disable for history-substitution (default: `'"`)
    - Be careful that `"!!"` is not substituted by DEFAULT.
    - To be compatible with 4.0, do `nyagos.antihistquot = [[']]`
* New variable: `nyagos.quotation`: characters for the completion (default: `"'`).
    - The first character of `nyagos.quotation` is the default quotation-mark.
    - The others characters are used when an user typed before completion.
    - When `nyagos.quotation=[["']]`
        - `C:\Prog[TAB]` is completed to `"C:\Program Files\ `  (`"` inserted)
        - `'C:\Prog[TAB]` is completed to `'C:\Program Files\ ` (`'` remains)
        - `"C:\Prog[TAB]` is completed to `"C:\Program Files\ ` (`"` remains)

NYAGOS 4.1-beta
================
* All Lua-callback function have thier own Lua-instances to avoid crash.
* Create the Lua-table `share[]` to share values between callback 
  functions and `.nyagos`.
* `*.wsf` is associated with cscript
* Warn on illeagal assign to nyagos[]
