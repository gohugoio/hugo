// Copyright 2024 The Hugo Authors. All rights reserved.
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

package pagesfromdata

import (
	"encoding/json"
	"fmt"
	"hash/maphash"
	"io"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
)

type PageData struct {
	skip map[uint64]bool

	SourceHash uint64 `json:"-"`

	pagemeta.PageConfig

	Content string
}

var errSkipPage = fmt.Errorf("skip page")

func (a *PageData) UnmarshalJSON(b []byte) error {
	a.calculateSourceHash(b)
	if a.skip != nil {
		if a.skip[a.SourceHash] {
			return errSkipPage
		}
	}
	return json.Unmarshal(b, &a.PageConfig)
}

// This is safe to use in parallel.
var sourceHashSeed = maphash.MakeSeed()

func (a *PageData) calculateSourceHash(b []byte) {
	var h maphash.Hash
	h.SetSeed(sourceHashSeed)
	h.Write(b)
	a.SourceHash = h.Sum64()
}

/*
 */
func PagesFromJSONFile(fim hugofs.FileMetaInfo, skip map[uint64]bool, handle func(p PageData, skip bool) error) error {
	f, err := fim.Meta().Open()
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	for {
		var p PageData
		p.skip = skip
		err := dec.Decode(&p)
		if err != nil {
			if err == io.EOF {
				break
			}
			if err != errSkipPage {
				return fmt.Errorf("parsing %q: %v", fim.Meta().Filename, err)
			}
		}
		p.skip = nil
		if err := handle(p, err == errSkipPage); err != nil {
			return err
		}

	}

	return nil
}
