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
	"github.com/zetamatta/go-box"
	"github.com/zetamatta/go-findfile"

	"../../dos"
)

const (
	O_STRIP_DIR     = 1
	O_LONG          = 2
	O_INDICATOR     = 4
	O_COLOR         = 8
	O_ALL           = 16
	O_TIME          = 32
	O_REVERSE       = 64
	O_RECURSIVE     = 128
	O_ONE           = 256
	O_HELP          = 512
	O_SIZESORT      = 1024
	O_HUMAN         = 2048
	O_NOT_RECURSIVE = 4096
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

func (this fileInfoT) Name() string { return this.name }

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

	if (perm & 4) > 0 {
		io.WriteString(out, "r")
	} else {
		io.WriteString(out, "-")
	}
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
		fmt.Fprint(out, stamp.Format(" Jan _2 15:04:05 "))
	} else {
		fmt.Fprint(out, stamp.Format(" Jan _2 2006     "))
	}
	fmt.Fprint(out, prefix)
	fmt.Fprint(out, name)
	fmt.Fprint(out, postfix)

	if (attr & dos.FILE_ATTRIBUTE_REPARSE_POINT) != 0 {
		indicator = "@"
	}
	if (flag & O_INDICATOR) > 0 {
		io.WriteString(out, indicator)
	}
	if indicator == "@" {
		path := dos.Join(folder, name)
		link_to, err := os.Readlink(path)
		if err == nil {
			fmt.Fprintf(out, " -> %s", link_to)
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
			(flag&O_INDICATOR) != 0 {
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
		if ctx != nil {
			select {
			case <-ctx.Done():
				return ErrCtrlC
			default:

			}
		}
	}
	return nil
}

func lsSimple(ctx context.Context, folder string, nodes []os.FileInfo, flag int, out io.Writer) error {
	for _, f := range nodes {
		fmt.Fprintln(out, f.Name())
		if ctx != nil {
			select {
			case <-ctx.Done():
				return ErrCtrlC
			default:

			}
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
	findfile.Walk(wildcard, func(f *findfile.FileInfo) bool {
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
	nodesArray.nodes = tmp
	sort.Sort(nodesArray)
	var err error
	if (flag & O_ONE) != 0 {
		err = lsSimple(ctx, folder_, nodesArray.nodes, O_STRIP_DIR|flag, out)
	} else if (flag & O_LONG) != 0 {
		err = lsLong(ctx, folder_, nodesArray.nodes, O_STRIP_DIR|flag, out)
	} else {
		err = lsBox(ctx, folder_, nodesArray.nodes, O_STRIP_DIR|flag, out)
	}
	if err != nil {
		return err
	}
	if folders != nil && len(folders) > 0 {
		for _, f1 := range folders {
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
		var nameStat string
		if rxDriveOnly.MatchString(name) {
			nameStat = name + "."
		} else {
			nameStat = name
		}
		status, err := os.Lstat(nameStat)
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
		if (flag & O_ONE) != 0 {
			err = lsSimple(ctx, ".", files, flag, out)
		} else if (flag & O_LONG) != 0 {
			err = lsLong(ctx, ".", files, flag, out)
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
	'S': func(flag *int) error {
		*flag |= O_SIZESORT
		return nil
	},
	'd': func(flag *int) error {
		*flag |= O_NOT_RECURSIVE
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
		message := make([]byte, 0, 80)
		message = append(message, "Usage: ls [-"...)
		for optKey, _ := range option {
			message = append(message, byte(optKey))
		}
		message = append(message, "] [PATH(s)]..."...)
		return errors.New(string(message))
	}
	if _, ok := out.(io.Closer); ok {
		// output is a not colorable instance.
		flag &^= O_COLOR
	}
	if (flag & O_COLOR) != 0 {
		io.WriteString(out, ANSI_END)
	}
	return lsCore(ctx, paths, flag, out, err)
}

// vim:set fenc=utf8 ts=4 sw=4 noet:
