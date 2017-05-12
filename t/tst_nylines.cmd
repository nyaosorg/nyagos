@more +1 "%~0" | "%~dp0..\\nyagos" & exit /b 0
lua_e "for line in nyagos.lines('tst_nylines.cmd') do print(line) end"
