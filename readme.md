The Nihongo Yet Another GOing Shell
===================================

English / [Japanese](./readme_ja.md)

NYAGOS is the commandline-shell for Windows written with the
Programming Language GO and Lua.

* UNIX-Like Shell
  * History (Ctrl-P and !-mark)
  * Alias
  * Filename/Command-name completion
* Support UNICODE
  * Can paste unicode charactor on clipboard and edit them.
  * Unicode-literal %U+XXXX%
  * Prompt Macro $Uxxxx
* Built-in ls
  * color support (-o option)
  * indicate junction-mark as @
* Customizing with Lua
  * built-in command written with Lua
  * filtering command-line
  * useful functions: ANSI-String & UTF8 convert , eval and so on.

Install
-------

The binary files can be downloaded on [Release](https://github.com/zetamatta/nyagos/releases).

    mkdir PATH\TO\INSTALLDIR
    cd PATH\TO\INSTALLDIR
    unzip PATH\TO\DOWNLOADDIR\nyagos-****.zip
    makeicon.cmd

The batchfile: `makeicon.cmd` makes icon on your desktop.

* [English Document](Doc/nyagos_en.md)
* [Japanese Document](Doc/nyagos_ja.md)

Uninstall
---------

Remove unzipped files and `%APPDATA%\NYAOS.ORG` and icon on the desktop.
NYAGOS.exe writes nothing on registry.

Build
-----

These sortwares are required.

* [go1.4.2 windows/386](http://golang.org)
* [LuaBinaries(5.3 for Win32)](http://sourceforge.net/projects/luabinaries/files/5.3/Tools%20Executables/lua-5.3_Win32_bin.zip)
* [NYOLE 0.0.0.5 or later](https://github.com/zetamatta/nyole/releases) (Optionally. Without it, some lua-extensions does not work, but nyagos.exe itself is available.)

On `%GOPATH%` folder,

    git clone https://github.com/zetamatta/nyagos nyagos
    cd nyagos

    unzip PATH\TO\lua-5.3_Win32_bin.zip lua53.dll

If you hava NYOLE.DLL,

    unzip PATH\TO\nyole-0.0.0.5.zip nyole.dll

Finally

    make.cmd get
    make.cmd
    make.cmd install INSTALLDIR

How to use make.cmd is shown with `make.cmd help`

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

Author
------

* HAYAMA\_Kaoru : [zetamatta](https://github.com/zetamatta) 

I dedicate this code to my late father.
