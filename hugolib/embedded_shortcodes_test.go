// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/tpl"
	"github.com/spf13/viper"
)

const (
	baseURL = "http://foo/bar"
)

func TestShortcodeCrossrefs(t *testing.T) {
	for _, relative := range []bool{true, false} {
		doTestShortcodeCrossrefs(t, relative)
	}
}

func doTestShortcodeCrossrefs(t *testing.T, relative bool) {
	var refShortcode string
	var expectedBase string

	if relative {
		refShortcode = "relref"
		expectedBase = "/bar"
	} else {
		refShortcode = "ref"
		expectedBase = baseURL
	}

	path := filepath.FromSlash("blog/post.md")
	in := fmt.Sprintf(`{{< %s "%s" >}}`, refShortcode, path)
	expected := fmt.Sprintf(`%s/simple/url/`, expectedBase)

	templ := tpl.New()
	p, _ := pageFromString(simplePageWithURL, path)
	p.Node.Site = &SiteInfo{
		AllPages: &(Pages{p}),
		BaseURL:  template.URL(helpers.SanitizeURLKeepTrailingSlash(baseURL)),
	}

	output, err := HandleShortcodes(in, p, templ)

	if err != nil {
		t.Fatal("Handle shortcode error", err)
	}

	if output != expected {
		t.Errorf("Got\n%q\nExpected\n%q", output, expected)
	}
}

func TestShortcodeHighlight(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	if !helpers.HasPygments() {
		t.Skip("Skip test as Pygments is not installed")
	}
	viper.Set("PygmentsStyle", "bw")
	viper.Set("PygmentsUseClasses", false)

	for i, this := range []struct {
		in, expected string
	}{
		{`
{{< highlight java >}}
void do();
{{< /highlight >}}`,
			"(?s)^\n<div class=\"highlight\" style=\"background: #ffffff\"><pre style=\"line-height: 125%\">.*?void</span> do().*?</pre></div>\n$",
		},
		{`
{{< highlight java "style=friendly" >}}
void do();
{{< /highlight >}}`,
			"(?s)^\n<div class=\"highlight\" style=\"background: #f0f0f0\"><pre style=\"line-height: 125%\">.*?void</span>.*?do</span>.*?().*?</pre></div>\n$",
		},
	} {
		templ := tpl.New()
		p, _ := pageFromString(simplePage, "simple.md")
		output, err := HandleShortcodes(this.in, p, templ)

		if err != nil {
			t.Fatalf("[%d] Handle shortcode error", i)
		}

		matched, err := regexp.MatchString(this.expected, output)

		if err != nil {
			t.Fatalf("[%d] Regexp error", i)
		}

		if !matched {
			t.Errorf("[%d] Hightlight mismatch, got %s\n", i, output)
		}
	}
}

func TestShortcodeFigure(t *testing.T) {
	for i, this := range []struct {
		in, expected string
	}{
		{
			`{{< figure src="/img/hugo-logo.png" >}}`,
			"(?s)^\n<figure >.*?<img src=\"/img/hugo-logo.png\" />.*?</figure>\n$",
		},
		{
			// set alt
			`{{< figure src="/img/hugo-logo.png" alt="Hugo logo" >}}`,
			"(?s)^\n<figure >.*?<img src=\"/img/hugo-logo.png\" alt=\"Hugo logo\" />.*?</figure>\n$",
		},
		// set title
		{
			`{{< figure src="/img/hugo-logo.png" title="Hugo logo" >}}`,
			"(?s)^\n<figure >.*?<img src=\"/img/hugo-logo.png\" />.*?<figcaption>.*?<h4>Hugo logo</h4>.*?</figcaption>.*?</figure>\n$",
		},
		// set attr and attrlink
		{
			`{{< figure src="/img/hugo-logo.png" attr="Hugo logo" attrlink="/img/hugo-logo.png" >}}`,
			"(?s)^\n<figure >.*?<img src=\"/img/hugo-logo.png\" />.*?<figcaption>.*?<p>.*?<a href=\"/img/hugo-logo.png\">.*?Hugo logo.*?</a>.*?</p>.*?</figcaption>.*?</figure>\n$",
		},
	} {
		templ := tpl.New()
		p, _ := pageFromString(simplePage, "simple.md")
		output, err := HandleShortcodes(this.in, p, templ)

		matched, err := regexp.MatchString(this.expected, output)

		if err != nil {
			t.Fatalf("[%d] Regexp error", i)
		}

		if !matched {
			t.Errorf("[%d] Hightlight mismatch, got %s\n", i, output)
		}
	}
}

func TestShortcodeSpeakerdeck(t *testing.T) {
	for i, this := range []struct {
		in, expected string
	}{
		{
			`{{< speakerdeck 4e8126e72d853c0060001f97 >}}`,
			"(?s)^<script async class='speakerdeck-embed' data-id='4e8126e72d853c0060001f97'.*?>.*?</script>$",
		},
	} {
		templ := tpl.New()
		p, _ := pageFromString(simplePage, "simple.md")
		output, err := HandleShortcodes(this.in, p, templ)

		matched, err := regexp.MatchString(this.expected, output)

		if err != nil {
			t.Fatalf("[%d] Regexp error", i)
		}

		if !matched {
			t.Errorf("[%d] Hightlight mismatch, got %s\n", i, output)
		}
	}
}

