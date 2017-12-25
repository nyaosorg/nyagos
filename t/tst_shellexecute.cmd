@more +1 "%~0" | "%~dp0..\nyagos.exe" - 2>nul & exit /b
lua_e "nyagos.shellexecute('runas',[[c:\windows\system32\cmd.exe]],'/k dir','')"
