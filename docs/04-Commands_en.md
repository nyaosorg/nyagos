[top](../readme.md) &gt; English / [Japanese](./04-Commands_ja.md)

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
* `cd shortcut.lnk` : move the target directory pointed shortcut.lnk

If the directory-path does not contain `:`,`/`,`\` and it does not
exists in the current directory, seek the directory to move in the
list of %CDPATH%.

### `chmod ooo FILE(s)`

### `env ENVVAR1=VAL1 ENVVAR2=VAL2 ... COMMAND ARG(s)`

While COMMAND is executed, change environment variables.

### `exit`

Quit NYAGOS.exe.

### foreach

`foreach` *VAR* *VAL1* *VAL2* ...
    STATEMENTS
`end`

### `history [N]`

Display the history. No arguments, the last ten are displayed.

### if

#### inline-if

`if` *COND* *THEN-STATEMENT*

#### block-if

`if` *COND* [`then`]
   *THEN-BLOCK*
`else`
   *ELSE-BLOCK*
`end`

* `endif` can be used as the alias of `end` for compatibility to nyaos-3000
* `then` can be ommited.

*COND* is:

* `not` *COND*
* `/i` *COND*
* *LEFT* `==` *RIGHT*
* `EXIST` *filename*
* `ERRORLEVEL` *n*

* if *COND* is true, execute *THEN-BLOCK* or *THEN-STATEMENT*
* if *COND* is false, execute *ELSE-BLOCK* or nothing.

### `kill PID`

Kill process specified by PID

### `killall NAME...`

Kill process by name

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
* `-?` Display help
* `-L` Show information for the file refernces rather than for the link it self.

Example of `ls -al`

```
drwx-sh    0 Feb 19 20:16:53 System Volume Information/
drwxa--    0 Sep  3 2016     TDM-GCC-64/
```

What FLAGs are meaning

- `d` - Directory
- `r` - Readable
- `w` - Writable (not read-only file)
- `x` - Executable
- `a` - Ready for archiving
- `s` - System file
- `h` - Hidden file


### `more`

Support both UTF8 and ANSI-text (auto detected)

### `ps`

Show a list of processes running.

### `pwd`

Print the current woking drive and directory.

* `pwd -N` (N:digit) : print the N-previous directory.
* `pwd -L` : use PWD from environment, even if it contains symlinks.(default)
* `pwd -P` : avoid symlinks.

### `set ENV=VAL`

Set the environment variable the value. When the value has any spaces,
you should `set "ENV=VAL"`.

* `PROMPT` ... The macro strings are compatible with CMD.EXE. Supported ANSI-ESCAPE SEQUENCE.
* `set ENV^=VAL` is same as `set ENV=VAL;%ENV%` but removes duplicated VAL.
* `set ENV+=VAL` is same as `set ENV=%ENV%;VAL` but removes duplicated VAL.

### `set -o OPTION-NAME`, `set +o OPTION-NAME`

`-o` makes OPTION true, `+o` false.

- `-o glob` enables the wildcard expansion on external commands also.
- `-o noclobber` overwriting the existing file by redirect is forbidden.
- `-o usesource` batchfiles can change the environment variable of nyagos.
- `+o usesource` you have to use `source BATCHFILE` to read the changes of the environment variables from batchfiles.
- `-o cleaup_buffer` clean up console input buffer before readline.

### `set -a "EQUATION"`, `set /a "EQUATION"`

Same as CMD.EXE. Evalute EQUATION

### `select FILENAME(s)...`

Open a file with dialog to select application.

### `touch [-t [CC[YY]MMDDhhmm[.ss]]] [-r ref_file ] FILENAME(s)`

If FILENAME exists, update its timestamp, otherwise create it.

### `which [-a] COMMAND-NAME`

Report which file is executed.

* `-a` - report all executable on %PATH%

### `copy SOURCE-FILENAME DESTINATE-FILENAME`
### `copy SOURCE-FILENAME(S)... DESINATE-DIRECTORY`
### `copy SOURCE-FILENAME(S)... SHORT-CUT(*.lnk)`
### `move OLD-FILENAME NEW-FILENAME`
### `move SOURCE-FILENAME(S)... DESITINATE-DIRECTORY`
### `move SOURCE-FILENAME(S)... SHORT-CUT(*.lnk)`
### `del FILE(S)...`
### `erase FILE(S)...`
### `mkdir [/p] NEWDIR(S)...`
### `rmdir [/s] DIR(S)...`
### `pushd [DIR]`
### `popd`
### `dirs`
### `diskfree`
### `diskused`

These built-in commands are always asking with prompt when files are override or removed.

### `source [-v] [-d] BATCHFILENAME`

Execute the batch-file(`*.cmd`,`*.bat`) by CMD.exe and
import the environment variables and working directory
which CMD.exe changed.

- We use `.` (one-period) as an alias of `source`.
- `source` makes a temporary file: `%TEMP%\nyagos-(PID).tmp`
    - It contains the new values of new current working directory and 
       the environemnt variables.
- With option -d, temporary files made by `source` is not to be removed.
- With option -v, `source` shows the temporari files to STDERR.

### `open FILE(s)`

Open the file with associated application.

### `clone`

Run another nyagos.exe on another console window.

### `su`

Run another nyagos.exe as Administrator.

## Commands implemented by Lua

### `abspath ARG(s)...` (nyagos.d\aliases.lua)

Print the absolute path of ARGs which are written in the relative path.

### `chompf FILE(s)` (nyagos.d\aliases.lua)

Output the contents of FILE(s) to STDOUT without the last CRLF before EOF.

### `lua_e "INLINE-LUA-COMMANDS"` (nyagos.d\aliases.lua)

Execute inline-lua-commands like 'lua.exe -e "..."'.

### `lua_f "LUA-SCRIPT-FILENAME" ARG(s)...` (nyagos.d\aliases.lua)

Execute lua-script.

### `trash FILE(S)` (nyagos.d\trash.lua)

It throws files into trashbox of Windows.

### `wildcard COMMAND ARG(s)...` (nyagos.d\aliases.lua)

Expand the wildcard included ARG(s) and call COMMAND.
