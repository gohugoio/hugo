// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

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
