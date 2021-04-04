package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/windows"

	"github.com/mattn/go-isatty"

	"github.com/zetamatta/go-box/v2"
	"github.com/zetamatta/go-findfile"
	"github.com/zetamatta/go-windows-shortcut"
	"github.com/zetamatta/nyagos/nodos"
)

const (
	_           = iota
	O_STRIP_DIR = (1 << iota)
	O_LONG
	O_INDICATOR
	O_COLOR
	O_ALL
	O_TIME
	O_REVERSE
	O_RECURSIVE
	O_ONE
	O_HELP
	O_SIZESORT
	O_HUMAN
	O_NOT_RECURSIVE
	O_DEREFERENCE
)

type fileInfoT struct {
	name        string
	os.FileInfo // anonymous
}

const (
	ANSI_EXEC     = "\x1B[35;1m"
	ANSI_DIR      = "\x1B[32;1m"
	ANSI_NORM     = "\x1B[37;1m"
	ANSI_READONLY = "\x1B[33;1m"
	ANSI_HIDDEN   = "\x1B[34;1m"
	ANSI_END      = "\x1B[0m"
)

func chkCancel(ctx context.Context) error {
	if ctx != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
	return nil
}

func (this fileInfoT) Name() string { return this.name }

func putFlag(value, flag uint32, c string, out io.Writer) {
	if (value & flag) != 0 {
		io.WriteString(out, c)
	} else {
		io.WriteString(out, "-")
	}
}

func lsOneLong(folder string, status os.FileInfo, flag int, width int, out io.Writer) {
	indicator := " "
	prefix := ""
	postfix := ""
	if (flag & O_COLOR) != 0 {
		prefix = ANSI_NORM
		postfix = ANSI_END
	}
	if status.IsDir() {
		io.WriteString(out, "d")
		indicator = "/"
		if (flag & O_COLOR) != 0 {
			prefix = ANSI_DIR
			postfix = ANSI_END
		}
	} else {
		io.WriteString(out, "-")
	}
	mode := status.Mode()
	perm := mode.Perm()
	name := status.Name()
	attr := findfile.GetFileAttributes(status)

	putFlag(uint32(perm), 4, "r", out)
	if (perm & 2) > 0 {
		io.WriteString(out, "w")
	} else {
		if (flag & O_COLOR) != 0 {
			prefix = ANSI_READONLY
			postfix = ANSI_END
		}
		io.WriteString(out, "-")
	}
	if (perm & 1) > 0 {
		io.WriteString(out, "x")
	} else if nodos.IsExecutableSuffix(filepath.Ext(name)) {
		io.WriteString(out, "x")
		indicator = "*"
		if (flag & O_COLOR) != 0 {
			prefix = ANSI_EXEC
			postfix = ANSI_END
		}
	} else {
		io.WriteString(out, "-")
	}
	putFlag(attr, windows.FILE_ATTRIBUTE_ARCHIVE, "a", out)
	putFlag(attr, windows.FILE_ATTRIBUTE_SYSTEM, "s", out)
	putFlag(attr, windows.FILE_ATTRIBUTE_HIDDEN, "h", out)

	if (attr&windows.FILE_ATTRIBUTE_HIDDEN) != 0 &&
		(flag&O_COLOR) != 0 {
		prefix = ANSI_HIDDEN
		postfix = ANSI_END
	}
	if (flag & O_STRIP_DIR) > 0 {
		name = filepath.Base(name)
	}
	if (flag & O_HUMAN) != 0 {
		fmt.Fprintf(out, " %*s", width, formatByHumanize(status.Size()))
	} else {
		fmt.Fprintf(out, " %*d", width, status.Size())
	}
	stamp := status.ModTime()
	now := time.Now()
	halflastyear := now.AddDate(0, -6, 0)
	if stamp.After(halflastyear) && !stamp.After(now) {
		io.WriteString(out, stamp.Format(" Jan _2 15:04:05 "))
	} else {
		io.WriteString(out, stamp.Format(" Jan _2 2006     "))
	}
	io.WriteString(out, prefix)
	io.WriteString(out, name)
	io.WriteString(out, postfix)

	var linkTo string
	if (attr & windows.FILE_ATTRIBUTE_REPARSE_POINT) != 0 {
		var err error
		path := filepath.Join(folder, name)
		linkTo, err = os.Readlink(path)
		if err == nil && linkTo != path {
			path, err = filepath.Abs(path)
			if err == nil && linkTo != path {
				indicator = "@"
			} else {
				linkTo = ""
			}
		} else {
			linkTo = ""
		}
	}
	if (flag & O_INDICATOR) > 0 {
		io.WriteString(out, indicator)
	}
	if linkTo != "" {
		fmt.Fprintf(out, " -> %s", linkTo)
	}
	if strings.HasSuffix(name, ".lnk") {
		path := nodos.Join(folder, name)
		shortcut, workdir, err := shortcut.Read(path)
		if err == nil && shortcut != "" {
			fmt.Fprintf(out, " -> %s", shortcut)
			if workdir != "" {
				fmt.Fprintf(out, " (%s)", workdir)
			}
		}
	}
	io.WriteString(out, "\n")
}

