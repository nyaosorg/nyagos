English / [Japanese](./01-Install_ja.md)

## Install

The binary files can be downloaded on [Release](https://github.com/zetamatta/nyagos/releases).

    mkdir PATH\TO\INSTALLDIR
    cd PATH\TO\INSTALLDIR
    unzip PATH\TO\DOWNLOADDIR\nyagos-****.zip
    makeicon.cmd

## Easy customizing

    notepad %USERPROFILE%\_nyagos

`_nyagos` is the dos-batchfile-like configuration file.
(Be careful that the filename starts with `_`(underscore)

For example:

    alias "grep=findstr"
    set "GOPATH=%USERPROFILE%\Share\GoSrc"
    suffix "awk=gawk -f"


## Customizing with Lua

    copy .nyagos "%USERPROFILE%\."
    notepad "%USERPROFILE%\.nyagos"

And please customize `%USERPROFILE%\.nyagos`

`.nyagos` is the configuration file written with Lua.
(Be careful that the filename starts with `.`(dot)
