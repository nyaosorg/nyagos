@echo off
setlocal
rem
rem For https://github.com/zetamatta/nyagos/pull/247
rem To occur error, git checkout 52b113f002300cc73503d861182fa7ecd95ab757~1
rem

set GOGC=1
set "ME=%~0"
set "ME=%ME:\=\\%"
nyagos.exe -e "for i=1,10000 do for line in nyagos.lines('%ME%') do end end"
endlocal
