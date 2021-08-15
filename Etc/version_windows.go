// +build windows

package etc

// for default icon
//go:generate cmd /c go run mkversioninfo.go version.txt version.txt < versioninfo.json > v.json && go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo -icon=nyagos.ico -o ..\nyagos.syso v.json && del v.json

// for second icon (disabled)
////go:generate cmd /c go run mkversioninfo.go version.txt version.txt < versioninfo.json > v.json && go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo -icon=nyagos32x32.ico -icon=nyagos16x16.ico -o ..\nyagos.syso v.json && del v.json
