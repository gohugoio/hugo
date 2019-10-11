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

package page

import "sort"

// Used in page binary search, the most common in front.
var pageLessFunctions = []func(p1, p2 Page) bool{
	DefaultPageSort,
	lessPageDate,
	lessPagePubDate,
	lessPageTitle,
	lessPageLinkTitle,
}

func searchPage(p Page, pages Pages) int {
	if len(pages) < 1000 {
		// For smaller data sets, doing a linear search is faster.
		return searchPageLinear(p, pages, 0)
	}

	less := isPagesProbablySorted(pages, pageLessFunctions...)
	if less == nil {
		return searchPageLinear(p, pages, 0)
	}

	i := searchPageBinary(p, pages, less)
	if i != -1 {
		return i
	}

	return searchPageLinear(p, pages, 0)
}

func searchPageLinear(p Page, pages Pages, start int) int {
	for i := start; i < len(pages); i++ {
		c := pages[i]
		if c.Eq(p) {
			return i
		}
	}
	return -1
}

func searchPageBinary(p Page, pages Pages, less func(p1, p2 Page) bool) int {
	n := len(pages)

	f := func(i int) bool {
		c := pages[i]
		isLess := less(c, p)
		return !isLess || c.Eq(p)
	}

	i := sort.Search(n, f)

	if i == n {
		return -1
	}

	return searchPageLinear(p, pages, i)

}

// isProbablySorted tests if the pages slice is probably sorted.
func isPagesProbablySorted(pages Pages, lessFuncs ...func(p1, p2 Page) bool) func(p1, p2 Page) bool {
	n := len(pages)
	step := 1
	if n > 500 {
		step = 50
	}

	is := func(less func(p1, p2 Page) bool) bool {
		samples := 0

		for i := n - 1; i > 0; i = i - step {
			if less(pages[i], pages[i-1]) {
				return false
			}
			samples++
			if samples >= 15 {
				return true
			}
		}
		return samples > 0
	}

	isReverse := func(less func(p1, p2 Page) bool) bool {
		samples := 0

		for i := 0; i < n-1; i = i + step {
			if less(pages[i], pages[i+1]) {
				return false
			}
			samples++

			if samples > 15 {
				return true
			}
		}
		return samples > 0
	}

	for _, less := range lessFuncs {
		if is(less) {
			return less
		}
		if isReverse(less) {
			return func(p1, p2 Page) bool {
				return less(p2, p1)
			}
		}
	}

	return nil
}
