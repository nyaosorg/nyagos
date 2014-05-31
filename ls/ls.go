package ls

import "fmt"
import "io"
import "os"
import "path"
import "strings"
import "time"

import "../box"

var exeSuffixes = map[string]bool{
	".bat": true,
	".cmd": true,
	".com": true,
	".exe": true,
}

const (
	O_STRIP_DIR = 1
	O_LONG      = 2
	O_INDICATOR = 4
	O_COLOR     = 8
)

type fileInfoT struct {
	name string
	info os.FileInfo
}

const (
	ANSI_EXEC     = "\x1B[1;35m"
	ANSI_DIR      = "\x1B[1;32m"
	ANSI_NORM     = "\x1B[1;37m"
	ANSI_READONLY = "\x1B[1;33m"
	ANSI_END      = "\x1B[39m"
)

func (this fileInfoT) Name() string       { return this.name }
func (this fileInfoT) Size() int64        { return this.info.Size() }
func (this fileInfoT) Mode() os.FileMode  { return this.info.Mode() }
func (this fileInfoT) ModTime() time.Time { return this.info.ModTime() }
func (this fileInfoT) IsDir() bool        { return this.info.IsDir() }
func (this fileInfoT) Sys() interface{}   { return this.info.Sys() }

func newMyFileInfoT(name string, info os.FileInfo) *fileInfoT {
	this := new(fileInfoT)
	this.name = name
	this.info = info
	return this
}

func lsOneLong(status os.FileInfo, flag int, out io.Writer) {
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
	} else if exeSuffixes[strings.ToLower(path.Ext(name))] {
		io.WriteString(out, "x")
		indicator = "*"
		if (flag & O_COLOR) != 0 {
			prefix = ANSI_EXEC
			postfix = ANSI_END
		}
	} else {
		io.WriteString(out, "-")
	}
	if (flag & O_STRIP_DIR) > 0 {
		name = path.Base(name)
	}
	io.WriteString(out, fmt.Sprintf(" %8d %s%s%s",
		status.Size(),
		prefix,
		name,
		postfix))
	if (flag & O_INDICATOR) > 0 {
		io.WriteString(out, indicator)
	}
	io.WriteString(out, "\n")
}

func lsBox(nodes []os.FileInfo, flag int, out io.Writer) {
	nodes_ := make([]string, len(nodes))
	for key, val := range nodes {
		prefix := ""
		postfix := ""
		if val.IsDir() {
			if (flag & O_COLOR) != 0 {
				prefix = ANSI_DIR
				postfix = ANSI_END
			}
			if (flag & O_INDICATOR) != 0 {
				postfix += "/"
			}
		} else if (val.Mode().Perm() & 2) == 0 {
			if (flag & O_COLOR) != 0 {
				prefix = ANSI_READONLY
				postfix = ANSI_END
			}
		} else if exeSuffixes[strings.ToLower(path.Ext(val.Name()))] {
			if (flag & O_COLOR) != 0 {
				prefix = ANSI_DIR
				postfix = ANSI_END
			}
			if (flag & O_INDICATOR) != 0 {
				postfix += "*"
			}
		}
		nodes_[key] = prefix + val.Name() + postfix
	}
	box.Print(nodes_, 80, out)
}

func lsLong(nodes []os.FileInfo, flag int, out io.Writer) {
	for _, finfo := range nodes {
		lsOneLong(finfo, flag, out)
	}
}

type fileInfoCollection struct {
	nodes []os.FileInfo
}

func (this *fileInfoCollection) Len() int {
	return len(this.nodes)
}
func (this *fileInfoCollection) Less(i, j int) bool {
	return this.nodes[i].Name() < this.nodes[j].Name()
}
func (this *fileInfoCollection) Swap(i, j int) {
	tmp := this.nodes[i]
	this.nodes[i] = this.nodes[j]
	this.nodes[j] = tmp
}

func lsFolder(folder string, flag int, out io.Writer) error {
	fd, err := os.Open(folder)
	if err != nil {
		return err
	}
	defer fd.Close()
	var nodesArray fileInfoCollection
	nodesArray.nodes, err = fd.Readdir(-1)
	if err != nil {
		return err
	}
	if (flag & O_LONG) > 0 {
		lsLong(nodesArray.nodes, O_STRIP_DIR|flag, out)
	} else {
		lsBox(nodesArray.nodes, O_STRIP_DIR|flag, out)
	}
	return nil
}

func lsCore(paths []string, flag int, out io.Writer) error {
	if len(paths) <= 0 {
		return lsFolder(".", flag, out)
	}
	dirs := make([]string, 0)
	printCount := 0
	files := make([]os.FileInfo, 0)
	for _, name := range paths {
		status, err := os.Stat(name)
		if err != nil {
			return err
		}
		if status.IsDir() {
			dirs = append(dirs, name)
		} else if (flag & O_LONG) != 0 {
			lsOneLong(newMyFileInfoT(name, status), flag, out)
			printCount += 1
		} else {
			files = append(files, newMyFileInfoT(name, status))
		}
	}
	if len(files) > 0 {
		lsBox(files, flag, out)
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
