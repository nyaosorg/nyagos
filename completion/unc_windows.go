package completion

import (
	"errors"
	"regexp"
	"strings"

	"github.com/zetamatta/nyagos/dos"
)

var rxUNCPattern1 = regexp.MustCompile(`^\\\\[^\\/]*$`)
var rxUNCPattern2 = regexp.MustCompile(`^(\\\\[^\\/]+)\\[^\\/]*$`)

var cacheServerNames []string
var cacheNodeNames = map[string][]string{}

func uncComplete(str string) ([]Element, error) {
	if rxUNCPattern1.MatchString(str) {
		server := strings.ToUpper(str)
		if cacheServerNames == nil {
			dos.EachMachine(func(n *dos.NetResource) bool {
				cacheServerNames = append(cacheServerNames, n.RemoteName())
				return true
			})
		}
		result := []Element{}
		for _, server1 := range cacheServerNames {
			if strings.HasPrefix(strings.ToUpper(server1), server) {
				result = append(result, Element1(server1+`\`))
			}
		}
		return result, nil
	}
	if m := rxUNCPattern2.FindStringSubmatch(str); m != nil {
		server := m[1]
		r, ok := cacheNodeNames[server]
		if !ok {
			r := []string{}
			dos.EachMachineNode(server, func(n *dos.NetResource) bool {
				r = append(r, n.RemoteName())
				return true
			})
			cacheNodeNames[server] = r
		}
		result := []Element{}
		path := strings.ToUpper(str)
		for _, node := range r {
			if strings.HasPrefix(strings.ToUpper(node), path) {
				result = append(result, Element1(node+`\`))
			}
		}
		return result, nil
	}
	return nil, errors.New("not support")
}
