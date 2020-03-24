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

package pagemeta

import (
	"github.com/mitchellh/mapstructure"
)

type URLPath struct {
	URL       string
	Permalink string
	Slug      string
	Section   string
}

const (
	Never       = "never"
	Always      = "always"
	ListLocally = "local"
)

var defaultBuildConfig = BuildConfig{
	List:             Always,
	Render:           true,
	PublishResources: true,
	set:              true,
}

// BuildConfig holds configuration options about how to handle a Page in Hugo's
// build process.
type BuildConfig struct {
	// Whether to add it to any of the page collections.
	// Note that the page can always be found with .Site.GetPage.
	// Valid values: never, always, local.
	// Setting it to 'local' means they will be available via the local
	// page collections, e.g. $section.Pages.
	// Note: before 0.57.2 this was a bool, so we accept those too.
	List string

	// Whether to render it.
	Render bool

	// Whether to publish its resources. These will still be published on demand,
	// but enabling this can be useful if the originals (e.g. images) are
	// never used.
	PublishResources bool

	set bool // BuildCfg is non-zero if this is set to true.
}

// Disable sets all options to their off value.
func (b *BuildConfig) Disable() {
	b.List = Never
	b.Render = false
	b.PublishResources = false
	b.set = true
}

func (b BuildConfig) IsZero() bool {
	return !b.set
}

func DecodeBuildConfig(m interface{}) (BuildConfig, error) {
	b := defaultBuildConfig
	if m == nil {
		return b, nil
	}

	err := mapstructure.WeakDecode(m, &b)

	// In 0.67.1 we changed the list attribute from a bool to a string (enum).
	// Bool values will become 0 or 1.
	switch b.List {
	case "0":
		b.List = Never
	case "1":
		b.List = Always
	case Always, Never, ListLocally:
	default:
		b.List = Always
	}

	return b, err
}
