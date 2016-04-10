@echo off
if exist version.cmd call version.cmd
if "%INSTALLDIR%" == "" set "INSTALLDIR=%~dp0"
if not exist "%INSTALLDIR%\nyagos.exe" set "INSTALLDIR=%~dp0"
if not "%2" == "" set "INSTALLDIR=%2"

for /F "usebackq" %%I in (`powershell "[Environment]::GetFolderPath('Desktop')"`) do set "DESKTOP=%%I"
cscript //nologo lnk.js "%INSTALLDIR%\nyagos.exe" "%DESKTOP%\Nihongo Yet Another GOing Shell.lnk" "WorkingDirectory=%USERPROFILE%" "HotKey=CTRL+ALT+N" "Description=Nihongo Yet Another GOing Shell"
