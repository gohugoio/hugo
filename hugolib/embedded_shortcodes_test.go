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
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"testing"

	"github.com/spf13/cast"

	"path/filepath"

	"github.com/gohugoio/hugo/deps"

	qt "github.com/frankban/quicktest"
)

const (
	testBaseURL = "http://foo/bar"
)

func TestShortcodeCrossrefs(t *testing.T) {
	t.Parallel()

	for _, relative := range []bool{true, false} {
		doTestShortcodeCrossrefs(t, relative)
	}
}

func doTestShortcodeCrossrefs(t *testing.T, relative bool) {
	var (
		cfg, fs = newTestCfg()
		c       = qt.New(t)
	)

	cfg.Set("baseURL", testBaseURL)

	var refShortcode string
	var expectedBase string

	if relative {
		refShortcode = "relref"
		expectedBase = "/bar"
	} else {
		refShortcode = "ref"
		expectedBase = testBaseURL
	}

	path := filepath.FromSlash("blog/post.md")
	in := fmt.Sprintf(`{{< %s "%s" >}}`, refShortcode, path)

	writeSource(t, fs, "content/"+path, simplePageWithURL+": "+in)

	expected := fmt.Sprintf(`%s/simple/url/`, expectedBase)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	c.Assert(len(s.RegularPages()), qt.Equals, 1)

	content, err := s.RegularPages()[0].Content()
	c.Assert(err, qt.IsNil)
	output := cast.ToString(content)

	if !strings.Contains(output, expected) {
		t.Errorf("Got\n%q\nExpected\n%q", output, expected)
	}
}

func TestShortcodeHighlight(t *testing.T) {
	t.Parallel()

	for _, this := range []struct {
		in, expected string
	}{
		{`{{< highlight java >}}
void do();
{{< /highlight >}}`,
			`(?s)<div class="highlight"><pre style="background-color:#fff;-moz-tab-size:4;-o-tab-size:4;tab-size:4"><code class="language-java"`,
		},
		{`{{< highlight java "style=friendly" >}}
void do();
{{< /highlight >}}`,
			`(?s)<div class="highlight"><pre style="background-color:#f0f0f0;-moz-tab-size:4;-o-tab-size:4;tab-size:4"><code class="language-java" data-lang="java">`,
		},
	} {

		var (
			cfg, fs = newTestCfg()
			th      = newTestHelper(cfg, fs, t)
		)

		cfg.Set("pygmentsStyle", "bw")
		cfg.Set("pygmentsUseClasses", false)

		writeSource(t, fs, filepath.Join("content", "simple.md"), fmt.Sprintf(`---
title: Shorty
---
%s`, this.in))
		writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), `{{ .Content }}`)

		buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

		th.assertFileContentRegexp(filepath.Join("public", "simple", "index.html"), this.expected)

	}
}

func TestShortcodeFigure(t *testing.T) {
	t.Parallel()

	for _, this := range []struct {
		in, expected string
	}{
		{
			`{{< figure src="/img/hugo-logo.png" >}}`,
			"(?s)<figure>.*?<img src=\"/img/hugo-logo.png\"/>.*?</figure>",
		},
		{
			// set alt
			`{{< figure src="/img/hugo-logo.png" alt="Hugo logo" >}}`,
			"(?s)<figure>.*?<img src=\"/img/hugo-logo.png\".+?alt=\"Hugo logo\"/>.*?</figure>",
		},
		// set title
		{
			`{{< figure src="/img/hugo-logo.png" title="Hugo logo" >}}`,
			"(?s)<figure>.*?<img src=\"/img/hugo-logo.png\"/>.*?<figcaption>.*?<h4>Hugo logo</h4>.*?</figcaption>.*?</figure>",
		},
		// set attr and attrlink
		{
			`{{< figure src="/img/hugo-logo.png" attr="Hugo logo" attrlink="/img/hugo-logo.png" >}}`,
			"(?s)<figure>.*?<img src=\"/img/hugo-logo.png\"/>.*?<figcaption>.*?<p>.*?<a href=\"/img/hugo-logo.png\">.*?Hugo logo.*?</a>.*?</p>.*?</figcaption>.*?</figure>",
		},
	} {

		var (
			cfg, fs = newTestCfg()
			th      = newTestHelper(cfg, fs, t)
		)

		writeSource(t, fs, filepath.Join("content", "simple.md"), fmt.Sprintf(`---
title: Shorty
---
%s`, this.in))
		writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), `{{ .Content }}`)

		buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

		th.assertFileContentRegexp(filepath.Join("public", "simple", "index.html"), this.expected)

	}
}

