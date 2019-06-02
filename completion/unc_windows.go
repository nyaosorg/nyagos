package completion

import (
	"errors"
	"regexp"
	"strings"

	"github.com/zetamatta/go-inline-animation"
	"github.com/zetamatta/nyagos/dos"
	//"github.com/zetamatta/go-outputdebug"
	//"time"
)

var rxUNCPattern1 = regexp.MustCompile(`^\\\\[^\\/]*$`)
var rxUNCPattern2 = regexp.MustCompile(`^(\\\\[^\\/]+)\\[^\\/]*$`)

var serverCache []string

func getServerCache() []string {
	if serverCache == nil {
		c := animation.Progress()
		defer c()
		serverCache = []string{}
		dos.EnumFileServer(func(n *dos.NetResource) bool {
			serverCache = append(serverCache, n.RemoteName())
			return true
		})
	}
	return serverCache
}

func hasServerCache() bool {
	return serverCache != nil
}

func uncComplete(str string, force bool) ([]Element, error) {
	if rxUNCPattern1.MatchString(str) {
		if !force && !hasServerCache() {
			return nil, ErrAskRetry
		}
		//outputdebug.String(`start complete \\server:` + time.Now().String())
		server := strings.ToUpper(str)
		result := []Element{}
		for _, server1 := range getServerCache() {
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

		if fs, err := dos.NewFileServer(server); err == nil {
			paths := []string{}
			fs.Enum(func(n *dos.NetResource) bool {
				paths = append(paths, n.RemoteName())
				return true
			})
			path := strings.ToUpper(str)
			for _, node := range paths {
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
