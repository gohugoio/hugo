// Copyright © 2013 Steve Francia <spf@spf13.com>.
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
	"log"
	"os/exec"
	"strings"
)

func Highlight(code string, lexer string) string {
	var pygmentsBin = "pygmentize"

	if _, err := exec.LookPath(pygmentsBin); err != nil {
		log.Print("Highlighting requries Pygments to be installed and in the path")
		return code
	}

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(pygmentsBin, "-l"+lexer, "-fhtml", "-O style=monokai,noclasses=true,encoding=utf-8")
	cmd.Stdin = strings.NewReader(code)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Print(stderr.String())
		return code
	}

	return out.String()
}
