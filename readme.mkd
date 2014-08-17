The Nihongo Yet Another GOing Shell
===================================

The binary files can be downloaded on [nyaos.org](http://www.nyaos.org/index.cgi?p=NYAGOS).

Required software 
-----------------

* [go1.3 windows/386](http://golang.org)
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
