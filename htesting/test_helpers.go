// Copyright 2019 The Hugo Authors. All rights reserved.
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

package htesting

import (
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/afero"
)

// CreateTempDir creates a temp dir in the given filesystem and
// returns the dirnam and a func that removes it when done.
func CreateTempDir(fs afero.Fs, prefix string) (string, func(), error) {
	tempDir, err := afero.TempDir(fs, "", prefix)
	if err != nil {
		return "", nil, err
	}

	_, isOsFs := fs.(*afero.OsFs)

	if isOsFs && runtime.GOOS == "darwin" && !strings.HasPrefix(tempDir, "/private") {
		// To get the entry folder in line with the rest. This its a little bit
		// mysterious, but so be it.
		tempDir = "/private" + tempDir
	}
	return tempDir, func() { fs.RemoveAll(tempDir) }, nil
}

// BailOut panics with a stack trace after the given duration. Useful for
// hanging tests.
func BailOut(after time.Duration) {
	time.AfterFunc(after, func() {
		buf := make([]byte, 1<<16)
		runtime.Stack(buf, true)
		panic(string(buf))
	})

}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandIntn(n int) int {
	return rnd.Intn(n)
}
