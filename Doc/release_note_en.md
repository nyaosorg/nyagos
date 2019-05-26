English / [Japanese](release_note_ja.md)

* (#233) Completion for `\\server\folder`
* (#238) copy: drawing progress
* Support `ENVVARNAME=VALUE COMMAND PARAMETERS..`
* Fix the problem temporary filenames on executing batchfile conflict
* (#277) Support set /A
* (#291) Print Process-ID when it is started with `&` mark like /bin/sh 
* (#361) Fix the problem GUI Application's STDOUT could not be redirected
* Implemented `.` and `source` for Linux (using /bin/sh)
* readline: fixed cursor-blink-off did not work while user wait.
* Fix: `mklink /J mountPt target-relative-path` made a broken junction.(Add change path absolute)
* Add options: `--chdir "DIR"` and `--netuse "X:=\\server\path"`
* Stop using CMD.EXE on `su` to set console-windows' icon nyagos'
* Completion: ask completion when more than 100 possibilities exist.
* ps: print `[self]` the line where nyagos.exe self exists.
* (#272) replace `!(historyNo)@` to the directory when the command was executed.
* (#130) Support Here Document
* ALT-O expands the path of shortcut(for example: SHORTCUT.lnk) to target-path
* (#368) Fix: Lua function: io.close() did not exist.
* (#332)(#369) Implement r+/w+/a+ mode for io.open

NYAGOS 4.4.3\_0
===============
on Apr.27,2019

* (#116) readline: implement UNDO on Ctrl-Z,Ctrl-`_`
* (#194) Update Console Window's ICON (LEFT-TOP-CORNER)
* Add alias for `date` and `time` in CMD.EXE
* Fix: the current-dir per each drive was mistaken after `cd RELATIVE-PATH`  
  ( `cd C:\x\y\z ; cd .. ; cd \\localhost\c$ ; c: ; pwd` -> `C:\x` (not `C:\x\y`) )

NYAGOS 4.4.2\_2
===============
on Apr.13,2019

* Implement Ctrl-RIGHT,ALT-F(forward-word) and Ctrl-LEFT,ALT-B(backward-word)
* Fix: wrong the count of backspaces to move top on starting incremental-search
* (#364) Fix: `ESC[0A` was used.

NYAGOS 4.4.2\_1
===============
on Apr.05,2019

* diskfree: trim spaces from the end of line
* Fix: on`~"\Program Files"`, the first quotation disappeared and `Files` was not contained in the argument.

NYAGOS 4.4.2\_0
===============
on Apr.02,2019

* Fix converting OLE-Object to Lua-Object causes panic on `VT_DATE` and some types.
* Fix: lua.LNumber was treated as integer. It should be as float64
* Lua: add function: `nyagos.to_ole_integer(n)` for `nyagos.d/trash.lua`
* Lua: support `for p in OLEObject:_iter() do ... end`
* Lua: add function: `OLEObject:_release()`
* Fix: trash.lua COM leak
* Fix: IUnknown instance created by `create_object` was not released.
* Implemented: expanding ~username
* Fix: exit status of executables (not batchfile) was not printed
* Fix: aliases using CMD.EXE (ren,mklink,dir...) did not work when %COMSPEC% is not defined.
* Fix: %U+3000% was regarded as a charactor of parameter separators
* (#359) -c and -k option can received multi arguments like CMD.EXE
* Fix: (not exist dir)\something [TAB] -> The system cannot find the path specified.(Thx! [tsuyoshicho](https://github.com/tsuyoshicho))
* (#360) Draw zero-width or surrogate paired characters as `<NNNNN>` (Thx! [tsuyoshicho](https://github.com/tsuyoshicho))
* Add the option --output-surrogate-pair to output them as it is (not `<NNNNN>`)
* su: network drives is not lost now after UNC-dialog
* (#197) `ln` makes Junction when the source-path is directory and -s is not given)
* Implemented built-in `mklink` command and remove the alias `mklink` as `CMD.exe /c mklink`
* Remove zero-bytes Lua files (cdlnk.lua, open.lua, su.lua, swapstdfunc.lua )
* (#262) `diskfree` shows volume label and filesystem
* Enabled to execute batch file even if UNC path is current directory.
* Fix rename,assoc,dir & for did not run when the current directory is UNC-path
* Fix (#363) Fix backquote did not work in nyagos.alias.COMMAND="string" (Thx! [tostos5963](https://github.com/tostos5963) & [sambatriste](https://github.com/sambatriste) )
* (#259) Implemented `select` command to open a file with dialog to select application.
* Fix the format of `diskfree`'s output

NYAGOS 4.4.1\_1
===============
on Feb.15,2019

* Made `print(nyagos.complete_for["COMMAND"])` work
* Fix (#356) `type` could output the last line which does not contain LF. (Thx! @spiegel-im-spiegel)
    * [zetamatta/go-texts](https://github.com/zetamatta/go-texts) v1.0.1 or laster is required
* Use `Go Modules` to build.
* Support completion for `killall` and `taskkill`.
* `kill` & `killall`: Forbide killing self process
* (#261) Set timeout(10sec) for completion and ls(1-folder)
* Fix: lua: ole object's setter(`__newindex`) did not work.
* (#357) Fix: on a french keyboard, AltGr + anykey did not work (Thx! @crile)
* (#358) Fix: When `foo.exe` and `foo.cmd` exist, typing `foo` calls `foo.cmd` rather than `foo.exe`

NYAGOS 4.4.1\_0
===============
on Feb.02,2019

* Support completion for `which`,`set`,`cd`,`pushd`,`rmdir` and `env` command. (Thx! [ChiyosukeF](https://twitter.com/ChiyosukeF))
* Fix (#353) Stopping OpenSSH with Ctrl-C on password prompt, Escape sequences and etc. are disabled. (Restore console mode for stdout after executing command) (Thx! [beepcap](https://twitter.com/beepcap))
* (#350) Stop calling os.Readlink on `ls -F` without `-l`
* Support `nyagos.complete_for["COMMANDNAME"] = function(args) ... end`
* Fix (#345) don't work git/svn/hg in subcomplete.lua (Thx! @tsuyoshicho)
* Fix io.popen(lua-function) did not work when redirect was used. (Thx! @tsuyoshicho)
* Fix (#354) box.lua: history completion did not start with C-X h (Thx! @fushihara)
* nyagos.d/catalog/subcomplete.lua supports completion for `hub` command. (Thx! @tsuyoshicho)

NYAGOS 4.4.0\_1
===============
on Jan.19,2019

* Abolished "--go-colorable" and "--enable-virtual-terminal-processing"
* Implemented `killall`
* Implemented `copy` and `move` for Linux
* (#351) Fix that `END` (and `F11`) key did not work 

NYAGOS 4.4.0\_0
===============
on Jan.12,2019

* To call a batchfile, stop to use `/V:ON` for CMD.EXE

NYAGOS 4.4.0\_beta
==================
on Jan.02,2019

* Support Linux (experimental)
* Fix the problem that current directories per drive were not inherited to child processes.
* Use the library "mattn/go-tty" instead of "zetamatta/go-getch"
* Stop using msvcrt.dll via "syscall" directly
* On linux, the filename NUL equals /dev/null
* Add lua-variable nyagos.goos
* (#341) Fix an unexpected space is inserted after wide characters
    * On Windows10, enable stdout virtual terminal processing always
    * If `git.exe push` disable virtual terminal processing, enable again.
* (#339) Fix that wildcard pattern `.??*` matches `..`
    * It requires github.com/zetamatta/go-findfile tagged 20181223-2

NYAGOS 4.3.3\_5
===============
on Dec.24,2018

* (#345) Fix subcomplete.lua don't work git (Thx! @tsuyoshicho)
* (#347) Fix the bug that STDOUT was closed after `dir 2>&1`.(Thx! @Matsuyanagi)
* (#348) Scrolling by mouse-wheel did not worked. (Thx! @tyochiai)
    * It requires github.com/zetamatta/go-getch tagged 20181223.

NYAGOS 4.3.3\_4
===============
on Dec.13,2018

* If stdin is not terminal, `more` command runs as `type`.
* On calling a batch file, `use CMD.EXE /V:ON /S /C "..."` for boosting code instead of temporary batchfile.
* (#340) Add lua variable `nyagos.histsize` to set the number of entries for history to save disk. (Thx! @crile)
* (#343) When %COMSPEC% is empty, use CMD.EXE (Thx! @orz--)

NYAGOS 4.3.3\_3
===============
on Oct.23,2018

* (#310) copy and move support shortcut files(`*.lnk`) as destination.
* (#313 reopened) Fix problem when `git blame FILES | type | gvim - &`, gvim starts with empty buffer.
* Fix: rmdir could not remove the broken junction
* Fix: Ctrl-C did not work in Lua-Script and some extern process
* (#267) `type` and `more` support UTF16 (requires go-texts package)
* (#336) Fix `io.write` did not work with -e and --lua-file
* (#337) Fix the crash the batchfile exit with -1 (Thx! @hogewest)

NYAGOS 4.3.3\_2
===============
on Sep.22,2018

* Append error message the filename on overwriting to existing file on redirect.
* Fix error for overwriting on redirect to `nul` when `noclobber` is set.
* diskused: continue counting how bytes disk used even if errors are found.
* ls: fixed `-l` option did not work with `-1` option
* ls: fixed: did not show one file per a line when output is not terminal.
* Not aliased builtin commands are able to be called as `\ls` like bash
* Fix the broken alias "for"
* Fix on completion the path separating characters were replaced to default one even if the word was not filepath for #334

NYAGOS 4.3.3\_1
===============
on Aug.29,2018

* #330,#331 Fix the original version of file:read incompatible behavior (Thx! @erw7)
* #332 stop buffering on io.open("w") (Thx! @spiegel-im-spiegel)
* #333 Fix file:seek() did not work on reading as expected (Thx! @erw7)
* #333 Fix file:close()'s return value was invalid. (Thx! @erw7)
* #319 Impl utf8.len()
* Fix: `which` reported files which has no suffixes
* `pwd` shows logical-path (=pwd -l) as default rather than phisical-path (=pwd -p)
* Fix: trash was left when incremental-search starts and some string exists on command-line.
* Shrink the executable with -lfdflags="-s -w"

NYAGOS 4.3.3\_0
===============
on Aug.14,2018

* #283 Omit the directory of path on completion by Ctrl-O
* #326 New option: `nyagos.option.tilde_expansion`
* Fix: `nyagos.option.xxxxxx = true` did not work
* Fix #328 `start https://...` fails (On CMD.EXE, it opens URL with Web Browser)
* Impl --read-stdin-as-file to read commands from stdin as a file for #327
* Fix: it sometimes failed to execute GUI application on symblic linked folder
* Fix: Commands with redirect (not pipeline) could not run on background
* Add lua-function: nyagos.fields(TEXT) which splits TEXT with spaces.
* #185 Add `ps` and `kill` command
* #329 Use `float64` instead of `int` for the number-type of Lua

NYAGOS 4.3.2\_0
===============
on Jul.23,2018

* #319 Support lua `bit32.*` all by github.com/BixData/gluabit32
* #323 Fix io.lines(), nyagos.lines() could not read from redirected stdin
* Fix: io.write() did not write to redirected stdout
* Replace `io.*` all with nyagos' own functions
* #324 Fix: Lua's print ignored --no-go-colorable (Thx @tignear)
* #325 Fix: `source` could not load the path which contains spaces.
* Add options: `--lua-first` and `--cmd-first`

NYAGOS 4.3.1\_3
===============
on Jun.19,2018

* #316 Fix: zero-length directory-name in %PATH% is regarded as the current directory
* #321 Fix: key function names `previous_history` & `next_history` were not registered.
* Add -h and --help option
* Lines starting with `@` of Lua script are now ignored to embed into batchfile.
* #322 Fix: change the encoding for batchfile's parameters from Thread Codepage to Console Codepage #322
* All of lua variables `nyagos.option.*` are now able to be set by nyagos.exe's command-line option.

NYAGOS 4.3.1\_2
===============
on Jun.12,2018

* #320: fix the imcompatibility: nyagos.rawexec & raweval did not expand tables in arguments.
* --show-version-only enables --norc automatically

NYAGOS 4.3.1\_1
===============
on Jun.11,2018

* Remove source code for lua53.dll
* #317: deadlock when `use "subcomplete"` is enabled and rclone.exe is found.
    - See also: https://github.com/yuin/gopher-lua/issues/181
* #318,#319: add compatible functions with lua 5.3
    - bit32.band/bitor/bxor
    - utf8.char/charpattern/codes

NYAGOS 4.3.1\_0
===============
on Jun.3,2018

* Support Windows10's native ESCAPE SEQUENCE processing with --no-go-colorable and --enable-virtual-terminal-processing
* For #304,#312, added options to search for the executable from the current directory
    * --look-curdir-first: do before %PATH% (compatible with CMD.EXE)
    * --look-curdir-last : do after %PATH% (compatible with PowerShell)
    * --look-curdir-never: never (compatible with UNIX Shells)
* nyagos.prompt can now be assigned string literal as prompt template directly.
* Fix #314 rmdir could not remove junctions.

NYAGOS 4.3.0\_4
===============
on May.12,2018

- Fix: #309 nyagos.getkey() raised panic (Thx @nocd5)
- Fix: error-message when command `lnk`'s target is not `*.lnk` nor exist.
- Fix: the cursor blink was switched to off on the child process.

NYAGOS 4.3.0\_3
===============
on May.9,2018

- Fix: forgot implement nyagos.setalias , nyagos.getalias (`alias { CMD=XXX}` did not work.)
- Fix: that the element [0] of the table value returned by alias-function was not used as the new command name to evaluate.
- Fix: `doc/09-Build_*.md` about how to download sourcefiles from github

NYAGOS 4.3.0\_2
===============
on May.7,2018

- #305: Fix issue that user's .nyagos was not loaded again (Thx! @erw7)

NYAGOS 4.3.0\_1
===============
on May.5,2018

- Fix: nyagos.d/start.lua did not worked because the member `rawargs` of alias-function's argument was not implemented.
- Fix: the return value of alias-function was not evaluted.
- Fix: for the script in -e option, arg[] was not assinged.
- Fix: On -f & -e option, warned as `getRegInt: could not find shell in Lua instanc
e`
- Fix: batchfile cound not return the value of `exit /b` as ERRORLEVEL

NYAGOS 4.3.0\_0
===============
on May.3,2018

- Add `ls -L` which shows information for the file refernces rather than for the link it self.

NYAGOS 4.3\_beta2
=================
on May.1,2018

- Fix: Typing C-o looks to raise hang up until Enter or ESCAPE is typed (on 4.3beta) #303
    - Fix the library: [go-box](https://github.com/zetamatta/go-box/commit/322b2318471f1ad3ce99a3531118b7095cdf3842)
- Fix: chcp did not work. (`chcp` was aliaes to update memory of screen width)

NYAGOS 4.3\_beta
==================
on Apr.30,2018

- Use Gopher-Lua instead of lua53.dll #300
    - nyagos.exe with lua53.dll can be built with `cd mains ; go build`
    - nyagos.exe with no Lua can be built with `cd ngs ; go build`
- Made `nyagos.option.cleanup_buffer` (default=false). When it is true, clean up console input buffer before readline.
- `set -o OPTION_NAME` and `set +o OPTION_NAME` (=`nyagos.option.OPTION_NAME=` on Lua)
- Buffer console-output ( go-colorable and bufio.Writer )

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
* Fix: #281 `cd \\server\folder ; open` -> `C:\Windows\system32` was open.
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

NYAGOS 4.1.9\_3
===============
on May 13,2017

* Fix #214: warned as `main/lua_cmd.go: cmdExec: not found interpreter object`

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
