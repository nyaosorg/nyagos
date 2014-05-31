# Maintenance file for nmake.exe

build:
	go build nyagos.go

fmt:
	for /R $(MAKEDIR) %%I IN (*.go) do go fmt %%I

clean :
	del nyagos.exe

nightly :
	zip -9 nyagos-%DATE:/=%.zip nyagos.cmd nyagos.exe nyagos_ja.txt