func lsBox(ctx context.Context, folder string, nodes []os.FileInfo, flag int, out io.Writer) error {
	nodes_ := make([]string, len(nodes))
	for key, val := range nodes {
		prefix := ""
		postfix := ""
		if (flag & O_COLOR) != 0 {
			prefix = ANSI_NORM
			postfix = ANSI_END
		}
		indicator := ""
		if val.IsDir() {
			if (flag & O_COLOR) != 0 {
				prefix = ANSI_DIR
				postfix = ANSI_END
			}
			if (flag & O_INDICATOR) != 0 {
				indicator = "/"
			}
		}
		if (val.Mode().Perm() & 2) == 0 {
			if (flag & O_COLOR) != 0 {
				prefix = ANSI_READONLY
				postfix = ANSI_END
			}
		}
		if !val.IsDir() && nodos.IsExecutableSuffix(filepath.Ext(val.Name())) {
			if (flag & O_COLOR) != 0 {
				prefix = ANSI_EXEC
				postfix = ANSI_END
			}
			if (flag & O_INDICATOR) != 0 {
				indicator = "*"
			}
		}
		attr := findfile.GetFileAttributes(val)
		if (attr&windows.FILE_ATTRIBUTE_HIDDEN) != 0 &&
			(flag&O_COLOR) != 0 {
			prefix = ANSI_HIDDEN
			postfix = ANSI_END
		}
		if (attr&windows.FILE_ATTRIBUTE_REPARSE_POINT) != 0 &&
			(flag&O_INDICATOR) != 0 {
			indicator = "@"
		}
		nodes_[key] = prefix + val.Name() + postfix + indicator
	}
	if !box.Print(ctx, nodes_, out) {
		return ctx.Err()
	}
	return nil
}

func keta(n int64) int {
	count := 1
	for n >= 10 {
		count++
		n /= 10
	}
	return count
}

func lsLong(ctx context.Context, folder string, nodes []os.FileInfo, flag int, out io.Writer) error {
	var width int = 0
	if (flag & O_HUMAN) != 0 {
		for _, finfo := range nodes {
			width1 := len(formatByHumanize(finfo.Size()))
			if width1 > width {
				width = width1
			}
		}
	} else {
		size := int64(1)
		for _, finfo := range nodes {
			if finfo.Size() > size {
				size = finfo.Size()
			}
		}
		width = keta(size)
	}
	for _, finfo := range nodes {
		lsOneLong(folder, finfo, flag, width, out)
		if err := chkCancel(ctx); err != nil {
			return err
		}
	}
	return nil
}

func lsSimple(ctx context.Context, folder string, nodes []os.FileInfo, flag int, out io.Writer) error {
	for _, f := range nodes {
		io.WriteString(out, f.Name())
		if (flag & O_INDICATOR) != 0 {
			if attr := findfile.GetFileAttributes(f); (attr & windows.FILE_ATTRIBUTE_REPARSE_POINT) != 0 {
				io.WriteString(out, "@")
			} else if f.IsDir() {
				io.WriteString(out, "/")
			} else if nodos.IsExecutableSuffix(filepath.Ext(f.Name())) {
				io.WriteString(out, "*")
			}
		}
		fmt.Fprintln(out)
		if err := chkCancel(ctx); err != nil {
			return err
		}
	}
	return nil
}

type fileInfoCollection struct {
	flag  int
	nodes []os.FileInfo
}

