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

package bufferpool

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestBufferPool(t *testing.T) {
	c := qt.New(t)

	buff := GetBuffer()
	buff.WriteString("do be do be do")
	c.Assert(buff.String(), qt.Equals, "do be do be do")
	PutBuffer(buff)

	c.Assert(buff.Len(), qt.Equals, 0)
}
