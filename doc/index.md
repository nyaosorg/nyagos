English / [Japanese](./index_ja.md)

## Welcome to the hybrid commandline shell

NYAGOS - Nihongo Yet Another GOing Shell is a command-line shell that combines bash-like command-line editing with seamless support for Windows file system paths and batch files. It also allows extensive customization of shell behavior using the Lua language.

![demo-animation](./demo.gif)

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
* Support Japanese input method editor: [SKK] \(Simple Kana Kanji conversion program\) - [How To Setup][SKKSetUpEn]
* Support OS:
  * Windows 7, 8.1, 10, 11, WindowsServer 2008 or later
  * Linux (experimental)

[SKK]: https://ja.wikipedia.org/wiki/SKK
[SKKSetUpEn]: 10-SetupSKK_en.md

[Video by @emisjerry](https://www.youtube.com/watch?v=WsfIrBWwAh0)

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
/ nu8 <!-- https://github.com/nu8 -->
/ [tomato3713](https://github.com/tomato3713)
/ tGqmJHoJKqgK <!-- https://github.com/tGqmJHoJKqgK -->
/ [juggler999](https://github.com/juggler999)
/ [zztkm](https://github.com/zztkm)
/ [8exBCYJi5ATL](https://github.com/8exBCYJi5ATL)
/ [ousttrue](https://github.com/ousttrue)
/ [kgasawa](https://github.com/kgasawa)
/ [HAYASHI-Masayuki](https://github.com/HAYASHI-Masayuki)
/ [naoyaikeda](https://github.com/naoyaikeda)
/ [emisjerry](https://github.com/emisjerry)

Author
------

* [hymkor - HAYAMA Kaoru](https://github.com/hymkor) (a.k.a zetamatta)
