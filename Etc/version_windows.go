//go:build windows

package etc

//go:generate cmd /c go1.20.14.exe run mkversioninfo.go > v.json && go1.20.14.exe run github.com/hymkor/goversioninfo/cmd/goversioninfo@master -icon=nyagos.ico,nyagos32x32.ico,nyagos16x16.ico -o ..\nyagos.syso v.json && del v.json
