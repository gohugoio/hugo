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

// Next returns the next page reletive to the given
func (p Pages) Next(cur Page) Page {
	x := searchPage(cur, p)
	if x <= 0 {
		return nil
	}
	return p[x-1]
}

// Prev returns the previous page reletive to the given
func (p Pages) Prev(cur Page) Page {
	x := searchPage(cur, p)

	if x == -1 || len(p)-x < 2 {
		return nil
	}

	return p[x+1]

}
