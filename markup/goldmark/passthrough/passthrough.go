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

// GOLDMARK-V2: The passthrough feature depends on the external module
// github.com/gohugoio/hugo-goldmark-extensions/passthrough which is built against
// goldmark v1 and has no v2 release yet. The whole implementation is therefore
// disabled for the v2 test upgrade; New returns nil so the feature is a no-op.
// See https://github.com/yuin/goldmark/discussions/559.
package passthrough

import (
	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"
	"github.com/yuin/goldmark/v2/parser"
)

// New returns nil: passthrough is disabled until the upstream extension supports
// goldmark v2.
func New(cfg goldmark_config.Passthrough) parser.Extension {
	return nil
}
