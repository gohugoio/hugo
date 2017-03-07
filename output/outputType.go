// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package output

import (
	"github.com/spf13/hugo/media"
)

var (
	HTMLType = Type{
		Name:      "HTML",
		MediaType: media.HTMLType,
		BaseName:  "index",
	}

	RSSType = Type{
		Name:      "RSS",
		MediaType: media.RSSType,
		BaseName:  "index",
	}
)

type Types []Type

// Type represents an output represenation, usually to a file on disk.
type Type struct {
	// The Name is used as an identifier. Internal output types (i.e. HTML and RSS)
	// can be overridden by providing a new definition for those types.
	Name string

	MediaType media.Type

	// Must be set to a value when there are two or more conflicting mediatype for the same resource.
	Path string

	// The base output file name used when not using "ugly URLs", defaults to "index".
	BaseName string

	// The protocol to use, i.e. "webcal://". Defaults to the protocol of the baseURL.
	Protocol string

	// IsPlainText decides whether to use text/template or html/template
	// as template parser.
	IsPlainText bool

	// Enable to ignore the global uglyURLs setting.
	NoUgly bool
}
