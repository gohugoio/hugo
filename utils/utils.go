package utils

import (
	"os"

	jww "github.com/spf13/jwalterweatherman"
)

func CheckErr(err error, s ...string) {
	if err != nil {
		if len(s) == 0 {
			jww.CRITICAL.Println(err)
		} else {
			for _, message := range s {
				jww.ERROR.Println(message)
			}
		}
	}
}

func StopOnErr(err error, s ...string) {
	if err != nil {
		if len(s) == 0 {
			jww.CRITICAL.Println(err)
		} else {
			for _, message := range s {
				jww.CRITICAL.Println(message)
			}
		}
		os.Exit(-1)
	}
}
