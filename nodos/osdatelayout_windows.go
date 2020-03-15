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
	return table.Replace(layout), nil
}

var table = strings.NewReplacer(
	"yyyy", "2006",
	"MM", "01",
	"dd", "02",
	"d", "2",
	"M", "1",
	"H", "15",
	"mm", "04",
	"ss", "05",
)
