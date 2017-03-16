// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"sync"

	"github.com/spf13/hugo/output"
)

// PageOutput represents one of potentially many output formats of a given
// Page.
type PageOutput struct {
	*Page

	// Pagination
	paginator     *Pager
	paginatorInit sync.Once

	// Keep this to create URL/path variations, i.e. paginators.
	targetPathDescriptor targetPathDescriptor

	outputFormat output.Format
}

func (p *PageOutput) targetPath(addends ...string) (string, error) {
	tp, err := p.createTargetPath(p.outputFormat, addends...)
	if err != nil {
		return "", err
	}
	return tp, nil

}

func newPageOutput(p *Page, createCopy bool, f output.Format) (*PageOutput, error) {
	if createCopy {
		p.initURLs()
		p = p.copy()
	}

	td, err := p.createTargetPathDescriptor(f)

	if err != nil {
		return nil, err
	}

	return &PageOutput{
		Page:                 p,
		outputFormat:         f,
		targetPathDescriptor: td,
	}, nil
}

// copy creates a copy of this PageOutput with the lazy sync.Once vars reset
// so they will be evaluated again, for word count calculations etc.
func (p *PageOutput) copy() *PageOutput {
	c, err := newPageOutput(p.Page, true, p.outputFormat)
	if err != nil {
		panic(err)
	}
	return c
}
