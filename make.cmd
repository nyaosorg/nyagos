@setlocal
@set PROMPT=$G$S
@if not "%1" == "" goto %1

:build
        echo %DATE:/=% > version.now
        goto _build

:release
        for /F %%I in (version.txt) do set "VERSION=-X main.version %%I"
        copy version.txt version.now

:_build
        if not exist nyagos.syso windres --output-format=coff -o nyagos.syso nyagos.rc
        for /F %%V in ('git log -1 --pretty^=format:%%H') do go build -ldflags "-X main.stamp %DATE% -X main.commit %%V %VERSION%"
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
        for /F %%I in (version.now) do zip -9 "nyagos-%%I%2.zip" nyagos.exe lua52.dll nyagos.lua nyagos_ja.md nyagos_en.md readme.md .nyagos nyagos.d\*.lua
        goto end

:install
        @if exist installdir (
            @start make.cmd install.
        ) else (
            @echo Please do 'mklink /J installdir PATH\TO\BIN'.
        )
        @goto end

:install.
        @echo Please close NYAGOS.exe and hit ENTER.
        @pause
        robocopy nyagos.d .\installdir\nyagos.d /mir
        copy nyagos.exe .\installdir\.
        copy nyagos.lua .\installdir\.
        copy lua52.dll  .\installdir\.
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
        echo   %0 (no arguments) == %0 build
        echo   %0 build .... Build nyagos.exe as snapshot (ignore version.txt)
        echo   %0 release .. Build nyagos.exe as release  (see version.txt)
        echo   %0 fmt ...... Format all source files with 'go fmt'
        echo   %0 clean .... Delete nyagos.exe and nyagos.syso
        echo   %0 install .. Copy binaries to %~dp0installdir
        echo   %0 package .. Make the package zip-file (see version.now)
        echo   %0 get ...... Do 'go get' for each github library
        echo   %0 upgrade .. Do 'git pull' for each github library
        echo   %0 help ..... Print help
:end
