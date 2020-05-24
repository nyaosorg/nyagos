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

type FileXmlT struct {
	XMLName xml.Name `xml:"File"`
	Id      string   `xml:"Id,attr"`      // 'NyagosExe'
	Name    string   `xml:"Name,attr"`    // 'nyagos.exe'
	DiskId  int      `xml:"DiskId,attr"`  // '1'
	Source  string   `xml:"Source,attr"`  // '../../cmd/amd64/nyagos.exe'
	KeyPath string   `xml:"KeyPath,attr"` // 'yes'
}

type ComponentXmlT struct {
	XMLName xml.Name `xml:"Component"`
	Id      string   `xml:"Id,attr"`
	Guid    string   `xml:"Guid,attr"`
	File    FileXmlT `xml:"File"`
}

type DirectoryXmlT struct {
	XMLName   xml.Name         `xml:"Directory"`
	Id        string           `xml:"Id,attr"`
	Name      string           `xml:"Name,attr"`
	File      []*ComponentXmlT `xml:"File,omitempty"`
	Directory []*DirectoryXmlT `xml:"Directory,omitempty"`
}

func (this *DirectoryXmlT) GetDirectory() []*DirectoryXmlT {
	return this.Directory
}
func (this *DirectoryXmlT) SetDirectory(value []*DirectoryXmlT) {
	this.Directory = value
}

func (this *DirectoryXmlT) GetFile() []*ComponentXmlT {
	return this.File
}

func (this *DirectoryXmlT) SetFile(value []*ComponentXmlT) {
	this.File = value
}

type IncludeXmlT struct {
	XMLName   xml.Name         `xml:"Include"`
	File      []*ComponentXmlT `xml:"File,omitempty"`
	Directory []*DirectoryXmlT `xml:"Directory"`
}

func (this *IncludeXmlT) GetDirectory() []*DirectoryXmlT {
	return this.Directory
}
func (this *IncludeXmlT) SetDirectory(value []*DirectoryXmlT) {
	this.Directory = value
}

func (this *IncludeXmlT) GetFile() []*ComponentXmlT {
	return this.File
}

func (this *IncludeXmlT) SetFile(value []*ComponentXmlT) {
	this.File = value
}

type DirectoryI interface {
	GetDirectory() []*DirectoryXmlT
	SetDirectory([]*DirectoryXmlT)
	GetFile() []*ComponentXmlT
	SetFile([]*ComponentXmlT)
}

var _ DirectoryI = &DirectoryXmlT{}
var _ DirectoryI = &IncludeXmlT{}

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

func setFiles(this DirectoryI, subdir []string, files []*ComponentXmlT) {
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
	d := &DirectoryXmlT{Id: makeId(subdir[0]) + "Dir", Name: subdir[0]}
	setFiles(d, subdir[1:], files)
	this.SetDirectory(append(this.GetDirectory(), d))
}

func SetFiles(this DirectoryI, subdir string, files []*ComponentXmlT) {
	subdirArray := strings.Split(filepath.ToSlash(subdir), "/")
	setFiles(this, subdirArray, files)
}

var rxTrimDots = regexp.MustCompile(`^(\.\./)+`)

func readerToFileObjs(uuidSeed uuid.UUID, in io.Reader) (*IncludeXmlT, error) {
	sc := bufio.NewScanner(in)
	fileXmls := []*ComponentXmlT{}
	directory := map[string][]*ComponentXmlT{}
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
		f := &ComponentXmlT{
			Id:   id,
			Guid: guid.String(),
			File: FileXmlT{
				Id:      "F" + id,
				Name:    filepath.Base(dstPath),
				DiskId:  len(fileXmls) + 1,
				Source:  filepath.FromSlash(sourcePath),
				KeyPath: "yes",
			},
		}

		dstDir := filepath.Dir(rxTrimDots.ReplaceAllString(dstPath, ""))

		directory[dstDir] = append(directory[dstDir], f)
	}
	root := &IncludeXmlT{}
	for dirName, files := range directory {
		SetFiles(root, dirName, files)
	}
	return root, nil
}

func fileXmlToBin(xml1 *IncludeXmlT) ([]byte, error) {
	return xml.MarshalIndent(xml1, "", "    ")
}

func MakeWxi(fileList io.Reader, uuidSeed string, wxsOutput io.Writer) error {
	seed, err := uuid.Parse(uuidSeed)
	if err != nil {
		return err
	}
	xml1, err := readerToFileObjs(seed, fileList)
	if err != nil {
		return err
	}
	bin, err := fileXmlToBin(xml1)
	if err != nil {
		return err
	}
	bw := bufio.NewWriter(wxsOutput)
	bw.Write([]byte(xml.Header))
	bw.Write(bin)
	bw.Write([]byte{'\n'})
	bw.Flush()
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr,
			"Usage: %s SeedUUID < FileList.txt > Files.wxi\n",
			os.Args[0])
		os.Exit(2)
	}
	if err := MakeWxi(os.Stdin, os.Args[1], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
