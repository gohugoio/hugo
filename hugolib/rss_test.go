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

	"github.com/spf13/viper"
)

func TestRSSOutput(t *testing.T) {
	testCommonResetState()

	rssURI := "customrss.xml"
	viper.Set("baseURL", "http://auth/bub/")
	viper.Set("rssURI", rssURI)
	viper.Set("title", "RSSTest")

	for _, s := range weightedSources {
		writeSource(t, filepath.Join("content", "sect", s.Name), string(s.Content))
	}

	if err := buildAndRenderSite(NewSiteDefaultLang()); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	// Home RSS
	assertFileContent(t, filepath.Join("public", rssURI), true, "<?xml", "rss version", "RSSTest")
	// Section RSS
	assertFileContent(t, filepath.Join("public", "sect", rssURI), true, "<?xml", "rss version", "Sects on RSSTest")
	// Taxonomy RSS
	assertFileContent(t, filepath.Join("public", "categories", "hugo", rssURI), true, "<?xml", "rss version", "Hugo on RSSTest")

}