func TestShortcodeYoutube(t *testing.T) {
	t.Parallel()

	for _, this := range []struct {
		in, expected string
	}{
		{
			`{{< youtube w7Ft2ymGmfc >}}`,
			"(?s)\n<div style=\".*?\">.*?<iframe src=\"https://www.youtube.com/embed/w7Ft2ymGmfc\" style=\".*?\" allowfullscreen title=\"YouTube Video\">.*?</iframe>.*?</div>\n",
		},
		// set class
		{
			`{{< youtube w7Ft2ymGmfc video>}}`,
			"(?s)\n<div class=\"video\">.*?<iframe src=\"https://www.youtube.com/embed/w7Ft2ymGmfc\" allowfullscreen title=\"YouTube Video\">.*?</iframe>.*?</div>\n",
		},
		// set class and autoplay (using named params)
		{
			`{{< youtube id="w7Ft2ymGmfc" class="video" autoplay="true" >}}`,
			"(?s)\n<div class=\"video\">.*?<iframe src=\"https://www.youtube.com/embed/w7Ft2ymGmfc\\?autoplay=1\".*?allowfullscreen title=\"YouTube Video\">.*?</iframe>.*?</div>",
		},
	} {
		var (
			cfg, fs = newTestCfg()
			th      = newTestHelper(cfg, fs, t)
		)

		writeSource(t, fs, filepath.Join("content", "simple.md"), fmt.Sprintf(`---
title: Shorty
---
%s`, this.in))
		writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), `{{ .Content }}`)

		buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

		th.assertFileContentRegexp(filepath.Join("public", "simple", "index.html"), this.expected)
	}

}

func TestShortcodeVimeo(t *testing.T) {
	t.Parallel()

	for _, this := range []struct {
		in, expected string
	}{
		{
			`{{< vimeo 146022717 >}}`,
			"(?s)\n<div style=\".*?\">.*?<iframe src=\"https://player.vimeo.com/video/146022717\" style=\".*?\" title=\"vimeo video\" webkitallowfullscreen mozallowfullscreen allowfullscreen>.*?</iframe>.*?</div>\n",
		},
		// set class
		{
			`{{< vimeo 146022717 video >}}`,
			"(?s)\n<div class=\"video\">.*?<iframe src=\"https://player.vimeo.com/video/146022717\" title=\"vimeo video\" webkitallowfullscreen mozallowfullscreen allowfullscreen>.*?</iframe>.*?</div>\n",
		},
		// set vimeo title
		{
			`{{< vimeo 146022717 video my-title >}}`,
			"(?s)\n<div class=\"video\">.*?<iframe src=\"https://player.vimeo.com/video/146022717\" title=\"my-title\" webkitallowfullscreen mozallowfullscreen allowfullscreen>.*?</iframe>.*?</div>\n",
		},
		// set class (using named params)
		{
			`{{< vimeo id="146022717" class="video" >}}`,
			"(?s)^<div class=\"video\">.*?<iframe src=\"https://player.vimeo.com/video/146022717\" title=\"vimeo video\" webkitallowfullscreen mozallowfullscreen allowfullscreen>.*?</iframe>.*?</div>",
		},
		// set vimeo title (using named params)
		{
			`{{< vimeo id="146022717" class="video" title="my vimeo video" >}}`,
			"(?s)^<div class=\"video\">.*?<iframe src=\"https://player.vimeo.com/video/146022717\" title=\"my vimeo video\" webkitallowfullscreen mozallowfullscreen allowfullscreen>.*?</iframe>.*?</div>",
		},
	} {
		var (
			cfg, fs = newTestCfg()
			th      = newTestHelper(cfg, fs, t)
		)

		writeSource(t, fs, filepath.Join("content", "simple.md"), fmt.Sprintf(`---
title: Shorty
---
%s`, this.in))
		writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), `{{ .Content }}`)

		buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

		th.assertFileContentRegexp(filepath.Join("public", "simple", "index.html"), this.expected)

	}
}

