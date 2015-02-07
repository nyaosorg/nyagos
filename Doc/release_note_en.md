NYAGOS 4.0.5\_0
================
* cd -N (N:digit): move the N-th previous directory.
* cd -h , cd ? : print current directory history.
* pwd -N (N:digit): print the N-th previous directory.
* %CD% , %ERRORLEVEL% are able to be completed.

Bugfix
------
* Lua-function 'include' didn't report error
* Command-name completion printed same-name in diffent directories.
* ReadLine sometimes left trash at replacing string
* Could not complete path including ./ #45
* Could not broken symbolic link with DEL #44
* Files contains '&' was not enclosed with ".." at completion.

Trivial fix
-----------
* make.cmd: add echo off a lot
* Let make.cmd on the top directory without arguments copy EXE top 

NYAGOS 4.0.4\_0
================
* Support Lua 5.3

NYAGOS 4.0.3\_2
===============
* Command-name completion supports alias and built-in commands.

BugFix
------
* `pwd` did not print correct UNC-Path (#37)
* `nyagos.gethistory( large-value )` crashes nyagos (#38 @1)
* `%APPDATA%/NYAOS_ORG/nyagos.history` did not be updated. (#39 @1)
* Crashed when `%APPDATA%/NYAOS_ORG/nyagos.history` did not exists or is empty.x (#40 @1)
* On French keyboard, keys shifted with AltGr could not be input. (#41)

@1 Probably these bugs appeared only on snapshot-build.
