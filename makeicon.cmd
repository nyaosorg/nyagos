@echo off
if exist version.cmd call version.cmd
if "%INSTALLDIR%" == "" set "INSTALLDIR=%~dp0"
if not exist "%INSTALLDIR%\nyagos.exe" set "INSTALLDIR=%~dp0"
if not "%2" == "" set "INSTALLDIR=%2"

powershell "$wsh=New-Object -Com WScript.Shell; $sc=$wsh.CreateShortcut([IO.Path]::Combine([Environment]::GetFolderPath('Desktop'),'NihongoYet Another GOing Shell.lnk')); $sc.TargetPath='%INSTALLDIR%\nyagos.exe'; $sc.WorkingDirectory='%USERPROFILE%' ; $sc.Description='Nihongo Yet Another GOing Shell'; $sc.HotKey='CTRL+ALT+N' ; $sc.Save();"
