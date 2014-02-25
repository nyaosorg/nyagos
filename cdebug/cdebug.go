package cdebug

import "os"
import "path/filepath"
import "strings"

var logDir string = ""
var w *os.File = nil

func prints(arg0 string, args []string) {
	if logDir == "" {
		logDir = filepath.Join(filepath.ToSlash(os.TempDir()), "DEBUG")
		fi, ok := os.Stat(logDir)
		if ok == nil && fi.IsDir() {
			fname := filepath.Base(filepath.ToSlash(os.Args[0]))
			lastPeriod := strings.LastIndex(fname, ".")
			if lastPeriod >= 0 {
				fname = fname[0:lastPeriod]
			}
			fullPath := filepath.Join(logDir, fname+".log")
			var err error
			w, err = os.OpenFile(filepath.ToSlash(fullPath), os.O_CREATE, 0666)
			if err != nil {
				w = nil
			}
		}
	}
	if w != nil {
		w.WriteString(arg0)
		for _, v := range args {
			w.WriteString(" ")
			w.WriteString(v)
		}
		w.Sync()
	}
}

func Print(arg0 string, args ...string) {
	prints(arg0, args)
}

func Println(arg0 string, args ...string) {
	prints(arg0, args)
	if w != nil {
		w.WriteString("\n")
	}
}
