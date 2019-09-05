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

package bundler

import (
	"testing"

	"github.com/gohugoio/hugo/helpers"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/hugio"
)

func TestMultiReadSeekCloser(t *testing.T) {
	c := qt.New(t)

	rc := newMultiReadSeekCloser(
		hugio.NewReadSeekerNoOpCloserFromString("A"),
		hugio.NewReadSeekerNoOpCloserFromString("B"),
		hugio.NewReadSeekerNoOpCloserFromString("C"),
	)

	for i := 0; i < 3; i++ {
		s1 := helpers.ReaderToString(rc)
		c.Assert(s1, qt.Equals, "ABC")
		_, err := rc.Seek(0, 0)
		c.Assert(err, qt.IsNil)
	}

}