func TestShortcodeYoutube(t *testing.T) {
	for i, this := range []struct {
		in, expected string
	}{
		{
			`{{< youtube w7Ft2ymGmfc >}}`,
			"(?s)^\n<div style=\".*?\">.*?<iframe src=\"//www.youtube.com/embed/w7Ft2ymGmfc\" style=\".*?\" allowfullscreen frameborder=\"0\">.*?</iframe>.*?</div>\n$",
		},
		// set class
		{
			`{{< youtube w7Ft2ymGmfc video>}}`,
			"(?s)^\n<div class=\"video\">.*?<iframe src=\"//www.youtube.com/embed/w7Ft2ymGmfc\" allowfullscreen frameborder=\"0\">.*?</iframe>.*?</div>\n$",
		},
		// set class and autoplay (using named params)
		{
			`{{< youtube id="w7Ft2ymGmfc" class="video" autoplay="true" >}}`,
			"(?s)^\n<div class=\"video\">.*?<iframe src=\"//www.youtube.com/embed/w7Ft2ymGmfc\\?autoplay=1\".*?allowfullscreen frameborder=\"0\">.*?</iframe>.*?</div>$",
		},
	} {
		templ := tpl.New()
		p, _ := pageFromString(simplePage, "simple.md")
		output, err := HandleShortcodes(this.in, p, templ)

		matched, err := regexp.MatchString(this.expected, output)

		if err != nil {
			t.Fatalf("[%d] Regexp error", i)
		}

		if !matched {
			t.Errorf("[%d] Hightlight mismatch, got %s\n", i, output)
		}
	}
}

func TestShortcodeVimeo(t *testing.T) {
	for i, this := range []struct {
		in, expected string
	}{
		{
			`{{< vimeo 146022717 >}}`,
			"(?s)^\n<div style=\".*?\">.*?<iframe src=\"//player.vimeo.com/video/146022717\" style=\".*?\" webkitallowfullscreen mozallowfullscreen allowfullscreen>.*?</iframe>.*?</div>\n$",
		},
		// set class
		{
			`{{< vimeo 146022717 video >}}`,
			"(?s)^\n<div class=\"video\">.*?<iframe src=\"//player.vimeo.com/video/146022717\" webkitallowfullscreen mozallowfullscreen allowfullscreen>.*?</iframe>.*?</div>\n$",
		},
		// set class (using named params)
		{
			`{{< vimeo id="146022717" class="video" >}}`,
			"(?s)^<div class=\"video\">.*?<iframe src=\"//player.vimeo.com/video/146022717\" webkitallowfullscreen mozallowfullscreen allowfullscreen>.*?</iframe>.*?</div>$",
		},
	} {
		templ := tpl.New()
		p, _ := pageFromString(simplePage, "simple.md")
		output, err := HandleShortcodes(this.in, p, templ)

		matched, err := regexp.MatchString(this.expected, output)

		if err != nil {
			t.Fatalf("[%d] Regexp error", i)
		}

		if !matched {
			t.Errorf("[%d] Hightlight mismatch, got %s\n", i, output)
		}
	}
}

func TestShortcodeGist(t *testing.T) {
	for i, this := range []struct {
		in, expected string
	}{
		{
			`{{< gist spf13 7896402 >}}`,
			"(?s)^<script src=\"//gist.github.com/spf13/7896402.js\"></script>$",
		},
		{
			`{{< gist spf13 7896402 "img.html" >}}`,
			"(?s)^<script src=\"//gist.github.com/spf13/7896402.js\\?file=img.html\"></script>$",
		},
	} {
		templ := tpl.New()
		p, _ := pageFromString(simplePage, "simple.md")
		output, err := HandleShortcodes(this.in, p, templ)

		matched, err := regexp.MatchString(this.expected, output)

		if err != nil {
			t.Fatalf("[%d] Regexp error", i)
		}

		if !matched {
			t.Errorf("[%d] Hightlight mismatch, got %s\n", i, output)
		}
	}
}

func TestShortcodeTweet(t *testing.T) {
	for i, this := range []struct {
		in, expected string
	}{
		{
			`{{< tweet 666616452582129664 >}}`,
			"(?s)^<blockquote class=\"twitter-tweet\"><p lang=\"en\" dir=\"ltr\">Hugo 0.15 will have 30%\\+ faster render times thanks to this commit <a href=\"https://t.co/FfzhM8bNhT\">https://t.co/FfzhM8bNhT</a>  <a href=\"https://twitter.com/hashtag/gohugo\\?src=hash\">#gohugo</a> <a href=\"https://twitter.com/hashtag/golang\\?src=hash\">#golang</a> <a href=\"https://t.co/ITbMNU2BUf\">https://t.co/ITbMNU2BUf</a></p>&mdash; Steve Francia \\(@spf13\\) <a href=\"https://twitter.com/spf13/status/666616452582129664\">November 17, 2015</a></blockquote>.*?<script async src=\"//platform.twitter.com/widgets.js\" charset=\"utf-8\"></script>$",
		},
	} {
		templ := tpl.New()
		p, _ := pageFromString(simplePage, "simple.md")
		cacheFileID := viper.GetString("CacheDir") + url.QueryEscape("https://api.twitter.com/1/statuses/oembed.json?id=666616452582129664")
		defer os.Remove(cacheFileID)
		output, err := HandleShortcodes(this.in, p, templ)

		matched, err := regexp.MatchString(this.expected, output)

		if err != nil {
			t.Fatalf("[%d] Regexp error", i)
		}

		if !matched {
			t.Errorf("[%d] Hightlight mismatch, got %s\n", i, output)
		}
	}
}
