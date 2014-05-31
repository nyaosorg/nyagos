@echo off
rem --- batchfile for test ---
set PROMPT=$e[34;40;1m$L$P$G$_$$ $e[37;1m
set OPTION=-a "ls=ls -oF"
if exist %~dp0nyagos.exe (
    %~dp0nyagos.exe %OPTION%
) else (
    go run nyagos.go %OPTION%
)
