English / [Japanese](history-4.2_ja.md)

NYAGOS 4.2.5\_1
===============
on Apr.14,2018

- Fix: `if [not] errorlevel N` did not work on block-if.
- Fix: that `ls -1F` did not show the indicator such as `/`,`*` or `@`.
- Fix: the problem that executables reparse-pointed but not symbolic-linked can not be found.
- Fix: `ls -F` marked '@' to files and directories which ar reparse-pointed but not symbolic-link nor junction
- Changed the error message when the command `history` is called in `_nyagos`

NYAGOS 4.2.5\_0
================
on Mar.31,2018

- Add Lua-flag: nyagos.option.usesource. When it is false, batchfiles can not change nyagos's environment variables and directory.(default:true)

NYAGOS 4.2.5\_beta2
===================
on Mar.27,2018

- Fix: #296 the batchfile could not be executed when the username contains multibyte-character.
    - Fix that the encoding of the temporary batchfile was UTF8.
    - Fix that the end of the each line of the temporary batchfile was LF not CRLF.
- Fix: #297 running the batchfile includes `exit` without `/b` option, an error occurs

NYAGOS 4.2.5\_beta
==================
on Mar.26,2018

- Read the value of environment variables that a batchfile changed like CMD.EXE.
- And refactored source files

NYAGOS 4.2.4\_0
===============
on Mar.9,2018

* lua: ole: `variable = OLE.property` is avaliable instead of `OLE:_get('property')`
* lua: ole: `OLE.property = value` is avaliable instead of `OLE:_set('property',value)`
* Load `nyagos.d/*.ny` as batchlike file
* #266: `lua_e "nyagos.option.noclobber = true"` forbides overwriting existing file by redirect.
* #269: `>| FILENAME` and `>! FILENAME` enable to overwrite the file already existing by redirect even if `nyagos.option.noclobber = true`
* #270: Console input buffer has been cleaned up when prompt is drawn.
* #228: Completion supports $ENV[TAB]... by native
* #275: Fix: history substitution like `!str:$` , `!?str?:$` did not work.
* The error `event not found` is caused when the event pointed !y does note exists.
* #285: Not wait GUI-process not using pipeline terminating like CMD.EXE (Call them with ShellExecute() instead of CreateProcess() )
* (Replaced `bytes.Buffer` to `strings.Builder` and Go 1.10 is required to build)
* When more than one are to be executed with `open` at once, display error: `open: ambiguous shellexecute`
* Fix that `nyagos.alias.NAME = nil` could not remove the alias.

NYAGOS 4.2.3\_4
===============
on Mar.4,2018

* `ls -?` for help instead of `ls -h`
* Building with `go build` instead of make.cmd, print version as `snapshot-GOARCH`
* Show an error when `type DIRECTORY` is executed.
* Made error message simple on `del NOTEXISTFILE`
* Fix: #279 Substitution on Environment variable (%VAR:OLD=NEW%) did not ignore case
* Fix: #281 `cd \\host-name\share-name ; open` -> `C:\Windows\system32` was open.
* Fix: #286 A tilde(~) after whitespace enclosed with double quotations was interpreted same as %USERPROFILE%
* #287 On the last entry of the history, do nothing for typing ARROW-DOWN

NYAGOS 4.2.3\_3
===============
on Jan.28,2018

* Fix: `print(nil,true,false)` outputs nothing..
* Fix the bug that `!notfoundstr` is replaced to `!n` only.
* #271: Fix Ctrl-O (box selector) does not work for the path contains %APPDATA% ( Fix zetamatta\go-findfile )
* On completion, don't append SPACE after PERCENT mark.
* #276 Fix that `source` did not execute a batch with stdout. (Thx @tyochiai )

NYAGOS 4.2.3\_2
===============
on Jan.6,2018

* Fix: #265 Type `ls` , SPACE and TAB -> command name completion starts. 

NYAGOS 4.2.3\_1
===============
on Dec.30,2017

