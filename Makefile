# Maintenance file for nmake.exe

build : nyagos.syso
	go build -ldflags "-X main.version [SNAPSHOT-%DATE:/=%]"

nyagos.syso : 
	windres --output-format=coff -o nyagos.syso nyagos.rc

fmt:
	for /R $(MAKEDIR) %%I IN (*.go) do @go fmt %%I

ver:
	yShowVer nyagos.exe $(HOME)\bin\nyagos.exe

clean :
	if exist nyagos.exe del nyagos.exe

snapshot :
	zip -9 nyagos-%DATE:/=%.zip nyagos.exe lua52.dll nyagos.lua nyagos_ja.mkd readme.mkd .nyagos
