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
	"github.com/gohugoio/hugo/resources/page"
)

// Sections returns the top level sections.
func (s *Site) Sections() page.Pages {
	s.CheckReady()
	return s.Home().Sections()
}

// Home is a shortcut to the home page, equivalent to .Site.GetPage "home".
func (s *Site) Home() page.Page {
	s.CheckReady()
	return s.s.home
}
