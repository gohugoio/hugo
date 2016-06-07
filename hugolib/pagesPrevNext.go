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

package hugolib

// Prev returns the previous page reletive to the given page.
func (p Pages) Prev(cur *Page) *Page {
	for x, c := range p {
		if c.UniqueID() == cur.UniqueID() {
			if x == 0 {
				return p[len(p)-1]
			}
			return p[x-1]
		}
	}
	return nil
}

// Next returns the next page reletive to the given page.
func (p Pages) Next(cur *Page) *Page {
	for x, c := range p {
		if c.UniqueID() == cur.UniqueID() {
			if x < len(p)-1 {
				return p[x+1]
			}
			return p[0]
		}
	}
	return nil
}
