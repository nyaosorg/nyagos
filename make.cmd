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
        call :fmt
        powershell -ExecutionPolicy RemoteSigned "cd '%~dp0main' ; . '%~dp0main\makesyso.ps1'"
        for /F "delims=" %%V in ('git log -1 --date^=short --pretty^=format:"-X main.stamp=%%ad -X main.commit=%%H"') do go build -o nyagos.exe -ldflags "%%V %X_VERSION%" .\main
        goto end

:fmt
        for /F %%I IN ('dir /s /b /aa *.go') do go fmt "%%I" & attrib -A "%%I"
        exit /b

:status
        nyagos -e "print(nyagos.version or 'Snapshot on '..nyagos.stamp)"
        goto end

:clean
        for %%I in (nyagos.exe nyagos.syso version.now) do if exist %%I del %%I
        powershell "ls -R | ?{ $_ -match '\.go$' } | %%{ [System.IO.Path]::GetDirectoryName($_.FullName)} | Sort-Object | Get-Unique | %%{ Write-Host 'go clean on',$_ ;  pushd $_ ; go clean ; popd }"
        for /R %%I in (*~ *.bak) do if exist %%I del %%I
        goto end

:sweep
        for /R %%I in (*~) do del %%I
        goto end

:get
        powershell "findstr /S github.com/ *.go | %%{ $_.Split()[-1] } | ?{ $_ -match 'github.com' } | Sort-Object | Get-Unique | %%{ Write-Host 'go get',$_ ; go get -u $_ }"
        goto end

:package
        for /F %%I in ('nyagos -e "print(nyagos.version or (string.gsub(nyagos.stamp,[[/]],[[]])))"') do set VERSION=%%I
        zip -9 "nyagos-%VERSION%-%GOARCH%.zip" nyagos.exe lua53.dll nyole.dll nyagos.lua .nyagos lnk.js makeicon.cmd nyagos.d\*.lua nyagos.d\catalog\*.lua license.txt readme_ja.md readme.md Doc\*.md
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
