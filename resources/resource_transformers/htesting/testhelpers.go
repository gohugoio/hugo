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
	"path/filepath"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/resources"
	"github.com/spf13/afero"
)

func NewResourceTransformerForSpec(spec *resources.Spec, filename, content string) (resources.ResourceTransformer, error) {
	filename = filepath.FromSlash(filename)

	fs := spec.Fs.Source
	if err := afero.WriteFile(fs, filename, []byte(content), 0o777); err != nil {
		return nil, err
	}

	var open hugio.OpenReadSeekCloser = func() (hugio.ReadSeekCloser, error) {
		return fs.Open(filename)
	}

	r, err := spec.NewResource(resources.ResourceSourceDescriptor{TargetPath: filepath.FromSlash(filename), OpenReadSeekCloser: open, GroupIdentity: identity.Anonymous})
	if err != nil {
		return nil, err
	}

	return r.(resources.ResourceTransformer), nil
}
