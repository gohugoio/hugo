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

package hugolib

import (
	"io"
	"path"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/hairyhenderson/go-codeowners"
	"github.com/spf13/afero"
)

var afs = afero.NewOsFs()

func findCodeOwnersFile(dir string) (io.Reader, error) {
	for _, p := range []string{".", "docs", ".github", ".gitlab"} {
		f := path.Join(dir, p, "CODEOWNERS")

		_, err := afs.Stat(f)
		if err != nil {
			if herrors.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		return afs.Open(f)
	}

	return nil, nil
}

type codeownerInfo struct {
	owners *codeowners.Codeowners
}

func (c *codeownerInfo) forPage(p page.Page) []string {
	return c.owners.Owners(p.File().Filename())
}

func newCodeOwners(workingDir string) (*codeownerInfo, error) {
	r, err := findCodeOwnersFile(workingDir)
	if err != nil || r == nil {
		return nil, err
	}

	owners, err := codeowners.FromReader(r, workingDir)
	if err != nil {
		return nil, err
	}

	return &codeownerInfo{owners: owners}, nil
}
