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
	"os"
	"path/filepath"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/spf13/afero"
)

// DefaultConfig returns the default configuration.
func DefaultConfig() *allconfig.Config {
	fs := afero.NewMemMapFs()
	all, err := allconfig.LoadConfig(allconfig.ConfigSourceDescriptor{Fs: fs, Environ: []string{"none"}})
	if err != nil {
		panic(err)
	}
	return all.Base
}

// ExampleConfig returns the some example configuration for documentation.
func ExampleConfig() (*allconfig.Config, error) {
	// Apply some example settings for the settings that does not come with a sensible default.
	configToml := `
title = 'My Blog'
baseURL = "https://example.com/"
disableKinds = ["term", "taxonomy"]

[outputs]
home = ['html', 'html', 'rss']
page = ['html']

[imaging]
bgcolor = '#ffffff'
hint = 'photo'
quality = 81
resamplefilter = 'CatmullRom'
[imaging.exif]
disableDate = true
disableLatLong = true
excludeFields = 'ColorSpace|Metering'

[params]
color = 'blue'
style = 'dark'


[languages]
[languages.ar]
languagedirection = 'rtl'
title = 'مدونتي'
weight = 2
[languages.en]
weight = 1
[languages.fr]
weight = 2
[languages.fr.params]
linkedin = 'https://linkedin.com/fr/whoever'
color = 'green'
[[languages.fr.menus.main]]
name = 'Des produits'
pageRef = '/products'
weight = 20

[menus]
[[menus.main]]
name = 'Home'
pageRef = '/'
weight = 10
[[menus.main]]
name = 'Products'
pageRef = '/products'
weight = 20
[[menus.main]]
name = 'Services'
pageRef = '/services'
weight = 30

[deployment]
order = [".jpg$", ".gif$"]
[[deployment.targets]]
name = "mydeployment"
url = "s3://mybucket?region=us-east-1"
cloudFrontDistributionID = "mydistributionid"
[[deployment.matchers]]
pattern = "^.+\\.(js|css|svg|ttf)$"
cacheControl = "max-age=31536000, no-transform, public"
gzip = true
[[deployment.matchers]]
pattern = "^.+\\.(png|jpg)$"
cacheControl = "max-age=31536000, no-transform, public"
gzip = false
[[deployment.matchers]]
pattern = "^sitemap\\.xml$"
contentType = "application/xml"
gzip = true
[[deployment.matchers]]
pattern = "^.+\\.(html|xml|json)$"
gzip = true

[permalinks]
posts = '/posts/:year/:month/:title/'

[taxonomies]
category = 'categories'
series = 'series'
tag = 'tags'
  
[module]
[module.hugoVersion]
min = '0.80.0'
[[module.imports]]
path = "github.com/bep/hugo-mod-misc/dummy-content"
ignoreconfig = true
ignoreimports = true
[[module.mounts]]
source = "content/blog"
target = "content"

[minify]
[minify.tdewolff]
[minify.tdewolff.json]
precision = 2

[[cascade]]
background = 'yosemite.jpg'
[cascade._target]
  kind = 'page'
  path = '/blog/**'
[[cascade]]
background = 'goldenbridge.jpg'
[cascade._target]
  kind = 'section'


`

	goMod := `
module github.com/bep/mymod
`

	cfg := config.New()

	tempDir := os.TempDir()
	cacheDir := filepath.Join(tempDir, "hugocache")
	if err := os.MkdirAll(cacheDir, 0o777); err != nil {
		return nil, err
	}
	cfg.Set("cacheDir", cacheDir)
	cfg.Set("workingDir", tempDir)
	defer func() {
		os.RemoveAll(tempDir)
	}()

	fs := afero.NewOsFs()

	if err := afero.WriteFile(fs, filepath.Join(tempDir, "hugo.toml"), []byte(configToml), 0o644); err != nil {
		return nil, err
	}

	if err := afero.WriteFile(fs, filepath.Join(tempDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		return nil, err
	}

	conf, err := allconfig.LoadConfig(allconfig.ConfigSourceDescriptor{Fs: fs, Flags: cfg, Environ: []string{"none"}})
	if err != nil {
		return nil, err
	}
	return conf.Base, err
}