* Fix: CR and LF did not work as the word seperator in the commandline.
* Fix: #264 Garbage appears on the screen when screen buffer width is too large.
    (You have to do `go get -u github.com/mattn/go-colorable`)

NYAGOS 4.2.3\_0
===============
on Dec.25,2017

* option --norc : not to load startup-scripts.
* #132 Support foreach and block-if
* Add option --lua-file which loads and runs Lua-Script even if the suffix of the filename is not .lua .
* Add members  to the parameter `c` of `nyagos.complete_hook(c)`
    * `c.field` : array split all commandline string with space.
    * `c.left` : string before cursor.
* Enable command-name completion even if it is after `|` , `&` , `;`
* #245 `print` of lua supports redirect.
* On incremental search, support Ctrl-S for backward search.
* Fix #258 Environment variable expanding does not work after backslash
* Add lua-function nyagos.msgbox(MESSAGE,TITLE)

NYAGOS 4.2.2\_2
===============
on Nov.26,2017

* #255 `start` command search the executable via %PATH%
* #254 Fix: -xxxx of `nyagos -f SCRIPT -xxxx` was treated as not SCRIPT's option but nyagos' option.
* Fix: Lua-stack overflow when arguments filter is not set

NYAGOS 4.2.2\_1
===============
on Oct.11,2017

* #250 Fix the crash in the built-in command `bindkey` without parameters. (Thx @masamitsu-murase)
* #252 Fix the problem that Shift/Ctrl keys typing cancels the screen-scroll. (Skip the output CURSOR-OFF/ON sequences at Shift/Ctrl keys typed) (Thx @masamitsu-murase)
* #253 Fix `nyagos-4.2.2_0-386` was built as a 64bit executable by make.cmd's bug (Thx @hazychill)

NYAGOS 4.2.2\_0
===============
on Oct.8,2017

* Append the new command commands by Lua: `abspath`,`chompf` and `wildcard`
* Append the forgotten builtin lua-commands reference: `lua_f`,`kill` and `killall`.
* #246 Reject conversion from userdata to Object. (Thx @masamitsu-murase)
    - To assign userdata(Lua) to `share[]` is forbidden
    - The global userdata(Lua) are not cloned on the forked Lua instance for the background goroutine to make pipelines.
* #247 Fixed the problem that Go's Garbage collector releases data refered by Lua and crashes (Thx @masamitsu-murase)
* #248 Completion hook can specify displayed-titles which differ from completed-strings.(Thx @masamitsu-murase)
* #249 Add `nyagos.completion_slash` option. When it is true, filename-completion uses a slash as the path-seperator as default. (Thx @masamitsu-murase)
* New building script(make.cmd) written in PowerShell

NYAGOS 4.2.1\_0
===============
on Aug.31,2017

* #241 Respect the item order in the list returned from `completion_hook` (Thx @masamitsu-murase)
* #242,#243 Support key combination for Alt+Backspace and Alt+"/". (Thx @masamitsu-murase)
* Remove built-in command `sudo`
* Add built-in command `more` (support color and unicode)
* readline: support C-q,C-v (`QUOTED_INSERT`)
* pwd: add options -L(use PWD from environment) and -P(avoid all symlinks)
* Output `nyagos.dump` if panic occurs.
* `diskused`: new command like du
* `rmdir` prints the progress as before.
* `diskfree`: new command like df

NYAGOS 4.2.0\_5
===============
on Aug.16,2017

* Fix: Building on Windows7, the version information was not written into the property of the executable because the script to make JSON for goversioninfo required the method ConvertTo-JSON of PowerShell 3.0 but Windows 7 does not support it.
* Fix: nyagos.box(LIST) ignored the order of LIST

NYAGOS 4.2.0\_4
===============
on Jul.29,2017

