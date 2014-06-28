# Maintenance file for nmake.exe

build :
	go build

icon : 
	rsrc -ico nyagos.ico -o nyagos.syso

fmt:
	for /R $(MAKEDIR) %%I IN (*.go) do go fmt %%I

clean :
	del nyagos.exe

snapshot :
	zip -9 nyagos-%DATE:/=%.zip nyagos.cmd nyagos.exe nyagos_ja.txt readme.mkd
