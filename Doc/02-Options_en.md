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


### `-e "SCRIPTCODE"`

Execute SCRIPTCODE with Lua interpretor and quit.

### `--norc`

Do not load the startup-scripts: `~\.nyagos` , `~\_nyagos` and `(BINDIR)\nyagos.d\*`.
