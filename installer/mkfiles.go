package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type FileXml struct {
	XMLName xml.Name `xml:"File"`
	Id      string   `xml:"Id,attr"`      // 'NyagosExe'
	Name    string   `xml:"Name,attr"`    // 'nyagos.exe'
	DiskId  int      `xml:"DiskId,attr"`  // '1'
	Source  string   `xml:"Source,attr"`  // '../../cmd/amd64/nyagos.exe'
	KeyPath string   `xml:"KeyPath,attr"` // 'yes'
}

type ComponentXml struct {
	XMLName xml.Name `xml:"Component"`
	Id      string   `xml:"Id,attr"`
	Guid    string   `xml:"Guid,attr"`
	File    FileXml  `xml:"File"`
}

type DirectoryXml struct {
	XMLName   xml.Name        `xml:"Directory"`
	Id        string          `xml:"Id,attr"`
	Name      string          `xml:"Name,attr"`
	File      []*ComponentXml `xml:"File,omitempty"`
	Directory []*DirectoryXml `xml:"Directory,omitempty"`
}

func (this *DirectoryXml) GetDirectory() []*DirectoryXml {
	return this.Directory
}
func (this *DirectoryXml) SetDirectory(value []*DirectoryXml) {
	this.Directory = value
}

func (this *DirectoryXml) GetFile() []*ComponentXml {
	return this.File
}

func (this *DirectoryXml) SetFile(value []*ComponentXml) {
	this.File = value
}

type IncludeXml struct {
	XMLName   xml.Name        `xml:"Include"`
	File      []*ComponentXml `xml:"File,omitempty"`
	Directory []*DirectoryXml `xml:"Directory"`
}

func (this *IncludeXml) GetDirectory() []*DirectoryXml {
	return this.Directory
}
func (this *IncludeXml) SetDirectory(value []*DirectoryXml) {
	this.Directory = value
}

func (this *IncludeXml) GetFile() []*ComponentXml {
	return this.File
}

func (this *IncludeXml) SetFile(value []*ComponentXml) {
	this.File = value
}

func (this *IncludeXml) WriteTo(out io.Writer) (int64, error) {
	bin, err := xml.MarshalIndent(this, "", "    ")
	if err != nil {
		return 0, err
	}
	n, err := out.Write(bin)
	return int64(n), err
}

type DirectoryI interface {
	GetDirectory() []*DirectoryXml
	SetDirectory([]*DirectoryXml)
	GetFile() []*ComponentXml
	SetFile([]*ComponentXml)
}

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

func setFiles(this DirectoryI, subdir []string, files []*ComponentXml) {
	if len(subdir) <= 0 || subdir[0] == "." {
		this.SetFile(append(this.GetFile(), files...))
		return
	}
	for _, d := range this.GetDirectory() {
		if d.Name == subdir[0] {
			setFiles(d, subdir[1:], files)
			return
		}
	}
	d := &DirectoryXml{Id: makeId(subdir[0]) + "Dir", Name: subdir[0]}
	setFiles(d, subdir[1:], files)
	this.SetDirectory(append(this.GetDirectory(), d))
}

func SetFiles(this DirectoryI, subdir string, files []*ComponentXml) {
	subdirArray := strings.Split(filepath.ToSlash(subdir), "/")
	setFiles(this, subdirArray, files)
}

var rxTrimDots = regexp.MustCompile(`^(\.\./)+`)

type ComponentRefXml struct {
	XMLName xml.Name `xml:"ComponentRef"`
	Id      string   `xml:"Id,attr"`
}

type Include2Xml struct {
	XMLName      xml.Name           `xml:"Include"`
	ComponentRef []*ComponentRefXml `xml:"ComponentRef"`
}

func (this *Include2Xml) WriteTo(out io.Writer) (int64, error) {
	bin, err := xml.MarshalIndent(this, "", "    ")
	if err != nil {
		return 0, err
	}
	n, err := out.Write(bin)
	return int64(n), err
}

func readerToFileObjs(uuidSeed uuid.UUID, in io.Reader) ([]io.WriterTo, error) {
	sc := bufio.NewScanner(in)
	fileXmls := []*ComponentXml{}
	directory := map[string][]*ComponentXml{}
	references := []*ComponentRefXml{}
	for sc.Scan() {
		field := strings.Fields(sc.Text())
		if len(field) <= 0 {
			continue
		}
		sourcePath := strings.TrimSpace(field[0])
		if sourcePath == "" {
			continue
		}
		sourcePath = filepath.ToSlash(sourcePath)
		if _, err := os.Stat(sourcePath); err != nil {
			return nil, fmt.Errorf("%s: %w", sourcePath, err)
		}
		dstPath := sourcePath
		if len(field) >= 2 {
			dstPath = filepath.ToSlash(strings.TrimSpace(field[1]))
		}

		id := makeId(sourcePath)
		guid := uuid.NewMD5(uuidSeed, []byte(id))
		f := &ComponentXml{
			Id:   id,
			Guid: guid.String(),
			File: FileXml{
				Id:      "F" + id,
				Name:    filepath.Base(dstPath),
				DiskId:  len(fileXmls) + 1,
				Source:  filepath.FromSlash(sourcePath),
				KeyPath: "yes",
			},
		}

		dstDir := filepath.Dir(rxTrimDots.ReplaceAllString(dstPath, ""))

		directory[dstDir] = append(directory[dstDir], f)

		references = append(references, &ComponentRefXml{Id: id})
	}
	root := &IncludeXml{}
	for dirName, files := range directory {
		SetFiles(root, dirName, files)
	}
	return []io.WriterTo{root, &Include2Xml{ComponentRef: references}}, nil
}

func outputWithXmlHeader(xml1 io.WriterTo,w io.Writer) error {
	bw := bufio.NewWriter(w)
	if _, err := bw.WriteString(xml.Header); err != nil {
		return err
	}
	if _, err := xml1.WriteTo(bw); err != nil {
		return err
	}
	if err := bw.WriteByte('\n'); err != nil {
		return err
	}
	return bw.Flush()
}


func MakeWxi(fileList io.Reader, uuidSeed string, out []io.Writer) error {
	seed, err := uuid.Parse(uuidSeed)
	if err != nil {
		return err
	}
	xmls, err := readerToFileObjs(seed, fileList)
	if err != nil {
		return err
	}
	for i, xml1 := range xmls {
		if i >= len(out) {
			break
		}
		if out[i] == nil {
			continue
		}
		if err := outputWithXmlHeader(xml1,out[i]) ; err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr,
			"Usage: %s SeedUUID < FileList.txt > Files.wxi\n",
			os.Args[0])
		os.Exit(2)
	}
	if err := MakeWxi(os.Stdin, os.Args[1], []io.Writer{os.Stdout}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
