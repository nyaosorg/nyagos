English / [Japanese](./02-Options_ja.md)


## Option for NYAGOS.EXE


### --cleanup-buffer (lua: `nyagos.option.cleanup_buffer=true`)
Clean up key buffer at prompt

## --clipboard (lua: `nyagos.option.clipboard=true`)
Enable clipboard integration for the kill buffer

### --cmd-first "COMMAND"
Execute "COMMAND" before processing any rcfiles and continue shell

### --completion-hidden (lua: `nyagos.option.completion_hidden=true`)
Include hidden files on completion

### --completion-slash (lua: `nyagos.option.completion_slash=true`)
use forward slash on completion

### --glob (lua: `nyagos.option.glob=true`)
Enable to expand wildcards

### --glob-slash (lua: `nyagos.option.glob_slash=true`,`set -o glob_slash`)
Use forward slash `/` on wildcard expansion instead of `\`

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

### --predict (lua: `nyagos.option.predict=true`) [default]
Enable the prediction on readlines (default: enabled)

### --no-cleanup-buffer (lua: `nyagos.option.cleanup_buffer=false`) [default]
Do not clean up key buffer at prompt

### --no-clipboard (lua: `nyagos.option.clipboard=false`)
Disable clipboard integration for the kill buffer

### --no-completion-hidden (lua: `nyagos.option.completion_hidden=false`) [default]
Do not include hidden files on completion

### --no-completion-slash (lua: `nyagos.option.completion_slash=false`) [default]
Do not use slash on completion

### --no-glob (lua: `nyagos.option.glob=false`) [default]
Disable to expand wildcards

### --no-glob-slash (lua: `nyagos.option.glob_slash=false`,`set +o glob_slash`)
Do not use slash on wildcard expansion.

### --no-output-surrogate-pair (lua: `nyagos.option.output_surrogate_pair=false`) [default]
Output surrogate pair characters like `<NNNNN>`

### --no-noclobber (lua: `nyagos.option.noclobber=false`) [default]
Do not forbide to overwrite files no redirect

### --no-read-stdin-as-file (lua: `nyagos.option.read_stdin_as_file=false`) [default]
Read commands from stdin as Windows Console(tty). (Enable to edit line)

### --no-tilde-expansion (lua: `nyagos.option.tilde_expansion=false`)
Disable Tilde Expansion

### --no-usesource (lua: `nyagos.option.usesource=false`)
forbide batchfile to change environment variables of nyagos

### --noclobber (lua: `nyagos.option.noclobber=true`)
forbide to overwrite files on redirect

### --no-predict (lua: `nyagos.option.predict=false`)
Disable the prediction on readlines (default: enabled)

### --norc
Do not load the startup-scripts: `~\.nyagos` , `(BINDIR)\.nyagos` , `(BINDIR)\nyagos.d\*.lua`, and `%APPDATA%\NYAOS_ORG\nyagos.d\*.lua`.

### --output-surrogate-pair (lua: `nyagos.option.output_surrogate_pair=true`)
Output surrogate pair characters as it is


### --read-stdin-as-file (lua: `nyagos.option.read_stdin_as_file=true`)
Read commands from stdin as a file stream (Disable to edit line)

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
