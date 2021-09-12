//go:build !windows
// +build !windows

package etc

// for default icon
//go:generate sh -c "go run mkversioninfo.go version.txt version.txt < versioninfo.json > v.json && go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo -icon=nyagos.ico -o ../nyagos.syso v.json && rm v.json"
