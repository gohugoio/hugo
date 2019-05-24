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

package tpl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractBaseof(t *testing.T) {
	assert := require.New(t)

	replaced := extractBaseOf(`failed: template: _default/baseof.html:37:11: executing "_default/baseof.html" at <.Parents>: can't evaluate field Parents in type *hugolib.PageOutput`)

	assert.Equal("_default/baseof.html", replaced)
	assert.Equal("", extractBaseOf("not baseof for you"))
	assert.Equal("blog/baseof.html", extractBaseOf("template: blog/baseof.html:23:11:"))
	assert.Equal("blog/baseof.ace", extractBaseOf("template: blog/baseof.ace:23:11:"))
}
