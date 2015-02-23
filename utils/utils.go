package utils

import (
	"os"
	"strings"

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
			newMessage := cutUsageMessage(err.Error())

			// Printing an empty string results in a error with
			// no message, no bueno.
			if newMessage != "" {
				jww.CRITICAL.Println(newMessage)
			}
		} else {
			for _, message := range s {
				message := cutUsageMessage(message)

				if message != "" {
					jww.CRITICAL.Println(message)
				}
			}
		}
		os.Exit(-1)
	}
}

// cutUsageMessage splits the incoming string on the beginning of the usage
// message text. Anything in the first element of the returned slice, trimmed
// of its Unicode defined spaces, should be returned. The 2nd element of the
// slice will have the usage message  that we wish to elide.
//
// This is done because Cobra already prints Hugo's usage message; not eliding
// would result in the usage output being printed twice, which leads to bug
// reports, more specifically: https://github.com/spf13/hugo/issues/374
func cutUsageMessage(s string) string {
	pieces := strings.Split(s, "Usage of")
	return strings.TrimSpace(pieces[0])
}
