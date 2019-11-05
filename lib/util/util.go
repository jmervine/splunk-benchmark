package util

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", 0)

func Verbose(s string) {
	if os.Getenv("VERBOSE") != "" && os.Getenv("VERY_VERBOSE") == "" {
		logger.Println(s)
	}
}

func Verbosef(s string, args ...interface{}) {
	if os.Getenv("VERBOSE") != "" && os.Getenv("VERY_VERBOSE") == "" {
		logger.Printf(s, args...)
	}
}

func Vverbose(s string) {
	if os.Getenv("VERY_VERBOSE") != "" {
		logger.Println(s)
	}
}

func Vverbosef(s string, args ...interface{}) {
	if os.Getenv("VERY_VERBOSE") != "" {
		logger.Printf(s, args...)
	}
}
