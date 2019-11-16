package completion

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/zetamatta/go-inline-animation"
	"github.com/zetamatta/nyagos/dos"
)

var rxUNCPattern1 = regexp.MustCompile(`^\\\\[^\\/]*$`)
var rxUNCPattern2 = regexp.MustCompile(`^(\\\\[^\\/]+)\\[^\\/]*$`)

func getServers() []string {
	servers := []string{}
	dos.EnumFileServer(func(n *dos.NetResource) bool {
		servers = append(servers, n.RemoteName())
		return true
	})
	return servers
}

func getCachePath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	dir = filepath.Join(dir, "nyaos_org")
	os.MkdirAll(dir, 0666)
	return filepath.Join(dir, "computers.txt")
}

var serverCache []string

func getServerCache() []string {
	if serverCache == nil {
		end := animation.Progress()
		serverCache = getServers()
		if cachePath := getCachePath(); cachePath != "" {
			ioutil.WriteFile(cachePath, []byte(strings.Join(serverCache, "\n")), 0666)
		}
		end()
	}
	return serverCache
}

func hasServerCache() bool {
	if serverCache != nil {
		return true
	}
	cachePath := getCachePath()
	if cachePath == "" {
		return false
	}
	fd, err := os.Open(cachePath)
	if err != nil {
		return false
	}
	serverCache = []string{}
	for sc := bufio.NewScanner(fd); sc.Scan(); {
		serverCache = append(serverCache, sc.Text())
	}
	stat, err := fd.Stat()
	fd.Close()
	//outputdebug.String(stat.ModTime().String())
	//outputdebug.String(stat.ModTime().Add(time.Hour * 2).String())

	if err == nil && stat.ModTime().Add(time.Hour*2).Before(time.Now()) {
		go func() {
			//outputdebug.String("begin searching servers at " + time.Now().String())
			tmp := getServers()
			ioutil.WriteFile(cachePath, []byte(strings.Join(tmp, "\n")), 0666)
			serverCache = tmp
			//outputdebug.String("update cache at " + time.Now().String())
		}()
	}
	return true
}

func uncComplete(str string, force bool) ([]Element, error) {
	if rxUNCPattern1.MatchString(str) {
		if !force && !hasServerCache() {
			return nil, ErrAskRetry
		}
		server := strings.ToUpper(str)
		result := []Element{}
		for _, server1 := range getServerCache() {
			if strings.HasPrefix(strings.ToUpper(server1), server) {
				result = append(result, Element1(server1+`\`))
			}
		}
		return result, nil
	}
	if m := rxUNCPattern2.FindStringSubmatch(str); m != nil {
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
		return result, nil
	}
	return nil, errors.New("not support")
}
