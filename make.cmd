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

:snapshot
	zip -9 nyagos-%DATE:/=%.zip nyagos.exe lua52.dll nyagos.lua nyagos_ja.mkd readme.mkd .nyagos

:end
