* `echo` prints quotations as same as cmd.exe.

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
