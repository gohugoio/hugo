// Copyright 2015 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
