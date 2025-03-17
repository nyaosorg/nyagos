English / [Japanese](./01-Install_ja.md)

## Install

### Download Binary

* https://github.com/nyaosorg/nyagos/releases

### Use "Scoop installer"

```
C:> scoop install nyagos
```

### Use "Chocolatey installer"

```
C:> choco install nyagos
```

### Customizing

    copy .nyagos "%USERPROFILE%\."
    notepad "%USERPROFILE%\.nyagos"

And please customize `%USERPROFILE%\.nyagos`

`.nyagos` is the configuration file written with Lua.
(Be careful that the filename starts with `.`(dot)

#### Setting environment variables

If you want to do `SET PATH="%PATH%;C:\BIN`,
write this in `%USERPROFILE%\.nyagos`

    nyagos.env.path = nyagos.env.path .. ";C:\\bin"

#### Setting aliases

If you want to use `lala` as `ls -al`:

    nyagos.alias.lala = 'ls -al $*'

You can use macros `$1`..`$9` and joined string macro `$*`.

You can use the same function defined by Lua for alias:

    nyagos.alias.lala = function(args)
        nyagos.exec{ "ls","-al", unpack(args) }
    end
