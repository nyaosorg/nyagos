@echo off
setlocal

if exist "%~dp0Misc\version.cmd" call "%~dp0Misc\version.cmd"

if not "%1" == "" goto %1

:build
        set "VERSION=%DATE:/=%"
        set X_VERSION=
        goto _build

:release
        for /F %%I in (%~dp0Misc\version.txt) do set "VERSION=%%I"
        set "X_VERSION=-X main.version %VERSION%"

:_build
        pushd "%~dp0main"
        for %%I in (windres.exe) do if not "%%~$PATH:I" == "" lua make_rc.lua | windres.exe --output-format=coff -o nyagos.syso
        for /F %%V in ('git log -1 --pretty^=format:%%H') do go build -o "%~dp0nyagos.exe" -ldflags "-X main.stamp %DATE% -X main.commit %%V %X_VERSION%"
        popd
        goto end

:status
        if exist %~dp0nyagos.exe for /F %%I in ('cscript //nologo %~dp0showver.js %~dp0nyagos.exe') do echo NYAGOS.EXE version='%%I'
        exit /b 0

:fmt
        for /R . %%I in (*~) do del %%I
        for /F %%I IN ('dir /s /b /aa *.go') do go fmt "%%I" & attrib -A "%%I"
        goto end

:clean
        for %%I in (nyagos.exe nyagos.syso version.now) do if exist %%I del %%I
        goto end

:sweep
        for /R %%I in (*~) do del %%I
        goto end

:get
        go get -u github.com/mattn/go-runewidth
        go get -u github.com/shiena/ansicolor 
        go get -u github.com/atotto/clipboard       
        goto end

:package
        zip -9 "nyagos-%VERSION%%2.zip" nyagos.exe lua53.dll lua.exe nyagos.lua .nyagos specialfolders.vbs lnk.vbs makeicon.cmd nyagos.d\*.lua readme.md
        zip -9j "nyagos-%VERSION%%2.zip" Doc\nyagos_*.md
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
        if not exist "%INSTALLDIR%\lua53.dll" copy lua53.dll "%INSTALLDIR%\."
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
        echo off
        ( echo @set "VERSION=%VERSION%"
          echo @set "X_VERSION=%X_VERSION%"
          echo @set "INSTALLDIR=%INSTALLDIR%"
        ) > "%~dp0Misc\version.cmd"
