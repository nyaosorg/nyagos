## Option for NYAGOS.EXE

### `-h`

Print Usage.

### `-c "COMMAND"`

Execute `COMMAND` and quit.

### `-k "COMMAND"`

Execute `COMMAND` and continue the command-line.

### `-f FILE ARG1 ARG2 ...`

If FILE's suffix is .lua, execute Lua-code on it.
(The script can refer arguments as `arg[]`)

Otherwise, read and execute commands on it.


### `-e "SCRIPTCODE"`

Execute SCRIPTCODE with Lua interpretor and quit.


