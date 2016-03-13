@setlocal
@SET PROMPT=$$$S

pushd "%~dp0"

if exist "%~dp0Misc\version.cmd" call "%~dp0Misc\version.cmd"

if exist goarch.txt for /F %%I in (goarch.txt) do set "GOARCH=%%I"
if "%GOARCH%" == "" for /F "delims=/ tokens=2" %%I in ('go version') do set "GOARCH=%%I"

if not "%1" == "" goto %1

goto build

:release
        for /F %%I in (%~dp0Misc\version.txt) do set "VERSION=%%I"
        set "X_VERSION=-X main.version=%VERSION%"

:build
        for /F %%I IN ('dir /s /b /aa *.go') do go fmt "%%I" & attrib -A "%%I"
        powershell -ExecutionPolicy RemoteSigned "cd '%~dp0main' ; . '%~dp0main\makesyso.ps1'"
        for /F "delims=" %%V in ('git log -1 --date^=short --pretty^=format:"-X main.stamp=%%ad -X main.commit=%%H"') do go build -o nyagos.exe -ldflags "%%V %X_VERSION%" .\main
        goto end

:fmt
        for /F %%I IN ('dir /s /b /aa *.go') do go fmt "%%I" & attrib -A "%%I"
        goto end

:status
        nyagos -e "print(nyagos.version or 'Snapshot on '..nyagos.stamp)"
        goto end

:clean
        for %%I in (nyagos.exe nyagos.syso version.now) do if exist %%I del %%I
        for %%I in (alias commands completion conio dos history interpreter lua main) do (cd "%%I" & go clean & cd ..)
        for /R %%I in (*~ *.bak) do if exist %%I del %%I
        goto end

:sweep
        for /R %%I in (*~) do del %%I
        goto end

:get
        go get -u github.com/mattn/go-runewidth
        go get -u github.com/shiena/ansicolor 
        go get -u github.com/atotto/clipboard       
        go get -u github.com/zetamatta/goutputdebugstring
        goto end

:package
        for /F %%I in ('nyagos -e "print(nyagos.version or (string.gsub(nyagos.stamp,[[/]],[[]])))"') do set VERSION=%%I
        zip -9 "nyagos-%VERSION%-%GOARCH%.zip" nyagos.exe lua53.dll nyole.dll nyagos.lua .nyagos specialfolders.js lnk.js makeicon.cmd nyagos.d\*.lua catalog.d\*.lua license.txt readme_ja.md readme.md Doc\*.md
        goto end

:install
        if not "%2" == "" set "INSTALLDIR=%2"
        if "%INSTALLDIR%" == "" (
            echo Please %0.cmd install PATH\TO\BIN, once
            goto end
        )
        if not exist "%INSTALLDIR%" (
            echo Please %0.cmd install EXIST\PATH\TO\BIN,  once
            goto end
        )
        start %~0 install.
        goto end

:install.
        robocopy nyagos.d "%INSTALLDIR%\nyagos.d" /E
        taskkill /F /im nyagos.exe
        copy nyagos.exe "%INSTALLDIR%\."
        copy nyagos.lua "%INSTALLDIR%\."
        copy nyole.dll "%INSTALLDIR%\."
        if not exist "%INSTALLDIR%\lua53.dll" copy lua53.dll "%INSTALLDIR%\."
        goto end

:install_catalog
        if not "%2" == "" set "INSTALLDIR=%2"
        if "%INSTALLDIR%" == "" (
            echo Please %0.cmd %1 PATH\TO\BIN, once
            goto end
        )
        robocopy catalog.d "%INSTALLDIR%\catalog.d" /E
        goto end

:icon
        makeicon.cmd
        goto end

:help
        echo Usage for make.cmd
        echo  %0          : Equals to '%0 build'
        echo  %0 build    : Build nyagos.exe as snapshot (ignore version.txt)
        echo  %0 release  : Build nyagos.exe as release  (see version.txt)
        echo  %0 fmt      : Format all source files with 'go fmt'
        echo  %0 clean    : Delete nyagos.exe and nyagos.syso
        echo  %0 package  : Make the package zip-file
        echo  %0 get      : Do 'go get' for each github library
        echo  %0 upgrade  : Do 'git pull' for each github library
        echo  %0 help     : Print help
        echo  %0 install INSTALLDIR 
        echo     : Copy binaries to INSTALLDIR
        echo  %0 install  
        echo     : Copy binaries to last INSTALLDIR
:end
        echo @set "INSTALLDIR=%INSTALLDIR%" > "%~dp0Misc\version.cmd"
        popd
