// Copyright 2016-present The Hugo Authors. All rights reserved.
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

package helpers

import (
	"strings"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/hugolib/paths"
)

// PathSpec holds methods that decides how paths in URLs and files in Hugo should look like.
type PathSpec struct {
	*paths.Paths
	*filesystems.BaseFs

	ProcessingStats *ProcessingStats

	// The file systems to use
	Fs *hugofs.Fs
}

// NewPathSpec creates a new PathSpec from the given filesystems and language.
// If an existing BaseFs is provided, parts of that is reused.
func NewPathSpec(fs *hugofs.Fs, cfg config.AllProvider, logger loggers.Logger, baseBaseFs *filesystems.BaseFs) (*PathSpec, error) {
	p, err := paths.New(fs, cfg)
	if err != nil {
		return nil, err
	}

	var options []func(*filesystems.BaseFs) error
	if baseBaseFs != nil {
		options = []func(*filesystems.BaseFs) error{
			filesystems.WithBaseFs(baseBaseFs),
		}
	}
	bfs, err := filesystems.NewBase(p, logger, options...)
	if err != nil {
		return nil, err
	}

	ps := &PathSpec{
		Paths:           p,
		BaseFs:          bfs,
		Fs:              fs,
		ProcessingStats: NewProcessingStats(p.Lang__()),
	}

	return ps, nil
}

// PermalinkForBaseURL creates a permalink from the given link and baseURL.
func (p *PathSpec) PermalinkForBaseURL(link, baseURL string) string {
	return baseURL + strings.TrimPrefix(link, "/")
}
