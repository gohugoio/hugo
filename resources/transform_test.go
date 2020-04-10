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

package resources

import (
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/gohugoio/hugo/htesting"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/internal"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
)

const gopher = `iVBORw0KGgoAAAANSUhEUgAAAEsAAAA8CAAAAAALAhhPAAAFfUlEQVRYw62XeWwUVRzHf2+OPbo9d7tsWyiyaZti6eWGAhISoIGKECEKCAiJJkYTiUgTMYSIosYYBBIUIxoSPIINEBDi2VhwkQrVsj1ESgu9doHWdrul7ba73WNm3vOPtsseM9MdwvvrzTs+8/t95ze/33sI5BqiabU6m9En8oNjduLnAEDLUsQXFF8tQ5oxK3vmnNmDSMtrncks9Hhtt/qeWZapHb1ha3UqYSWVl2ZmpWgaXMXGohQAvmeop3bjTRtv6SgaK/Pb9/bFzUrYslbFAmHPp+3WhAYdr+7GN/YnpN46Opv55VDsJkoEpMrY/vO2BIYQ6LLvm0ThY3MzDzzeSJeeWNyTkgnIE5ePKsvKlcg/0T9QMzXalwXMlj54z4c0rh/mzEfr+FgWEz2w6uk8dkzFAgcARAgNp1ZYef8bH2AgvuStbc2/i6CiWGj98y2tw2l4FAXKkQBIf+exyRnteY83LfEwDQAYCoK+P6bxkZm/0966LxcAAILHB56kgD95PPxltuYcMtFTWw/FKkY/6Opf3GGd9ZF+Qp6mzJxzuRSractOmJrH1u8XTvWFHINNkLQLMR+XHXvfPPHw967raE1xxwtA36IMRfkAAG29/7mLuQcb2WOnsJReZGfpiHsSBX81cvMKywYZHhX5hFPtOqPGWZCXnhWGAu6lX91ElKXSalcLXu3UaOXVay57ZSe5f6Gpx7J2MXAsi7EqSp09b/MirKSyJfnfEEgeDjl8FgDAfvewP03zZ+AJ0m9aFRM8eEHBDRKjfcreDXnZdQuAxXpT2NRJ7xl3UkLBhuVGU16gZiGOgZmrSbRdqkILuL/yYoSXHHkl9KXgqNu3PB8oRg0geC5vFmLjad6mUyTKLmF3OtraWDIfACyXqmephaDABawfpi6tqqBZytfQMqOz6S09iWXhktrRaB8Xz4Yi/8gyABDm5NVe6qq/3VzPrcjELWrebVuyY2T7ar4zQyybUCtsQ5Es1FGaZVrRVQwAgHGW2ZCRZshI5bGQi7HesyE972pOSeMM0dSktlzxRdrlqb3Osa6CCS8IJoQQQgBAbTAa5l5epO34rJszibJI8rxLfGzcp1dRosutGeb2VDNgqYrwTiPNsLxXiPi3dz7LiS1WBRBDBOnqEjyy3aQb+/bLiJzz9dIkscVBBLxMfSEac7kO4Fpkngi0ruNBeSOal+u8jgOuqPz12nryMLCniEjtOOOmpt+KEIqsEdocJjYXwrh9OZqWJQyPCTo67LNS/TdxLAv6R5ZNK9npEjbYdT33gRo4o5oTqR34R+OmaSzDBWsAIPhuRcgyoteNi9gF0KzNYWVItPf2TLoXEg+7isNC7uJkgo1iQWOfRSP9NR11RtbZZ3OMG/VhL6jvx+J1m87+RCfJChAtEBQkSBX2PnSiihc/Twh3j0h7qdYQAoRVsRGmq7HU2QRbaxVGa1D6nIOqaIWRjyRZpHMQKWKpZM5feA+lzC4ZFultV8S6T0mzQGhQohi5I8iw+CsqBSxhFMuwyLgSwbghGb0AiIKkSDmGZVmJSiKihsiyOAUs70UkywooYP0bii9GdH4sfr1UNysd3fUyLLMQN+rsmo3grHl9VNJHbbwxoa47Vw5gupIqrZcjPh9R4Nye3nRDk199V+aetmvVtDRE8/+cbgAAgMIWGb3UA0MGLE9SCbWX670TDy1y98c3D27eppUjsZ6fql3jcd5rUe7+ZIlLNQny3Rd+E5Tct3WVhTM5RBCEdiEK0b6B+/ca2gYU393nFj/n1AygRQxPIUA043M42u85+z2SnssKrPl8Mx76NL3E6eXc3be7OD+H4WHbJkKI8AU8irbITQjZ+0hQcPEgId/Fn/pl9crKH02+5o2b9T/eMx7pKoskYgAAAABJRU5ErkJggg==`

