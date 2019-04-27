package completion

import (
	"errors"
	"regexp"
	"strings"

	"github.com/zetamatta/nyagos/dos"
)

var rxUNCPattern1 = regexp.MustCompile(`^\\\\[^\\/]*$`)
var rxUNCPattern2 = regexp.MustCompile(`^(\\\\[^\\/]+)\\[^\\/]*$`)

func uncComplete(str string) ([]Element, error) {
	if rxUNCPattern1.MatchString(str) {
		server := strings.ToUpper(str)
		result := []Element{}
		err := dos.EachMachine(func(n *dos.NetResource) bool {
			server1 := n.RemoteName()
			if strings.HasPrefix(strings.ToUpper(server1), server) {
				result = append(result, Element1(server1+`\`))
			}
			return true
		})
		return result, err
	}
	if m := rxUNCPattern2.FindStringSubmatch(str); m != nil {
		server := m[1]
		path := strings.ToUpper(str)
		result := []Element{}
		err := dos.EachMachineNode(server, func(n *dos.NetResource) bool {
			remoteName := n.RemoteName()
			if strings.HasPrefix(strings.ToUpper(remoteName), path) {
				result = append(result, Element1(remoteName+`\`))
			}
			return true
		})
		return result, err
	}
	return nil, errors.New("not support")
}
