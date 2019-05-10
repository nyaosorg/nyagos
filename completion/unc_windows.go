package completion

import (
	"errors"
	"regexp"
	"strings"

	"github.com/zetamatta/nyagos/dos"
	//"github.com/zetamatta/go-outputdebug"
	//"time"
)

var rxUNCPattern1 = regexp.MustCompile(`^\\\\[^\\/]*$`)
var rxUNCPattern2 = regexp.MustCompile(`^(\\\\[^\\/]+)\\[^\\/]*$`)

type _ServerCache struct {
	dos.NetResource
	Path []string
}

var serverCache map[string]*_ServerCache

func getServerCache() map[string]*_ServerCache {
	if serverCache == nil {
		serverCache = make(map[string]*_ServerCache)
		dos.EachMachine(func(n *dos.NetResource) bool {
			serverCache[n.RemoteName()] = &_ServerCache{
				NetResource: *n,
			}
			return true
		})
	}
	return serverCache
}

func uncComplete(str string) ([]Element, error) {
	if rxUNCPattern1.MatchString(str) {
		//outputdebug.String(`start complete \\server:` + time.Now().String())
		server := strings.ToUpper(str)
		result := []Element{}
		for server1 := range getServerCache() {
			if strings.HasPrefix(strings.ToUpper(server1), server) {
				result = append(result, Element1(server1+`\`))
			}
		}
		//outputdebug.String(`end complete \\server:` + time.Now().String())
		return result, nil
	}
	if m := rxUNCPattern2.FindStringSubmatch(str); m != nil {
		//outputdebug.String(`start complete \\server\path:` + time.Now().String())
		server := m[1]
		result := []Element{}

		cache := getServerCache()
		root, ok := cache[server]
		if ok {
			if root.Path == nil {
				root.NetResource.Enum(func(n *dos.NetResource) bool {
					root.Path = append(root.Path, n.RemoteName())
					return true
				})
			}
			path := strings.ToUpper(str)
			for _, node := range root.Path {
				if strings.HasPrefix(strings.ToUpper(node), path) {
					result = append(result, Element1(node+`\`))
				}
			}
		}
		//outputdebug.String(`end complete \\server\path:` + time.Now().String())
		return result, nil
	}
	return nil, errors.New("not support")
}
