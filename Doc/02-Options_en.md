English / [Japanese](./02-Options_ja.md)

## Option for NYAGOS.EXE

### `-h`

Print Usage.

### `-c "COMMAND"`

Execute `COMMAND` and quit.

### `-k "COMMAND"`

Execute `COMMAND` and continue the command-line.

### `-b "BASE64edCOMMAND"`

Decode and execute the command which is encoded with Base64.

### `-f FILE ARG1 ARG2 ...`

If FILE's suffix is .lua, execute Lua-code on it.
(The script can refer arguments as `arg[]`)

Otherwise, read and execute commands on it.

### `--lua-file FILE ARG1 ARG2 ...`

Execute FILE as Lua Script.

### `-e "SCRIPTCODE"`

Execute SCRIPTCODE with Lua interpretor and quit.

### `--norc`

Do not load the startup-scripts: `~\.nyagos` , `~\_nyagos` and `(BINDIR)\nyagos.d\*`.

### `--no-go-colorable`

Do not use the ESCAPE SEQUENCE emulation with go-colorable library.

### `--enable-virtual-terminal-processing`

Enable Windows10's native ESCAPE SEQUENCE. It should be used with `--no-go-colorable`.

### `--look-curdir-first`

Search for the executable from the current directory before %PATH%.
(compatible with CMD.EXE)

### `--look-curdir-last`

Search for the executable from the current directory after %PATH%.
(compatible with PowerShell)

### `--look-curdir-never`

Never search for the executable from the current directory unless %PATH%
contains. (compatible with UNIX Shells)