func (this fileInfoCollection) Len() int {
	return len(this.nodes)
}

func (this fileInfoCollection) Less(i, j int) bool {
	var result bool
	if (this.flag & O_TIME) != 0 {
		result = this.nodes[i].ModTime().After(this.nodes[j].ModTime())
		if !result && !this.nodes[i].ModTime().Before(this.nodes[j].ModTime()) {
			result = (this.nodes[i].Name() < this.nodes[j].Name())
		}
	} else if (this.flag & O_SIZESORT) != 0 {
		diff := this.nodes[i].Size() - this.nodes[j].Size()
		if diff != 0 {
			result = (diff < 0)
		} else {
			result = (this.nodes[i].Name() < this.nodes[j].Name())
		}
	} else {
		result = (this.nodes[i].Name() < this.nodes[j].Name())
	}
	if (this.flag & O_REVERSE) != 0 {
		result = !result
	}
	return result
}
func (this fileInfoCollection) Swap(i, j int) {
	this.nodes[i], this.nodes[j] = this.nodes[j], this.nodes[i]
}

func lsFolder(ctx context.Context, folder string, flag int, out io.Writer) error {
	var folder_ string
	if rxDriveOnly.MatchString(folder) {
		folder_ = folder + "."
	} else {
		folder_ = folder
	}
	nodesArray := fileInfoCollection{flag: flag}
	var folders []string = nil
	if (flag & O_RECURSIVE) != 0 {
		folders = make([]string, 0)
	}
	tmp := make([]os.FileInfo, 0)

	var wildcard string
	if folder == "" {
		wildcard = "*"
	} else {
		wildcard = nodos.Join(folder, "*")
	}
	_ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	var canceled error
	findfile.Walk(wildcard, func(f *findfile.FileInfo) bool {
		if err := chkCancel(_ctx); err != nil {
			canceled = err
			return false
		}
		if (flag & O_ALL) == 0 {
			if strings.HasPrefix(f.Name(), ".") {
				return true
			}
			attr := findfile.GetFileAttributes(f)
			if (attr & windows.FILE_ATTRIBUTE_HIDDEN) != 0 {
				return true
			}
		}
		if f.IsDir() && folders != nil && f.Name() != "." && f.Name() != ".." {
			folders = append(folders, f.Name())
		} else {
			tmp = append(tmp, f)
		}
		return true
	})
	if canceled != nil {
		cancel()
		return canceled
	}
	nodesArray.nodes = tmp
	sort.Sort(nodesArray)
	var err error
	if (flag & O_LONG) != 0 {
		err = lsLong(_ctx, folder_, nodesArray.nodes, O_STRIP_DIR|flag, out)
	} else if (flag & O_ONE) != 0 {
		err = lsSimple(_ctx, folder_, nodesArray.nodes, O_STRIP_DIR|flag, out)
	} else {
		err = lsBox(_ctx, folder_, nodesArray.nodes, O_STRIP_DIR|flag, out)
	}
	cancel()
	if err != nil {
		return err
	}
	if folders != nil && len(folders) > 0 {
		for _, f1 := range folders {
			if err := chkCancel(ctx); err != nil {
				return err
			}
			f1fullpath := nodos.Join(folder, f1)
			fmt.Fprintf(out, "\n%s:\n", f1fullpath)
			if err := lsFolder(ctx, f1fullpath, flag, out); err != nil {
				return err
			}
		}
	}
	return nil
}

var rxDriveOnly = regexp.MustCompile("^[a-zA-Z]:$")

