The Nihongo Yet Another GOing Shell
===================================

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

How to Install
--------------

The binary files can be downloaded on [Release](https://github.com/zetamatta/nyagos/releases).

    > mkdir PATH\TO\INSTALLDIR
    > cd PATH\TO\INSTALLDIR
    > unzip PATH\TO\DOWNLOADDIR\nyagos-****.zip
    > makeicon.cmd

The batchfile: `makeicon.cmd` makes icon on your desktop.

* [English Document](nyagos_en.md)
* [Japanese Document](nyagos_ja.md)

How to Build
------------

These sortwares are required.

* [go1.4 windows/386](http://golang.org)
* [Lua 5.3](http://www.lua.org)
* [MinGW](http://www.mingw.org) for building Lua 5.3

On `%GOPATH%` folder,

    git clone https://github.com/zetamatta/nyagos nyagos
    cd nyagos\src

When you have lua53.dll

    copy PATH\TO\lua53.dll lua\.

Otherwise,

    tar zxvf PATH\TO\lua-5.3.0.tar.gz
    cd lua-5.3.0
    mingw32-make.exe mingw
    copy src\lua53.dll ..\..
    cd ..

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

Author
------

* [zetamatta](https://github.com/zetamatta)

I dedicate this code to my late father.