func TestShortcodeGist(t *testing.T) {
	t.Parallel()

	for _, this := range []struct {
		in, expected string
	}{
		{
			`{{< gist spf13 7896402 >}}`,
			"(?s)^<script type=\"application/javascript\" src=\"https://gist.github.com/spf13/7896402.js\"></script>",
		},
		{
			`{{< gist spf13 7896402 "img.html" >}}`,
			"(?s)^<script type=\"application/javascript\" src=\"https://gist.github.com/spf13/7896402.js\\?file=img.html\"></script>",
		},
	} {
		var (
			cfg, fs = newTestCfg()
			th      = newTestHelper(cfg, fs, t)
		)

		writeSource(t, fs, filepath.Join("content", "simple.md"), fmt.Sprintf(`---
title: Shorty
---
%s`, this.in))
		writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), `{{ .Content }}`)

		buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

		th.assertFileContentRegexp(filepath.Join("public", "simple", "index.html"), this.expected)

	}
}

func TestShortcodeTweet(t *testing.T) {
	t.Parallel()

	for i, this := range []struct {
		privacy            map[string]interface{}
		in, resp, expected string
	}{
		{
			map[string]interface{}{
				"twitter": map[string]interface{}{
					"simple": true,
				},
			},
			`{{< tweet 666616452582129664 >}}`,
			`{"url":"https:\/\/twitter.com\/spf13\/status\/666616452582129664","author_name":"Steve Francia","author_url":"https:\/\/twitter.com\/spf13","html":"\u003Cblockquote class=\"twitter-tweet\"\u003E\u003Cp lang=\"en\" dir=\"ltr\"\u003EHugo 0.15 will have 30%+ faster render times thanks to this commit \u003Ca href=\"https:\/\/t.co\/FfzhM8bNhT\"\u003Ehttps:\/\/t.co\/FfzhM8bNhT\u003C\/a\u003E  \u003Ca href=\"https:\/\/twitter.com\/hashtag\/gohugo?src=hash\"\u003E#gohugo\u003C\/a\u003E \u003Ca href=\"https:\/\/twitter.com\/hashtag\/golang?src=hash\"\u003E#golang\u003C\/a\u003E \u003Ca href=\"https:\/\/t.co\/ITbMNU2BUf\"\u003Ehttps:\/\/t.co\/ITbMNU2BUf\u003C\/a\u003E\u003C\/p\u003E&mdash; Steve Francia (@spf13) \u003Ca href=\"https:\/\/twitter.com\/spf13\/status\/666616452582129664\"\u003ENovember 17, 2015\u003C\/a\u003E\u003C\/blockquote\u003E\n\u003Cscript async src=\"\/\/platform.twitter.com\/widgets.js\" charset=\"utf-8\"\u003E\u003C\/script\u003E","width":550,"height":null,"type":"rich","cache_age":"3153600000","provider_name":"Twitter","provider_url":"https:\/\/twitter.com","version":"1.0"}`,
			`.twitter-tweet a`,
		},
		{
			map[string]interface{}{
				"twitter": map[string]interface{}{
					"simple": false,
				},
			},
			`{{< tweet 666616452582129664 >}}`,
			`{"url":"https:\/\/twitter.com\/spf13\/status\/666616452582129664","author_name":"Steve Francia","author_url":"https:\/\/twitter.com\/spf13","html":"\u003Cblockquote class=\"twitter-tweet\"\u003E\u003Cp lang=\"en\" dir=\"ltr\"\u003EHugo 0.15 will have 30%+ faster render times thanks to this commit \u003Ca href=\"https:\/\/t.co\/FfzhM8bNhT\"\u003Ehttps:\/\/t.co\/FfzhM8bNhT\u003C\/a\u003E  \u003Ca href=\"https:\/\/twitter.com\/hashtag\/gohugo?src=hash\"\u003E#gohugo\u003C\/a\u003E \u003Ca href=\"https:\/\/twitter.com\/hashtag\/golang?src=hash\"\u003E#golang\u003C\/a\u003E \u003Ca href=\"https:\/\/t.co\/ITbMNU2BUf\"\u003Ehttps:\/\/t.co\/ITbMNU2BUf\u003C\/a\u003E\u003C\/p\u003E&mdash; Steve Francia (@spf13) \u003Ca href=\"https:\/\/twitter.com\/spf13\/status\/666616452582129664\"\u003ENovember 17, 2015\u003C\/a\u003E\u003C\/blockquote\u003E\n\u003Cscript async src=\"\/\/platform.twitter.com\/widgets.js\" charset=\"utf-8\"\u003E\u003C\/script\u003E","width":550,"height":null,"type":"rich","cache_age":"3153600000","provider_name":"Twitter","provider_url":"https:\/\/twitter.com","version":"1.0"}`,
			`(?s)^<blockquote class="twitter-tweet"><p lang="en" dir="ltr">Hugo 0.15 will have 30%. faster render times thanks to this commit <a href="https://t.co/FfzhM8bNhT">https://t.co/FfzhM8bNhT</a>  <a href="https://twitter.com/hashtag/gohugo.src=hash">#gohugo</a> <a href="https://twitter.com/hashtag/golang.src=hash">#golang</a> <a href="https://t.co/ITbMNU2BUf">https://t.co/ITbMNU2BUf</a></p>&mdash; Steve Francia .@spf13. <a href="https://twitter.com/spf13/status/666616452582129664">November 17, 2015</a></blockquote>.*?<script async src="//platform.twitter.com/widgets.js" charset="utf-8"></script>`,
		},
	} {
		// overload getJSON to return mock API response from Twitter
		tweetFuncMap := template.FuncMap{
			"getJSON": func(urlParts ...interface{}) interface{} {
				var v interface{}
				err := json.Unmarshal([]byte(this.resp), &v)
				if err != nil {
					t.Fatalf("[%d] unexpected error in json.Unmarshal: %s", i, err)
					return err
				}
				return v
			},
		}

		var (
			cfg, fs = newTestCfg()
			th      = newTestHelper(cfg, fs, t)
		)

		cfg.Set("privacy", this.privacy)

		writeSource(t, fs, filepath.Join("content", "simple.md"), fmt.Sprintf(`---
title: Shorty
---
%s`, this.in))
		writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), `{{ .Content }}`)

		buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg, OverloadedTemplateFuncs: tweetFuncMap}, BuildCfg{})

		th.assertFileContentRegexp(filepath.Join("public", "simple", "index.html"), this.expected)

	}
}

