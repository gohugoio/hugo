package meta_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugolib"
)

func TestMeta(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[imaging.meta]
fields = ['**']
sources = ['exif', 'iptc', 'xmp']
-- assets/sunset.jpg --
sourcefilename: ../../testdata/sunset.jpg
-- layouts/home.html --
{{ $img := resources.Get "sunset.jpg" }}
{{ $meta := $img.Meta }}
{{ with $meta }}
Lat: {{ .Lat }}
Long: {{ .Long }}
Date: {{ .Date.Format "2006-01-02" }}
ExifMake: {{ .Exif.Make }}
ExifModel: {{ .Exif.Model }}
ExifISO: {{ .Exif.ISO }}
ExifArtist: {{ .Exif.Artist }}
IPTCCountry: {{ index .IPTC "Country-PrimaryLocationName" }}
IPTCProvince: {{ index .IPTC "Province-State" }}
IPTCKeywords: {{ .IPTC.Keywords }}
XMPCity: {{ .XMP.City }}
XMPCountry: {{ .XMP.Country }}
XMPCreator: {{ .XMP.Creator }}

{{ $exifAndIPTC := merge .Exif .IPTC }}
MergedExifIPTC_Make: {{ $exifAndIPTC.Make }}
MergedExifIPTC_Model: {{ $exifAndIPTC.Model }}
MergedExifIPTC_Keywords: {{ $exifAndIPTC.Keywords }}
MergedExifIPTC_ProvinceState: {{ index $exifAndIPTC "Province-State" }}
Same Type: {{ eq (printf "%T" $exifAndIPTC)  (printf "%T" .Exif) }}
{{ $all := merge .XMP .Exif .IPTC }}
MergedAll_Make: {{ $all.Make }}
MergedAll_City: {{ $all.City }}
MergedAll_Keywords: {{ $all.Keywords }}
MergedAll_Creator: {{ $all.Creator }}
{{ end }}

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"Lat: 36.597441",
		"Long: -4.50846",
		"Date: 2017-10-27",
		"ExifMake: RICOH IMAGING COMPANY, LTD.",
		"ExifModel: PENTAX K-3 II",
		"ExifISO: 100",
		"ExifArtist: bjorn.erik.pedersen@gmail.com",
		"IPTCCountry: Spain",
		"IPTCProvince: Andalucía",
		"IPTCKeywords: [Malaga Torremolinos]",
		"XMPCity: Benalmádena",
		"XMPCountry: Spain",
		"XMPCreator: bjorn.erik.pedersen@gmail.com",
		// merge .Exif .IPTC: contains keys from both, IPTC values take precedence
		"MergedExifIPTC_Make: RICOH IMAGING COMPANY, LTD.",
		"MergedExifIPTC_Model: PENTAX K-3 II",
		"MergedExifIPTC_Keywords: [Malaga Torremolinos]",
		"MergedExifIPTC_ProvinceState: Andalucía",
		"Same Type: true",
		// merge .XMP .Exif .IPTC: contains keys from all, rightmost (IPTC) wins
		"MergedAll_Make: RICOH IMAGING COMPANY, LTD.",
		"MergedAll_City: Benalmádena",
		"MergedAll_Keywords: [Malaga Torremolinos]",
		"MergedAll_Creator: bjorn.erik.pedersen@gmail.com",
	)
}

