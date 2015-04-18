NYAGOS 4.0.7\_2
===============
* Fix:on completion, / was always replaced to \ .(Thx @nocd5)
* Fix:nyagos.shellexecute() did not report some errors.
* To use COM on Lua scripts, include and use [NYOLE.DLL](https://github.com/zetamatta/nyole).

NYAGOS 4.0.7\_1
===============
* Set arg[..] in nyagos -e "LUA-CODE".

Bugfix
------
* In nyagos -f "LUA-FILE", arg[i]'s elements were shifted.

NYAGOS 4.0.7\_0
===============
* Support incremental-search(C-r)
* Add option -e "LUA-CODE" to nyagos.exe
* Set executable's property the version-number
* Change error-message when files do not exists like bash.

NYAGOS 4.0.6\_0
===============

* Add built-in command: pushd/popd/dirs
* Add the method boxprint(),firstword(),lastword() to nyagos.bindkey's first argument
* Update Document about uninstalling
* Add `nyagos.completion_hook`
* Sub-command completion for git, Subversion and Mercurial.

Bugfix
------
* Completion failed when 0001 is typed where 0001.txt and "0001 copy.txt" were.

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
