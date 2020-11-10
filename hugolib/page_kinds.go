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
	"github.com/gohugoio/hugo/resources/page/pagekinds"
)

// This is all the kinds we can expect to find in .Site.Pages.
var allKindsInPages = []string{pagekinds.Page, pagekinds.Home, pagekinds.Section, pagekinds.Term, pagekinds.Taxonomy}

const (
	pageResourceType = "page"
)
