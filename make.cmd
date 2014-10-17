@setlocal
@set PROMPT=$G$S
@if not "%1" == "" goto %1

:build
        if not exist nyagos.syso windres --output-format=coff -o nyagos.syso nyagos.rc
        for /F %%V in ('git log -1 --pretty^=format:%%H') do go build -ldflags "-X main.stamp %DATE% -X main.commit %%V %VERSION%"
        goto end

:release
        for /F %%I in (version.txt) do set "VERSION=-X main.version %%I"
        goto build

:fmt
        for /R . %%I IN (*.go) do go fmt %%I
        for /R . %%I in (.*~) do del %%I
        goto end

:clean
        for %%I in (nyagos.exe nyagos.syso) do if exist %%I del %%I
        goto end

:get
        go get github.com/mattn/go-runewidth
        go get github.com/shiena/ansicolor 
        go get github.com/atotto/clipboard       
        goto end

:snapshot
        zip -9 nyagos-%DATE:/=%%2.zip nyagos.exe lua52.dll nyagos.lua nyagos_ja.md nyagos_en.md readme.md .nyagos nyagos.d\*.lua
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
        robocopy nyagos.d .\installdir\nyagos.d /E
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

:end
