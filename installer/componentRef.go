package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func makeId(sourcePath string) string {
	baseName := filepath.Base(sourcePath)
	suffix := filepath.Ext(sourcePath)

	left := baseName[:len(baseName)-len(suffix)]
	if len(left) <= 0 {
		left = "Dot"
	} else {
		left = strings.Title(left)
	}
	return left + strings.Title(strings.TrimLeft(suffix, "."))
}

type ComponentRef struct {
	XMLName xml.Name `xml:"ComponentRef"`
	Id      string   `xml:"Id,attr"`
}

type Include struct {
	XMLName xml.Name `xml:"Include"`
	ComponentRef []*ComponentRef `xml:"ComponentRef"`
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	list := []*ComponentRef{}
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) <= 0 {
			continue
		}
		id := makeId(f[0])
		list = append(list, &ComponentRef{Id: id})
	}
	bin, err := xml.MarshalIndent(&Include{ ComponentRef:list } , "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	bw := bufio.NewWriter(os.Stdout)
	bw.Write([]byte(xml.Header))
	bw.Write(bin)
	bw.Write([]byte{'\n'})
	bw.Flush()
}
