// Copyright 2024 The Hugo Authors. All rights reserved.
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

package httpcache_test

import (
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

func TestConfigCustom(t *testing.T) {
	files := `
-- hugo.toml --
[httpcache]
[httpcache.cache.for]
includes = ["**gohugo.io**"]
[[httpcache.polls]]
low = "5s"
high = "32s"
[httpcache.polls.for]
includes = ["**gohugo.io**"]
		
	
`

	b := hugolib.Test(t, files)

	httpcacheConf := b.H.Configs.Base.HTTPCache
	compiled := b.H.Configs.Base.C.HTTPCache

	b.Assert(httpcacheConf.Cache.For.Includes, qt.DeepEquals, []string{"**gohugo.io**"})
	b.Assert(httpcacheConf.Cache.For.Excludes, qt.IsNil)

	pc := compiled.PollConfigFor("https://gohugo.io/foo.jpg")
	b.Assert(pc.Config.Low, qt.Equals, 5*time.Second)
	b.Assert(pc.Config.High, qt.Equals, 32*time.Second)
	b.Assert(compiled.PollConfigFor("https://example.com/foo.jpg").IsZero(), qt.IsTrue)
}

func TestConfigDefault(t *testing.T) {
	files := `
-- hugo.toml --
`
	b := hugolib.Test(t, files)

	compiled := b.H.Configs.Base.C.HTTPCache

	b.Assert(compiled.For("https://gohugo.io/posts.json"), qt.IsFalse)
	b.Assert(compiled.For("https://gohugo.io/foo.jpg"), qt.IsFalse)
	b.Assert(compiled.PollConfigFor("https://gohugo.io/foo.jpg").Config.Disable, qt.IsTrue)
}
