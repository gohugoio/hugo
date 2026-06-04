package warpc_test

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

func TestAvifBasic(t *testing.T) {
	files := `
-- hugo.toml --
-- assets/sunset.avif --
sourcefilename: ../../resources/testdata/bep/dock-75-hdr.avif
-- layouts/home.html --
{{ $image := resources.Get "sunset.avif" }}
Width/Height: {{ $image.Width }}/{{ $image.Height }}|
Decode avif: {{ $jpeg := $image.Process "jpeg" }}|{{ $jpeg.RelPermalink }}|
Encode avif from JPEG: {{ $avif := $jpeg.Process "avif" }}|{{ $avif.RelPermalink }}|
Encode avif from avif: {{ $avif := $image.Process "avif" }}|{{ $avif.RelPermalink }}|
`

	b, err := hugolib.TestE(t, files)
	b.Assert(err, qt.IsNil)
	b.AssertFileContent("public/index.html", "Width/Height: 1024/683")
}

// giphy.avif is an animated AVIF transcoded from giphy.gif (14 frames @ 200ms,
// infinite loop) via avifenc. Verifies that decode preserves animation, the
// pipeline runs every frame through resize, the GIF encoder writes a multi-frame
// GIF, and the libavif "repetitionCount" → Go "LoopCount" inversion (libavif
// -1 = infinite, image/gif 0 = infinite) is handled at the AVIF decoder boundary.
func TestAvifAnimatedToGif(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["page", "section", "taxonomy", "term", "sitemap", "robotsTXT", "404"]
-- assets/giphy.avif --
sourcefilename: ../../resources/testdata/giphy.avif
-- layouts/home.html --
{{ $img := resources.Get "giphy.avif" }}
{{ $gif := $img.Resize "x60 gif" }}
{{ $gif.Publish }}
gif:{{ $gif.RelPermalink }}
`

	b := hugolib.Test(t, files)
	rel := b.FileContent("public/index.html")

	var gifFile string
	for line := range strings.SplitSeq(rel, "\n") {
		if rest, ok := strings.CutPrefix(line, "gif:"); ok {
			gifFile = "public" + rest
			break
		}
	}
	b.Assert(gifFile, qt.Not(qt.Equals), "")

	durations := make([]int, 14)
	for i := range durations {
		durations[i] = 200
	}
	b.ImageHelper(gifFile).
		AssertFormat("gif").
		AssertIsAnimated(true).
		AssertLoopCount(0).
		AssertFrameDurations(durations)
}

// See issue 14987.
func TestAvifEncodeHintSubsampling(t *testing.T) {
	files := `
-- hugo.toml --
-- assets/logo.png --
sourcefilename: ../../resources/testdata/gohugoio24.png
-- layouts/home.html --
{{ $img := resources.Get "logo.png" }}
default:{{ ($img.Process "avif photo").RelPermalink }}|
photo:{{ ($img.Process "avif photo").RelPermalink }}|
text:{{ ($img.Process "avif text").RelPermalink }}|
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html",
		`
default:/logo_hu_72e150ce03376cfd.avif|
photo:/logo_hu_72e150ce03376cfd.avif|
text:/logo_hu_9de8f282912bc925.avif|
`)
}

// See issue 14985.
func TestAvifEncodeOutOfMemory(t *testing.T) {
	files := `
-- assets/gopher.png --
sourcefilename: ../../resources/testdata/bw-gopher.png
-- layouts/home.html --
{{ $img := resources.Get "gopher.png" }}
{{ $r := try ($img.Resize "3000x3000 avif") }}
{{ with $r.Err }}BigErr: {{ . }}|{{ else }}BigOK|{{ end }}
{{ $small := $img.Resize "32x32 avif" }}
SmallAfter: {{ $small.RelPermalink }}|
`

	b := hugolib.Test(t, files, hugolib.TestOptWithConfig(func(c *hugolib.IntegrationTestConfig) {
		c.WarpcMemory = 8
	}))

	// The big resize must fail gracefully (caught by try) ...
	b.AssertFileContent("public/index.html",
		"BigErr:",
		"out of memory allocating",
		"for blob data",
		// ... and the dispatcher must still process images afterwards.
		"SmallAfter: /gopher_",
	)
}
