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

import (
	"bytes"
	"testing"

	"text/template"

	"github.com/stretchr/testify/require"
)

func TestPageData(t *testing.T) {
	assert := require.New(t)

	data := make(Data)

	assert.Nil(data.Pages())

	pages := Pages{
		&testPage{title: "a1"},
		&testPage{title: "a2"},
	}

	data["pages"] = pages

	assert.Equal(pages, data.Pages())

	data["pages"] = func() Pages {
		return pages
	}

	assert.Equal(pages, data.Pages())

	templ, err := template.New("").Parse(`Pages: {{ .Pages }}`)

	assert.NoError(err)

	var buff bytes.Buffer

	assert.NoError(templ.Execute(&buff, data))

	assert.Contains(buff.String(), "Pages(2)")
}
