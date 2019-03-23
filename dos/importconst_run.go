// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// CSource is the filename of the temporary file *.c
const CSource = "makeconst.cpp"

// CC is the Compiler command
const CC = "gcc"

// GoSource is the filename of the temporary file *.go
const GoSource = "const.go"

var clean = flag.Bool("c", false, "clean output")
var debug = flag.Bool("d", false, "debug flag")
var packageName = flag.String("p", "main", "package name")

func makeCSource(csrcname string, headers []string, vars []string) {
	fd, err := os.Create(csrcname)
	if err != nil {
		fmt.Fprintf(fd, "%s: can not create %s\n", os.Args[0], csrcname)
		return
	}
	defer fd.Close()

	for _, header1 := range headers {
		fmt.Fprintf(fd, "#include %s\n", header1)
	}
	fmt.Fprint(fd, `
void p(const char *name,const char *s){
	printf("const %s=\"%s\"\n",name,s);
}
void p(const char *name,int n){
	printf("const %s=%d\n",name,n);
}
void p(const char *name,long n){
	printf("const %s=%ld\n",name,n);
}
void p(const char *name,unsigned long n){
	printf("const %s=%ld\n",name,n);
}
void p(const char *name,double n){
	printf("const %s=%lf\n",name,n);
}

int main()
{
`)
	fmt.Fprintln(fd, `    printf("package `+*packageName+`\n\n");`)

	for _, name1 := range vars {
		fmt.Fprintf(fd, "    p(\"%s\",%s);\n", name1, name1)
	}
	fmt.Fprintln(fd, "    return 0;\n}\n")
}

func compile() error {
	var gcc exec.Cmd
	gcc.Args = []string{
		CC,
		CSource,
	}
	fn, err := exec.LookPath(CC)
	if err != nil {
		return err
	}
	gcc.Path = fn
	gcc.Stdout = os.Stdout
	gcc.Stderr = os.Stderr
	return gcc.Run()
}

func nameOfExecutable() string {
	if runtime.GOOS == "windows" {
		return "a.exe"
	} else {
		return "a.out"
	}
}

func aexe() (string, error) {
	constC, err := os.Create(GoSource)
	if err != nil {
		return "", err
	}
	defer constC.Close()

	fname := nameOfExecutable()
	aexe := exec.Cmd{
		Args:   []string{fname},
		Path:   fname,
		Stdout: constC,
		Stderr: os.Stderr,
	}
	return fname, aexe.Run()
}

func gofmt() error {
	var gofmt exec.Cmd
	gofmt.Args = []string{
		"go",
		"fmt",
		GoSource,
	}
	fn, err := exec.LookPath("go")
	if err != nil {
		return err
	}
	gofmt.Path = fn
	gofmt.Stdout = os.Stdout
	gofmt.Stderr = os.Stderr
	return gofmt.Run()
}

func main1() error {
	err := compile()
	if err != nil {
		return err
	}
	fname, err := aexe()
	if err != nil {
		return err
	}
	os.Remove(fname)
	err = gofmt()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()

	if *clean {
		os.Remove(CSource)
		os.Remove(nameOfExecutable())
		os.Remove(GoSource)
		return
	}

	headers := []string{"<cstdio>"}
	if runtime.GOOS == "windows" {
		headers = append(headers, "<windows.h>")
	}
	vars := make([]string, 0)

	for _, arg1 := range flag.Args() {
		if len(arg1) > 0 && arg1[0] == '<' {
			headers = append(headers, arg1)
		} else if strings.HasSuffix(arg1, ".h") {
			headers = append(headers, fmt.Sprintf(`"%s"`, arg1))
		} else {
			vars = append(vars, arg1)
		}
	}
	makeCSource(CSource, headers, vars)

	if !*debug {
		defer os.Remove(CSource)
	}
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
