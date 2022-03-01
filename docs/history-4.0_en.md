[top](../readme.md) &gt; English / [Japanese](history-4.0_ja.md)

NYAGOS 4.0.9\_11
================
on Dec 12,2015

* Build with go 1.5.2
* Readline: support Ctrl-C as interrupt on isearch (#98 Thx @hattya)
* Add suppression options (#95,#97 Thx @hattya)
    - Add /q option to copy,move,del,erase and rmdir.
* ls: calculate width of file size field for each directory (#94 Thx @hattya)
* move: be able to override files (api: MoveFileW -> MoveFileExW) (#93 Thx @hattya)
* Fixed lua-function 'set' did not expand env-var including `_` (#92) 
* Fixed typographical error(s)
    - Changed charactor to character in README (#90: @orthographic-pedant)
* `nyagos.command_not_found` became deprecated.
    - Removed from Document (it gets deprecated)
    - Removed catalog.d/ezoe.lua
* To avoid crash (#89)
    - Forbade nyagos.prompt to execute background
    - Forbade using Lua on background or not 1st command of pipeline

NYAGOS 4.0.9\_10
================
on Sep 17,2015

* Build with go 1.5.1
* Fixed #88 Ctrl-U (`UNIX_LINE_DISCARD`) did not consider scroll.
* Fixed #87 Querying not-starting process's errorlevel which occured crash.
* Do not quit and do warn only when nyagos.prompt returns NaN.
* Fixed file-handle created with REDIRECT was not closed.
* Improved make.cmd
* Fixed interpreter instance was not given to nyagos.prompt.
* `which` command supports -a option
* Support C-w (unix-word-rubout) like nyaos (#85)
* Fixed wrong month in prompt 'd' (Thx @Matsuyanagi)

NYAGOS 4.0.9\_9
===============
on Sep 07,2015

* Fixed #80 built-in commands and aliases could not return errorlevel
* Let nyagos.exec return %ERRORLEVEL% and error-message.
* Fixed #83 panic: calling `nyagos.exec` on `nyagos.on_command_not_found`
* Fixed #82 panic: using pipeline between replacing-type-aliases.
* Fixed #81 No errors reported on rmdir NOT directory before prompt.

NYAGOS 4.0.9\_8
===============
on Aug 20,2015

* Build with go 1.5
* Set go-version to Lua-variable: `nyagos.goversion` and print at startup.
* Support to build both 32bit and 64bit executable.

NYAGOS 4.0.9\_7
===============
on Jul 29,2015

* Fixed a panic occurs when only '||' command-line typed.

NYAGOS 4.0.9\_6
===============
on Jul 21,2015

* Fixed nyagos.stat(nil) caused a panic.
* Fixed not all stack trace was printed when a panic was recovered.

NYAGOS 4.0.9\_5
===============
on Jul 15,2015

* Add ls-option -S (Sort by Size) and -h (Print size with human-readable format)
* Add nyagos.rawexec and nyagos.raweval
* Forbade to use more than one lua-command on the same pipeline to avoid crash.
* Forbade to run Lua background to avoid crash.
* Fixed #77: echo "{a,b}" -> "a b" is printed.({a,b} should be printed)

NYAGOS 4.0.9\_4
===============
on Jun 24,2015

* Fixed `ls (NO_MATCHING_WILDCARD)` worked as if `ls` with (NO-ARGUMENTS)
* Fixed `A ; B` worked as if `A & B`.

NYAGOS 4.0.9\_3
===============
on Jun 15,2015

* Fixed nyagos.stat did not work on reparsepoints.
* Fixed #74 "ls -a" never prints "." and ".." (current and parent directory)
* Disabled ls-color when stdout is not console even if -o option exists.

NYAGOS 4.0.9\_2
===============
on Jun 12,2015

* Fixed the limiter of alias subsutitution was only one. Changed upto 5
* Fixed more than one nyagos.argsfiler could run at once on pipelines. It caused to crash the process.

NYAGOS 4.0.9\_1
===============
on Jun 03,2015

* Fixed nyagos.exec() crashed when it was redirected.

NYAGOS 4.0.9\_0
===============
on Jun 02,2015

* Add the lua-table `nyagos.env`
* Add the lua-hook `nyagos.on_command_not_found`
* Add the lua-function `nyagos.getviewwidth` (Thx @nocd5)
* Supported Surrogate-Pair on GetKey (Thx @rururutan)
* Add the lua-function `nyagos.stat`
* Fixed #72 an empty argument was inserted before SPACE & REDIRECT-MARK (Thx @hogewest)
* Made lua-script catalog folder 'catalog.d'(not loaded automatically)
* Associated the suffix .py to either of IronPython or CPython automatically(Thx @hattya)

NYAGOS 4.0.8\_0
===============
on May 28,2015

* Supported single-quatations like UNIX-Shell
* Added lua-function `nyagos.getkey()`/`nyagos.getalias()`
* Renamed lua-function `nyagos.alias` to `nyagos.setalias()`
* Made alias-defining-table `nyagos.alias`. (nyagos.lua)
* Changed the output of `x("COMMAND")` from stdout to stderr.
* Added `cdlnk.lua` which enables `cd SHORTCUT.LNK`.
* Improved `su` and `clone` to retry as `%COMSPEC% /c NYAGOS.EXE` when error occued.
* Improved lua-function nyagos.glob to receives more than one wildcards.
* Added `trash.lua` which provides `trash` command which throws files Windows' trashbox.
* Replaced the bundled NYOLE.DLL 0.0.0.4 to 0.0.0.5 which trash.lua requires.
* Enabled NYAGOS.EXE run without NYOLE.DLL (trash.lua & cdlnk.lua are disabled)
* Replaced lua53.dll to LuaBinaries' version.
* Unbundled lua.exe from zip-file.

Bugfix
-------
* Fixed #66 `echo a>a` makes a empty file named `aa`
* Fixed suffix.lua problem (#69 Thx hattya)
* Fixed panic when nyagos.argsfilter returns empty array(#68 Thx hattya)
* Reset to default color after ls (#67 Thx @hattya)

NYAGOS 4.0.7\_5
===============
on May 09,2015

* Fix #64 invalid cursor position after Ctrl-T typed.(Not compatible with other shells)

NYAGOS 4.0.7\_4
===============
on May 05,2015

* Fix bug that filename-completion replacing slash all to backslash.
* Fix #63 ESCAPE-Key let clipboard empty-string. (Thx hokorobi)
* Let Ctrl-U copy erased string to clipboard.

NYAGOS 4.0.7\_3
===============
on May 05,2015

* `SET VAR=` removes environment variable `VAR` (Thx @pine613)
* lnk.js with one parameter shows its linked file. (#59 Thx @NSP-0123456)
* Support Ctrl-T (swapchar) (#62)

BugFix
------
* missed a double-quatation after a percent (#57)
* wildcard's case was not ignored (#58)
* completion removed .\ (dot and shash) (#61)
* `open ARGUMENT(s)` did not work

NYAGOS 4.0.7\_2
===============
on Apr 18,2015

* Fix:on completion, / was always replaced to \ .(Thx @nocd5)
* Fix:nyagos.shellexecute() did not report some errors.
* To use COM on Lua scripts, include and use [NYOLE.DLL](https://github.com/zetamatta/nyole).

NYAGOS 4.0.7\_1
===============
on Apr 13,2015

* Set arg[..] in nyagos -e "LUA-CODE".

Bugfix
------
* In nyagos -f "LUA-FILE", arg[i]'s elements were shifted.

NYAGOS 4.0.7\_0
===============
on Apr 15,2015

* Support incremental-search(C-r)
* Add option -e "LUA-CODE" to nyagos.exe
* Set executable's property the version-number
* Change error-message when files do not exists like bash.

NYAGOS 4.0.6\_0
===============
on Mar 19,2015

* Add built-in command: pushd/popd/dirs
* Add the method boxprint(),firstword(),lastword() to nyagos.bindkey's first argument
* Update Document about uninstalling
* Add `nyagos.completion_hook`
* Sub-command completion for git, Subversion and Mercurial.

Bugfix
------
* Completion failed when 0001 is typed where 0001.txt and "0001 copy.txt" were.

NYAGOS 4.0.5\_0
================
on Feb 07,2015

* cd -N (N:digit): move the N-th previous directory.
* cd -h , cd ? : print current directory history.
* pwd -N (N:digit): print the N-th previous directory.
* %CD% , %ERRORLEVEL% are able to be completed.

Bugfix
------
* Lua-function 'include' didn't report error
* Command-name completion printed same-name in diffent directories.
* ReadLine sometimes left trash at replacing string
* Could not complete path including ./ #45
* Could not broken symbolic link with DEL #44
* Files contains '&' was not enclosed with ".." at completion.

Trivial fix
-----------
* make.cmd: add echo off a lot
* Let make.cmd on the top directory without arguments copy EXE top 

NYAGOS 4.0.4\_0
================
on Jan 19,2015

* Support Lua 5.3

NYAGOS 4.0.3\_2
===============
on Jan 18,2015

* Command-name completion supports alias and built-in commands.

BugFix
------
* `pwd` did not print correct UNC-Path (#37)
* `nyagos.gethistory( large-value )` crashes nyagos (#38 @1)
* `%APPDATA%/NYAOS_ORG/nyagos.history` did not be updated. (#39 @1)
* Crashed when `%APPDATA%/NYAOS_ORG/nyagos.history` did not exists or is empty.x (#40 @1)
* On French keyboard, keys shifted with AltGr could not be input. (#41)

@1 Probably these bugs appeared only on snapshot-build.
