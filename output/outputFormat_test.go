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

package output

import (
	"sort"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/media"
)

func TestDefaultTypes(t *testing.T) {
	c := qt.New(t)
	c.Assert(CalendarFormat.Name, qt.Equals, "calendar")
	c.Assert(CalendarFormat.MediaType, qt.Equals, media.Builtin.CalendarType)
	c.Assert(CalendarFormat.Protocol, qt.Equals, "webcal://")
	c.Assert(CalendarFormat.Path, qt.HasLen, 0)
	c.Assert(CalendarFormat.IsPlainText, qt.Equals, true)
	c.Assert(CalendarFormat.IsHTML, qt.Equals, false)

	c.Assert(CSSFormat.Name, qt.Equals, "css")
	c.Assert(CSSFormat.MediaType, qt.Equals, media.Builtin.CSSType)
	c.Assert(CSSFormat.Path, qt.HasLen, 0)
	c.Assert(CSSFormat.Protocol, qt.HasLen, 0) // Will inherit the BaseURL protocol.
	c.Assert(CSSFormat.IsPlainText, qt.Equals, true)
	c.Assert(CSSFormat.IsHTML, qt.Equals, false)

	c.Assert(CSVFormat.Name, qt.Equals, "csv")
	c.Assert(CSVFormat.MediaType, qt.Equals, media.Builtin.CSVType)
	c.Assert(CSVFormat.Path, qt.HasLen, 0)
	c.Assert(CSVFormat.Protocol, qt.HasLen, 0)
	c.Assert(CSVFormat.IsPlainText, qt.Equals, true)
	c.Assert(CSVFormat.IsHTML, qt.Equals, false)
	c.Assert(CSVFormat.Permalinkable, qt.Equals, false)

	c.Assert(HTMLFormat.Name, qt.Equals, "html")
	c.Assert(HTMLFormat.MediaType, qt.Equals, media.Builtin.HTMLType)
	c.Assert(HTMLFormat.Path, qt.HasLen, 0)
	c.Assert(HTMLFormat.Protocol, qt.HasLen, 0)
	c.Assert(HTMLFormat.IsPlainText, qt.Equals, false)
	c.Assert(HTMLFormat.IsHTML, qt.Equals, true)
	c.Assert(AMPFormat.Permalinkable, qt.Equals, true)

	c.Assert(AMPFormat.Name, qt.Equals, "amp")
	c.Assert(AMPFormat.MediaType, qt.Equals, media.Builtin.HTMLType)
	c.Assert(AMPFormat.Path, qt.Equals, "amp")
	c.Assert(AMPFormat.Protocol, qt.HasLen, 0)
	c.Assert(AMPFormat.IsPlainText, qt.Equals, false)
	c.Assert(AMPFormat.IsHTML, qt.Equals, true)
	c.Assert(AMPFormat.Permalinkable, qt.Equals, true)

	c.Assert(RSSFormat.Name, qt.Equals, "rss")
	c.Assert(RSSFormat.MediaType, qt.Equals, media.Builtin.RSSType)
	c.Assert(RSSFormat.Path, qt.HasLen, 0)
	c.Assert(RSSFormat.IsPlainText, qt.Equals, false)
	c.Assert(RSSFormat.NoUgly, qt.Equals, true)
	c.Assert(CalendarFormat.IsHTML, qt.Equals, false)

	c.Assert(len(DefaultFormats), qt.Equals, 11)
}

func TestGetFormatByName(t *testing.T) {
	c := qt.New(t)
	formats := Formats{AMPFormat, CalendarFormat}
	tp, _ := formats.GetByName("AMp")
	c.Assert(tp, qt.Equals, AMPFormat)
	_, found := formats.GetByName("HTML")
	c.Assert(found, qt.Equals, false)
	_, found = formats.GetByName("FOO")
	c.Assert(found, qt.Equals, false)
}

func TestGetFormatByExt(t *testing.T) {
	c := qt.New(t)
	formats1 := Formats{AMPFormat, CalendarFormat}
	formats2 := Formats{AMPFormat, HTMLFormat, CalendarFormat}
	tp, _ := formats1.GetBySuffix("html")
	c.Assert(tp, qt.Equals, AMPFormat)
	tp, _ = formats1.GetBySuffix("ics")
	c.Assert(tp, qt.Equals, CalendarFormat)
	_, found := formats1.GetBySuffix("not")
	c.Assert(found, qt.Equals, false)

	// ambiguous
	_, found = formats2.GetBySuffix("html")
	c.Assert(found, qt.Equals, false)
}

func TestGetFormatByFilename(t *testing.T) {
	c := qt.New(t)
	noExtNoDelimMediaType := media.Builtin.TextType
	noExtNoDelimMediaType.Delimiter = ""

	noExtMediaType := media.Builtin.TextType

	var (
		noExtDelimFormat = Format{
			Name:      "NEM",
			MediaType: noExtNoDelimMediaType,
			BaseName:  "_redirects",
		}
		noExt = Format{
			Name:      "NEX",
			MediaType: noExtMediaType,
			BaseName:  "next",
		}
	)

	formats := Formats{AMPFormat, HTMLFormat, noExtDelimFormat, noExt, CalendarFormat}
	f, found := formats.FromFilename("my.amp.html")
	c.Assert(found, qt.Equals, true)
	c.Assert(f, qt.Equals, AMPFormat)
	_, found = formats.FromFilename("my.ics")
	c.Assert(found, qt.Equals, true)
	f, found = formats.FromFilename("my.html")
	c.Assert(found, qt.Equals, true)
	c.Assert(f, qt.Equals, HTMLFormat)
	f, found = formats.FromFilename("my.nem")
	c.Assert(found, qt.Equals, true)
	c.Assert(f, qt.Equals, noExtDelimFormat)
	f, found = formats.FromFilename("my.nex")
	c.Assert(found, qt.Equals, true)
	c.Assert(f, qt.Equals, noExt)
	_, found = formats.FromFilename("my.css")
	c.Assert(found, qt.Equals, false)
}

func TestSort(t *testing.T) {
	c := qt.New(t)
	c.Assert(DefaultFormats[0].Name, qt.Equals, "html")
	c.Assert(DefaultFormats[1].Name, qt.Equals, "amp")

	json := JSONFormat
	json.Weight = 1

	formats := Formats{
		AMPFormat,
		HTMLFormat,
		json,
	}

	sort.Sort(formats)

	c.Assert(formats[0].Name, qt.Equals, "json")
	c.Assert(formats[1].Name, qt.Equals, "html")
	c.Assert(formats[2].Name, qt.Equals, "amp")
}
