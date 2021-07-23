[top](..) &gt; English / [Japanese](release_note_ja.md)

* Fix: the replacing result for %DATE% was not compatible with CMD.EXE's output on codepage 437
* go-readline-ny v0.4.13: Support Mathematical Bold Capital (U+1D400 - U+1D7FF) on the Windows Terminal
* Even if the option -o is given, remove ESC[0m from the output of redirected ls
* Fix: (#411) the document that English part and Japanese were written in inverse places (Thx! @tomato3713)
* Organize and automate the test-codes

NYAGOS 4.4.10\_1
================
on Jul.02,2021

* Fix: When the folder `./dll` existed and the folder `DLL` existed on CDPATH, typed path dll was replaced to `DLL` with completion (The problem is that the path's cases were changed)
* Fix: clone: the current directory was not kept when the current directory name has a space

NYAGOS 4.4.10\_0
===============
on Jun.25,2021

* nyagos.d/aliases.lua: abspath can be given wildcard names
* Raise a suitable error when `OLEOBJECT._release()` is used instead of `OLEOBJECT:_release()` on Lua script by updating the package glua-ole.
* `echo` outputs a double-quatation same as CMD.EXE does
* Fix that `"..\"..\".."` was changed to `"..\"..".."` as raw-string on parsing.
* (#410) Stop to ignore SIGTERM to shutdown immediately when terminal's close button is pressed (Thx @nocd5)
* On WindowsTerminal (1.8~) , `clone` starts a new shell in the same Window's other tab
* Fix: clone command did not work on WindowsTerminal because wt.exe was not available via apppath
* Fix: su: the network drives were not connected
* Use Makefile(GNU Make) to build instead of PowerShell(make.cmd)
* Can build on Linux

NYAGOS 4.4.9\_7
===============
on May.22,2021

* (#409) Fix: the wildcards expansion with `set -o glob` or `nyagos.option.glob=true` did not work for command alias (Thx @juggler999)

NYAGOS 4.4.9\_6
===============
on May.07,2021

* (#406) Fix: nyagos.argsfilter did not work, raw arguments were not converted and suffix command did not work as expected. (Thx @tGqmJHoJKqgK)

NYAGOS 4.4.9\_5
===============
on May.03,2021

* go-readline-ny v0.4.10: Fix that Yes/No's answer:Y is inserted in the next commandline.
* go-readline-ny v0.4.11: Support Emoji Moifier Sequence (skin tone)
* Improve colored-ls's speed on Windows8.1 when CPU load is high.
* ( Do not use "io/ioutil" )
* Prevent from the message break the commandline when the process finishes which is started by `open` command.
* go-readline-ny v0.4.12: Disable emoji editing in the terminal of VisualStudioCode
* (#403) Support -S,-C and -K options like CMD.EXE
* (#403) Fix: some irregular double-quotations in the commandline were removed when the parameter is sent to the external commands.
* (#405) Add fuzzyfinder catalog module (Thx @tsuyoshicho)

NYAGOS 4.4.9\_4
===============
on Mar.06,2021

* (#400) add check the existance of commands for subcomplete.lua (Thx @tsuyoshicho )
* (#401) add subcompletion choco/chocolaty (Thx @tsuyoshicho )
* Fix: the layout on ls and Ctrl-O selection was broken in WindowsTerminal
* go-readline-ny v0.4.4: Fix: When ANYCHARACTER + Ctrl-B + ZeroWidthJoin Sequence are typed, view was broken
* go-readline-ny v0.4.5: Variation Selector Sequence can include ZeroWidthJoinerSequence
* go-readline-ny v0.4.6: Editing Combining Enclosing Keycap after Variation Selector (&#x0023;&#xFE0F;&#x20E3;)
* Fix: (#402) "echo !xxx" causes unexpected end of nyagos (Thx @masamitsu-murase)
* go-readline-ny v0.4.7: REGIONAL INDICATOR (U+1F1E6..U+1F1FF)
* go-readline-ny v0.4.8: WAVING WHITE FLAG and its variations (U+1F3F3 U+FE0F?)
* go-readline-ny v0.4.9: RAINBOW FLAG (U+1F3F3 U+200D U+1F308)

NYAGOS 4.4.9\_3
===============
on Feb.20,2021

* readline: on WindowsTerminal: Support Variation Selectors of Unicode
* (#397) Add scoop subcommand completion ( `use "subcomplete.lua"` ) (Thx @tomato3713)
* Completion: make English case of completed word to the shortest candidate
* (#398) Fix: the default value for io.popen's second parameter did not work (Thx @ironsand)
* (#399) Improve utf8 offset (Thx @masamitsu-murase)
* Support ALT-/ key bind (Thx @masamitsu-murase) https://github.com/zetamatta/go-readline-ny/pull/1
* readline: Fix the problem that emoji and circled digits could not be input in WindowsTerminal 1.5

NYAGOS 4.4.9\_2
===============
on Jan.8,2021

* (#342) Stop killing child process on Ctrl-C pressed.

NYAGOS 4.4.9\_1
===============
on Dec.21,2020

* Fix: the first `make install` without path parameter fails.
* (#396) Fix: panic on Ctrl-W when left-scrool are required.
* Fix: sometimes more, clip & type from console did not echo input
* (#342) Improve Ctrl-C Interrupt handling to prevent from crash

NYAGOS 4.4.9\_0
===============
on Dec.5,2020

* (#390) Support zero-width-join sequence for WindowsTerminal
* Fix the cursor position broken on VARIATION SELECTOR-1..16  
  ( VARIATION SELECTOR-1..16 are shown like &lt;FE0F&gt; )
* su and clone commands supports WindowsTerminal
* The background process start/end message is no longer displayed during editing.
* C-r: Incremental Search: compare case-insensitively
* Fix: Command-name completion did not work after && and ||.
* C-y: Trim the last CRLF on pasting
* Fix: (#393) the first key after terminal-window activated was input twice. (Thanks to @tostos5963)
* Stop using upx.exe because antivirus software sometimes disjudges as a virus.

NYAGOS 4.4.8\_0
===============
on Oct.3,2020

* git.lua: completion for `git add`: 
    - unquoto filenames like "\343\201\202\343\201\257\343\201\257"
    - support files under untracked directories.
* diskused: shows size like `ls -h`
* Fix: diskused can not be stoped by Ctrl-C
* Implement: environment variable substring like %ENV:~10,5%
* (#308) Fix: error: `The operation was canceled by the user` when the executable is placed on the network folder which is not written by true-UNC-Path
* Fix: `clone` command showed security error dialog when nyagos exists on network.
* (#389) su: keep drive mounting by SUBST
* (#390) Fixed: Some unicode character from U+2000 to U+2FFF could not be input
* (#390) Fixed: Characters represented by Surrogate pair could not be input
* box.lua: Fix: C-o and ESCAPE erased the user-input-word.
* (#391) subcommand.lua: add gh first level subcommand rule (Thanks to @tsuyoshicho)

NYAGOS 4.4.7\_0
===============
on Jul.18,2020

* cd,push and their completion supports %CDPATH% like bash
* load scripts on `%APPDATA%\NYAOS_ORG\nyagos.d`
* On WindowsTerminal, print surrogate-paired unicode by not escaped like &lt;nnnnn&gt;
* Change the directory put binary from Cmd to bin
* catalog/subcomplete.lua
    - Use new completion api:`nyagos.complete_for`
    - Caching subcommands to complete to speed-up nyagos starting. 
    - Implement `clear_subcomands_cache` to clear cache.
    - Subcompletion for `fsutil` and `go`
* catalog/git.lua
    - load `subcomplete.lua` automatically
    - Complete commit-hash like branch-name
    - Complete commit-hash,branch-name and modified filenames after `git checkout`
* (#386) Fix the file size output of `ls -h` to be displayed in units.(Thx! [@Matsuyanagi](https://github.com/Matsuyanagi))
* Fix: `nyagos.exec{ ALIAS-COMMAND-USING $@ }` causes panic
* Add: `nyagos.complete_for_files `(which returns table of completable files)`

NYAGOS 4.4.6\_2
===============
on Jun.09,2020

* Fix: Ctrl-C terminated nyagos.exe like Ctrl-D (which is made on fixing #383 at `4.4.6_0`)

NYAGOS 4.4.6\_1
===============
on May.31,2020

* (#385) Fix: Can not move to any folder in the other drive whose last folder is removed.
* Fix: cd's history did not ignore filepath's case.(see `cd -h`)
* Fix: change drive(`x:`) did not push the last directory to directory history
* Fix: The last element of `nyagos.rawexec{...}` was ignored.

NYAGOS 4.4.6\_0
===============
on May.08,2020

* Implement: %DATE% and %TIME%
* nyagos.envdel now returns removed directories.
* use github.com/zetamatta/go-windows-netresource instead of `dos/net*.go`
* (#379) Add: nyagos.preexechook & postexechook
* (#383) Fixed the bug to go into an infinite loop when the terminal crashes
* Tab-key after `start` completes as a command name as `which`
* Fix: when `cd x:\y\z` failed, the current directory is moved to x:\ (root)

NYAGOS 4.4.5\_4
===============
on Mar.13,2020

* box.lua: fixed C-xC-r, C-xC-h & C-xC-g did not work because github.com/BixData/gluabit32 disappeared.
* (#319) Add bit32.band , bor & bxor again.
* (#378) nyagos.d/catalog/subcomplete.lua: ignore command's case and suffixes by default.
* (#377) Fix: Escape sequence does not work after `git gui` installed by scoop
* make.cmd: do not compress executable on every building by upx.exe. Use it only on making packages.

NYAGOS 4.4.5\_3
===============
on Mar.08,2020

* UNC Path Cache is saved to `~/appdata/local/nyaos_org/computers.txt` rather than `~/appdata/local/nyaos.org/computers.txt` because other features use `nyaos_org` folder.
* Sub command completion(`complete_for`) now matches command-name ignoring its suffix
* Compressed the executable by upx.exe
* Lua function `bit32.*` are not available because github.com/BixData/gluabit32 is 404.
* Use Windows10's native ansi-escape-sequence through mattn/go-colorable
* Fix that `echo $(gawk "BEGIN{ print \"\x22\x22\" }")` could not print double-quatations

NYAGOS 4.4.5\_2
===============
on Oct.26,2019

* (#375) Fix: `~randomstring` causes panic
* (#374) Fix: `ls -l` for future timestamp's files do not print year.

NYAGOS 4.4.5\_1
===============
on Oct.20,2019

* Fix that built-in `box` command did not support selecting multi-items.
* Do not move cursor at drawing [PID] when process starts or shutdowns
* Ctrl-O: do not append a quotation after last backslach (NG: `"Program Files\"` -> OK:`"Program Files\`)
* nyagos.stat/access can understand ~ and %ENV% now.

NYAGOS 4.4.5\_0
===============
on Sep.01,2019

* Implement `nyagos.dirname()` as a Lua function.
* C-o supports selecting multi-files by Space,BackSpace,Shift-H/J/K/L,Ctrl-Left/Right/Down/Up
* Alt-Y(paste with quotes) puts quotes around CRLF
* C-o: append `\`(on Windows) or `/`(on Linux) after choices when they are directories.
* Implement `nyagos.envadd("ENVNAME","DIR")`,`nyagos.envdel("ENVNAME","PATTERN")`
* `nyagos.pathjoin()` now expands `%ENVNAME%` and `~\`,`~/`

NYAGOS 4.4.4\_3
===============
on Jul.07,2019

* (#371) Could not execute `foo.bar.exe` as `foo.bar`
* diskfree shows UNCPath assigned network drive

NYAGOS 4.4.4\_2
===============
on Jun.14,2019

* Speed up the completion for `\\host-name` by updating the cache on background.

NYAGOS 4.4.4\_1
===============
on May.30,2019

* Fix: executable for Linux could not be built.

NYAGOS 4.4.4\_0 the Reiwa edition
=================================
on May.27,2019

* (#233) Completion for `\\host-name\share-name`
* (#238) copy: drawing progress
* Support `ENVVARNAME=VALUE COMMAND PARAMETERS..`
* Fix the problem temporary filenames on executing batchfile conflict
* (#277) Support set /A
* (#291) Print Process-ID when it is started with `&` mark like /bin/sh 
* (#361) Fix the problem GUI Application's STDOUT could not be redirected
* Implemented `.` and `source` for Linux (using /bin/sh)
* readline: fixed cursor-blink-off did not work while user wait.
* Fix: `mklink /J mountPt target-relative-path` made a broken junction.(Add change path absolute)
* Add options: `--chdir "DIR"` and `--netuse "X:=\\host-name\share-name"`
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
