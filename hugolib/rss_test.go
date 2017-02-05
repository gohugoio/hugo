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
	"path/filepath"
	"testing"

	"github.com/spf13/hugo/deps"
)

func TestRSSOutput(t *testing.T) {
	t.Parallel()
	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg}
	)

	rssURI := "customrss.xml"

	cfg.Set("baseURL", "http://auth/bub/")
	cfg.Set("rssURI", rssURI)
	cfg.Set("title", "RSSTest")

	for _, src := range weightedSources {
		writeSource(t, fs, filepath.Join("content", "sect", src.Name), string(src.Content))
	}

	buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	// Home RSS
	th.assertFileContent(t, fs, filepath.Join("public", rssURI), true, "<?xml", "rss version", "RSSTest")
	// Section RSS
	th.assertFileContent(t, fs, filepath.Join("public", "sect", rssURI), true, "<?xml", "rss version", "Sects on RSSTest")
	// Taxonomy RSS
	th.assertFileContent(t, fs, filepath.Join("public", "categories", "hugo", rssURI), true, "<?xml", "rss version", "Hugo on RSSTest")

}
