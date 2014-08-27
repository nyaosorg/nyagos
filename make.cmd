@setlocal
@set PROMPT=$G$S
@if not "%1" == "" goto %1

:build
        if not exist nyagos.syso windres --output-format=coff -o nyagos.syso nyagos.rc
        for /F %%V in ('git log -1 --pretty^=format:%%H') do go build -ldflags "-X main.stamp %DATE% -X main.commit %%V"
        goto end

:fmt
        for /R . %%I IN (*.go) do go fmt %%I
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
        zip -9 nyagos-%DATE:/=%.zip nyagos.exe lua52.dll nyagos.lua nyagos_ja.mkd readme.mkd .nyagos
        goto end

:install
        if exist installdir (
            start make.cmd install_
        ) else (
            @echo Please do 'mklink /J installdir PATH\TO\BIN' for example.
        )
        @goto end

:install_
        @echo Please close NYAGOS.exe and hit ENTER.
        @pause
        copy nyagos.exe .\installdir\.
        copy nyagos.lua .\installdir\.
        copy lua52.dll  .\installdir\.
        goto end
:end
