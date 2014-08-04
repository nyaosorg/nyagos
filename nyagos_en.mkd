# NYAGOS - Nihongo Yet Another GOing Shell

NYAGOS is the commandline-shell for Windows written with the
Programming Language GO and Lua. NYAGOS supports UNICODE.

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

## Built-in commnads

### `alias ALIAS=DEFINE`

Define the alias. Macros in DEFINE:

* `$n` which is replaced the n'th parameter
* `$*` which is replaced to all parameters

When DEFINE is empty, the alias is removed.
Not found the mark =, display the define of the alias.

No arguments, list up the all aliases.

These aliaes are defined in nyagos.lua:

* `assoc=%COMSPEC% /c assoc`
* `attrib=%COMSPEC% /c attrib`
* `copy=%COMSPEC% /c copy`
* `del=%COMSPEC% /c del`
* `dir=%COMSPEC% /c dir`
* `for=%COMSPEC% /c for`
* `md=%COMSPEC% /c md`
* `mkdir=%COMSPEC% /c mkdir`
* `mklink=%COMSPEC% /c mklink`
* `move=%COMSPEC% /c move`
* `rd=%COMSPEC% /c rd`
* `ren=%COMSPEC% /c ren`
* `rename=%COMSPEC% /c rename`
* `rmdir=%COMSPEC% /c rmdir`
* `start=%COMSPEC% /c start`
* `type=%COMSPEC% /c type`

### `cd DRIVE:DIRECTORY`

Change the current working drive and directory.
No arguments, move to %HOME% or %USERPROFILE%.

### `echo STRING`

Print STRING.

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

#### Special Variable

* `PROMPT` ... The macro strings are compatible with CMD.EXE. Supported ANSI-ESCAPE SEQUENCE.

### `source BATCHFILENAME`

## On the Startup

## Substitution

### History

### Environment variable

## Lua Extension
