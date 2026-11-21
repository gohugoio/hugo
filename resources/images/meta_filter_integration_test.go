// Copyright 2026 The Hugo Authors. All rights reserved.
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

package images_test

import (
	"math"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/bep/imagemeta"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

func TestMetaFilter(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
timeZone = "Europe/Oslo"
[imaging.meta]
fields = ["**"]
[params]
photoDate = 2024-05-01
-- data/meta.yaml --
Copyright: "© Data"
-- assets/sunset.jpg --
sourcefilename: ../testdata/sunset.jpg
-- assets/giphy.gif --
sourcefilename: ../testdata/giphy.gif
-- assets/giphy.avif --
sourcefilename: ../testdata/giphy.avif
-- layouts/home.html --
{{ $img := resources.Get "sunset.jpg" }}
{{ $filters := slice
  images.AutoOrient
  (images.Process "resize 64x")
  (images.Meta (dict "copyright" "© Bjørn Erik" "creator" "Bjørn Erik" "description" "First"))
  (images.Meta (dict "title" "Solnedgang" "description" "Second" "date" $img.Meta.Date "lat" $img.Meta.Lat "long" $img.Meta.Long))
}}
{{ range slice "jpg" "png" "webp" "avif" "tiff" }}
{{ with $img.Filter $filters (images.Process .) }}{{ .MediaType.SubType }}: {{ .RelPermalink }}|{{ end }}
{{ end }}
{{ $plain := $img.Filter (images.Process "resize 64x") }}
{{ $noGPS := $img.Filter (images.Process "resize 64x") (images.Meta (dict "lat" 0 "long" 0 "creator" "Bjørn Erik")) }}
noGPS: {{ $noGPS.RelPermalink }}|{{ ne $plain.RelPermalink $noGPS.RelPermalink }}|
{{ $nilEmpty := $img.Filter (images.Process "resize 64x") (images.Meta nil) (images.Meta (dict)) (images.Meta (dict "lat" nil "long" "" "date" "" "copyright" nil "creator" "Bjørn Erik")) }}
nilEmpty: {{ $nilEmpty.RelPermalink }}|
{{ $dateLocal := $img.Filter (images.Process "resize 64x") (images.Meta (dict "date" "2024-05-01 12:30:00")) }}
dateLocal: {{ $dateLocal.RelPermalink }}|
{{ $dateTOML := $img.Filter (images.Process "resize 64x") (images.Meta (dict "date" site.Params.photoDate)) }}
dateTOML: {{ $dateTOML.RelPermalink }}|
{{ $fromData := $img.Filter (images.Process "resize 64x") (images.Meta hugo.Data.meta) }}
fromData: {{ $fromData.RelPermalink }}|dataKeyPreserved: {{ isset hugo.Data.meta "Copyright" }}|
{{ with (resources.Get "giphy.avif").Filter (images.Process "resize 32x") (images.Meta (dict "copyright" "Still AVIF")) }}animAVIF: {{ .RelPermalink }}|{{ end }}
{{ with (resources.Get "giphy.gif").Filter (images.Meta (dict "copyright" "Bjørn Erik")) }}GIF: {{ .MediaType.SubType }}|{{ end }}
`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	b.AssertFileContent("public/index.html", "GIF: gif|", "dataKeyPreserved: true|")
	b.AssertLogContains("Writing metadata to GIF images is not supported", "warning-image-meta-not-written")
	b.Assert(b.LogString(), qt.Not(qt.Contains), "animated AVIF")

	links := map[string]string{}
	for _, m := range regexp.MustCompile(`(\w+): (/\S+?)\|`).FindAllStringSubmatch(b.FileContent("public/index.html"), -1) {
		links[m[1]] = m[2]
	}

	formats := map[string]imagemeta.ImageFormat{
		"jpeg": imagemeta.JPEG, "png": imagemeta.PNG, "webp": imagemeta.WebP, "avif": imagemeta.AVIF, "tiff": imagemeta.TIFF,
		"noGPS": imagemeta.JPEG, "nilEmpty": imagemeta.JPEG, "dateLocal": imagemeta.JPEG, "dateTOML": imagemeta.JPEG,
		"fromData": imagemeta.JPEG, "animAVIF": imagemeta.AVIF,
	}
	b.Assert(links, qt.HasLen, len(formats))

	decode := func(name string) *imagemeta.Tags {
		var tags imagemeta.Tags
		_, err := imagemeta.Decode(imagemeta.Options{
			R:               strings.NewReader(b.FileContent("public" + links[name])),
			ImageFormat:     formats[name],
			ShouldHandleTag: func(imagemeta.TagInfo) bool { return true },
			HandleTag: func(ti imagemeta.TagInfo) error {
				tags.Add(ti)
				return nil
			},
		})
		b.Assert(err, qt.IsNil)
		return &tags
	}

	for _, name := range []string{"jpeg", "png", "webp", "avif", "tiff"} {
		b.Run(name, func(c *qt.C) {
			tags := decode(name)
			exif, xmp := tags.EXIF(), tags.XMP()
			c.Assert(exif["Copyright"].Value, qt.Equals, "© Bjørn Erik")
			c.Assert(exif["Artist"].Value, qt.Equals, "Bjørn Erik")
			c.Assert(exif["ImageDescription"].Value, qt.Equals, "Second")
			c.Assert(xmp["Title"].Value, qt.Equals, "Solnedgang")
			d, err := tags.GetDateTime()
			c.Assert(err, qt.IsNil)
			c.Assert(d.Format("2006-01-02T15:04:05"), qt.Equals, "2017-10-27T08:38:52")
			lat, long, err := tags.GetLatLong()
			c.Assert(err, qt.IsNil)
			c.Assert(math.Abs(lat-36.597441) < 1e-6, qt.IsTrue, qt.Commentf("%v", lat))
			c.Assert(math.Abs(long+4.50846) < 1e-6, qt.IsTrue, qt.Commentf("%v", long))
		})
	}

	for _, name := range []string{"noGPS", "nilEmpty"} {
		tags := decode(name)
		b.Assert(tags.EXIF()["Artist"].Value, qt.Equals, "Bjørn Erik")
		_, hasGPS := tags.EXIF()["GPSLatitude"]
		b.Assert(hasGPS, qt.IsFalse)
	}
	b.AssertFileContent("public/index.html", "|true|")

	d, err := decode("dateLocal").GetDateTime()
	b.Assert(err, qt.IsNil)
	b.Assert(d.Format(time.RFC3339), qt.Equals, "2024-05-01T12:30:00+02:00")

	d, err = decode("dateTOML").GetDateTime()
	b.Assert(err, qt.IsNil)
	b.Assert(d.Format(time.RFC3339), qt.Equals, "2024-05-01T00:00:00+02:00")

	b.Assert(decode("fromData").EXIF()["Copyright"].Value, qt.Equals, "© Data")
	b.Assert(decode("animAVIF").EXIF()["Copyright"].Value, qt.Equals, "Still AVIF")
}

func TestMetaFilterErrors(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		meta string
		want string
	}{
		{`(dict "copyrigth" "Bjørn")`, `unknown meta field "copyrigth"`},
		{`(dict "lat" 60)`, `lat and long must be set together`},
		{`(dict "lat" nil "long" 10)`, `lat and long must be set together`},
		{`(dict "lat" "" "long" 10)`, `lat and long must be set together`},
		{`(dict "lat" 91 "long" 10)`, `latitude 91 out of range`},
		{`(dict "date" "not a date")`, `invalid date`},
		{`(dict "creator" (slice "a" "b"))`, `invalid creator`},
		{`"not a map"`, `invalid meta options`},
	} {
		t.Run(test.want, func(t *testing.T) {
			t.Parallel()
			files := `
-- assets/sunset.jpg --
sourcefilename: ../testdata/sunset.jpg
-- layouts/home.html --
{{ with (resources.Get "sunset.jpg").Filter (images.Meta ` + test.meta + `) }}{{ .RelPermalink }}{{ end }}
`
			_, err := hugolib.TestE(t, files)
			qt.Assert(t, err, qt.ErrorMatches, `(?s).*`+test.want+`.*`)
		})
	}
}
