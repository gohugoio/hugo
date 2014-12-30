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
	"fmt"
	"os/exec"
	"strings"

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

	var out bytes.Buffer
	var stderr bytes.Buffer
	style := viper.GetString("PygmentsStyle")

	noclasses := "true"
	if viper.GetBool("PygmentsUseClasses") {
		noclasses = "false"
	}

	cmd := exec.Command(pygmentsBin, "-l"+lexer, "-fhtml", "-O",
		fmt.Sprintf("style=%s,noclasses=%s,encoding=utf8", style, noclasses))
	cmd.Stdin = strings.NewReader(code)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		jww.ERROR.Print(stderr.String())
		return code
	}

	return out.String()
}
