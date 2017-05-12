@more +1 "%~0" | "%~dp0..\nyagos.exe" - 2>nul & exit /b
true
if errorlevel 1 echo NG(1)
true
if not errorlevel 1 echo OK(2)
false
if errorlevel 1 echo OK(3)
false
if not errorlevel 1 echo NG(4)
