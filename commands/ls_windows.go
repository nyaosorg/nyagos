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

	"github.com/nyaosorg/go-windows-shortcut"
	"github.com/nyaosorg/nyagos/nodos"
	"github.com/zetamatta/go-box/v2"
	"github.com/zetamatta/go-findfile"
)

const (
	_              = iota
	optionStripDir = (1 << iota)
	optionLong
	optionIndicator
	optionColor
	optionAll
	optionTime
	optionReserve
	optionRecursive
	optionOne
	optionHelp
	optionSizeSort
	optionHuman
	optionNotRecursive
	optionDereference
)

type fileInfoT struct {
	name        string
	os.FileInfo // anonymous
}

const (
	ansiExec     = "\x1B[35;1m"
	ansiDir      = "\x1B[32;1m"
	ansiNorm     = "\x1B[37;1m"
	ansiReadOnly = "\x1B[33;1m"
	ansiHidden   = "\x1B[34;1m"
	ansiEnd      = "\x1B[0m"
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

func (f fileInfoT) Name() string { return f.name }

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
	if (flag & optionColor) != 0 {
		prefix = ansiNorm
		postfix = ansiEnd
	}
	if status.IsDir() {
		io.WriteString(out, "d")
		indicator = "/"
		if (flag & optionColor) != 0 {
			prefix = ansiDir
			postfix = ansiEnd
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
		if (flag & optionColor) != 0 {
			prefix = ansiReadOnly
			postfix = ansiEnd
		}
		io.WriteString(out, "-")
	}
	if (perm & 1) > 0 {
		io.WriteString(out, "x")
	} else if nodos.IsExecutableSuffix(filepath.Ext(name)) {
		io.WriteString(out, "x")
		indicator = "*"
		if (flag & optionColor) != 0 {
			prefix = ansiExec
			postfix = ansiEnd
		}
	} else {
		io.WriteString(out, "-")
	}
	putFlag(attr, windows.FILE_ATTRIBUTE_ARCHIVE, "a", out)
	putFlag(attr, windows.FILE_ATTRIBUTE_SYSTEM, "s", out)
	putFlag(attr, windows.FILE_ATTRIBUTE_HIDDEN, "h", out)

	if (attr&windows.FILE_ATTRIBUTE_HIDDEN) != 0 &&
		(flag&optionColor) != 0 {
		prefix = ansiHidden
		postfix = ansiEnd
	}
	if (flag & optionStripDir) > 0 {
		name = filepath.Base(name)
	}
	if (flag & optionHuman) != 0 {
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
	if (flag & optionIndicator) > 0 {
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
	_nodes := make([]string, len(nodes))
	for key, val := range nodes {
		prefix := ""
		postfix := ""
		if (flag & optionColor) != 0 {
			prefix = ansiNorm
			postfix = ansiEnd
		}
		indicator := ""
		if val.IsDir() {
			if (flag & optionColor) != 0 {
				prefix = ansiDir
				postfix = ansiEnd
			}
			if (flag & optionIndicator) != 0 {
				indicator = "/"
			}
		}
		if (val.Mode().Perm() & 2) == 0 {
			if (flag & optionColor) != 0 {
				prefix = ansiReadOnly
				postfix = ansiEnd
			}
		}
		if !val.IsDir() && nodos.IsExecutableSuffix(filepath.Ext(val.Name())) {
			if (flag & optionColor) != 0 {
				prefix = ansiExec
				postfix = ansiEnd
			}
			if (flag & optionIndicator) != 0 {
				indicator = "*"
			}
		}
		attr := findfile.GetFileAttributes(val)
		if (attr&windows.FILE_ATTRIBUTE_HIDDEN) != 0 &&
			(flag&optionColor) != 0 {
			prefix = ansiHidden
			postfix = ansiEnd
		}
		if (attr&windows.FILE_ATTRIBUTE_REPARSE_POINT) != 0 &&
			(flag&optionIndicator) != 0 {
			indicator = "@"
		}
		if indicator != "" {
			_nodes[key] = prefix + val.Name() + postfix + indicator
		} else {
			_nodes[key] = prefix + val.Name()
		}
	}
	isSucceeded := box.Print(ctx, _nodes, out)
	if (flag & optionColor) != 0 {
		io.WriteString(out, ansiEnd)
	}
	if !isSucceeded {
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
	if (flag & optionHuman) != 0 {
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
		if (flag & optionIndicator) != 0 {
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

func (f fileInfoCollection) Len() int {
	return len(f.nodes)
}

func (f fileInfoCollection) Less(i, j int) bool {
	var result bool
	if (f.flag & optionTime) != 0 {
		result = f.nodes[i].ModTime().After(f.nodes[j].ModTime())
		if !result && !f.nodes[i].ModTime().Before(f.nodes[j].ModTime()) {
			result = (f.nodes[i].Name() < f.nodes[j].Name())
		}
	} else if (f.flag & optionSizeSort) != 0 {
		diff := f.nodes[i].Size() - f.nodes[j].Size()
		if diff != 0 {
			result = (diff < 0)
		} else {
			result = (f.nodes[i].Name() < f.nodes[j].Name())
		}
	} else {
		result = (f.nodes[i].Name() < f.nodes[j].Name())
	}
	if (f.flag & optionReserve) != 0 {
		result = !result
	}
	return result
}
func (f fileInfoCollection) Swap(i, j int) {
	f.nodes[i], f.nodes[j] = f.nodes[j], f.nodes[i]
}

func lsFolder(ctx context.Context, folder string, flag int, out io.Writer) error {
	var _folder string
	if rxDriveOnly.MatchString(folder) {
		_folder = folder + "."
	} else {
		_folder = folder
	}
	nodesArray := fileInfoCollection{flag: flag}
	var folders []string = nil
	if (flag & optionRecursive) != 0 {
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
		if (flag & optionAll) == 0 {
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
	if (flag & optionLong) != 0 {
		err = lsLong(_ctx, _folder, nodesArray.nodes, optionStripDir|flag, out)
	} else if (flag & optionOne) != 0 {
		err = lsSimple(_ctx, _folder, nodesArray.nodes, optionStripDir|flag, out)
	} else {
		err = lsBox(_ctx, _folder, nodesArray.nodes, optionStripDir|flag, out)
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
		if (flag & optionDereference) != 0 {
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
		} else if status.IsDir() && (flag&optionNotRecursive) == 0 {
			dirs = append(dirs, name)
		} else {
			files = append(files, &fileInfoT{filepath.Clean(name), status})
		}
	}
	if len(files) > 0 {
		nodesArray := fileInfoCollection{flag: flag, nodes: files}
		sort.Sort(nodesArray)
		var err error
		if (flag & optionLong) != 0 {
			err = lsLong(ctx, ".", files, flag, out)
		} else if (flag & optionOne) != 0 {
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
		*flag |= optionLong
		return nil
	},
	'F': func(flag *int) error {
		*flag |= optionIndicator
		return nil
	},
	'o': func(flag *int) error {
		*flag |= optionColor
		return nil
	},
	'a': func(flag *int) error {
		*flag |= optionAll
		return nil
	},
	't': func(flag *int) error {
		*flag |= optionTime
		return nil
	},
	'r': func(flag *int) error {
		*flag |= optionReserve
		return nil
	},
	'R': func(flag *int) error {
		*flag |= optionRecursive
		return nil
	},
	'1': func(flag *int) error {
		*flag |= optionOne
		return nil
	},
	'h': func(flag *int) error {
		*flag |= optionHuman
		return nil
	},
	'?': func(flag *int) error {
		*flag |= optionHelp
		return nil
	},
	'S': func(flag *int) error {
		*flag |= optionSizeSort
		return nil
	},
	'd': func(flag *int) error {
		*flag |= optionNotRecursive
		return nil
	},
	'L': func(flag *int) error {
		*flag |= optionDereference
		return nil
	},
}

// OptionError is the error when the given option does not exist in the specification.
type OptionError struct {
	Option rune
}

func (err OptionError) Error() string {
	return fmt.Sprintf("-%c: No such option", err.Option)
}

func cmdLs(ctx context.Context, cmd Param) (int, error) {
	return Ls(ctx, cmd.Args(), cmd.Out(), cmd.Err(), cmd.Term())
}

func Ls(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer, term io.Writer) (int, error) {
	flag := 0
	paths := make([]string, 0)
	for _, arg := range args[1:] {
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
	if (flag & optionHelp) != 0 {
		var message strings.Builder
		message.WriteString("Usage: ls [-")
		for optKey := range option {
			message.WriteRune(optKey)
		}
		message.WriteString("] [PATH(s)]...")
		return 1, errors.New(message.String())
	}

	if file, ok := stdout.(*os.File); ok && !isatty.IsTerminal(file.Fd()) {
		flag |= optionOne
		flag &^= optionColor
	}

	// cmd.Term() is colorableTerminal which is not fast.
	if (flag & optionColor) == 0 {
		_out := bufio.NewWriter(stdout)
		defer _out.Flush()
		stdout = _out
	} else if stdout == os.Stdout {
		_out := bufio.NewWriter(term)
		defer _out.Flush()
		stdout = _out
	}
	if (flag & optionColor) != 0 {
		io.WriteString(stdout, ansiEnd)
	}
	return 0, lsCore(ctx, paths, flag, stdout, stderr)
}

// vim:set fenc=utf8 ts=4 sw=4 noet:
