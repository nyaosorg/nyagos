@echo off
setlocal
set "PROMPT=$$ "
for /F %%I in (..\Etc\version.txt) do set "VER=%%I"
call :"%1"
endlocal
exit /b

:""
    call :"version"
    call :"amd64"
    call :"386"
    exit /b

:"amd64"
    call :mkmsi nyagos-amd64 "-arch x64"
    ren nyagos-amd64.msi nyagos-%VER%-windows-amd64.msi
    exit /b

:"386"
    call :mkmsi nyagos-386
    ren nyagos-386.msi nyagos-%VER%-windows-386.msi
    exit /b

:mkmsi
    @echo on
        candle %~2 "%~1.wxs" || exit /b 1
        light  "%~1.wixobj" || exit /b 1
        del "%~1.wixobj"
        del "%~1.wixpdb"
    @echo off
    exit /b 0

:"status"
    @echo off
    for %%I in ("C:\Program Files" "C:\Program Files (x86)") do (
        if exist "%%~I\NyaosOrg\nyagos.exe" (
            echo Found: %%~I\NyaosOrg
        ) else (
            echo Not Found: %%~I\NyaosOrg
        )
    )
    exit /b

:"install"
    @echo on
        msiexec /i nyagos-amd64.msi
    @echo off
    exit /b 0

:"uninstall"
    @echo on
    msiexec /x nyagos-amd64.msi
    @echo off
    exit /b 0

:"clean"
    @echo on
    del *.msi *.wixobj *.bak *.wixpdb
    @echo off
    exit /b 0

:"files"
    @echo on
    if not exist mkfwxi.exe (
        pushd mkfwxi
        go build
        move mkfwxi.exe ..
        popd
    )
    mkfwxi.exe -c files.wxi -r refer.wxi "dada523c-cb49-4e4e-a9cb-d509c50631b9" < files.txt
    @echo off
    exit /b 0

:"version"
    set "VER=%VER:_=.%"
    ( echo ^<?xml version="1.0" encoding="UTF-8"?^>
      echo ^<Include^>
      echo ^<?define Version=%VER% ?^>
      echo ^</Include^> ) > version.wxi
    exit /b
