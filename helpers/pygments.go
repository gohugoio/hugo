// Copyright © 2013-14 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helpers

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

const pygmentsBin = "pygmentize"

// HasPygments checks to see if Pygments is installed and available
// on the system.
func HasPygments() bool {
	if _, err := exec.LookPath(pygmentsBin); err != nil {
		return false
	}
	return true
}

// Highlight takes some code and returns highlighted code.
func Highlight(code string, lexer string) string {

	if !HasPygments() {
		jww.WARN.Println("Highlighting requires Pygments to be installed and in the path")
		return code
	}

	fs := hugofs.OsFs

	style := viper.GetString("PygmentsStyle")

	noclasses := "true"
	if viper.GetBool("PygmentsUseClasses") {
		noclasses = "false"
	}

	// Try to read from cache first
	hash := sha1.New()
	io.WriteString(hash, lexer)
	io.WriteString(hash, code)
	io.WriteString(hash, style)
	io.WriteString(hash, noclasses)

	cachefile := filepath.Join(viper.GetString("CacheDir"), fmt.Sprintf("pygments-%x", hash.Sum(nil)))
	exists, err := Exists(cachefile, fs)
	if err != nil {
		jww.ERROR.Print(err.Error())
		return code
	}
	if exists {
		f, err := fs.Open(cachefile)
		if err != nil {
			jww.ERROR.Print(err.Error())
			return code
		}

		s, err := ioutil.ReadAll(f)
		if err != nil {
			jww.ERROR.Print(err.Error())
			return code
		}

		return string(s)
	}

	// No cache file, render and cache it
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(pygmentsBin, "-l"+lexer, "-fhtml", "-O",
		fmt.Sprintf("style=%s,noclasses=%s,encoding=utf8", style, noclasses))
	cmd.Stdin = strings.NewReader(code)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		jww.ERROR.Print(stderr.String())
		return code
	}

	// Write cache file
	if err := WriteToDisk(cachefile, bytes.NewReader(out.Bytes()), fs); err != nil {
		jww.ERROR.Print(stderr.String())
	}

	return out.String()
}
