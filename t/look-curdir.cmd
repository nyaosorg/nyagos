@echo off

REM Test program for option
REM   --look-curdir-first
REM   --look-curdir-last
REM   --look-curdir-none

setlocal
pushd

cd "%TEMP%"
( echo @echo CALLED_CURDIR
  echo @exit /b ) > hogehoge.cmd

mkdir subdir
set "PATH=%TEMP%\subdir;%PATH%"

( echo @echo CALLED_PATH
  echo @exit /b ) > subdir\hogehoge.cmd

call :test --look-curdir-first CALLED_CURDIR
call :test --look-curdir-last  CALLED_PATH
del subdir\hogehoge.cmd
rmdir subdir
call :test --look-curdir-never NOT-FOUND

del hogehoge.cmd

popd
endlocal
exit /b

:test
    setlocal
    set FOUND=NOT-FOUND
    ( for /F %%I in ('%~dp0..\nyagos.exe %1 -c "cd %TEMP% ; hogehoge.cmd"') do set "FOUND=%%I" ) 2>nul
    IF "%FOUND%" == "%2" (
        echo [OK] %1: %FOUND%
    ) else (
        echo [NG] %1: expected %2 but got %FOUND%
    )
    endlocal
    exit /b
