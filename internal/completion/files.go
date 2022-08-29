package completion

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/nyaosorg/go-windows-findfile"
	"github.com/nyaosorg/nyagos/internal/nodos"
)

const (
	stdSlash = string(os.PathSeparator)
	optSlash = "/"
)

var IncludeHidden = false

func ListUpFiles(ctx context.Context, ua UncCompletion, str string) ([]Element, error) {
	return listUpWithFilter(ctx, str, ua, func(*findfile.FileInfo) bool { return true })
}
func listUpDirs(ctx context.Context, ua UncCompletion, str string) ([]Element, error) {
	return listUpWithFilter(ctx, str, ua, func(fd *findfile.FileInfo) bool {
		return fd.IsDir() || strings.HasSuffix(strings.ToLower(fd.Name()), ".lnk")
	})
}

type UncCompletion int

const (
	DoNotUncCompletion UncCompletion = iota
	AskDoUncCompletion
	DoUncCompletion
)

var errAskRetry = errors.New("Complete Network Path ?")

func listUpWithFilter(ctx context.Context, str string, ua UncCompletion, filter func(*findfile.FileInfo) bool) ([]Element, error) {
	if ua != DoNotUncCompletion {
		if r, err := uncComplete(str, ua == DoUncCompletion); err == nil {
			return r, nil
		} else if err == errAskRetry {
			return nil, err
		}
	}
	orgSlash := stdSlash[0]
	if UseSlash {
		orgSlash = optSlash[0]
	}
	if pos := strings.IndexAny(str, stdSlash+optSlash); pos >= 0 {
		orgSlash = str[pos]
	}
	str = strings.Replace(strings.Replace(str, optSlash, stdSlash, -1), `"`, "", -1)
	directory := DirName(str)
	wildcard := join(findfile.ExpandEnv(directory), "*")

	// Drive letter
	cutprefix := 0
	if strings.HasPrefix(directory, stdSlash) {
		wd, _ := os.Getwd()
		directory = wd[0:2] + directory
		cutprefix = 2
	}
	commons := make([]Element, 0)
	STR := strings.ToUpper(str)
	var canceled error = nil
	err := findfile.WalkContext(ctx, wildcard, func(fd *findfile.FileInfo) bool {
		if err := checkTimeout(ctx); err != nil {
			canceled = ctx.Err()
			return false
		}
		if fd.Name() == "." || fd.Name() == ".." {
			return true
		}
		if !IncludeHidden && fd.IsHidden() {
			return true
		}
		if !filter(fd) {
			return true
		}
		listname := fd.Name()
		name := join(directory, fd.Name())
		if fd.IsDir() {
			name += stdSlash
			listname += optSlash
		}
		if cutprefix > 0 {
			name = name[2:]
		}
		nameUpr := strings.ToUpper(name)
		if strings.HasPrefix(nameUpr, STR) {
			if orgSlash != stdSlash[0] {
				name = strings.Replace(name, stdSlash, optSlash, -1)
			}
			element := Element2{name, listname}
			commons = append(commons, element)
		}
		return true
	})
	if canceled != nil {
		return commons, canceled
	}
	if os.IsNotExist(err) {
		return commons, nil
	}
	return commons, err
}

func join(dir, name string) string {
	return nodos.Join(dir, name)
}
