English / [Japanese](./02-Options_ja.md)

## Options

### -b "BASE64edCOMMAND"
Decode and execute the command encoded with Base64.

### -c "COMMAND"
Execute `COMMAND` and quit.

### --clipboard / --no-clipboard
(Lua: `nyagos.option.clipboard=true` / `false`)

Enable / Disable clipboard integration for the kill buffer.

### --cmd-first "COMMAND"
Execute "COMMAND" before processing any rcfiles and continue the shell.

### --completion-hidden / --no-completion-hidden
(Lua: `nyagos.option.completion_hidden = true` / `false`)

Include hidden files in completion.  
Use `--no-completion-hidden` to disable it.

### --completion-slash / --no-completion-slash
(Lua: `nyagos.option.completion_slash = true` / `false`)

Use forward slash in completion.  
Use `--no-completion-slash` to disable it.

### -e "SCRIPTCODE"
Execute SCRIPTCODE with the Lua interpreter and quit.

### -f FILE ARG1 ARG2 ...
If FILE's suffix is `.lua`, execute it as a Lua script.  
The script can access arguments as `arg[]`.  
Otherwise, read and execute commands from the file.

### --glob / --no-glob
(Lua: `nyagos.option.glob=true` / `false`)

Enable wildcard expansion.  
Use `--no-glob` to disable it.

### --glob-slash / --no-glob-slash
(Lua: `nyagos.option.glob_slash=true`, `set -o glob_slash`)

Use forward slash `/` for wildcard expansion instead of `\`.  
Use `--no-glob-slash` to disable it.

### -h , --help
Print this usage information.

### -k "COMMAND"
Execute "COMMAND" and continue the command-line.

### --look-curdir-first
Search for executables in the current directory before checking `%PATH%`.  
(Compatible with CMD.EXE)

### --look-curdir-last
Search for executables in the current directory after checking `%PATH%`.  
(Compatible with PowerShell)

### --look-curdir-never
Never search for executables in the current directory  
unless `%PATH%` contains it.  
(Compatible with UNIX shells)

### --lua-file FILE ARG1 ARG2...
Execute FILE as a Lua script, even if its suffix is not `.lua`.  
The script can access arguments as `arg[]`.  
Lines starting with `@` are ignored for batch file embedding.

### --lua-first "LUACODE"
Execute "LUACODE" before processing any rcfiles and continue the shell.

### --noclobber / --no-noclobber
(Lua: `nyagos.option.noclobber=true`)

Prevent overwriting files during redirection.  
Use `--no-noclobber` to disable it.

### --norc
Do not load startup scripts:  
`~\.nyagos`, `(BINDIR)\.nyagos`, `(BINDIR)\nyagos.d\*.lua`, and `%APPDATA%\NYAOS_ORG\nyagos.d\*.lua`.

### --output-surrogate-pair / --no-output-surrogate-pair
(Lua: `nyagos.option.output_surrogate_pair = true` / `false`)

Output surrogate pair characters as they are.  
Use `--no-output-surrogate-pair` to output surrogate pair characters as `<NNNNN>`.

### --predict / --no-predict
(Lua: `nyagos.option.predict=true` / `false`)

Enable prediction in readline.  
Use `--no-predict` to disable it.

### --read-stdin-as-file / --no-read-stdin-as-file
(Lua: `nyagos.option.read_stdin_as_file = true` / `false`)

Read commands from stdin as a file stream. (Disables line editing.)  
Use `--no-read-stdin-as-file` to disable it.

### --show-version-only
Show version only.

### --tilde-expansion / --no-tilde-expansion
(Lua: `nyagos.option.tilde_expansion=true`)

Enable tilde expansion.  
Use `--no-tilde-expansion` to disable it.

### --usesource / --no-usesource
(Lua: `nyagos.option.usesource=true`)

Allow batch files to change the environment variables of nyagos.  
Use `--no-usesource` to disable it.

