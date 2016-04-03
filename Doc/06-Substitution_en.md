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

#### Variables

* `nyagos.histchar`: header-character for substitution (default:`!`)
* `nyagos.antihistquot`: characters to disable (default: `'"`)

### Environment variable

* `~` (tilde) are replaced to `%HOME%` or `%USERPROFILE%`.

### Unicode Literal

* `%u+XXXX%` are replaced to Unicode charactor (XXXX is hexadecimal number.)

