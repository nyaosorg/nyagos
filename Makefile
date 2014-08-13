# Maintenance file for nmake.exe

build : nyagos.syso
	for /F %%V IN ('git log -1 --pretty^=format:%%H') DO go build -ldflags "-X main.version [SNAPSHOT-%DATE:/=%-%%V]"

nyagos.syso : 
	windres --output-format=coff -o nyagos.syso nyagos.rc

fmt:
	for /R $(MAKEDIR) %%I IN (*.go) do @go fmt %%I

ver:
	yShowVer nyagos.exe $(HOME)\bin\nyagos.exe

clean :
	for %%I in (nyagos.exe nyagos.syso) do if exist %%I del %%I

snapshot :
	zip -9 nyagos-%DATE:/=%.zip nyagos.exe lua52.dll nyagos.lua nyagos_ja.mkd readme.mkd .nyagos
