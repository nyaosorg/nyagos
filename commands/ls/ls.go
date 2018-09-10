package ls

import (
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

	"github.com/dustin/go-humanize"
	"github.com/mattn/go-isatty"

	"github.com/zetamatta/go-box"
	"github.com/zetamatta/go-findfile"
	"github.com/zetamatta/nyagos/dos"
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

var ErrCtrlC = errors.New("C-c")

func isCancel(ctx context.Context) bool {
	if ctx != nil {
		select {
		case <-ctx.Done():
			return true
		default:
		}
	}
	return false
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
	} else if dos.IsExecutableSuffix(filepath.Ext(name)) {
		io.WriteString(out, "x")
		indicator = "*"
		if (flag & O_COLOR) != 0 {
			prefix = ANSI_EXEC
			postfix = ANSI_END
		}
	} else {
		io.WriteString(out, "-")
	}
	putFlag(attr, dos.FILE_ATTRIBUTE_ARCHIVE, "a", out)
	putFlag(attr, dos.FILE_ATTRIBUTE_SYSTEM, "s", out)
	putFlag(attr, dos.FILE_ATTRIBUTE_HIDDEN, "h", out)

	if (attr&dos.FILE_ATTRIBUTE_HIDDEN) != 0 &&
		(flag&O_COLOR) != 0 {
		prefix = ANSI_HIDDEN
		postfix = ANSI_END
	}
	if (flag & O_STRIP_DIR) > 0 {
		name = filepath.Base(name)
	}
	if (flag & O_HUMAN) != 0 {
		fmt.Fprintf(out, " %*s", width, humanize.Comma(status.Size()))
	} else {
		fmt.Fprintf(out, " %*d", width, status.Size())
	}
	stamp := status.ModTime()
	onelastyear := time.Now().AddDate(0, -11, 0)
	if stamp.After(onelastyear) {
		io.WriteString(out, stamp.Format(" Jan _2 15:04:05 "))
	} else {
		io.WriteString(out, stamp.Format(" Jan _2 2006     "))
	}
	io.WriteString(out, prefix)
	io.WriteString(out, name)
	io.WriteString(out, postfix)

	var linkTo string
	if (attr & dos.FILE_ATTRIBUTE_REPARSE_POINT) != 0 {
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
		path := dos.Join(folder, name)
		shortcut, workdir, err := dos.ReadShortcut(path)
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
		if !val.IsDir() && dos.IsExecutableSuffix(filepath.Ext(val.Name())) {
			if (flag & O_COLOR) != 0 {
				prefix = ANSI_EXEC
				postfix = ANSI_END
			}
			if (flag & O_INDICATOR) != 0 {
				indicator = "*"
			}
		}
		attr := findfile.GetFileAttributes(val)
		if (attr&dos.FILE_ATTRIBUTE_HIDDEN) != 0 &&
			(flag&O_COLOR) != 0 {
			prefix = ANSI_HIDDEN
			postfix = ANSI_END
		}
		if (attr&dos.FILE_ATTRIBUTE_REPARSE_POINT) != 0 &&
			(flag&O_INDICATOR) != 0 &&
			hasLink(folder, val.Name()) {

			indicator = "@"
		}
		nodes_[key] = prefix + val.Name() + postfix + indicator
	}
	if !box.Print(ctx, nodes_, out) {
		return ErrCtrlC
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
	size := int64(1)
	for _, finfo := range nodes {
		if finfo.Size() > size {
			size = finfo.Size()
		}
	}
	width := keta(size)
	if (flag & O_HUMAN) != 0 {
		width = width * 4 / 3
	}
	for _, finfo := range nodes {
		lsOneLong(folder, finfo, flag, width, out)
		if isCancel(ctx) {
			return ErrCtrlC
		}
	}
	return nil
}

func hasLink(folder, name string) bool {
	origpath := filepath.Join(folder, name)
	fullpath, err := os.Readlink(origpath)
	if err != nil || fullpath == "" {
		return false
	}
	origpath, err = filepath.Abs(origpath)
	return err == nil && origpath != fullpath
}

func lsSimple(ctx context.Context, folder string, nodes []os.FileInfo, flag int, out io.Writer) error {
	for _, f := range nodes {
		io.WriteString(out, f.Name())
		if (flag & O_INDICATOR) != 0 {
			if attr := findfile.GetFileAttributes(f); (attr&dos.FILE_ATTRIBUTE_REPARSE_POINT) != 0 && hasLink(folder, f.Name()) {
				io.WriteString(out, "@")
			} else if f.IsDir() {
				io.WriteString(out, "/")
			} else if dos.IsExecutableSuffix(filepath.Ext(f.Name())) {
				io.WriteString(out, "*")
			}
		}
		fmt.Fprintln(out)
		if isCancel(ctx) {
			return ErrCtrlC
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
		wildcard = dos.Join(folder, "*")
	}
	canceled := false
	findfile.Walk(wildcard, func(f *findfile.FileInfo) bool {
		if isCancel(ctx) {
			canceled = true
			return false
		}
		if (flag & O_ALL) == 0 {
			if strings.HasPrefix(f.Name(), ".") {
				return true
			}
			attr := findfile.GetFileAttributes(f)
			if (attr & dos.FILE_ATTRIBUTE_HIDDEN) != 0 {
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
	if canceled {
		return ErrCtrlC
	}
	nodesArray.nodes = tmp
	sort.Sort(nodesArray)
	var err error
	if (flag & O_LONG) != 0 {
		err = lsLong(ctx, folder_, nodesArray.nodes, O_STRIP_DIR|flag, out)
	} else if (flag & O_ONE) != 0 {
		err = lsSimple(ctx, folder_, nodesArray.nodes, O_STRIP_DIR|flag, out)
	} else {
		err = lsBox(ctx, folder_, nodesArray.nodes, O_STRIP_DIR|flag, out)
	}
	if err != nil {
		return err
	}
	if folders != nil && len(folders) > 0 {
		for _, f1 := range folders {
			if isCancel(ctx) {
				return ErrCtrlC
			}
			f1fullpath := dos.Join(folder, f1)
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
		if isCancel(ctx) {
			return ErrCtrlC
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

// ls 機能のエントリ:引数をオプションとパスに分離する
func Main(ctx context.Context, args []string, out io.Writer, err io.Writer) error {
	flag := 0
	paths := make([]string, 0)
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			for _, o := range arg[1:] {
				setter, ok := option[o]
				if !ok {
					return OptionError{Option: o}
				}
				if err := setter(&flag); err != nil {
					return err
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
		return errors.New(message.String())
	}
	if _, ok := out.(io.Closer); ok {
		// output is a not colorable instance.
		flag &^= O_COLOR
	}
	if (flag & O_COLOR) != 0 {
		io.WriteString(out, ANSI_END)
	}
	if file, ok := out.(*os.File); ok && !isatty.IsTerminal(file.Fd()) {
		flag |= O_ONE
	}
	return lsCore(ctx, paths, flag, out, err)
}

// vim:set fenc=utf8 ts=4 sw=4 noet:
