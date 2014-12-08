@echo off
if exist version.cmd call version.cmd
if "%INSTALLDIR%" == "" set "INSTALLDIR=%~dp0"
if not exist "%INSTALLDIR%\nyagos.exe" set "INSTALLDIR=%~dp0"
if not "%2" == "" set "INSTALLDIR=%2"

for /F %%I in ('cscript /nologo specialfolders.vbs desktop') do cscript /nologo lnk.vbs "%INSTALLDIR%\nyagos.exe" "%%I"
