// Copyright 2021 The Hugo Authors. All rights reserved.
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

package htesting

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestExtractMinorVersionFromGoTag(t *testing.T) {
	c := qt.New(t)

	c.Assert(extractMinorVersionFromGoTag("go1.17"), qt.Equals, 17)
	c.Assert(extractMinorVersionFromGoTag("go1.16.7"), qt.Equals, 16)
	c.Assert(extractMinorVersionFromGoTag("go1.17beta1"), qt.Equals, 17)
	c.Assert(extractMinorVersionFromGoTag("asdfadf"), qt.Equals, -1)
}
