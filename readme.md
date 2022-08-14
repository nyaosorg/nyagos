[![Build status](https://ci.appveyor.com/api/projects/status/bh7866s6oasvchpj?svg=true)](https://ci.appveyor.com/project/zetamatta/nyagos)
[![GoDoc](https://godoc.org/github.com/nyaosorg/nyagos?status.svg)](https://godoc.org/github.com/nyaosorg/nyagos)
[![Go Report Card](https://goreportcard.com/badge/github.com/nyaosorg/nyagos)](https://goreportcard.com/report/github.com/nyaosorg/nyagos)
[![Github latest Releases](https://img.shields.io/github/downloads/nyaosorg/nyagos/latest/total.svg)](https://github.com/nyaosorg/nyagos/releases/latest)

The Nihongo Yet Another GOing Shell
===================================

English
/ [Japanese](./readme_ja.md)

NYAGOS is the commandline-shell written with the Programming Language GO and Lua.

![demo-animation](./demo.gif)

There are some shells in Windows compatible with ones in UNIX.
Most of them do not support Windows's traditional PATH-style
like `X:\DIR\FILE.EXT` that many applications require as arguments.

So, I created a new shell like below:

* UNIX-Like Shell
  * Keybinding
    * Features are bound to keys like Bash on default
    * Customized like
        * `nyagos.key.c_u = "KILL_WHOLE_LINE"` on %USERPROFILE%\\.nyagos ([Lua](https://github.com/yuin/gopher-lua))
    * A lua-functions can be bound to a key like
        * `nyagos.key.escape = function(this) nyagos.exec("start vim.exe") end`
  * History (Ctrl-P and !-mark)
  * Alias
    * like DOSKEY
        * `nyagos.alias["g++"]="g++.exe -std=gnu++17 $*"`
    * ones implemented by Lua functions
        * `nyagos.alias["lala"]=function(args) nyagos.exec("ls","-al",unpack(args)) end`
  * Custom completions
```lua
            nyagos.complete_for["go"] = function(args)
                if #args == 2 then
                    return {
                        "bug","doc","fmt","install","run","version",
                        "build","env","generate","list","test","vet",
                        "clean","fix","get","mod","tool" }
                else
                    return nil -- files completion
                end
            end
```
* Shell that follows the Windows' style like CMD.EXE
  * Windows' path format `C:\path\to\file` are able to be used.
  * Each drive has its own current directory.
  * `copy`,`move` and some dos-like built-in commands work.
  * No additional DLL are required.
  * Registry are not used.
* Color command-line
* Support Unicode. Windows unicode APIs are used.
  * Can paste unicode character on clipboard and edit them.
  * Unicode-literal %U+XXXX%
  * Prompt Macro $Uxxxx
* Built-in ls
  * color support (-o option)
  * print hard-link,symbolic-link and junction's target-path
* Support OS:
  * Windows 8.1 or later
  * Linux (experimental)

Download Binary
---------------

* https://github.com/nyaosorg/nyagos/releases

Contents
--------

### Release note

- [4.4.x](docs/release_note_en.md)
- [4.3.x](docs/history-4.3_en.md)
- [4.2.x](docs/history-4.2_en.md)
- [4.1.x](docs/history-4.1_en.md)
- [4.0.x](docs/history-4.0_en.md)

### Documents

1. [Install](docs/01-Install_en.md)
2. [Option for NYAGOS](docs/02-Options_en.md)
3. [Editor](docs/03-Readline_en.md)
4. [Built-in commands](docs/04-Commands_en.md)
5. [What is done on the Startup](docs/05-Startup_en.md)
6. [Substitution](docs/06-Substitution_en.md)
7. [Lua functions extenteded by NYAGOS](docs/07-LuaFunctions_en.md)
8. [Uninstall](docs/08-Uninstall_en.md)
9. [How To Build](docs/09-Build_en.md)

License
-------

You can use, copy and modify under the New BSD License.

Acknowledgement
---------------

[nocd5](https://github.com/nocd5)
/ [mattn](https://github.com/mattn)
/ [hattya](https://github.com/hattya)
/ [shiena](https://github.com/shiena)
/ [atotto](https://github.com/atotto)
/ [ironsand](https://github.com/ironsand)
/ [kardianos](https://github.com/kardianos)
/ [malys](https://github.com/malys)
/ [pine613](https://github.com/pine613)
/ [NSP-0123456](https://github.com/NSP-0123456)
/ [hokorobi](https://github.com/hokorobi)
/ [amuramatsu](https://github.com/amuramatsu)
/ [spiegel-im-spiegel](https://github.com/spiegel-im-spiegel)
/ [rururutan](https://github.com/rururutan/)
/ [hogewest](https://github.com/hogewest)
/ [cagechi](https://github.com/cagechi)
/ [Matsuyanagi](https://github.com/Matsuyanagi)
/ [Shougo](https://github.com/Shougo)
/ [orthographic-pedant](https://github.com/orthographic-pedant)
/ HABATA Katsuyuki
/ [hisomura](https://github.com/hisomura)
/ [tsuyoshicho](https://github.com/tsuyoshicho)
/ [rane-hs](https://github.com/rane-hs)
/ [hami-jp](https://github.com/hami-jp)
/ [3bch](https://github.com/3bch)
/ [AoiMoe](https://github.com/aoimoe)
/ [DeaR](https://github.com/DeaR)
/ [gracix](https://github.com/gracix)
/ [orz--](https://github.com/orz--)
/ [zkangaroo](https://github.com/zkangaroo)
/ [maskedw](https://github.com/maskedw)
/ [tyochiai](https://github.com/tyochiai)
/ [masamitsu-murase](https://github.com/masamitsu-murase)
/ [hazychill](https://github.com/hazychill)
/ [erw7](https://github.com/erw7)
/ [tignear](https://github.com/tignear)
/ [crile](https://github.com/crile)
/ [fushihara](https://github.com/fushihara)
/ [ChiyosukeF](https://twitter.com/ChiyosukeF)
/ [beepcap](https://twitter.com/beepcap)
/ [tostos5963](https://github.com/tostos5963)
/ [sambatriste](https://github.com/sambatriste)
/ [terepanda](https://github.com/terepanda)
/ [Takmg](https://github.com/Takmg)
/ [nu8](https://github.com/nu8)
/ [tomato3713](https://github.com/tomato3713)
/ [tGqmJHoJKqgK](https://github.com/tGqmJHoJKqgK)
/ [juggler999](https://github.com/juggler999)
/ [zztkm](https://github.com/zztkm)

Author
------

* [hymkor - HAYAMA Kaoru](https://github.com/hymkor) (a.k.a zetamatta)
