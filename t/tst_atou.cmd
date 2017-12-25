@more +1 "%~0" | "%~dp0..\nyagos.exe" - 2>nul & exit /b
lua_e "nyagos.write( nyagos.atou('SHIFTJIS\x95\xB6\x8E\x9A\x97\xF1') , '\n' )"
