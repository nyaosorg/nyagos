pushd "%~dp0"Src 
call make.cmd %*
if "%1" == "" copy nyagos.exe ..
popd
