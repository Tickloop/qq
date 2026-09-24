package utils

import (
	"log"
	"os"
)

var debug *log.Logger

func init() {
	if os.Getenv("QQ_DEBUG") != "" {
		debug = log.New(os.Stdout, "[qq]", log.Ltime)
	}
}

func Dbg(format string, args ...any) {
	if debug != nil {
		debug.Printf(format, args...)
	}
}
