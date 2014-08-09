# Maintenance file for nmake.exe

build :
	go build -ldflags "-X main.version [SNAPSHOT-%DATE:/=%]" nyagos.go

icon : 
	windres --output-format=coff -o nyagosico.syso nyagosico.rc

fmt:
	for /R $(MAKEDIR) %%I IN (*.go) do @go fmt %%I

ver:
	yShowVer nyagos.exe $(HOME)\bin\nyagos.exe

clean :
	if exist nyagos.exe del nyagos.exe

snapshot :
	zip -9 nyagos-%DATE:/=%.zip nyagos.exe lua52.dll nyagos.lua nyagos_ja.mkd readme.mkd .nyagos
