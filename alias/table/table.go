package aliasTable

var Table = map[string]string{
	"assoc":  "%COMSPEC% /c assoc",
	"attrib": "%COMSPEC% /c attrib",
	"copy":   "%COMSPEC% /c copy",
	"del":    "%COMSPEC% /c del",
	"dir":    "%COMSPEC% /c dir",
	"for":    "%COMSPEC% /c for",
	"md":     "%COMSPEC% /c md",
	"mkdir":  "%COMSPEC% /c mkdir",
	"mklink": "%COMSPEC% /c mklink",
	"move":   "%COMSPEC% /c move",
	"open":   "%COMSPEC% /c for %I in ($*) do @start \"%I\"",
	"rd":     "%COMSPEC% /c rd",
	"ren":    "%COMSPEC% /c ren",
	"rename": "%COMSPEC% /c rename",
	"rmdir":  "%COMSPEC% /c rmdir",
	"start":  "%COMSPEC% /c start",
	"type":   "%COMSPEC% /c type",
}
