[![Build status](https://ci.appveyor.com/api/projects/status/bh7866s6oasvchpj?svg=true)](https://ci.appveyor.com/project/zetamatta/nyagos)
[![GoDoc](https://godoc.org/github.com/zetamatta/nyagos?status.svg)](https://godoc.org/github.com/zetamatta/nyagos)
[![Go Report Card](https://goreportcard.com/badge/github.com/zetamatta/nyagos)](https://goreportcard.com/report/github.com/zetamatta/nyagos)

The Nihongo Yet Another GOing Shell
===================================

English
/ [Japanese](./readme_ja.md)

NYAGOS is the commandline-shell for Windows written with the
Programming Language GO and Lua.

* UNIX-Like Shell
  * Keybinding like Emacs.
  * History (Ctrl-P and !-mark)
  * Alias
  * Filename/Command-name completion
* Support UNICODE
  * Can paste unicode character on clipboard and edit them.
  * Unicode-literal %U+XXXX%
  * Prompt Macro $Uxxxx
* Built-in ls
  * color support (-o option)
  * print hard-link,symbolic-link and junction's target-path
* Customizing with [GopherLua](https://github.com/yuin/gopher-lua)
  * built-in command written with Lua
  * filtering command-line
  * useful functions: ANSI-String & UTF8 convert , eval and so on.
  * Support to call COM(OLE)

Download Binary
---------------

* https://github.com/zetamatta/nyagos/releases

Contents
--------

### Release note and history 

- [Current Release note](Doc/release_note_en.md)
- [History ~4.0](Doc/history_4.0_en.md)
- [What is new since 4.1](Doc/since_4.1_en.md)

### Documents

1. [Install](Doc/01-Install_en.md)
2. [Option for NYAGOS.EXE](Doc/02-Options_en.md)
3. [Editor](Doc/03-Readline_en.md)
4. [Built-in commands](Doc/04-Commands_en.md)
5. [What is done on the Startup](Doc/05-Startup_en.md)
6. [Substitution](Doc/06-Substitution_en.md)
7. [Lua functions extenteded by NYAGOS](Doc/07-LuaFunctions_en.md)
8. [Uninstall](Doc/08-Uninstall_en.md)
9. [How To Build](Doc/09-Build_en.md)

License
-------

You can use, copy and modify under the New BSD License.

Acknowledgement
---------------

* [nocd5](https://github.com/nocd5)
* [mattn](https://github.com/mattn)
* [hattya](https://github.com/hattya)
* [shiena](https://github.com/shiena)
* [atotto](https://github.com/atotto)
* [ironsand](https://github.com/ironsand)
* [kardianos](https://github.com/kardianos)
* [malys](https://github.com/malys)
* [pine613](https://github.com/pine613)
* [NSP-0123456](https://github.com/NSP-0123456)
* [hokorobi](https://github.com/hokorobi)
* [amuramatsu](https://github.com/amuramatsu)
* [spiegel-im-spiegel](https://github.com/spiegel-im-spiegel)
* [rururutan](https://github.com/rururutan/)
* [hogewest](https://github.com/hogewest)
* [cagechi](https://github.com/cagechi)
* [Matsuyanagi](https://github.com/Matsuyanagi)
* [Shougo](https://github.com/Shougo)
* [orthographic-pedant](https://github.com/orthographic-pedant)
* HABATA Katsuyuki
* [hisomura](https://github.com/hisomura)
* [tsuyoshicho](https://github.com/tsuyoshicho)
* [rane-hs](https://github.com/rane-hs)
* [hami-jp](https://github.com/hami-jp)
* [3bch](https://github.com/3bch)
* [AoiMoe](https://github.com/aoimoe)
* [DeaR](https://github.com/DeaR)
* [gracix](https://github.com/gracix)
* [orz--](https://github.com/orz--)
* [zkangaroo](https://github.com/zkangaroo)
* [maskedw](https://github.com/maskedw)
* [tyochiai](https://github.com/tyochiai)
* [masamitsu-murase](https://github.com/masamitsu-murase)
* [hazychill](https://github.com/hazychill)
* [erw7](https://github.com/erw7)
* [tignear](https://github.com/tignear)
* [crile](https://github.com/crile)

Author
------

* [zetamatta](https://github.com/zetamatta)
