package nodos

import (
	"golang.org/x/sys/windows/registry"
	"strings"
)

func international(key string) (string, error) {
	k, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Control Panel\International`,
		registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	val, _, err := k.GetStringValue(key)
	return val, err
}

func osDateLayout() (string, error) {
	layout, err := international("sShortDate")
	if err != nil {
		return "", err
	}
	if strings.HasSuffix(layout, "d") {
		return tableForMMDD.Replace(layout), nil
	} else {
		return tableForDDMM.Replace(layout), nil
	}
}

var tableForMMDD = strings.NewReplacer(
	"yyyy", "2006",
	"MM", "01",
	"dd", "02",
	"d", "02 Mon",
	"M", "01",
	"H", "15",
	"mm", "04",
	"ss", "05",
)

var tableForDDMM = strings.NewReplacer(
	"yyyy", "2006",
	"MM", "01",
	"dd", "02",
	"d", "Mon 02",
	"M", "01",
	"H", "15",
	"mm", "04",
	"ss", "05",
)
