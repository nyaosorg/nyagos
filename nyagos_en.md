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
* Customizing with Lua
  * built-in command written with Lua
  * filtering command-line
  * useful functions: ANSI-String & UTF8 convert , eval and so on.

## Install

Put `nyagos.exe`, `nyagos.lua`, `lua52.dll` into the one 
directory pointed with `%PATH%`.

Put .nyagos into the directory pointed with %USERPROFILE%
or %HOME% and modify as you like.

## Option for NYAGOS.EXE

### `-c "COMMAND"`

Execute `COMMAND` and quit.

### `-k "COMMAND"`

Execute `COMMAND` and continue the command-line.

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


## Built-in commnads

### `alias ALIAS=DEFINE`

Define the alias. Macros in DEFINE:

* `$n` which is replaced the n'th parameter
* `$*` which is replaced to all parameters

When DEFINE is empty, the alias is removed.
Not found the mark =, display the define of the alias.

No arguments, list up the all aliases.

### `cd DRIVE:DIRECTORY`

Change the current working drive and directory.
No arguments, move to %HOME% or %USERPROFILE%.

### `echo STRING`

Print STRING to Standard-output with UTF-8

### `exit`

Quit NYAGOS.exe.

### `history [N]`

Display the history. No arguments, the last ten are displayed.

### `ls -OPTION FILES`

List the directory. Supported options are below:

* `-l` Long format
* `-F` Mark `/` after directories' name. `*' after executables' name.
* `-o` Enable color
* `-a` Print all files.
* `-R` Print Subdirectories recursively.

### `pwd`

Print the current woking drive and directory.

### `set ENV=VAR`

Set the environment variable the value. When the value has any spaces,
you should `set "ENV=VAR"`.

* `PROMPT` ... Escape sequences are avaliable.

#### Special Variable

* `PROMPT` ... The macro strings are compatible with CMD.EXE. Supported ANSI-ESCAPE SEQUENCE.

### `source BATCHFILENAME`

Execute the batch-file(*.cmd,*.bat) for CMD.exe.
The batch-file is executed with CMD.exe , so it can not use
NYAGOS' built-in commands, but the environment variables
CMD.exe changes are imported to NYAGOS.exe.

We use . (one-period) as an alias of source.

## On the Startup

1. NYAGOS.exe loads and execute nyagos.lua where NYAGOS.exe are put. NYAGOS.lua is wrritten with the programming language Lua, and it load .nyagos on HOME directory(=%USERPROFILE% or %HOME%). Users can customize with the .nyagos .
2. History are recorded on %APPDATA%\NYAOS_ORG\nyagos.history

## Substitution

### History

* `!!` previous input string
* `!n` n'th input string
* `!-n` n'th previous input string

These suffix are available.

* `:0` command name
* `:m` m'th argument
* `^` first argument
* `$` last argument
* `\*` all argument

### Environment variable

* `~` (tilde) are replaced to `%HOME%` or `%USERPROFILE%`.

### Unicode Literal

* `%u+XXXX%` are replaced to Unicode charactor (XXXX is hexadecimal number.)

## Lua extension

### `nyagos.alias("NAME","REPLACED-STRING")`

It defines an alias. These macros are available on "REPLACE-STRING".

* $1,$2,$3 ... the number's argument.
* $* ... all arguments

### `nyagos.setenv("NAME","VALUE")`

It changes the environment variable.

### `nyagos.exec("COMMAND")`

It executes "COMMAND" as shell command.

### `OUTPUT = nyagos.eval("COMMAND")`

It executes "COMMAND" and set its standard output into the lua-variable OUTPUT.
When error occures, OUTPUT is set `nil`.

### `nyagos.echo("TEXT")`

It output "TEXT" to the standard output with UTF8.

### `ANSISTRING = nyagos.utoa(UTF8STRING)`

It converts UTF8 string to the current code page multibyte string.

### `UTF8STRING = nyagos.atou(ANSISTRING)`

It converts the current codepage multibyte string to UTF8 string.

### `files = nyagos.glob("WILDCARD-PATTERN")`

It expands the wildcard pattern to table contains filenames.

### `nyagos.bindkey("KEYNAME","FUNCNAME")`

KEYNAME are:

	"BACKSPACE" "CTRL" "C_A" "C_B" "C_C" "C_D" "C_E" "C_F" "C_G" "C_H"
	"C_I" "C_J" "C_K" "C_L" "C_M" "C_N" "C_O" "C_P" "C_Q" "C_R" "C_S"
	"C_T" "C_U" "C_V" "C_W" "C_X" "C_Y" "C_Z" "DEL" "DOWN" "END"
	"ENTER" "ESCAPE" "HOME" "LEFT" "RIGHT" "SHIFT" "UP"

FUNCNAME are:

	"BACKSPACE" "BACKWORD" "CLEAR" "DELETE" "DELETE_OR_ABORT"
	"ENTER" "ERASEAFTER" "ERASEBEFORE" "FORWARD" "HEAD"
	"PASS" "PASTE" "REPAINT" "TAIL" "HISTORY_UP" "HISTORY_DOWN"
        "COMPLETE"

If it succeeded, it returns true only. Failed, it returns nil and error-message.

### `nyagos.filter = function(cmdline) ... end`

`nyagos.filter` can modify user input command-line.
If it returns string, NYAGOS.exe replace the command-line-string it.

### `nyagos.argsfilter = function(args) ... end`

`nyagos.argsfilter` is like `nyaos.filter`, but its argument are
not a string but a table as string array which has each command 
arguments.

## Misc.

You can get the source files from
    https://github.com/zetamatta/nyagos
the binary files from
    http://www.nyaos.org/index.cgi?p=NYAGOS .

On the New BSD-licenses, you can use NYAGOS.

To build nyagos.exe , these softwares are required:

- [go1.3 windows/386](http://golang.org)
- [Mingw-Gcc 4.8.1-4](http://mingw.org/)
- [LuaBinaries 5.2.3 for Win32 and MinGW](http://luabinaries.sourceforge.net/index.html)
- http://github.com/mattn/go-runewidth
- http://github.com/shiena/ansicolor
- http://github.com/atotto/clipboard

Thanks to the authors of them.
