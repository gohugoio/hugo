// Copyright 2023 The Hugo Authors. All rights reserved.
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

package sass

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestIsUnquotedCSSValue(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		in  any
		out bool
	}{
		{"24px", true},
		{"1.5rem", true},
		{"10%", true},
		{"hsl(0, 0%, 100%)", true},
		{"calc(24px + 36px)", true},
		{"24xxx", true}, // a false positive.
		{123, true},
		{123.12, true},
		{"#fff", true},
		{"#ffffff", true},
		{"#ffffffff", false},
	} {
		c.Assert(isTypedCSSValue(test.in), qt.Equals, test.out)
	}

}
