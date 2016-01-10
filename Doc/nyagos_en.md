# NYAGOS - Nihongo Yet Another GOing Shell

NYAGOS is the commandline-shell for Windows written with the
Programming Language GO and Lua.

* Support UNICODE
  * Can paste unicode charactor on clipboard and edit them.
  * Unicode-literal %U+XXXX%
  * Prompt Macro $Uxxxx
* Built-in ls
  * color support (-o option)
  * indicate junction-mark as @
* UNIX-Like Shell
  * History (Ctrl-P and !-mark)
  * Alias
  * Filename/Command-name completion
  * Single quatations enclosing a command-line parameter.
* Customizing with Lua
  * built-in command written with Lua
  * filtering command-line
  * useful functions: ANSI-String & UTF8 convert , eval and so on.

## Install

Put files:`nyagos.exe`, `nyagos.lua` and `lua53.dll`, and
directory `nyagos.d` into the one directory pointed with `%PATH%`.

Put .nyagos into the directory pointed with %USERPROFILE%
or %HOME% and modify as you like.

## Option for NYAGOS.EXE

### `-h`

Print Usage.

### `-c "COMMAND"`

Execute `COMMAND` and quit.

### `-k "COMMAND"`

Execute `COMMAND` and continue the command-line.

### `-f SCRIPTFILE ARG1 ARG2 ...`

Execute SCRIPTFILE with Lua interpretor and quit.
The script can refer arguments as `arg[]`.

### `-e "SCRIPTCODE"`

Execute SCRIPTCODE with Lua interpretor and quit.

## Editor