* Fix: error's line number was not displayed when `.nyagos` has an error.
* Fix: `.nyagos` cache errors when executable architecture (amd64 or 386) changes previous'run
* Fix: `ls | more` outputs `ESC[0K`. (fixed by go-box)
* (internal) follow the change in go-colorable's `ESC[%dC` & `ESC[%dD`'s behaviour ( https://github.com/mattn/go-colorable/commit/3fa8c76f , thanks to @tyochiai )
* Fix: on default `_nyagos`, `suffix "lua=nyagos"` was wrong. Added `.exe -f`
* Error if scripts on `nyagos.d` are executed by not nyagos.exe
* Do not insert interpreter-name when user-typed-command-name does not have a suffix to fix #237 that `cd nyagos.d` and `suffix` -> new nyagos.exe processes start infinitely.
* Fix #240: on empty dir, C-o -> `bad argument # 1 to 'find' (string expected, got nil)`

NYAGOS 4.2.0\_3
==================
on Jul 13,2017

* Fix: panic occurs when `box` Enter & Ctrl-C pressed.
* Fix: panic occurs when `lua_e "nyagos.box({})"`
* Fix: box: cursor disappear at scrolling (go-box's fix)
* `box`: decrease flickering (go-box's fix)
* Fix: #235 .nyagos on the same directory with nyaos.exe wasn't read on startup.
* completion: enclose with "" when ! mark is found.
* Fix: `suffix ps1` => `?:-1: attempt to concatenate a table value`

NYAGOS 4.2.0\_2
===============
on Jun 16,2017

* Fix the problem that `lnk . ~` failed.
* Fix the problem executables on the folder symbolic-linked to network one and to be elevated are unable to be called. (ShellExecute with physical path)
* Fix: readline: isearch: BACKSPACE-KEY did not redraw a found commandline
* Fix: crash that `index out of range` occurs when an empty string in the global variable in Lua exists and pipeline was used. (#232)

NYAGOS 4.2.0\_1
==================
on Jun 06,2017

* Fix: sample of `_nyagos` was forgotten to bundle into the package. (#230)
* Implemented `chmod`. (#199)
* nyagos.d/catalog/dollar.lua: support completion $TEMP\xxxx format. (#228)
* nyagos.d/catalog/ezoe.lua: revival

NYAGOS 4.2.0\_0
===============
on May 29,2017

* **Improved the restriction that the Lua-variables not in `share[]` are not shared in all Lua-instances.(#210,#208)**
    * Do not create a new lua instance except when background thread begins to run.
    * Copy global variables all including ones not in `share[]` from the Lua-instance in the forwaground thread to the new instance for the new background thread.
    * To print prompt, use the same Lua instance loading ~/.nyagos

New feature
-----------
* `nyagos.completion_hidden`: If it is set true, hidden filenames are also completed.
* Add built-in command `env`
* #189 Support `nyagos.history[..]` and `#nyagos.history`
* Make `type` as built-in command.
* Make `clip` as built-in command which read/write both UTF8/MBCS (#202)
* Support `del /f`: delete even if it is a readonly file. (#198)
* Add command to make shortcut(`lnk FILENAME SHORTCUT WORKING-DIRECTORY`)
* Add `attrib` as built-in command. (#199)
* Support `$(  )` format to quote command-output by backquote.lua
* `ls -l`: Show shortcut's target and working directory.
* New lua-function: `nyagos.box()`

Trivial fix
-----------
* Support `-b BASE64edCOMMANDSTRING` as startup option (#200)
* Rewrote `cd/push SHORTCUT.lnk` from Lua(nyagos.d/cdlnk.lua) to Golang-native
* `nyagos.alias.grep = "findstr.exe"`

Bugfix
------
* Fix `\` in `%USERPROFILE:\=/%` were replaced once only
* Fix: `ll` was aliased to non-colored version on default `_nyagos`
* Fix the problem that C-o could not complete filenames which has ` ` and `~`.
* Fix Ctrl-O (filename-completion) causes panic. (#204)
* Never cut double-quotations of parameters which users wrote for FIND.EXE and so on #218,#222
* Fix: Executing commands requiring elevation causes the error `The requested operation requires elevation`. Now UAC elevation dialog is shown. #227
* Fix: executed `FOO.123.EXE` even when `FOO` was typed #229
