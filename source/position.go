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

package source

import "fmt"

// Position holds a source position.
type Position struct {
	Filename     string // filename, if any
	Offset       int    // byte offset, starting at 0
	LineNumber   int    // line number, starting at 1
	ColumnNumber int    // column number, starting at 1 (character count per line)
}

func (pos Position) String() string {
	filename := pos.Filename
	if filename == "" {
		filename = "<stream>"
	}
	return fmt.Sprintf("%s:%d:%d", filename, pos.LineNumber, pos.ColumnNumber)

}
