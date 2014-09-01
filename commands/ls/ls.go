package ls

import "fmt"
import "io"
import "os"
import "path/filepath"
import "regexp"
import "sort"
import "strings"

import "../../conio"
import "../../dos"

const (
	O_STRIP_DIR = 1
	O_LONG      = 2
	O_INDICATOR = 4
	O_COLOR     = 8
	O_ALL       = 16
	O_TIME      = 32
	O_REVERSE   = 64
	O_RECURSIVE = 128
)

type fileInfoT struct {
	name        string
	os.FileInfo // anonymous
}

const (
	ANSI_EXEC     = "\x1B[1;35m"
	ANSI_DIR      = "\x1B[1;32m"
	ANSI_NORM     = "\x1B[1;37m"
	ANSI_READONLY = "\x1B[1;33m"
	ANSI_HIDDEN   = "\x1B[1;34m"
	ANSI_END      = "\x1B[39m"
)

func (this fileInfoT) Name() string { return this.name }

func newMyFileInfoT(name string, info os.FileInfo) *fileInfoT {
	return &fileInfoT{name, info}
}

func lsOneLong(folder string, status os.FileInfo, flag int, out io.Writer) {
	indicator := " "
	prefix := ""
	postfix := ""
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
	attr, attrErr := dos.NewFileAttr(dos.Join(folder, status.Name()))
	if attrErr == nil && attr.IsReparse() {
		indicator = "@"
	}
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
	if attr != nil && attr.IsHidden() && (flag&O_COLOR) != 0 {
		prefix = ANSI_HIDDEN
		postfix = ANSI_END
	}
	if (flag & O_STRIP_DIR) > 0 {
		name = filepath.Base(name)
	}
	stamp := status.ModTime()
	fmt.Fprintf(out, " %8d %04d-%02d-%02d %02d:%02d %s%s%s",
		status.Size(),
		stamp.Year(),
		stamp.Month(),
		stamp.Day(),
		stamp.Hour(),
		stamp.Minute(),
		prefix,
		name,
		postfix)
	if (flag & O_INDICATOR) > 0 {
		io.WriteString(out, indicator)
	}
	io.WriteString(out, "\n")
}

func lsBox(folder string, nodes []os.FileInfo, flag int, out io.Writer) {
	nodes_ := make([]string, len(nodes))
	for key, val := range nodes {
		prefix := ""
		postfix := ""
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
		if dos.IsExecutableSuffix(filepath.Ext(val.Name())) {
			if (flag & O_COLOR) != 0 {
				prefix = ANSI_EXEC
				postfix = ANSI_END
			}
			if (flag & O_INDICATOR) != 0 {
				indicator = "*"
			}
		}
		attr, attrErr := dos.NewFileAttr(dos.Join(folder, val.Name()))
		if attrErr == nil && attr.IsHidden() && (flag&O_COLOR) != 0 {
			prefix = ANSI_HIDDEN
			postfix = ANSI_END
		}
		if attrErr == nil && attr.IsReparse() && (flag&O_INDICATOR) != 0 {
			indicator = "@"
		}
		nodes_[key] = prefix + val.Name() + postfix + indicator
	}
	conio.BoxPrint(nodes_, out)
}

func lsLong(folder string, nodes []os.FileInfo, flag int, out io.Writer) {
	for _, finfo := range nodes {
		lsOneLong(folder, finfo, flag, out)
	}
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
	} else {
		result = (this.nodes[i].Name() < this.nodes[j].Name())
	}
	if (this.flag & O_REVERSE) != 0 {
		result = !result
	}
	return result
}
func (this fileInfoCollection) Swap(i, j int) {
	tmp := this.nodes[i]
	this.nodes[i] = this.nodes[j]
	this.nodes[j] = tmp
}

func lsFolder(folder string, flag int, out io.Writer) error {
	var folder_ string
	if rxDriveOnly.MatchString(folder) {
		folder_ = folder + "."
	} else {
		folder_ = folder
	}
	fd, err := os.Open(folder_)
	if err != nil {
		return err
	}
	var nodesArray fileInfoCollection
	nodesArray.nodes, err = fd.Readdir(-1)
	fd.Close()
	nodesArray.flag = flag
	if err != nil {
		return err
	}
	tmp := make([]os.FileInfo, 0)
	var folders []string = nil
	if (flag & O_RECURSIVE) != 0 {
		folders = make([]string, 0)
	}
	for _, f := range nodesArray.nodes {
		attr, attrErr := dos.NewFileAttr(dos.Join(folder_, f.Name()))
		if (strings.HasPrefix(f.Name(), ".") || (attrErr == nil && attr.IsHidden())) && (flag&O_ALL) == 0 {
			continue
		}
		if f.IsDir() && folders != nil {
			folders = append(folders, f.Name())
		} else {
			tmp = append(tmp, f)
		}
	}
	nodesArray.nodes = tmp
	sort.Sort(nodesArray)
	if (flag & O_LONG) > 0 {
		lsLong(folder_, nodesArray.nodes, O_STRIP_DIR|flag, out)
	} else {
		lsBox(folder_, nodesArray.nodes, O_STRIP_DIR|flag, out)
	}
	if folders != nil && len(folders) > 0 {
		for _, f1 := range folders {
			f1fullpath := dos.Join(folder, f1)
			fmt.Fprintf(out, "\n%s:\n", f1fullpath)
			lsFolder(f1fullpath, flag, out)
		}
	}
	return nil
}

var rxDriveOnly = regexp.MustCompile("^[a-zA-Z]:$")

func lsCore(paths []string, flag int, out io.Writer) error {
	if len(paths) <= 0 {
		return lsFolder(".", flag, out)
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
		status, err := os.Stat(nameStat)
		if err != nil {
			continue
		} else if status.IsDir() {
			dirs = append(dirs, name)
		} else if (flag & O_LONG) != 0 {
			lsOneLong(".", newMyFileInfoT(name, status), flag, out)
			printCount += 1
		} else {
			files = append(files, newMyFileInfoT(name, status))
		}
	}
	if len(files) > 0 {
		lsBox(".", files, flag, out)
		printCount = len(files)
	}
	for _, name := range dirs {
		if len(paths) > 1 {
			if printCount > 0 {
				io.WriteString(out, "\n")
			}
			io.WriteString(out, name)
			io.WriteString(out, ":\n")
		}
		err := lsFolder(name, flag, out)
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
}

// 存在しないオプションに関するエラー
type OptionError struct {
	Option rune
}

func (this OptionError) Error() string {
	return fmt.Sprintf("-%c: No such option", this.Option)
}

// ls 機能のエントリ:引数をオプションとパスに分離する
func Main(args []string, out io.Writer) error {
	flag := 0
	paths := make([]string, 0)
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			for _, o := range arg[1:] {
				setter, ok := option[o]
				if !ok {
					var err OptionError
					err.Option = o
					return err
				}
				err := setter(&flag)
				if err != nil {
					return err
				}
			}
		} else {
			paths = append(paths, arg)
		}
	}
	return lsCore(paths, flag, out)
}

// vim:set fenc=utf8 ts=4 sw=4 noet:
