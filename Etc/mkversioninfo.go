//go:build run
// +build run

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/josephspurrier/goversioninfo"
)

var de = regexp.MustCompile(`[\._]`)

func newVersion(versionString string) (*goversioninfo.FileVersion, error) {
	v := de.Split(versionString, -1)
	if len(v) < 4 {
		return nil, fmt.Errorf("%s: too short version string", versionString)
	}
	var err error
	var fv goversioninfo.FileVersion

	if fv.Major, err = strconv.Atoi(v[0]); err != nil {
		return nil, fmt.Errorf("%s: invalid major version", versionString)
	}
	if fv.Minor, err = strconv.Atoi(v[1]); err != nil {
		return nil, fmt.Errorf("%s: invalid minor version", versionString)
	}
	if fv.Patch, err = strconv.Atoi(v[2]); err != nil {
		return nil, fmt.Errorf("%s: invalid patch version", versionString)
	}
	if fv.Build, err = strconv.Atoi(v[3]); err != nil {
		return nil, fmt.Errorf("%s: invalid build version", versionString)
	}
	return &fv, nil
}

func main1() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("Usage: %s FileVerFile ProdVerFile < base-json > final-json", os.Args[0])
	}

	var v goversioninfo.VersionInfo
	if err := json.NewDecoder(os.Stdin).Decode(&v); err != nil {
		return err
	}
	fileVerBin, err := os.ReadFile(os.Args[1])
	if err != nil {
		return err
	}
	fileVer := strings.TrimSpace(string(fileVerBin))

	prodVerBin, err := os.ReadFile(os.Args[2])
	if err != nil {
		return err
	}
	prodVer := strings.TrimSpace(string(prodVerBin))

	v.StringFileInfo.FileVersion = fileVer
	v.StringFileInfo.ProductVersion = prodVer
	nv, err := newVersion(fileVer)
	if err != nil {
		return err
	}
	v.FixedFileInfo.FileVersion = *nv

	nv, err = newVersion(prodVer)
	if err != nil {
		return err
	}
	v.FixedFileInfo.ProductVersion = *nv
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(&v)
	return nil
}

func main() {
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