func TestMetaConfig(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[imaging]
[imaging.meta]
fields = ['! *{Model,ColorSpace,Metering}*']
sources = ['exif', 'iptc']
-- assets/sunset.jpg --
sourcefilename: ../../testdata/sunset.jpg
-- layouts/home.html --
{{ $img := resources.Get "sunset.jpg" }}
{{ $meta := $img.Meta }}
{{ with $meta }}
Lat: {{ .Lat }}
Long: {{ .Long }}
Date: {{ .Date.Format "2006-01-02" }}
ExifMake: {{ .Exif.Make }}
ExifModel: {{ with .Exif.Model }}{{ . }}{{ else }}EXCLUDED{{ end }}
IPTCCountry: {{ index .IPTC "Country-PrimaryLocationName" }}
XMPCity: {{ with .XMP.City }}{{ . }}{{ else }}NOT_IN_SOURCE{{ end }}
{{ end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"Lat: 36.597441",
		"Long: -4.50846",
		"Date: 2017-10-27",
		// Exif is included
		"ExifMake: RICOH IMAGING COMPANY, LTD.",
		// Model is excluded by fields pattern '! *Model*'
		"ExifModel: EXCLUDED",
		// IPTC is included
		"IPTCCountry: Spain",
		// XMP is not in sources, so should be empty
		"XMPCity: NOT_IN_SOURCE",
	)
}

func TestMetaFieldsFilterDateAndGPS(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[imaging.meta]
fields = ['! GPS*', '! *Date*', '! *Time*']
-- assets/sunset.jpg --
sourcefilename: ../../testdata/sunset.jpg
-- layouts/home.html --
{{ $img := resources.Get "sunset.jpg" }}
{{ $meta := $img.Meta }}
{{ with $meta }}
Lat: {{ .Lat }}
Long: {{ .Long }}
Date: {{ .Date.Format "2006-01-02" }}
ExifMake: {{ .Exif.Make }}
{{ end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"Lat: 0",
		"Long: 0",
		"Date: 0001-01-01",
		"ExifMake: RICOH IMAGING COMPANY, LTD.",
	)
}

// TestMetaXMPOnly verifies behavior when only XMP is configured as a source.
// Note: Date and Lat/Long are extracted from whichever configured source contains them.
// The test image has GPS/Date in EXIF only, so with XMP-only sources they will be empty.
// If the image had GPS/Date in XMP, they would be extracted (imagemeta v0.14.0+).
func TestMetaXMPOnly(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[imaging.meta]
sources = ['xmp']
-- assets/sunset.jpg --
sourcefilename: ../../testdata/sunset.jpg
-- layouts/home.html --
{{ $img := resources.Get "sunset.jpg" }}
{{ $meta := $img.Meta }}
{{ with $meta }}
Lat: {{ .Lat }}
Long: {{ .Long }}
Date: {{ .Date.Format "2006-01-02" }}
XMPCity: {{ .XMP.City }}
XMPCountry: {{ .XMP.Country }}
ExifMake: {{ with .Exif.Make }}{{ . }}{{ else }}NOT_IN_SOURCE{{ end }}
{{ end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		// The test image has Date/GPS only in EXIF, not in XMP.
		// With XMP-only sources, these will be empty/zero.
		"Lat: 0",
		"Long: 0",
		"Date: 0001-01-01",
		// XMP values are present
		"XMPCity: Benalmádena",
		"XMPCountry: Spain",
		// EXIF is not in sources
		"ExifMake: NOT_IN_SOURCE",
	)
}

func TestExifIsDeprecated(t *testing.T) {
	// This cannot be parallel.
	files := `
-- hugo.toml --
-- assets/sunset.jpg --
sourcefilename: ../../testdata/sunset.jpg
-- layouts/home.html --
Home.
{{ $img := resources.Get "sunset.jpg" }}
{{ $exif := $img.Exif }}
`
	b := hugolib.Test(t, files, hugolib.TestOptInfo())

	b.AssertLogContains("deprecated: Image.Exif")
}

func TestMetaInvalidSource(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[imaging.meta]
sources = ['foo']
-- assets/sunset.jpg --
sourcefilename: ../../testdata/sunset.jpg
-- layouts/home.html --
Home.
{{ $img := resources.Get "sunset.jpg" }}
{{ $exif := $img.Meta }}
`
	b, err := hugolib.TestE(t, files, hugolib.TestOptInfo())
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, `invalid metadata source "foo" in imaging.meta.sources config; must be one of [exif iptc xmp]`)
}

func TestAVIFMetaWidthAndHeight(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[imaging.meta]
fields = ['**']
sources = ['exif', 'iptc', 'xmp']
-- assets/sunset.avif --
sourcefilename: ../../testdata/sunset.avif
-- assets/sunset.jpg --
sourcefilename: ../../testdata/sunset.jpg
-- assets/giphy.gif --
sourcefilename: ../../testdata/giphy.gif
-- assets/pix.bmp --
sourcefilename: ../../testdata/pix.bmp
-- assets/mytext.txt --
This is a text file, not an image.
-- assets/mysvg.svg --
<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
  <rect width="100" height="100" fill="blue" />
</svg>
-- layouts/home.html --
{{ $txt := resources.Get "mytext.txt" }}
{{ $svg := resources.Get "mysvg.svg" }}
{{ $avif := resources.Get "sunset.avif" }}
{{ $jpg := resources.Get "sunset.jpg" }}
{{ $gif := resources.Get "giphy.gif" }}
{{ $bmp := resources.Get "pix.bmp" }}
{{ $ic := images.Config "/assets/sunset.avif" }}

$avif.Width/Height: {{ $avif.Width }}x{{ $avif.Height }}
$ic.Width/Height: {{ $ic.Width }}x{{ $ic.Height }}

{{ template "is-meta-etc" dict "what" "AVIF" "dot" $avif -}}
{{ template "is-meta-etc" dict "what" "JPG" "dot" $jpg -}}
{{ template "is-meta-etc" dict "what" "TXT" "dot" $txt -}}
{{ template "is-meta-etc" dict "what" "SVG" "dot" $svg -}}
{{ template "is-meta-etc" dict "what" "GIF" "dot" $gif -}}
{{ template "is-meta-etc" dict "what" "BMP" "dot" $bmp -}}
 
{{ $meta := $avif.Meta }}
Num Exif tags: {{ $meta.Exif | len }}|
{{ define "is-meta-etc"}}
IsImageResource {{ .what }}: {{ if reflect.IsImageResource .dot }}true{{ else }}false{{ end }}
IsImageResourceWithMeta {{ .what }}: {{ if reflect.IsImageResourceWithMeta .dot }}true{{ else }}false{{ end }}
IsImageResourceProcessable {{ .what }}: {{ if reflect.IsImageResourceProcessable .dot }}true{{ else }}false{{ end }}
{{ if reflect.IsImageResourceWithMeta .dot }}
Has width {{ .what }}: {{ if .dot.Width }}true{{ else }}false{{ end }}
Has meta {{ .what }}: {{ if .dot.Meta }}true{{ else }}false{{ end }}
{{ end }}
{{ end }}

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		`
$avif.Width/Height: 900x562
$ic.Width/Height: 900x562

IsImageResource AVIF: true
IsImageResourceWithMeta AVIF: true
IsImageResourceProcessable AVIF: false
Has width AVIF: true
Has meta AVIF: true

IsImageResource JPG: true
IsImageResourceWithMeta JPG: true
IsImageResourceProcessable JPG: true

Has width JPG: true
Has meta JPG: true

IsImageResource TXT: false
IsImageResourceWithMeta TXT: false
IsImageResourceProcessable TXT: false

IsImageResource SVG: true
IsImageResourceWithMeta SVG: false
IsImageResourceProcessable SVG: false

IsImageResource GIF: true
IsImageResourceWithMeta GIF: true
IsImageResourceProcessable GIF: true

Has width GIF: true
Has meta GIF: false

IsImageResource BMP: true
IsImageResourceWithMeta BMP: true
IsImageResourceProcessable BMP: true

Has width BMP: true
Has meta BMP: false

Num Exif tags: 52|
`,
	)
}
