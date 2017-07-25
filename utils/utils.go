// Copyright 2016 The Hugo Authors. All rights reserved.
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

// CheckErr logs the messages given and then the error.
// TODO(bep) Remove this package.
func CheckErr(logger *jww.Notepad, err error, s ...string) {
	if err == nil {
		return
	}
	if len(s) == 0 {
		logger.CRITICAL.Println(err)
		return
	}
	for _, message := range s {
		logger.ERROR.Println(message)
	}
	logger.ERROR.Println(err)
}

// StopOnErr exits on any error after logging it.
func StopOnErr(logger *jww.Notepad, err error, s ...string) {
	if err == nil {
		return
	}

	defer os.Exit(-1)

	if len(s) == 0 {
		newMessage := err.Error()
		// Printing an empty string results in a error with
		// no message, no bueno.
		if newMessage != "" {
			logger.CRITICAL.Println(newMessage)
		}
	}
	for _, message := range s {
		if message != "" {
			logger.CRITICAL.Println(message)
		}
	}
}
