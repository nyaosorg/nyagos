@setlocal
@set PROMPT=$G$S

@if exist "%~dp0version.cmd" call "%~dp0version.cmd"

@if not "%1" == "" goto %1

:build
        @set "VERSION=%DATE:/=%"
        @set X_VERSION=
        @goto _build

:release
        for /F %%I in (version.txt) do set VERSION=%%I
        set "X_VERSION=-X main.version %VERSION%"

:_build
        if not exist nyagos.syso windres --output-format=coff -o nyagos.syso nyagos.rc
        for /F %%V in ('git log -1 --pretty^=format:%%H') do go build -ldflags "-X main.stamp %DATE% -X main.commit %%V %X_VERSION%"
        goto end

:fmt
        for /R . %%I IN (*.go) do go fmt %%I
        for /R . %%I in (.*~) do del %%I
        goto end

:clean
        for %%I in (nyagos.exe nyagos.syso version.now) do if exist %%I del %%I
        goto end

:get
        go get github.com/mattn/go-runewidth
        go get github.com/shiena/ansicolor 
        go get github.com/atotto/clipboard       
        goto end

:package
        zip -9 "nyagos-%VERSION%%2.zip" nyagos.exe lua52.dll nyagos.lua nyagos_ja.md nyagos_en.md readme.md .nyagos nyagos.d\*.lua
        goto end

:install
        @echo off
        if not "%2" == "" set "INSTALLDIR=%2"
        if "%INSTALLDIR%" == "" (
            @echo Please %0.cmd install PATH\TO\BIN, once
            goto end
        )
        if exist "%INSTALLDIR%" (
            @start make.cmd install.
        ) else (
            @echo Please %0.cmd install EXIST\PATH\TO\BIN,  once
        )
        goto end

:install.
        @echo Please close NYAGOS.exe and hit ENTER.
        @pause
        robocopy nyagos.d "%INSTALLDIR%\nyagos.d" /E
        copy nyagos.exe "%INSTALLDIR%\."
        copy nyagos.lua "%INSTALLDIR%\."
        copy lua52.dll  "%INSTALLDIR%\."
        goto end

:upgrade
        for %%I in (mattn\go-runewidth shiena\ansicolor atotto\clipboard ) do (
            cd %GOPATH%\Src\github.com\%%I
            git pull origin master:master
            go build
        )
        goto end

:help
        @echo off
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
        ( echo set "VERSION=%VERSION%"
          echo set "X_VERSION=%X_VERSION%"
          echo set "INSTALLDIR=%INSTALLDIR%"
        ) > "%~dp0version.cmd"
