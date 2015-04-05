@echo off
if exist version.cmd call version.cmd
if "%INSTALLDIR%" == "" set "INSTALLDIR=%~dp0"
if not exist "%INSTALLDIR%\nyagos.exe" set "INSTALLDIR=%~dp0"
if not "%2" == "" set "INSTALLDIR=%2"

for /F %%I in ('cscript /nologo specialfolders.js desktop') do echo Desktop=%%I & cscript //nologo lnk.js "%INSTALLDIR%\nyagos.exe" "%%I\Nihongo Yet Another GOing Shell.lnk" "WorkingDirectory=%USERPROFILE%" "HotKey=CTRL+ALT+N" "Description=Nihongo Yet Another GOing Shell"
