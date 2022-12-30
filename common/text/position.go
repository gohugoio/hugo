// Copyright 2018 The Hugo Authors. All rights reserved.
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

package text

import (
	"fmt"
	"os"
	"strings"

	"github.com/gohugoio/hugo/common/terminal"
)

// Positioner represents a thing that knows its position in a text file or stream,
// typically an error.
type Positioner interface {
	// Position returns the current position.
	// Useful in error logging, e.g. {{ errorf "error in code block: %s" .Position }}.
	Position() Position
}

// Position holds a source position in a text file or stream.
type Position struct {
	Filename     string // filename, if any
	Offset       int    // byte offset, starting at 0. It's set to -1 if not provided.
	LineNumber   int    // line number, starting at 1
	ColumnNumber int    // column number, starting at 1 (character count per line)
}

func (pos Position) String() string {
	if pos.Filename == "" {
		pos.Filename = "<stream>"
	}
	return positionStringFormatfunc(pos)
}

// IsValid returns true if line number is > 0.
func (pos Position) IsValid() bool {
	return pos.LineNumber > 0
}

var positionStringFormatfunc func(p Position) string

func createPositionStringFormatter(formatStr string) func(p Position) string {
	if formatStr == "" {
		formatStr = "\":file::line::col\""
	}

	identifiers := []string{":file", ":line", ":col"}
	var identifiersFound []string

	for i := range formatStr {
		for _, id := range identifiers {
			if strings.HasPrefix(formatStr[i:], id) {
				identifiersFound = append(identifiersFound, id)
			}
		}
	}

	replacer := strings.NewReplacer(":file", "%s", ":line", "%d", ":col", "%d")
	format := replacer.Replace(formatStr)

	f := func(pos Position) string {
		args := make([]any, len(identifiersFound))
		for i, id := range identifiersFound {
			switch id {
			case ":file":
				args[i] = pos.Filename
			case ":line":
				args[i] = pos.LineNumber
			case ":col":
				args[i] = pos.ColumnNumber
			}
		}

		msg := fmt.Sprintf(format, args...)

		if terminal.PrintANSIColors(os.Stdout) {
			return terminal.Notice(msg)
		}

		return msg
	}

	return f
}

func init() {
	positionStringFormatfunc = createPositionStringFormatter(os.Getenv("HUGO_FILE_LOG_FORMAT"))
}
