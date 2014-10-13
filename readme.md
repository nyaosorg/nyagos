The Nihongo Yet Another GOing Shell
===================================

The binary files can be downloaded on [Release](https://github.com/zetamatta/nyagos/releases).

* [English Document](./nyagos_en.md)
* [Japanese Document](./nyagos_ja.md)

Required software 
-----------------

* [go1.3.3 windows/386](http://golang.org)
* [Mingw-Gcc 4.8.1-4](http://mingw.org)
* [LuaBinaries 5.2.3 for Win32 and MinGW](http://luabinaries.sourceforge.net/index.html)

How to Build
------------

On `%GOPATH%` folder,

    git clone https://github.com/zetamatta/nyagos nyagos
    cd nyagos\lua
    unzip PATH\TO\lua-5.2.3_Win32_dllw4_lib.zip 
    copy lua52.dll ..\..
    cd ..
    make.cmd get
    make.cmd

Special Thanks
--------------

* [nocd5](https://github.com/nocd5/)
* [mattn](https://github.com/mattn/)
* [shiena](https://github.com/shiena/)
* [atotto](https://github.com/atotto/)

License
-------

You can use, copy and modify under the New BSD License.

Author
------

[zetamatta](https://github.com/zetamatta/)
