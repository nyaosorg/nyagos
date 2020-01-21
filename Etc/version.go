package etc

//go:generate cmd /c go run mkversioninfo.go version.txt version.txt < versioninfo.json > v.json && go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo -icon=nyagos.ico -o ..\nyagos.syso v.json && del v.json
