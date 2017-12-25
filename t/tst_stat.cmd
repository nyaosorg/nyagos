@more +1 "%~0" | "%~dp0..\nyagos.exe" - 2>nul & exit /b
lua_e "s=nyagos.stat('tst_stat.cmd') ; print(s.name,s.size,s.isdir,s.mtime.year,s.mtime.month,s.mtime.day,s.mtime.hour,s.mtime.minute,s.mtime.second)"
ls -l tst_stat.cmd
