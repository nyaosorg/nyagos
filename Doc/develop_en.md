- To call a batchfile, stop to use `/V:ON` for CMD.EXE

4.4.0\_beta
===========

- Support Linux (experimental)
- Fix the problem that current directories per drive were not inherited to child processes.
- Use the library "mattn/go-tty" instead of "zetamatta/go-getch"
- Stop using msvcrt.dll via "syscall" directly
- On linux, the filename NUL equals /dev/null
- Add lua-variable nyagos.goos
- (#341) Fix an unexpected space is inserted after wide characters
    * On Windows10, enable stdout virtual terminal processing always
    * If `git.exe push` disable virtual terminal processing, enable again.
- (#339) Fix that wildcard pattern `.??*` matches `..`
    * It requires github.com/zetamatta/go-findfile tagged 20181223-2