func gopherPNG() io.Reader { return base64.NewDecoder(base64.StdEncoding, strings.NewReader(gopher)) }

func TestTransform(t *testing.T) {
	c := qt.New(t)

	createTransformer := func(spec *Spec, filename, content string) Transformer {
		filename = filepath.FromSlash(filename)
		fs := spec.Fs.Source
		afero.WriteFile(fs, filename, []byte(content), 0777)
		r, _ := spec.New(ResourceSourceDescriptor{Fs: fs, SourceFilename: filename})
		return r.(Transformer)
	}

	createContentReplacer := func(name, old, new string) ResourceTransformation {
		return &testTransformation{
			name: name,
			transform: func(ctx *ResourceTransformationCtx) error {
				in := helpers.ReaderToString(ctx.From)
				in = strings.Replace(in, old, new, 1)
				ctx.AddOutPathIdentifier("." + name)
				fmt.Fprint(ctx.To, in)
				return nil
			},
		}
	}

	// Verify that we publish the same file once only.
	assertNoDuplicateWrites := func(c *qt.C, spec *Spec) {
		c.Helper()
		d := spec.Fs.Destination.(hugofs.DuplicatesReporter)
		c.Assert(d.ReportDuplicates(), qt.Equals, "")
	}

	assertShouldExist := func(c *qt.C, spec *Spec, filename string, should bool) {
		c.Helper()
		exists, _ := helpers.Exists(filepath.FromSlash(filename), spec.Fs.Destination)
		c.Assert(exists, qt.Equals, should)
	}

	c.Run("All values", func(c *qt.C) {
		c.Parallel()

		spec := newTestResourceSpec(specDescriptor{c: c})

		transformation := &testTransformation{
			name: "test",
			transform: func(ctx *ResourceTransformationCtx) error {
				// Content
				in := helpers.ReaderToString(ctx.From)
				in = strings.Replace(in, "blue", "green", 1)
				fmt.Fprint(ctx.To, in)

				// Media type
				ctx.OutMediaType = media.CSVType

				// Change target
				ctx.ReplaceOutPathExtension(".csv")

				// Add some data to context
				ctx.Data["mydata"] = "Hugo Rocks!"

				return nil
			},
		}

		r := createTransformer(spec, "f1.txt", "color is blue")

		tr, err := r.Transform(transformation)
		c.Assert(err, qt.IsNil)
		content, err := tr.(resource.ContentProvider).Content()
		c.Assert(err, qt.IsNil)

		c.Assert(content, qt.Equals, "color is green")
		c.Assert(tr.MediaType(), eq, media.CSVType)
		c.Assert(tr.RelPermalink(), qt.Equals, "/f1.csv")
		assertShouldExist(c, spec, "public/f1.csv", true)

		data := tr.Data().(map[string]interface{})
		c.Assert(data["mydata"], qt.Equals, "Hugo Rocks!")

		assertNoDuplicateWrites(c, spec)
	})

	c.Run("Meta only", func(c *qt.C) {
		c.Parallel()

		spec := newTestResourceSpec(specDescriptor{c: c})

		transformation := &testTransformation{
			name: "test",
			transform: func(ctx *ResourceTransformationCtx) error {
				// Change media type only
				ctx.OutMediaType = media.CSVType
				ctx.ReplaceOutPathExtension(".csv")

				return nil
			},
		}

		r := createTransformer(spec, "f1.txt", "color is blue")

		tr, err := r.Transform(transformation)
		c.Assert(err, qt.IsNil)
		content, err := tr.(resource.ContentProvider).Content()
		c.Assert(err, qt.IsNil)

		c.Assert(content, qt.Equals, "color is blue")
		c.Assert(tr.MediaType(), eq, media.CSVType)

		// The transformed file should only be published if RelPermalink
		// or Permalink is called.
		n := htesting.RandIntn(3)
		shouldExist := true
		switch n {
		case 0:
			tr.RelPermalink()
		case 1:
			tr.Permalink()
		default:
			shouldExist = false
		}

		assertShouldExist(c, spec, "public/f1.csv", shouldExist)
		assertNoDuplicateWrites(c, spec)
	})

	c.Run("Memory-cached transformation", func(c *qt.C) {
		c.Parallel()

		spec := newTestResourceSpec(specDescriptor{c: c})

		// Two transformations with same id, different behaviour.
		t1 := createContentReplacer("t1", "blue", "green")
		t2 := createContentReplacer("t1", "color", "car")

		for i, transformation := range []ResourceTransformation{t1, t2} {
			r := createTransformer(spec, "f1.txt", "color is blue")
			tr, _ := r.Transform(transformation)
			content, err := tr.(resource.ContentProvider).Content()
			c.Assert(err, qt.IsNil)
			c.Assert(content, qt.Equals, "color is green", qt.Commentf("i=%d", i))

			assertShouldExist(c, spec, "public/f1.t1.txt", false)
		}

		assertNoDuplicateWrites(c, spec)
	})

	c.Run("File-cached transformation", func(c *qt.C) {
		c.Parallel()

		fs := afero.NewMemMapFs()

		for i := 0; i < 2; i++ {
			spec := newTestResourceSpec(specDescriptor{c: c, fs: fs})

			r := createTransformer(spec, "f1.txt", "color is blue")

			var transformation ResourceTransformation

			if i == 0 {
				// There is currently a hardcoded list of transformations that we
				// persist to disk (tocss, postcss).
				transformation = &testTransformation{
					name: "tocss",
					transform: func(ctx *ResourceTransformationCtx) error {
						in := helpers.ReaderToString(ctx.From)
						in = strings.Replace(in, "blue", "green", 1)
						ctx.AddOutPathIdentifier("." + "cached")
						ctx.OutMediaType = media.CSVType
						ctx.Data = map[string]interface{}{
							"Hugo": "Rocks!",
						}
						fmt.Fprint(ctx.To, in)
						return nil
					},
				}
			} else {
				// Force read from file cache.
				transformation = &testTransformation{
					name: "tocss",
					transform: func(ctx *ResourceTransformationCtx) error {
						return herrors.ErrFeatureNotAvailable
					},
				}
			}

			msg := qt.Commentf("i=%d", i)

			tr, _ := r.Transform(transformation)
			c.Assert(tr.RelPermalink(), qt.Equals, "/f1.cached.txt", msg)
			content, err := tr.(resource.ContentProvider).Content()
			c.Assert(err, qt.IsNil)
			c.Assert(content, qt.Equals, "color is green", msg)
			c.Assert(tr.MediaType(), eq, media.CSVType)
			c.Assert(tr.Data(), qt.DeepEquals, map[string]interface{}{
				"Hugo": "Rocks!",
			})

			assertNoDuplicateWrites(c, spec)
			assertShouldExist(c, spec, "public/f1.cached.txt", true)

		}
	})

	c.Run("Access RelPermalink first", func(c *qt.C) {
		c.Parallel()

		spec := newTestResourceSpec(specDescriptor{c: c})

		t1 := createContentReplacer("t1", "blue", "green")

		r := createTransformer(spec, "f1.txt", "color is blue")

		tr, _ := r.Transform(t1)

		relPermalink := tr.RelPermalink()

		content, err := tr.(resource.ContentProvider).Content()
		c.Assert(err, qt.IsNil)

		c.Assert(relPermalink, qt.Equals, "/f1.t1.txt")
		c.Assert(content, qt.Equals, "color is green")
		c.Assert(tr.MediaType(), eq, media.TextType)

		assertNoDuplicateWrites(c, spec)
		assertShouldExist(c, spec, "public/f1.t1.txt", true)
	})

	c.Run("Content two", func(c *qt.C) {
		c.Parallel()

		spec := newTestResourceSpec(specDescriptor{c: c})

		t1 := createContentReplacer("t1", "blue", "green")
		t2 := createContentReplacer("t1", "color", "car")

		r := createTransformer(spec, "f1.txt", "color is blue")

		tr, _ := r.Transform(t1, t2)
		content, err := tr.(resource.ContentProvider).Content()
		c.Assert(err, qt.IsNil)

		c.Assert(content, qt.Equals, "car is green")
		c.Assert(tr.MediaType(), eq, media.TextType)

		assertNoDuplicateWrites(c, spec)
	})

	c.Run("Content two chained", func(c *qt.C) {
		c.Parallel()

		spec := newTestResourceSpec(specDescriptor{c: c})

		t1 := createContentReplacer("t1", "blue", "green")
		t2 := createContentReplacer("t2", "color", "car")

		r := createTransformer(spec, "f1.txt", "color is blue")

		tr1, _ := r.Transform(t1)
		tr2, _ := tr1.Transform(t2)

		content1, err := tr1.(resource.ContentProvider).Content()
		c.Assert(err, qt.IsNil)
		content2, err := tr2.(resource.ContentProvider).Content()
		c.Assert(err, qt.IsNil)

		c.Assert(content1, qt.Equals, "color is green")
		c.Assert(content2, qt.Equals, "car is green")

		assertNoDuplicateWrites(c, spec)
	})

	c.Run("Content many", func(c *qt.C) {
		c.Parallel()

		spec := newTestResourceSpec(specDescriptor{c: c})

		const count = 26 // A-Z

		transformations := make([]ResourceTransformation, count)
		for i := 0; i < count; i++ {
			transformations[i] = createContentReplacer(fmt.Sprintf("t%d", i), fmt.Sprint(i), string(rune(i+65)))
		}

		var countstr strings.Builder
		for i := 0; i < count; i++ {
			countstr.WriteString(fmt.Sprint(i))
		}

		r := createTransformer(spec, "f1.txt", countstr.String())

		tr, _ := r.Transform(transformations...)
		content, err := tr.(resource.ContentProvider).Content()
		c.Assert(err, qt.IsNil)

		c.Assert(content, qt.Equals, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")

		assertNoDuplicateWrites(c, spec)
	})

	c.Run("Image", func(c *qt.C) {
		c.Parallel()

		spec := newTestResourceSpec(specDescriptor{c: c})

		transformation := &testTransformation{
			name: "test",
			transform: func(ctx *ResourceTransformationCtx) error {
				ctx.AddOutPathIdentifier(".changed")
				return nil
			},
		}

		r := createTransformer(spec, "gopher.png", helpers.ReaderToString(gopherPNG()))

		tr, err := r.Transform(transformation)
		c.Assert(err, qt.IsNil)
		c.Assert(tr.MediaType(), eq, media.PNGType)

		img, ok := tr.(resource.Image)
		c.Assert(ok, qt.Equals, true)

		c.Assert(img.Width(), qt.Equals, 75)
		c.Assert(img.Height(), qt.Equals, 60)

		// RelPermalink called.
		resizedPublished1, err := img.Resize("40x40")
		c.Assert(err, qt.IsNil)
		c.Assert(resizedPublished1.Height(), qt.Equals, 40)
		c.Assert(resizedPublished1.RelPermalink(), qt.Equals, "/gopher.changed_hu2e827f5a78333ebc04166dd643235dea_1462_40x40_resize_linear_2.png")
		assertShouldExist(c, spec, "public/gopher.changed_hu2e827f5a78333ebc04166dd643235dea_1462_40x40_resize_linear_2.png", true)

		// Permalink called.
		resizedPublished2, err := img.Resize("30x30")
		c.Assert(err, qt.IsNil)
		c.Assert(resizedPublished2.Height(), qt.Equals, 30)
		c.Assert(resizedPublished2.Permalink(), qt.Equals, "https://example.com/gopher.changed_hu2e827f5a78333ebc04166dd643235dea_1462_30x30_resize_linear_2.png")
		assertShouldExist(c, spec, "public/gopher.changed_hu2e827f5a78333ebc04166dd643235dea_1462_30x30_resize_linear_2.png", true)

		// Not published because none of RelPermalink or Permalink was called.
		resizedNotPublished, err := img.Resize("50x50")
		c.Assert(err, qt.IsNil)
		c.Assert(resizedNotPublished.Height(), qt.Equals, 50)
		//c.Assert(resized.RelPermalink(), qt.Equals, "/gopher.changed_hu2e827f5a78333ebc04166dd643235dea_1462_50x50_resize_linear_2.png")
		assertShouldExist(c, spec, "public/gopher.changed_hu2e827f5a78333ebc04166dd643235dea_1462_50x50_resize_linear_2.png", false)

		assertNoDuplicateWrites(c, spec)

	})

	c.Run("Concurrent", func(c *qt.C) {
		spec := newTestResourceSpec(specDescriptor{c: c})

		transformers := make([]Transformer, 10)
		transformations := make([]ResourceTransformation, 10)

		for i := 0; i < 10; i++ {
			transformers[i] = createTransformer(spec, fmt.Sprintf("f%d.txt", i), fmt.Sprintf("color is %d", i))
			transformations[i] = createContentReplacer("test", strconv.Itoa(i), "blue")
		}

		var wg sync.WaitGroup

		for i := 0; i < 13; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < 23; j++ {
					id := (i + j) % 10
					tr, err := transformers[id].Transform(transformations[id])
					c.Assert(err, qt.IsNil)
					content, err := tr.(resource.ContentProvider).Content()
					c.Assert(err, qt.IsNil)
					c.Assert(content, qt.Equals, "color is blue")
					c.Assert(tr.RelPermalink(), qt.Equals, fmt.Sprintf("/f%d.test.txt", id))
				}
			}(i)
		}
		wg.Wait()

		assertNoDuplicateWrites(c, spec)
	})
}

type testTransformation struct {
	name      string
	transform func(ctx *ResourceTransformationCtx) error
}

func (t *testTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey(t.name)
}

func (t *testTransformation) Transform(ctx *ResourceTransformationCtx) error {
	return t.transform(ctx)
}
