English / [Japanese](./02-Options_ja.md)


## Option for NYAGOS.EXE


### --cleanup-buffer (lua: `nyagos.option.cleanup_buffer=true`)
Clean up key buffer at prompt

### --cmd-first "COMMAND"
Execute "COMMAND" before processing any rcfiles and continue shell

### --completion-hidden (lua: `nyagos.option.completion_hidden=true`)
Include hidden files on completion

### --completion-slash (lua: `nyagos.option.completion_slash=true`)
use forward slash on completion

### --disable-virtual-terminal-processing
Do not use Windows10's native ESCAPE SEQUENCE.

### --enable-virtual-terminal-processing
Enable Windows10's native ESCAPE SEQUENCE.
It should be used with `--no-go-colorable`.

### --glob (lua: `nyagos.option.glob=true`)
Enable to expand wildcards

### --go-colorable
Use the ESCAPE SEQUENCE emulation with go-colorable library.

### --help
Print this usage

### --look-curdir-first
Search for the executable from the current directory before %PATH%.
(compatible with CMD.EXE)

### --look-curdir-last
Search for the executable from the current directory after %PATH%.
(compatible with PowerShell)

### --look-curdir-never
Never search for the executable from the current directory
unless %PATH% contains.
(compatible with UNIX Shells)

### --lua-file FILE ARG1 ARG2...
Execute FILE as Lua Script even if FILE's suffix is not .lua .
The script can refer arguments as `arg[]`.
Lines starting with `@` are ignored to embed into batchfile.

### --lua-first "LUACODE"
Execute "LUACODE" before processing any rcfiles and continue shell

### --no-cleanup-buffer (lua: `nyagos.option.cleanup_buffer=false`) [default]
Do not clean up key buffer at prompt

### --no-completion-hidden (lua: `nyagos.option.completion_hidden=false`) [default]
Do not include hidden files on completion

### --no-completion-slash (lua: `nyagos.option.completion_slash=false`) [default]
Do not use slash on completion

### --no-glob (lua: `nyagos.option.glob=false`) [default]
Disable to expand wildcards

### --no-go-colorable
Do not use the ESCAPE SEQUENCE emulation with go-colorable library.

### --no-noclobber (lua: `nyagos.option.noclobber=false`) [default]
Do not forbide to overwrite files no redirect

### --no-tilde-expansion (lua: `nyagos.option.tilde_expansion=false`)
Disable Tilde Expansion

### --no-usesource (lua: `nyagos.option.usesource=false`)
forbide batchfile to change environment variables of nyagos

### --noclobber (lua: `nyagos.option.noclobber=true`)
forbide to overwrite files on redirect

### --norc
Do not load the startup-scripts: `~\.nyagos` , `~\_nyagos`
and `(BINDIR)\nyagos.d\*`.

### --show-version-only
show version only

### --tilde-expansion (lua: `nyagos.option.tilde_expansion=true`) [default]
Enable Tilde Expansion

### --usesource (lua: `nyagos.option.usesource=true`) [default]
allow batchfile to change environment variables of nyagos

### -b "BASE64edCOMMAND"
Decode and execute the command which is encoded with Base64.

### -c "COMMAND"
Execute `COMMAND` and quit.

### -e "SCRIPTCODE"
Execute SCRIPTCODE with Lua interpreter and quit.

### -f FILE ARG1 ARG2 ...
If FILE's suffix is .lua, execute Lua-code on it.
The script can refer arguments as `arg[]`.
Otherwise, read and execute commands on it.

### -h
Print this usage

### -k "COMMAND"
Execute "COMMAND" and continue the command-line.
