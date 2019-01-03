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

package hugolib

import (
	"github.com/gohugoio/hugo/resources/page"
)

// This is sent to the shortcodes. They cannot access the content
// they're a part of. It would cause an infinite regress.
//
// Go doesn't support virtual methods, so this careful dance is currently (I think)
// the best we can do.
type pageWithoutContent struct {
	page.PageWithoutContent
	page.ContentProvider
}

func (p pageWithoutContent) page() page.Page {
	return p.PageWithoutContent.(page.Page)
}

func newPageWithoutContent(p page.Page) page.Page {
	return pageWithoutContent{
		PageWithoutContent: p,
		ContentProvider:    page.NopPage,
	}
}