* BackSpace , Ctrl-H : Delete a charactor left of cursor
* Enter , Ctrl-M     : Execute commandline
* Del                : Delete a charactor on cursor
* Home , Ctrl-A      : Move cursor to top
* Left , Ctrl-B      : Move cursor to left
* Ctrl-D             : Delete a charactor on cursor or quit
* End , Ctrl-E       : Move cursor to the tail of commandline
* Right , Ctrl-F     : Move cursor right
* Ctrl-K             : Remove text from cursor to tail
* Ctrl-L             : Repaint screen
* Ctrl-U             : Remove text from top to cursor
* Ctrl-Y             : Paste text from clipboard
* Esc , Ctrl-[       : Remove all-commandline
* UP , Ctrl-P        : Replace commandline to previous input one
* DOWN , Ctrl-N      : Replace commnadline to next input one
* TAB , Ctrl-I       : Complete file or command-name
* Ctrl-C             : Drop text all
* Ctrl-R             : Incremental search
* Ctrl-W             : Remove current word.

## Built-in commnads

These commands have their alias. For example, `ls` => `__ls__`.

### `cd DRIVE:DIRECTORY`

Change the current working drive and directory.
No arguments, move to %HOME% or %USERPROFILE%.

* `cd -` : move the previous directory.
* `cd -N` (N:digit) : move the N-previous directory.
* `cd -h` , `cd ?` : listing directories stayed.

### `exit`

Quit NYAGOS.exe.

### `history [N]`

Display the history. No arguments, the last ten are displayed.

### `ln [-s] SRC DST`

Make hardlink or symbolic-link.
The alias 'lns' defined on `nyagos.d\lns.lua` shows UAC-dialog
and calls `ln -s`.

### `ls -OPTION FILES`

List the directory. Supported options are below:

* `-l` Long format
* `-F` Mark `/` after directories' name. `*' after executables' name.
* `-o` Enable color
* `-a` Print all files.
* `-R` Print Subdirectories recursively.
* `-1` Print filename only.
* `-t` Sort with last modified time.
* `-r` Revert sort order.
* `-h` With -l, print sizes in human readable format (e.g., 1K 234M 2G)
* `-S` Sort by file size

### `pwd`

Print the current woking drive and directory.

* `pwd -N` (N:digit) : print the N-previous directory.

### `set ENV=VAR`

Set the environment variable the value. When the value has any spaces,
you should `set "ENV=VAR"`.

* `PROMPT` ... The macro strings are compatible with CMD.EXE. Supported ANSI-ESCAPE SEQUENCE.

### `which [-a] COMMAND-NAME`

Report which file is executed.

* `-a` - report all executable on %PATH%

### `copy SOURCE-FILENAME DESTINATE-FILENAME`
### `copy SOURCE-FILENAME(S)... DESINATE-DIRECTORY`
### `move OLD-FILENAME NEW-FILENAME`
### `move SOURCE-FILENAME(S)... DESITINATE-DIRECTORY`
### `del FILE(S)...`
### `erase FILE(S)...`
### `mkdir [/p] NEWDIR(S)...`
### `rmdir [/s] DIR(S)...`
### `pushd [DIR]`
### `popd`
### `dirs`

These built-in commands are always asking with prompt when files are override or removed.

### `source BATCHFILENAME`

Execute the batch-file(*.cmd,*.bat) by CMD.exe and
import the environment variables and working directory
which CMD.exe changed.

We use . (one-period) as an alias of source.

## On the Startup

1. NYAGOS.exe loads and execute nyagos.lua where NYAGOS.exe are put. NYAGOS.lua is wrritten with the programming language Lua, and it load .nyagos on HOME directory(=%USERPROFILE% or %HOME%). Users can customize with the .nyagos .
2. History are recorded on %APPDATA%\NYAOS_ORG\nyagos.history

## Substitution

### History

* `!!` previous input string
* `!n` n'th input string
* `!-n` n'th previous input string
* `!STR` input string starting with STR
* `!?STR?` input string containing STR

These suffix are available.

* `:0` command name
* `:m` m'th argument
* `^` first argument
* `$` last argument
* `\*` all argument

#### Variables

* `nyagos.histchar`: header-character for substitution (default:`!`)
* `nyagos.antihistquot`: characters to disable (default: `'"`)

### Environment variable

* `~` (tilde) are replaced to `%HOME%` or `%USERPROFILE%`.

### Unicode Literal

* `%u+XXXX%` are replaced to Unicode charactor (XXXX is hexadecimal number.)

## Lua extension

### `nyagos.alias.NAME = "REPLACED-STRING"`

It defines an alias. These macros are available on "REPLACE-STRING".

* `$1`,`$2`,`$3` ... the number's argument.
* `$\*` ... all arguments

### `nyagos.alias.NAME = function(ARGS)...end`

It assigns the function to the command-name `"NAME"`.
`ARGS` is the table: { 1stArgument,2nd,3rd,... }

When an error occures, the function should return
the number(integer) for %ERRORLEVEL% and error-message.
(No 'return' equals 'return 0,nil')

When the return-value is a string(or string-table), nyagos.exe
executes the string(-table) as a new commandline.

Aliases run on the other Lua-instance and can not access variables
assigned on .nyagos but `share[]`. You can use share[] as you like.
Only the member of the table `share[]` are shared on all Lua-instances 
of nyagos.

### `nyagos.env.NAME`

It is linked to the the environment variable, which are able 
to be refered and assigned.

### `errorlevel,errormessage = nyagos.exec("COMMAND")`

It executes "COMMAND" as shell command.
It returns the integer-value for %ERRORLEVEL% and the error-message.
With no error, they are 0 and nil.

### `errorlevel,errormessage = nyagos.rawexec('COMMAND-NAME','ARG-1','ARG-2'...)`

It executes "COMMAND-NAME" with ARGs. COMMAND-NAME is not interpreted as
a built-in command nor an alias. The difference with os.execute is that
the errormessage is written with utf8.

### `OUTPUT = nyagos.eval("COMMAND")`

It executes "COMMAND" and set its standard output into the lua-variable OUTPUT.
When error occures, OUTPUT is set `nil`.

### `OUTPUT,ERR = nyagos.raweval('COMMAND-NAME','ARG-1','ARG-2'...)`

It executes "COMMAND-NAME" with ARGs and returns commands' standard-output.
COMMAND-NAME is not intepreted as a built-in command nor an alias.

### `WD = nyaos.getwd()`

Get current working directory.

### `nyagos.write("TEXT")`

It output "TEXT" to the standard output with UTF8.

### `nyagos.writerr("TEXT")`

It output "TEXT" to the standard error with UTF8.

### `ANSISTRING = nyagos.utoa(UTF8STRING)`

It converts UTF8 string to the current code page multibyte string.

### `UTF8STRING = nyagos.atou(ANSISTRING)`

It converts the current codepage multibyte string to UTF8 string.

### `FILES = nyagos.glob("WILDCARD-PATTERN1","WILDCARD-PATTERN2"...)`

It returns the table which includes files matching the wildcard pattern(s).

### `path = nyagos.pathjoin('path','to','where'...)`

It makes parts of path-string join.

### `nyagos.bindkey("KEYNAME","FUNCNAME")`

KEYNAME are:
        "C_A" "C_B" ... "C_Z" "M_A" "M_B" ... "M_Z"
        "F1" "F2" ..."F24"
        "BACKSPACE" "CTRL" "DEL" "DOWN" "END"
        "ENTER" "ESCAPE" "HOME" "LEFT" "RIGHT" "SHIFT" "UP"
        "C_BREAK" "CAPSLOCK" "PAGEUP", "PAGEDOWN" "PAUSE"

FUNCNAME are:

        "BACKWARD_DELETE_CHAR" "BACKWARD_CHAR" "CLEAR_SCREEN" "DELETE_CHAR"
        "DELETE_OR_ABORT" "ACCEPT_LINE" "KILL_LINE" "UNIX_LINE_DISCARD"
        "FORWARD_CHAR" "BEGINNING_OF_LINE" "PASS" "YANK" "KILL_WHOLE_LINE"
        "END_OF_LINE" "COMPLETE" "PREVIOUS_HISTORY" "NEXT_HISTORY" "INTR"
        "ISEARCH_BACKWARD"

If it succeeded, it returns true only. Failed, it returns nil and error-message.
Cases are ignores and, the character '-' is same as '\_'.

### `nyagos.bindkey("KEYNAME",function(this)...end)`

When the key is pressed, call the function.

`this` is the table which have these members.

* `this.pos` ... cursor position counted with bytes (==1 when beginning of line)
* `this.text` ... all text represented with utf8
* `this:call("FUNCNAME")` ... call function like `this:call("BACKWARD_DELETE_CHAR")`
* `this:insert("TEXT")` ... insert TEXT at the cursor position.
* `this:firstword()` ... get the first word(=command-name) on the command-line.
* `this:lastword()` ... get the last word and its position on the command-line.
* `this:boxprint({...})` ... listing table values like completion-list.

The return value of function is used as below

* When it is a string, it is inserted into cursor position.
* When it is `true`, accept line as same as Enter is pressed.
* When it is `false`, drop line as same as Ctrl-C is pressed.
* When it is `nil`, it is ignored.

### `nyagos.filter = function(cmdline) ... end`

`nyagos.filter` can modify user input command-line.
If it returns string, NYAGOS.exe replace the command-line-string it.

### `nyagos.argsfilter = function(args) ... end`

`nyagos.argsfilter` is like `nyaos.filter`, but its argument are
not a string but a table as string array which has each command
arguments.

### `length = nyagos.prompt(template)`

`nyagos.prompt` is assigned function which draw prompt.
You can swap the prompt-function as below.

    share.backup_prompt = nyagos.prompt
    nyagos.prompt = function(template)
        nyagos.echo("xxxxx")
        return share.backup_prompt(template)
    end

### `nyagos.gethistory(N)`

Get the n-th command-line history. When N < 0, last (-N)-th history.
With no arguments, get the count of the command-line history.

### `nyagos.access(PATH,MODE)`

Returns the boolean value whether the PATH can be access with MODE.
It equals the access function of the programming language C.

### `nyagos.completion_hook = function(c) ... end`

This is the Hook for completion. It should be assigned a function.
The argument `c` is the table which has these members.

    c.list[1] .. c.list[#c.list] - command/filename completion result
    c.word - original word without double-quotations.
    c.rawword - original word which may has double-quotations.
    c.pos - position word exists.
    c.text - all command-line text.

`nyagos.completion_hook` should return updated list(table) or `nil`.
Returning nil equals to returning c.list with no change.

### `nyagos.on_command_not_found = function(args) ... end`

It is called when the command which user typed is not found.
The command-name and parameters are set to args[0]...args[#args].
If the function returns nil or false, nyagos.exe prints errors of
usual.

Since the function runs the other Lua-instance, accesss to variables
assigned on .nyagos have the same restriction with aliases. 

### `nyagos.getkey()`

It returns three values : typed key's UNICODE,SCANCODE and SHIFT-Status.

### `WIDTH,HEIGHT=nyagos.getviewwidth()`

It returns the width and height of the terminal.

### `STAT = nyagos.stat(FILENAME)`

It returns the file's information.
If the file exists, the table STAT has these members.

    STAT.name
    STAT.isdir (is set true when the file is directory, otherwise false)
    STAT.size  (bytes)
    STAT.mtime.year
    STAT.mtime.month
    STAT.mtime.day
    STAT.mtime.hour
    STAT.mtime.minute
    STAT.mtime.second

If the file does not exist, STAT is nil.

### `nyagos.goversion`

Go-version string to build nyagos.exe
(for example, "go1.5")

### `nyagos.goarch`

The string compilation architecture of nyagos.exe.
(for example, "386" or "amd64" )

### `nyagos.exe`

This string variable has the value of the fullpath of nyagos.exe.

## Misc.

You can get NYAGOS's package from https://github.com/zetamatta/nyagos

On the New BSD-licenses, you can use NYAGOS.

To build nyagos.exe , these softwares are required:

* [go 1.5.2 for windows](http://golang.org)
* [LuaBinaries 5.3 - Release 1 for Win32/64](http://luabinaries.sourceforge.net/download.html)
- http://github.com/mattn/go-runewidth
- http://github.com/shiena/ansicolor
- http://github.com/atotto/clipboard

Thanks to the authors of them.
