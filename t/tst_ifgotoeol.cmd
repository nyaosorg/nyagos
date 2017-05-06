@more +1 "%~0" | "%~dp0..\nyagos.exe" - 2>nul & exit /b
set FLAG=OK
if not 1 == 1 echo NG(1) ; set "FLAG=NG(2)"
echo %FLAG%
