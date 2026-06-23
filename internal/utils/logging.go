package utils

import (
	"log"
	"os"
)

var debug *log.Logger

func init() {
	debug = log.New(os.Stdout, "[qq]", log.Ltime)
}

func Dbg(format string, args ...any) {
	if debug == nil {
		debug = log.New(os.Stderr, "[qq] ", log.Ltime)
	}
	debug.Printf(format, args...)
}
