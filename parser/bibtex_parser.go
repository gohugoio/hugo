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

package parser

import (
	"bytes"
	"github.com/nickng/bibtex"
)

func HandleBibtexData(datum []byte) (interface{}, error) {

	bib, err := bibtex.Parse(bytes.NewReader(datum))

	m := make([]map[string]string, len(bib.Entries))

	for idx, entry := range bib.Entries {
		e := make(map[string]string)
		e["type"] = entry.Type
		e["key"] = entry.CiteName
		for key, val := range entry.Fields {
			e[key] = val.String()
		}
		m[idx] = e
	}

  if (err != nil) {
    return m, err  
  }

	return m, nil
}
