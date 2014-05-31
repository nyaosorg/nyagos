@echo off
rem --- batchfile for test ---
set PROMPT=$e[34;40;1m$L$P$G$_$$ $e[37;1m
set OPTION=-a "ls=ls -oF" %1 %2 %3 %4 %5 %6 %7 %8 %9
if exist %~dp0nyagos.exe (
    %~dp0nyagos.exe %OPTION%
) else (
    go run nyagos.go %OPTION%
)