func lsCore(ctx context.Context, paths []string, flag int, out io.Writer, errout io.Writer) error {
	if len(paths) <= 0 {
		return lsFolder(ctx, "", flag, out)
	}
	dirs := make([]string, 0)
	printCount := 0
	files := make([]os.FileInfo, 0)
	for _, name := range paths {
		if err := chkCancel(ctx); err != nil {
			return err
		}
		var nameStat string
		if rxDriveOnly.MatchString(name) {
			nameStat = name + "."
		} else {
			nameStat = name
		}
		var status os.FileInfo
		var err error
		if (flag & O_DEREFERENCE) != 0 {
			status, err = os.Stat(nameStat)
		} else {
			status, err = os.Lstat(nameStat)
		}
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(errout, "ls: %s not exist.\n", nameStat)
			} else if _, ok := err.(*os.PathError); ok {
				fmt.Fprintf(errout, "ls: %s: Path Error.\n", nameStat)
			} else {
				fmt.Fprintf(errout, "ls: %s\n", err.Error())
			}
			continue
		} else if status.IsDir() && (flag&O_NOT_RECURSIVE) == 0 {
			dirs = append(dirs, name)
		} else {
			files = append(files, &fileInfoT{filepath.Clean(name), status})
		}
	}
	if len(files) > 0 {
		nodesArray := fileInfoCollection{flag: flag, nodes: files}
		sort.Sort(nodesArray)
		var err error
		if (flag & O_LONG) != 0 {
			err = lsLong(ctx, ".", files, flag, out)
		} else if (flag & O_ONE) != 0 {
			err = lsSimple(ctx, ".", files, flag, out)
		} else {
			err = lsBox(ctx, ".", files, flag, out)
		}
		if err != nil {
			return err
		}
		printCount = len(files)
	}
	for _, name := range dirs {
		if len(paths) > 1 {
			if printCount > 0 {
				io.WriteString(out, "\n")
			}
			fmt.Fprintf(out, "%s:\n", name)
		}
		err := lsFolder(ctx, name, flag, out)
		if err != nil {
			return err
		}
		printCount++
	}
	return nil
}

var option = map[rune](func(*int) error){
	'l': func(flag *int) error {
		*flag |= O_LONG
		return nil
	},
	'F': func(flag *int) error {
		*flag |= O_INDICATOR
		return nil
	},
	'o': func(flag *int) error {
		*flag |= O_COLOR
		return nil
	},
	'a': func(flag *int) error {
		*flag |= O_ALL
		return nil
	},
	't': func(flag *int) error {
		*flag |= O_TIME
		return nil
	},
	'r': func(flag *int) error {
		*flag |= O_REVERSE
		return nil
	},
	'R': func(flag *int) error {
		*flag |= O_RECURSIVE
		return nil
	},
	'1': func(flag *int) error {
		*flag |= O_ONE
		return nil
	},
	'h': func(flag *int) error {
		*flag |= O_HUMAN
		return nil
	},
	'?': func(flag *int) error {
		*flag |= O_HELP
		return nil
	},
	'S': func(flag *int) error {
		*flag |= O_SIZESORT
		return nil
	},
	'd': func(flag *int) error {
		*flag |= O_NOT_RECURSIVE
		return nil
	},
	'L': func(flag *int) error {
		*flag |= O_DEREFERENCE
		return nil
	},
}

// 存在しないオプションに関するエラー
type OptionError struct {
	Option rune
}

func (this OptionError) Error() string {
	return fmt.Sprintf("-%c: No such option", this.Option)
}

func cmdLs(ctx context.Context, cmd Param) (int, error) {
	flag := 0
	paths := make([]string, 0)
	for _, arg := range cmd.Args()[1:] {
		if strings.HasPrefix(arg, "-") {
			for _, o := range arg[1:] {
				setter, ok := option[o]
				if !ok {
					return 1, OptionError{Option: o}
				}
				if err := setter(&flag); err != nil {
					return 1, err
				}
			}
		} else {
			paths = append(paths, arg)
		}
	}
	if (flag & O_HELP) != 0 {
		var message strings.Builder
		message.WriteString("Usage: ls [-")
		for optKey := range option {
			message.WriteRune(optKey)
		}
		message.WriteString("] [PATH(s)]...")
		return 1, errors.New(message.String())
	}

	out := cmd.Out()
	err := cmd.Err()

	if file, ok := out.(*os.File); ok && !isatty.IsTerminal(file.Fd()) {
		flag |= O_ONE
	}

	// cmd.Term() is colorableTerminal which is not fast.
	if (flag & O_COLOR) == 0 {
		_out := bufio.NewWriter(cmd.Out())
		defer _out.Flush()
		out = _out
	} else if out == os.Stdout {
		_out := bufio.NewWriter(cmd.Term())
		defer _out.Flush()
		out = _out
	}
	if (flag & O_COLOR) != 0 {
		io.WriteString(out, ANSI_END)
	}
	return 0, lsCore(ctx, paths, flag, out, err)
}

// vim:set fenc=utf8 ts=4 sw=4 noet:
