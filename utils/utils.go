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
			jww.ERROR.Println(err)
		}
	}
}

func StopOnErr(err error, s ...string) {
	if err != nil {
		if len(s) == 0 {
			newMessage := err.Error()

			// Printing an empty string results in a error with
			// no message, no bueno.
			if newMessage != "" {
				jww.CRITICAL.Println(newMessage)
			}
		} else {
			for _, message := range s {
				if message != "" {
					jww.CRITICAL.Println(message)
				}
			}
		}
		os.Exit(-1)
	}
}
