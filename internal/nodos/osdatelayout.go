package nodos

import (
	"time"
)

func TimeFormatOsLayout(t time.Time) (string, error) {
	return timeFormatOsLayout(t)
}
