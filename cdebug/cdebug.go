package cdebug

import "os"
import "path/filepath"
import "strings"

var logDir string = ""
var w *os.File = nil

func prints(args []string) {
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
		w.WriteString( strings.Join(args," ") )
	}
}

func Print(args ...string) {
	prints(args)
	if w != nil {
		w.Sync()
	}
}

func Println(args ...string) {
	prints(args)
	if w != nil {
		w.WriteString("\n")
		w.Sync()
	}
}
