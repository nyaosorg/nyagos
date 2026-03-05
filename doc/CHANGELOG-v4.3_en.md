Changelog v4.3
==============
( English / [Japanese](history-4.3_ja.md) )

v4.3.3\_5
---------
Dec 24, 2018

* (#345) Fix subcomplete.lua don't work git (Thx! @tsuyoshicho)
* (#347) Fix the bug that STDOUT was closed after `dir 2>&1`.(Thx! @Matsuyanagi)
* (#348) Scrolling by mouse-wheel did not worked. (Thx! @tyochiai)
    * It requires github.com/zetamatta/go-getch tagged 20181223.

v4.3.3\_4
---------
Dec 13, 2018

* If stdin is not terminal, `more` command runs as `type`.
* On calling a batch file, `use CMD.EXE /V:ON /S /C "..."` for boosting code instead of temporary batchfile.
* (#340) Add lua variable `nyagos.histsize` to set the number of entries for history to save disk. (Thx! @crile)
* (#343) When %COMSPEC% is empty, use CMD.EXE (Thx! @orz--)

v4.3.3\_3
---------
Oct 23, 2018

* (#310) copy and move support shortcut files(`*.lnk`) as destination.
* (#313 reopened) Fix problem when `git blame FILES | type | gvim - &`, gvim starts with empty buffer.
* Fix: rmdir could not remove the broken junction
* Fix: Ctrl-C did not work in Lua-Script and some extern process
* (#267) `type` and `more` support UTF16 (requires go-texts package)
* (#336) Fix `io.write` did not work with -e and --lua-file
* (#337) Fix the crash the batchfile exit with -1 (Thx! @hogewest)

v4.3.3\_2
---------
Sep 22, 2018

* Append error message the filename on overwriting to existing file on redirect.
* Fix error for overwriting on redirect to `nul` when `noclobber` is set.
* diskused: continue counting how bytes disk used even if errors are found.
* ls: fixed `-l` option did not work with `-1` option
* ls: fixed: did not show one file per a line when output is not terminal.
* Not aliased builtin commands are able to be called as `\ls` like bash
* Fix the broken alias "for"
* Fix on completion the path separating characters were replaced to default one even if the word was not filepath for #334

v4.3.3\_1
---------
Aug 29, 2018

* #330,#331 Fix the original version of file:read incompatible behavior (Thx! @erw7)
* #332 stop buffering on io.open("w") (Thx! @spiegel-im-spiegel)
* #333 Fix file:seek() did not work on reading as expected (Thx! @erw7)
* #333 Fix file:close()'s return value was invalid. (Thx! @erw7)
* #319 Impl utf8.len()
* Fix: `which` reported files which has no suffixes
* `pwd` shows logical-path (=pwd -l) as default rather than phisical-path (=pwd -p)
* Fix: trash was left when incremental-search starts and some string exists on command-line.
* Shrink the executable with -lfdflags="-s -w"

v4.3.3\_0
---------
Aug 14, 2018

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

v4.3.2\_0
---------
Jul 23, 2018

* #319 Support lua `bit32.*` all by github.com/BixData/gluabit32
* #323 Fix io.lines(), nyagos.lines() could not read from redirected stdin
* Fix: io.write() did not write to redirected stdout
* Replace `io.*` all with nyagos' own functions
* #324 Fix: Lua's print ignored --no-go-colorable (Thx @tignear)
* #325 Fix: `source` could not load the path which contains spaces.
* Add options: `--lua-first` and `--cmd-first`

v4.3.1\_3
---------
Jun 19, 2018

* #316 Fix: zero-length directory-name in %PATH% is regarded as the current directory
* #321 Fix: key function names `previous_history` & `next_history` were not registered.
* Add -h and --help option
* Lines starting with `@` of Lua script are now ignored to embed into batchfile.
* #322 Fix: change the encoding for batchfile's parameters from Thread Codepage to Console Codepage #322
* All of lua variables `nyagos.option.*` are now able to be set by nyagos.exe's command-line option.

v4.3.1\_2
---------
Jun 12, 2018

* #320: fix the imcompatibility: nyagos.rawexec & raweval did not expand tables in arguments.
* --show-version-only enables --norc automatically

v4.3.1\_1
---------
Jun 11, 2018

* Remove source code for lua53.dll
* #317: deadlock when `use "subcomplete"` is enabled and rclone.exe is found.
    - See also: https://github.com/yuin/gopher-lua/issues/181
* #318,#319: add compatible functions with lua 5.3
    - bit32.band/bitor/bxor
    - utf8.char/charpattern/codes

v4.3.1\_0
---------
Jun 3, 2018

* Support Windows10's native ESCAPE SEQUENCE processing with --no-go-colorable and --enable-virtual-terminal-processing
* For #304,#312, added options to search for the executable from the current directory
    * --look-curdir-first: do before %PATH% (compatible with CMD.EXE)
    * --look-curdir-last : do after %PATH% (compatible with PowerShell)
    * --look-curdir-never: never (compatible with UNIX Shells)
* nyagos.prompt can now be assigned string literal as prompt template directly.
* Fix #314 rmdir could not remove junctions.

v4.3.0\_4
---------
May 12, 2018

- Fix: #309 nyagos.getkey() raised panic (Thx @nocd5)
- Fix: error-message when command `lnk`'s target is not `*.lnk` nor exist.
- Fix: the cursor blink was switched to off on the child process.

v4.3.0\_3
---------
May 9, 2018

- Fix: forgot implement nyagos.setalias , nyagos.getalias (`alias { CMD=XXX}` did not work.)
- Fix: that the element [0] of the table value returned by alias-function was not used as the new command name to evaluate.
- Fix: `doc/09-Build_*.md` about how to download sourcefiles from github

v4.3.0\_2
---------
May 7, 2018

- #305: Fix issue that user's .nyagos was not loaded again (Thx! @erw7)

v4.3.0\_1
---------
May 5, 2018

- Fix: nyagos.d/start.lua did not worked because the member `rawargs` of alias-function's argument was not implemented.
- Fix: the return value of alias-function was not evaluted.
- Fix: for the script in -e option, arg[] was not assinged.
- Fix: On -f & -e option, warned as `getRegInt: could not find shell in Lua instanc
e`
- Fix: batchfile cound not return the value of `exit /b` as ERRORLEVEL

v4.3.0\_0
---------
May 3, 2018

- Add `ls -L` which shows information for the file refernces rather than for the link it self.

v4.3\_beta2
-----------
May 1, 2018

- Fix: Typing C-o looks to raise hang up until Enter or ESCAPE is typed (on 4.3beta) #303
    - Fix the library: [go-box](https://github.com/zetamatta/go-box/commit/322b2318471f1ad3ce99a3531118b7095cdf3842)
- Fix: chcp did not work. (`chcp` was aliaes to update memory of screen width)

v4.3\_beta
----------
Apr 30, 2018

- Use Gopher-Lua instead of lua53.dll #300
    - nyagos.exe with lua53.dll can be built with `cd mains ; go build`
    - nyagos.exe with no Lua can be built with `cd ngs ; go build`
- Made `nyagos.option.cleanup_buffer` (default=false). When it is true, clean up console input buffer before readline.
- `set -o OPTION_NAME` and `set +o OPTION_NAME` (=`nyagos.option.OPTION_NAME=` on Lua)
- Buffer console-output ( go-colorable and bufio.Writer )
