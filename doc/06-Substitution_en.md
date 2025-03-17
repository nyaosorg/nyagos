English / [Japanese](./06-Substitution_ja.md)

## Substitution

### History

* `!!` previous input string
* `!n` n'th input string
* `!-n` n'th previous input string
* `!STR` input string starting with STR
* `!?STR?` input string containing STR

These suffix are available.

* `:0` command name
* `:m` m'th argument
* `^` first argument
* `$` last argument
* `\*` all argument
* `@` the directory

#### Variables

* `nyagos.histchar`: header-character for substitution (default:`!`)
* `nyagos.antihistquot`: characters to disable (default: `'"`)

### Environment variable

* `~` (tilde) are replaced to `%HOME%` or `%USERPROFILE%`.

### Unicode Literal

* `%u+XXXX%` are replaced to Unicode charactor (XXXX is hexadecimal number.)

### Command Substitution (nyagos.d\backquote.lua)

    `COMMAND`
  OR
    $(COMMAND)

is replaced to what COMMAND print to standard output.

### Brace Expansion (nyagos.d\brace.lua)

    echo a{b,c,d}e

is replaced to

    echo abe ace ade

### Inserting Interpreter-name (nyagos.d\suffix.lua)

- `FOO.pl  ...` is replaced to `perl   FOO.pl ...`
- `FOO.py  ...` is replaced to `ipy FOO.py ...` , `py FOO.py` or `python FOO.py ...` (inserted the first found interpreter's name)
- `FOO.rb  ...` is replaced to `ruby   FOO.rb ...`
- `FOO.lua ...` is replaced to `lua    FOO.lua ...`
- `FOO.awk ...` is replaced to `awk -f FOO.awk ...`
- `FOO.js  ...` is replaced to `cscript //nologo FOO.js ...`
- `FOO.vbs ...` is replaced to `cscript //nologo FOO.vbs ...`
- `FOO.ps1 ...` is replaced to `powershell -file FOO.ps1 ...`

To append the new associtation between the suffix and interpreter,
write

    suffix.xxx = "INTERPRETERNAME"
    suffix.xxx = { "INTERPRETERNAME","OPTION" ... }
    suffix[".xxx] = "INTERPRETERNAME"
    suffix[".xxx] = { "INTERPRETERNAME","OPTION" ... }
    suffix(".xxx","INTERPRETERNAME")
    suffix(".xxx",{ "INTERPRETERNAME","OPTION" ... })

in `%USERPROFILE%\\.nyagos`
