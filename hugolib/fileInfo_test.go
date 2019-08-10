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
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/cast"
)

func TestFileInfo(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)
		fi := &fileInfo{}
		_, err := cast.ToStringE(fi)
		c.Assert(err, qt.IsNil)
	})
}
