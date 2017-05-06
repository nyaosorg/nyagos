@more +1 "%~0" | "%~dp0..\nyagos.exe" - 2>nul & exit /b
lua_e "foo = 'OK'"
lua_e "nyagos.alias.xxx = function() assert(print(foo or 'NG')) end"
xxx | more
