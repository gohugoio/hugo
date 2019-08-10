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

	qt "github.com/frankban/quicktest"
)

func TestPageData(t *testing.T) {
	c := qt.New(t)

	data := make(Data)

	c.Assert(data.Pages(), qt.IsNil)

	pages := Pages{
		&testPage{title: "a1"},
		&testPage{title: "a2"},
	}

	data["pages"] = pages

	c.Assert(data.Pages(), eq, pages)

	data["pages"] = func() Pages {
		return pages
	}

	c.Assert(data.Pages(), eq, pages)

	templ, err := template.New("").Parse(`Pages: {{ .Pages }}`)

	c.Assert(err, qt.IsNil)

	var buff bytes.Buffer

	c.Assert(templ.Execute(&buff, data), qt.IsNil)

	c.Assert(buff.String(), qt.Contains, "Pages(2)")

}
