@set "PROMPT=$G"
for /F %%I in ('cscript /nologo specialfolders.vbs desktop') do cscript /nologo lnk.vbs "%~dp0\nyagos.exe" "%%I"
