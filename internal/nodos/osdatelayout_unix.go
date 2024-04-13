//go:build !windows
// +build !windows

package nodos

import (
	"time"
)

func timeFormatOsLayout(t time.Time) (string, error) {
	return t.Format("Jan.02,2006"), nil
}
