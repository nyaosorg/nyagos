@setlocal
@SET PROMPT=$$$S

@pushd "%~dp0"

if exist "%~dp0Misc\version.cmd" call "%~dp0Misc\version.cmd"

if exist goarch.txt for /F %%I in (goarch.txt) do set "GOARCH=%%I"
if "%GOARCH%" == "" for /F "delims=/ tokens=2" %%I in ('go version') do set "GOARCH=%%I"

call :"%~1"
@popd
@endlocal
@exit /b

:""
        for /F %%I in ('git describe --tags') do set "X_VERSION=-X main.version=%%I"
        call :"build"
        @exit /b 0

:"debug"
        set "TAGS=-tags=debug"
        @exit /b

:"release"
        for /F %%I in (%~dp0Misc\version.txt) do set "VERSION=%%I"
        set "X_VERSION=-X main.version=%VERSION%"
        call :"build"
        @exit /b

:"build"
        call :"fmt"
        call :"goversioninfo"
        for /F %%I in ('dir /b /s /aa nyagos.d') do attrib -A "%%I" & if exist main\bindata.go del main\bindata.go
        if not exist main\bindata.go call :"bindata"
        for /F "delims=" %%V in ('git log -1 --date^=short --pretty^=format:"-X main.stamp=%%ad -X main.commit=%%H"') do go build -o nyagos.exe -ldflags "%%V %X_VERSION%" %TAGS% .\main
        @exit /b

:"fmt"
        for /F %%I IN ('dir /s /b /aa *.go') do go fmt "%%I" & attrib -A "%%I"
        @exit /b

:"status"
        nyagos -e "print(nyagos.version or 'Snapshot on '..nyagos.stamp)"
        @exit /b

:"clean"
        for %%I in (nyagos.exe main\nyagos.syso version.now main\bindata.go) do if exist %%I del %%I
        powershell "ls -R | ?{ $_ -match '\.go$' } | %%{ [System.IO.Path]::GetDirectoryName($_.FullName)} | Sort-Object | Get-Unique | %%{ Write-Host 'go clean on',$_ ;  pushd $_ ; go clean ; popd }"

:"sweep"
        for /R %%I in (*~ *.bak) do if exist %%I del %%I
        if exist main\main.exe del main\main.exe
        @exit /b

:"get"
        powershell "Get-ChildItem . -Recurse | ?{ $_.Extension -eq '.go' } | %%{  Get-Content $_.FullName | %%{ ($_ -replace '\s*//.*$','').Split()[-1] } | ?{ $_ -match 'github.com/' } } | Sort-Object | Get-Unique | %%{ Write-Host $_ ; go get -u $_ }"
        @exit /b

:getbindata
        go get "github.com/jteeuwen/go-bindata"
        pushd "%GOPATH%\src\github.com\jteeuwen\go-bindata\go-bindata"
        go build
        copy go-bindata.exe "%~dp0\."
        popd
        @exit /b

:"bindata"
        if not exist go-bindata.exe call :getbindata
        go-bindata.exe -o "main\bindata.go" "nyagos.d/..."
        @exit /b

:getgoversioninfo
        go get "github.com/josephspurrier/goversioninfo"
        pushd "%GOPATH%\src\github.com\josephspurrier\goversioninfo\cmd\goversioninfo"
        go build
        copy goversioninfo.exe "%~dp0\."
        popd
        @exit /b

:"goversioninfo"
        if not exist goversioninfo.exe call :getgoversioninfo
        powershell -ExecutionPolicy RemoteSigned -File "%~dp0main\makejson.ps1" > "%~dp0Misc\version.json"
        goversioninfo.exe -icon main\nyagos.ico -o main\nyagos.syso "%~dp0Misc\version.json"
        @exit /b

:"const"
        for /F %%I in ('dir /b /s makeconst.cmd') do pushd %%~dpI & call %%I & popd
        @exit /b

:"package"
        for /F %%I in ('nyagos -e "print(nyagos.version or (string.gsub(nyagos.stamp,[[/]],[[]])))"') do set VERSION=%%I
        zip -9 "nyagos-%VERSION%-%GOARCH%.zip" nyagos.exe lua53.dll nyagos.lua .nyagos makeicon.cmd nyagos.d\*.lua nyagos.d\catalog\*.lua license.txt readme_ja.md readme.md Doc\*.md
        @exit /b

:"install"
        if not "%2" == "" set "INSTALLDIR=%2" & echo @set "INSTALLDIR=%2" > "%~dp0Misc\version.cmd"
        if "%INSTALLDIR%" == "" (
            echo Please %0.cmd install PATH\TO\BIN, once
            @exit /b
        )
        if not exist "%INSTALLDIR%" (
            echo Please %0.cmd install EXIST\PATH\TO\BIN,  once
            @exit /b
        )

        robocopy nyagos.d "%INSTALLDIR%\nyagos.d" /E
        if not exist "%INSTALLDIR%\lua53.dll" copy lua53.dll "%INSTALLDIR%\."
        copy nyagos.lua "%INSTALLDIR%\."
        copy /-Y _nyagos "%INSTALLDIR%\."
        copy nyagos.exe "%INSTALLDIR%\."
        if errorlevel 1 (start "" "%~dpfx0" install_ & @exit /b)
        @exit /b

:"install_"
        taskkill /F /im nyagos.exe
        copy nyagos.exe "%INSTALLDIR%\."
        timeout /T 3
        @exit %ERRORLEVEL%

:"icon"
        makeicon.cmd
        @exit /b

:"help"
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
        @exit /b
