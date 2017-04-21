English / [Japanese](release_note_ja.md)

* Add built-in command `env`
* #189 Support `nyagos.history[..]` and `#nyagos.history`
* Make `type` as built-in command.
* To print prompt, use the same Lua instance loading ~/.nyagos
* Fix the problem that C-o could not complete filenames which has ` ` and `~`.
* Make `clip` as built-in command which read/write both UTF8/MBCS (#202)
* Support `del /f`: delete even if it is a readonly file. (#198)

NYAGOS 4.1.9\_2
===============
on Apr 3,2017

* Fix #191 the option`-c` printed `option parse error`.
* A new Lua function: `nyagos.elevated()` which returns true on elevated mode.
* The default title bar prints `(admin)` on elevated mode.

NYAGOS 4.1.9\_1
===============
on Mar 28,2017

* Fix: sometimes cursor disappears on readline on some environment on 4.1.9\_0.
* Be able to use the escape sequence `\033]0;(title)\007` to change the title of the command-prompt by the new go-colorable's feature.

NYAGOS 4.1.9\_0
===============
on Mar 27,2017

* Fix: `open http(s)://...` did not work.
* Support `cd file:///...`
* ALT-y: if string on clipboard has space, paste it with double-quatations.
* Listing filenames to completion, cut dirname of the fullpath
* Fix: `history`-command did not display ID to use !-mark
* Command-name on %NYAGOSPATH% are completed with TAB.
* Not expand environment variables on filename-completion.
* Not enclose ~/ & ~\ with double-quotations.
* At completion, ignore the string before `;`,`=` (for set command)
* Speedup print working directory on Prompt by not to fix filename case.
* `cd C:\Program Files` works without double-quotations.(#182)
* `cd /D` works. ignore /D option for compatibility with CMD.EXE.(#182)
* Sort `history`'s output by time.
* Remove file existance check on `open` to `open regedit`
* `clone`,`su`,`sudo`: ShellExecute with the destinate paths of the symbolic links not to fail on network folders.(#122)
* set: be compatible with CMD.EXE's (set FOO=A B is same FOO="A B")
* Fix #184 Backquotation does not work in `_nyagos`
* `_nyagos`: support `bindkey KEYNAME FUNCNAME`
* Support %ENVNAME:FROM=TO% like CMD.EXE
* On incremental search, bind ESCAPE-KEY to quit search-mode.
* New completions by new built-in command `box`
    * Ctrl-O          : Insert filename to select by Cursor (box.lua)
    * Ctrl-XR , Alt-R : Insert history to select by Cursor (box.lua)
    * Ctrl-XG , Alt-G : Insert Git-revision to select by Cursor (box.lua)
    * Ctrl-XH , Alt-H : Insert `CD`ed directory to select by Cursor (box.lua)
* Support `lua_e "nyagos.key = function(this) end"`

NYAGOS 4.1.8\_0
===============
on Feb 15,2017

* Add new customizing file `_nyagos`(command.com-batchlike)
* Fix #173 could not stop `ls` and built-in commands with Ctrl-C
* ls -h: display size with COMMA not Kilo,Mega or Giga
* Support nyagos.lines(FILENAME,"n") but value is not float but int #147
* Add %NYAGOSPATH% which works like %PATH% only in nyagos.exe not childprocess
* Support SET VAR+=VALUE , VAR^=VALUE like vim
* Fix #176 Bug on `gawk "BEGIN{ print substr(""%01"",2) }"`
* Use github.com/josephspurrier/goversioninfo instead of windres.exe to attach icon
* Support `if` compatible with command.com's one (`==`,`not`,`errorlevel`,`/I`)
* New alias macro `$~1` `$~2` ... `$~*` which remove double quotations.
* Record current directories, times and process-id as history (#112)
* ls -l: change timestamp format to 'Jan 2 15:04:05' or 'Jan 2 2006'
* When lua53.dll is not found, display not a stacktrace but a readable error.
* '#' became a comment mark.
* open,clone,su,sudo : rewrite with Go (from Lua)

NYAGOS 4.1.7\_0
===============
on Nov 29,2016

* Abolished nyagos.lua, which role nyagos.exe do itself.
* Caching ~/.nyagos with `%APPDATA%\NYAOS_ORG/dotnyagos.luac`
* `nyagos.d/*` are bundled with nyagos.exe self.
* Fix #167 Could not call executable symbolic-linked to relative path
* Fix `ls -l` could not display `@` and linked path for symbolic-linked-executables
* Fix su.lua: clone/su displayed broken path.
* Fix #168 `ls RELATIVE-SYMLINKED-FILEPATH` occured error.
* Fix Widths for filesize in `ls -lh` were broken
* Set default alias ls="ls -oFh" (add -h) 
* `history` outputs history lines all when stdout is not a terminal.
* `open` prints a prompt for each files if more than one parameters are given.
* `use "cho"` -> powered by [cho](https://github.com/mattn/cho)
        * C-r: History
        * C-o: Filename completion
        * M-h: Directory history
        * M-g: Git-revision
* Fix: brace expansion "{a,b,c}" worked even in quotated strings

NYAGOS 4.1.6\_1
===============
on Sep 7,2016

* Fix: the package zip did not have lua53.dll

NYAGOS 4.1.6\_0
===============
on Sep 7,2016

* Use "\x1B[0K" as ERASELINE instead of " " & Backspace
* Use "\x1B[mC as m-times of Backspace
* Fix #159: Stop to print prompt again when terminal window resized
* Fix #164: `cd --history` changed the current directory to home.
* copy and move always regard the desitinate path matching with `[\\/:]\.{0,2}$` as a directory wheter it fails or not to stat the path.

NYAGOS 4.1.5\_1
===============
on Jul 31,2016

* Fix #157++: Overflow line on the text appended after screen resized.
* Error when it the upvalue named as 'prompter' is used on closures(nyagos.prompt) for invalid ~/.nyagos of 4.0.x on default for #155,#158

NYAGOS 4.1.5\_0
===============
on Jul 31,2016

* `cd --history` outputs the current directory at first to prevent peco(M-h) fro
m hangup with no cd histories.
* On lua, `nyagos.option.glob = true` enables the wildcard expansion on external commands also.(#150)
* Tried to improve the compatibility of `source`
* Support nyagos.lines(FILENAME,X) X='a','l','L',Number for #147
* Fix #156: %U+0000% causes panic
* Fix #152 ls -ld Downloads\ -> Downloads\/ printed.
* Fix #157 Reset the readline-width on the console window resized.
* Moved some packages to the other repositories.

NYAGOS 4.1.4\_1
===============
on Jun 12,2016

* Fix #151 `&&` and `||` work same as ` ;`
* Add nyagos.d/catalog/autocd.lua & autols.lua (#149 Thx @DeaR)

NYAGOS 4.1.4\_0
===============
on May 29,2016

* Implemented built-in tiny OLE interface and nyole.dll is not necessary now.
* Define default-prompt function as `nyagos.default_prompt` and it can change
console-title(second parameter)
* Fix: nyagos.lines() did not remove CRLF #144
* Fix: Lua's default file handles(STDIN/STDOUT) were opened by binary-mode. #146
* nyagos.d/catalog/peco.lua: C-r: revert order of display and improved speed.

NYAGOS 4.1.3\_1
===============
on May 8,2016

* Fix: %APPDATA%\nyaos.org\nyagos.history was not updated (#138)
* Fix: when nyagos.history was deleted, warnings are displayed until `exit` was typed.
* Fix: nyagos.d/catalog/peco.lua: when nyagos.history does not exist, peco hangs

NYAGOS 4.1.3\_0
===============
on May 5,2016

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
* Fix: `make install DIR` did not save DIR to Misc/version.cmd
* Fix: nyagos.exe could not load nyagos.lua when nyagos.exe exists on non-ascii-path (#133)
* Fix: nyagos.d/catalog/subcomplete.lua does not work after 4.1 (#135)
* Switch escape sequence emulater to github.com/mattn/go-colorable (#137)
* Fix: `ls -ltr *` was not sorted by modified time. (#136)
* Support: nyagos -f NOT-Lua-Script(COMMANDS-Script)

(Add forgotten change on May 17,2016)
-------------------------------------

* Not to confuse whether the encoding is ANSI or UTF8 string , stop to print('UTF8-String with ESCAPE-SEQUENCE'). Now print remains to be the bundled one of lua53.dll. ( #129 )

NYAGOS 4.1.2\_0
===============
on Mar 29,2016

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
on Feb 17,2016

* Fix the miss to convert filename UTF8 to ANSI for loadfile() of Lua (#110,Thx Mr.HABATA)

NYAGOS 4.1.1\_1
===============
on Feb 16,2016

* Force to insert a line feed when prompt is too wide (#104)
* Fix the error message when no files matches with a given wildcard (#108)
* Fix the environment variable like %ProgramFiles(x86)% were not expanded. (#109 Thx @hattya)

NYAGOS 4.1.1\_0
===============
on Jan 15,2016

* Support UTF-16 surrogate pair on getkey
* `mkdir` suppports /p: make parent directories as needed.

NYAGOS 4.1.0\_0
===============
on Jan 3,2016

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
on Dec 13,2015

* All Lua-callback function have thier own Lua-instances to avoid crash.
* Create the Lua-table `share[]` to share values between callback 
  functions and `.nyagos`.
* `*.wsf` is associated with cscript
* Warn on illeagal assign to nyagos[]
