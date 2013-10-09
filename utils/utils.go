package utils

import (
	"log"
	"os"
)

func CheckErr(err error, s ...string) {
	if err != nil {
		for _, message := range s {
			log.Fatalf(message)
		}
		log.Fatalf("Fatal Error: %v", err)
	}
}

func CheckErrExit(err error, s ...string) {
	if err != nil {
		CheckErr(err, s...)
		os.Exit(-1)
	}
}
