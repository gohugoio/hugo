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

package main

import (
	"runtime"

	"os"

	"github.com/gohugoio/hugo/commands"
	jww "github.com/spf13/jwalterweatherman"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	commands.Execute()

	if jww.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError) > 0 {
		os.Exit(-1)
	}

	if commands.Hugo != nil {
		if commands.Hugo.Log.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError) > 0 {
			os.Exit(-1)
		}
	}
}
