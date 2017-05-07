@more +1 "%~0" | "%~dp0..\..\nyagos.exe" 2>nul & exit /b 0
lua_e "nyagos.alias.yyy=function() print 'OK' end"
yyy | more
