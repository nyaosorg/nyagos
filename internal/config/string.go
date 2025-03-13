package config

import (
	"github.com/nyaosorg/nyagos/internal/go-ignorecase-sorted"
)

var PredictColor = [...]string{"\x1B[3;22;34m", "\x1B[0m"}

type StringPtr = ConfigPtr[string]
type StringFunc = ConfigFunc[string]

type String interface {
	Usage() string
	NoUsage() string
	Set(value string)
	Get() string
}

var Strings = ignoreCaseSorted.MapToDictionary(map[string]String{
	"predict_color": &StringPtr{
		ptr:     &PredictColor[0],
		usage:   "predict color",
		noUsage: "predict color",
	},
})
