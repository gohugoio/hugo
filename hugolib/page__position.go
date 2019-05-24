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
	"github.com/gohugoio/hugo/lazy"
	"github.com/gohugoio/hugo/resources/page"
)

func newPagePosition(n *nextPrev) pagePosition {
	return pagePosition{nextPrev: n}
}

func newPagePositionInSection(n *nextPrev) pagePositionInSection {
	return pagePositionInSection{nextPrev: n}

}

type nextPrev struct {
	init     *lazy.Init
	prevPage page.Page
	nextPage page.Page
}

func (n *nextPrev) next() page.Page {
	n.init.Do()
	return n.nextPage
}

func (n *nextPrev) prev() page.Page {
	n.init.Do()
	return n.prevPage
}

type pagePosition struct {
	*nextPrev
}

func (p pagePosition) Next() page.Page {
	return p.next()
}

func (p pagePosition) NextPage() page.Page {
	return p.Next()
}

func (p pagePosition) Prev() page.Page {
	return p.prev()
}

func (p pagePosition) PrevPage() page.Page {
	return p.Prev()
}

type pagePositionInSection struct {
	*nextPrev
}

func (p pagePositionInSection) NextInSection() page.Page {
	return p.next()
}

func (p pagePositionInSection) PrevInSection() page.Page {
	return p.prev()
}
