English / [Japanese](./04-Commands_ja.md)

## Built-in commands

These commands have their alias. For example, `ls` => `__ls__`.

### `bindkey KEYNAME FUNCNAME`

Customize the key-binding for line-editing.

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
        "ISEARCH_BACKWARD" "REPAINT_ON_NEWLINE"

### `cd DRIVE:DIRECTORY`

Change the current working drive and directory.
No arguments, move to %HOME% or %USERPROFILE%.

* `cd -` : move the previous directory.
* `cd -N` (N:digit) : move the N-previous directory.
* `cd -h` , `cd ?` : listing directories stayed.
* `cd --history` : listing directories stayed all with no decoration

### `env ENVVAR1=VAL1 ENVVAR2=VAL2 ... COMMAND ARG(s)`

While COMMAND is executed, change environment variables.

### `exit`

Quit NYAGOS.exe.

### `history [N]`

Display the history. No arguments, the last ten are displayed.

### `ln [-s] SRC DST`

Make hardlink or symbolic-link.
The alias 'lns' defined on `nyagos.d\lns.lua` shows UAC-dialog
and calls `ln -s`.

### `lnk FILENAME SHORTCUT [WORKING-DIRECTORY]`

Make shortcut.

### `ls -OPTION FILES`

List the directory. Supported options are below:

* `-l` Long format
* `-F` Mark `/` after directories' name. `*` after executables' name.
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

### `set ENV=VAL`

Set the environment variable the value. When the value has any spaces,
you should `set "ENV=VAL"`.

* `PROMPT` ... The macro strings are compatible with CMD.EXE. Supported ANSI-ESCAPE SEQUENCE.
* `set ENV^=VAL` is same as `set ENV=VAL;%ENV%` but removes duplicated VAL.
* `set ENV+=VAL` is same as `set ENV=%ENV%;VAL` but removes duplicated VAL.

### `touch [-t [CC[YY]MMDDhhmm[.ss]]] [-r ref_file ] FILENAME(s)`

If FILENAME exists, update its timestamp, otherwise create it.

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

Execute the batch-file(`*.cmd`,`*.bat`) by CMD.exe and
import the environment variables and working directory
which CMD.exe changed.

We use . (one-period) as an alias of source.

### `open FILE(s)`

Open the file with associated application.

### `clone`

Run another nyagos.exe on another console window.

### `su`

Run another nyagos.exe as Administrator.

### `sudo COMMAND ARGS(s)...`

Run COMMAND as Administrator

## Commands implemented by Lua

### `lua_e "INLINE-LUA-COMMANDS"` (nyagos.d\aliased.lua)

Execute inline-lua-commands like 'lua.exe -e "..."'.

### `cd SHORTCUT.LNK` (nyagos.d\cdlnk.lua)

`cd` cat change the current directory pointed with shortcut-file(`*.lnk`).

### `trash FILE(S)` (nyagos.d\trash.lua)

It throws files into trashbox of Windows.
