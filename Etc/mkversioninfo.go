//go:build run
// +build run

package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var jsonTemplate = `{
	"FixedFileInfo":
	{
		"FileVersion": {
			"Major": %d,
			"Minor": %d,
			"Patch": %d,
			"Build": %d
		},
		"ProductVersion": {
			"Major": %d,
			"Minor": %d,
			"Patch": %d,
			"Build": %d
		},
		"FileFlagsMask": "3f",
		"FileFlags ": "00",
		"FileOS": "040004",
		"FileType": "01",
		"FileSubType": "00"
	},
	"StringFileInfo":
	{
		"Comments": "",
		"CompanyName": "NYAOS.ORG",
		"FileDescription": "Nihongo Yet Another GOing Shell",
		"FileVersion": "%s",
		"InternalName": "",
		"LegalCopyright": "Copyright (C) 2014-2025 HAYAMA_Kaoru",
		"LegalTrademarks": "",
		"OriginalFilename": "NYAGOS.EXE",
		"PrivateBuild": "",
		"ProductName": "Nihongo Yet Another GOing Shell",
		"ProductVersion": "%s",
		"SpecialBuild": ""
	},
	"VarFileInfo":
	{
		"Translation": {
			"LangID": "0411",
			"CharsetID": "04E4"
		}
	}
}
`

var de = regexp.MustCompile(`[-\._]`)

func versionStrToNum(versionString string) ([]int, error) {
	v := de.Split(versionString, -1)
	if len(v) < 4 {
		return nil, fmt.Errorf("%s: too short version string", versionString)
	}

	var vn [4]int
	var err error

	if vn[0], err = strconv.Atoi(v[0]); err != nil {
		return nil, fmt.Errorf("%s: invalid major version", versionString)
	}
	if vn[1], err = strconv.Atoi(v[1]); err != nil {
		return nil, fmt.Errorf("%s: invalid minor version", versionString)
	}
	if vn[2], err = strconv.Atoi(v[2]); err != nil {
		return nil, fmt.Errorf("%s: invalid patch version", versionString)
	}
	if vn[3], err = strconv.Atoi(v[3]); err != nil {
		return nil, fmt.Errorf("%s: invalid build version", versionString)
	}
	return vn[:], nil
}

func getVersionData() (string, []int, error) {
	bin, err := exec.Command("git", "describe", "--tags").Output()
	if err != nil {
		return "", nil, fmt.Errorf("Could not get version string from git (%w)", err)
	}
	str := strings.TrimSpace(string(bin))
	num, err := versionStrToNum(str)
	return str, num, err
}

func main() {
	fileVerStr, fileVerNum, err := getVersionData()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf(jsonTemplate,
		fileVerNum[0],
		fileVerNum[1],
		fileVerNum[2],
		fileVerNum[3],
		fileVerNum[0],
		fileVerNum[1],
		fileVerNum[2],
		fileVerNum[3],
		fileVerStr,
		fileVerStr)
}
