@echo off
setlocal
set "PROMPT=$ "
candle.exe | findstr "version.2"
if errorlevel 1 (
    set "WIX=wix3"
) else (
    set "WIX=wix2"
)
call :"%1"
endlocal
exit /b

:""
    call :"amd64"
    call :"386"
    exit /b

:"amd64"
    upx ..\cmd\amd64\nyagos.exe
    call :mkmsi nyagos-amd64-%WIX%
    exit /b

:"386"
    upx ..\cmd\386\nyagos.exe
    call :mkmsi nyagos-386-%WIX%
    exit /b

:mkmsi
    candle "%~1.wxs" || exit /b 1
    light  "%~1.wixobj" || exit /b 1
    del "%~1.wixobj"
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
    msiexec /i nyagos-amd64-%WIX%.msi
    exit /b 0

:"uninstall"
    msiexec /x nyagos-amd64-%WIX%.msi
    exit /b 0

:"clean"
    del *.msi *.wixobj *.bak *.wixpdb
    exit /b 0