func TestShortcodeInstagram(t *testing.T) {
	t.Parallel()

	for i, this := range []struct {
		in, hidecaption, resp, expected string
	}{
		{
			`{{< instagram BMokmydjG-M >}}`,
			`0`,
			`{"provider_url": "https://www.instagram.com", "media_id": "1380514280986406796_25025320", "author_name": "instagram", "height": null, "thumbnail_url": "https://scontent-amt2-1.cdninstagram.com/t51.2885-15/s640x640/sh0.08/e35/15048135_1880160212214218_7827880881132929024_n.jpg?ig_cache_key=MTM4MDUxNDI4MDk4NjQwNjc5Ng%3D%3D.2", "thumbnail_width": 640, "thumbnail_height": 640, "provider_name": "Instagram", "title": "Today, we\u2019re introducing a few new tools to help you make your story even more fun: Boomerang and mentions. We\u2019re also starting to test links inside some stories.\nBoomerang lets you turn everyday moments into something fun and unexpected. Now you can easily take a Boomerang right inside Instagram. Swipe right from your feed to open the stories camera. A new format picker under the record button lets you select \u201cBoomerang\u201d mode.\nYou can also now share who you\u2019re with or who you\u2019re thinking of by mentioning them in your story. When you add text to your story, type \u201c@\u201d followed by a username and select the person you\u2019d like to mention. Their username will appear underlined in your story. And when someone taps the mention, they'll see a pop-up that takes them to that profile.\nYou may begin to spot \u201cSee More\u201d links at the bottom of some stories. This is a test that lets verified accounts add links so it\u2019s easy to learn more. From your favorite chefs\u2019 recipes to articles from top journalists or concert dates from the musicians you love, tap \u201cSee More\u201d or swipe up to view the link right inside the app.\nTo learn more about today\u2019s updates, check out help.instagram.com.\nThese updates for Instagram Stories are available as part of Instagram version 9.7 available for iOS in the Apple App Store, for Android in Google Play and for Windows 10 in the Windows Store.", "html": "\u003cblockquote class=\"instagram-media\" data-instgrm-captioned data-instgrm-version=\"7\" style=\" background:#FFF; border:0; border-radius:3px; box-shadow:0 0 1px 0 rgba(0,0,0,0.5),0 1px 10px 0 rgba(0,0,0,0.15); margin: 1px; max-width:658px; padding:0; width:99.375%; width:-webkit-calc(100% - 2px); width:calc(100% - 2px);\"\u003e\u003cdiv style=\"padding:8px;\"\u003e \u003cdiv style=\" background:#F8F8F8; line-height:0; margin-top:40px; padding:50.0% 0; text-align:center; width:100%;\"\u003e \u003cdiv style=\" background:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACwAAAAsCAMAAAApWqozAAAABGdBTUEAALGPC/xhBQAAAAFzUkdCAK7OHOkAAAAMUExURczMzPf399fX1+bm5mzY9AMAAADiSURBVDjLvZXbEsMgCES5/P8/t9FuRVCRmU73JWlzosgSIIZURCjo/ad+EQJJB4Hv8BFt+IDpQoCx1wjOSBFhh2XssxEIYn3ulI/6MNReE07UIWJEv8UEOWDS88LY97kqyTliJKKtuYBbruAyVh5wOHiXmpi5we58Ek028czwyuQdLKPG1Bkb4NnM+VeAnfHqn1k4+GPT6uGQcvu2h2OVuIf/gWUFyy8OWEpdyZSa3aVCqpVoVvzZZ2VTnn2wU8qzVjDDetO90GSy9mVLqtgYSy231MxrY6I2gGqjrTY0L8fxCxfCBbhWrsYYAAAAAElFTkSuQmCC); display:block; height:44px; margin:0 auto -44px; position:relative; top:-22px; width:44px;\"\u003e\u003c/div\u003e\u003c/div\u003e \u003cp style=\" margin:8px 0 0 0; padding:0 4px;\"\u003e \u003ca href=\"https://www.instagram.com/p/BMokmydjG-M/\" style=\" color:#000; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px; text-decoration:none; word-wrap:break-word;\" target=\"_blank\"\u003eToday, we\u2019re introducing a few new tools to help you make your story even more fun: Boomerang and mentions. We\u2019re also starting to test links inside some stories. Boomerang lets you turn everyday moments into something fun and unexpected. Now you can easily take a Boomerang right inside Instagram. Swipe right from your feed to open the stories camera. A new format picker under the record button lets you select \u201cBoomerang\u201d mode. You can also now share who you\u2019re with or who you\u2019re thinking of by mentioning them in your story. When you add text to your story, type \u201c@\u201d followed by a username and select the person you\u2019d like to mention. Their username will appear underlined in your story. And when someone taps the mention, they\u0026#39;ll see a pop-up that takes them to that profile. You may begin to spot \u201cSee More\u201d links at the bottom of some stories. This is a test that lets verified accounts add links so it\u2019s easy to learn more. From your favorite chefs\u2019 recipes to articles from top journalists or concert dates from the musicians you love, tap \u201cSee More\u201d or swipe up to view the link right inside the app. To learn more about today\u2019s updates, check out help.instagram.com. These updates for Instagram Stories are available as part of Instagram version 9.7 available for iOS in the Apple App Store, for Android in Google Play and for Windows 10 in the Windows Store.\u003c/a\u003e\u003c/p\u003e \u003cp style=\" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; line-height:17px; margin-bottom:0; margin-top:8px; overflow:hidden; padding:8px 0 7px; text-align:center; text-overflow:ellipsis; white-space:nowrap;\"\u003eA photo posted by Instagram (@instagram) on \u003ctime style=\" font-family:Arial,sans-serif; font-size:14px; line-height:17px;\" datetime=\"2016-11-10T15:02:28+00:00\"\u003eNov 10, 2016 at 7:02am PST\u003c/time\u003e\u003c/p\u003e\u003c/div\u003e\u003c/blockquote\u003e\n\u003cscript async defer src=\"//platform.instagram.com/en_US/embeds.js\"\u003e\u003c/script\u003e", "width": 658, "version": "1.0", "author_url": "https://www.instagram.com/instagram", "author_id": 25025320, "type": "rich"}`,
			`(?s)<blockquote class="instagram-media" data-instgrm-captioned data-instgrm-version="7" .*defer src="//platform.instagram.com/en_US/embeds.js"></script>`,
		},
		{
			`{{< instagram BMokmydjG-M hidecaption >}}`,
			`1`,
			`{"provider_url": "https://www.instagram.com", "media_id": "1380514280986406796_25025320", "author_name": "instagram", "height": null, "thumbnail_url": "https://scontent-amt2-1.cdninstagram.com/t51.2885-15/s640x640/sh0.08/e35/15048135_1880160212214218_7827880881132929024_n.jpg?ig_cache_key=MTM4MDUxNDI4MDk4NjQwNjc5Ng%3D%3D.2", "thumbnail_width": 640, "thumbnail_height": 640, "provider_name": "Instagram", "title": "Today, we\u2019re introducing a few new tools to help you make your story even more fun: Boomerang and mentions. We\u2019re also starting to test links inside some stories.\nBoomerang lets you turn everyday moments into something fun and unexpected. Now you can easily take a Boomerang right inside Instagram. Swipe right from your feed to open the stories camera. A new format picker under the record button lets you select \u201cBoomerang\u201d mode.\nYou can also now share who you\u2019re with or who you\u2019re thinking of by mentioning them in your story. When you add text to your story, type \u201c@\u201d followed by a username and select the person you\u2019d like to mention. Their username will appear underlined in your story. And when someone taps the mention, they'll see a pop-up that takes them to that profile.\nYou may begin to spot \u201cSee More\u201d links at the bottom of some stories. This is a test that lets verified accounts add links so it\u2019s easy to learn more. From your favorite chefs\u2019 recipes to articles from top journalists or concert dates from the musicians you love, tap \u201cSee More\u201d or swipe up to view the link right inside the app.\nTo learn more about today\u2019s updates, check out help.instagram.com.\nThese updates for Instagram Stories are available as part of Instagram version 9.7 available for iOS in the Apple App Store, for Android in Google Play and for Windows 10 in the Windows Store.", "html": "\u003cblockquote class=\"instagram-media\" data-instgrm-version=\"7\" style=\" background:#FFF; border:0; border-radius:3px; box-shadow:0 0 1px 0 rgba(0,0,0,0.5),0 1px 10px 0 rgba(0,0,0,0.15); margin: 1px; max-width:658px; padding:0; width:99.375%; width:-webkit-calc(100% - 2px); width:calc(100% - 2px);\"\u003e\u003cdiv style=\"padding:8px;\"\u003e \u003cdiv style=\" background:#F8F8F8; line-height:0; margin-top:40px; padding:50.0% 0; text-align:center; width:100%;\"\u003e \u003cdiv style=\" background:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACwAAAAsCAMAAAApWqozAAAABGdBTUEAALGPC/xhBQAAAAFzUkdCAK7OHOkAAAAMUExURczMzPf399fX1+bm5mzY9AMAAADiSURBVDjLvZXbEsMgCES5/P8/t9FuRVCRmU73JWlzosgSIIZURCjo/ad+EQJJB4Hv8BFt+IDpQoCx1wjOSBFhh2XssxEIYn3ulI/6MNReE07UIWJEv8UEOWDS88LY97kqyTliJKKtuYBbruAyVh5wOHiXmpi5we58Ek028czwyuQdLKPG1Bkb4NnM+VeAnfHqn1k4+GPT6uGQcvu2h2OVuIf/gWUFyy8OWEpdyZSa3aVCqpVoVvzZZ2VTnn2wU8qzVjDDetO90GSy9mVLqtgYSy231MxrY6I2gGqjrTY0L8fxCxfCBbhWrsYYAAAAAElFTkSuQmCC); display:block; height:44px; margin:0 auto -44px; position:relative; top:-22px; width:44px;\"\u003e\u003c/div\u003e\u003c/div\u003e\u003cp style=\" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; line-height:17px; margin-bottom:0; margin-top:8px; overflow:hidden; padding:8px 0 7px; text-align:center; text-overflow:ellipsis; white-space:nowrap;\"\u003e\u003ca href=\"https://www.instagram.com/p/BMokmydjG-M/\" style=\" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px; text-decoration:none;\" target=\"_blank\"\u003eA photo posted by Instagram (@instagram)\u003c/a\u003e on \u003ctime style=\" font-family:Arial,sans-serif; font-size:14px; line-height:17px;\" datetime=\"2016-11-10T15:02:28+00:00\"\u003eNov 10, 2016 at 7:02am PST\u003c/time\u003e\u003c/p\u003e\u003c/div\u003e\u003c/blockquote\u003e\n\u003cscript async defer src=\"//platform.instagram.com/en_US/embeds.js\"\u003e\u003c/script\u003e", "width": 658, "version": "1.0", "author_url": "https://www.instagram.com/instagram", "author_id": 25025320, "type": "rich"}`,
			`(?s)<blockquote class="instagram-media" data-instgrm-version="7" style=" background:#FFF; border:0; .*<script async defer src="//platform.instagram.com/en_US/embeds.js"></script>`,
		},
	} {
		// overload getJSON to return mock API response from Instagram
		instagramFuncMap := template.FuncMap{
			"getJSON": func(urlParts ...string) interface{} {
				var v interface{}
				err := json.Unmarshal([]byte(this.resp), &v)
				if err != nil {
					t.Fatalf("[%d] unexpected error in json.Unmarshal: %s", i, err)
					return err
				}
				return v
			},
		}

		var (
			cfg, fs = newTestCfg()
			th      = newTestHelper(cfg, fs, t)
		)

		writeSource(t, fs, filepath.Join("content", "simple.md"), fmt.Sprintf(`---
title: Shorty
---
%s`, this.in))
		writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), `{{ .Content | safeHTML }}`)

		buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg, OverloadedTemplateFuncs: instagramFuncMap}, BuildCfg{})

		th.assertFileContentRegexp(filepath.Join("public", "simple", "index.html"), this.expected)

	}
}
