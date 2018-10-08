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

package collections

import (
	"errors"

	"github.com/gohugoio/hugo/common/collections"
)

// Append appends the arguments up to the last one to the slice in the last argument.
// This construct allows template constructs like this:
//     {{ $pages = $pages | append $p2 $p1 }}
// Note that with 2 arguments where both are slices of the same type,
// the first slice will be appended to the second:
//     {{ $pages = $pages | append .Site.RegularPages }}
func (ns *Namespace) Append(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("need at least 2 arguments to append")
	}

	to := args[len(args)-1]
	from := args[:len(args)-1]

	return collections.Append(to, from...)

}
