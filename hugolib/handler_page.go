// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import "github.com/spf13/hugo/source"

func init() {
	RegisterHandler(markdownHandler)
	RegisterHandler(htmlHandler)
}

var markdownHandler = Handle{
	extensions: []string{"mdown", "markdown", "md"},
	read: func(f *source.File, s *Site, results HandleResults) {
		page, err := NewPage(f.Path())
		if err != nil {
			results <- HandledResult{file: f, err: err}
		}

		if err := page.ReadFrom(f.Contents); err != nil {
			results <- HandledResult{file: f, err: err}
		}

		page.Site = &s.Info
		page.Tmpl = s.Tmpl

		results <- HandledResult{file: f, page: page, err: err}
	},
	pageConvert: func(p *Page, s *Site, results HandleResults) {
		p.ProcessShortcodes(s.Tmpl)
		err := p.Convert()
		if err != nil {
			results <- HandledResult{err: err}
		}

		results <- HandledResult{err: err}
	},
}

var htmlHandler = Handle{
	extensions: []string{"html", "htm"},
	read: func(f *source.File, s *Site, results HandleResults) {
		page, err := NewPage(f.Path())
		if err != nil {
			results <- HandledResult{file: f, err: err}
		}

		if err := page.ReadFrom(f.Contents); err != nil {
			results <- HandledResult{file: f, err: err}
		}

		page.Site = &s.Info
		page.Tmpl = s.Tmpl

		results <- HandledResult{file: f, page: page, err: err}
	},
	pageConvert: func(p *Page, s *Site, results HandleResults) {
		p.ProcessShortcodes(s.Tmpl)
		err := p.Convert()
		if err != nil {
			results <- HandledResult{err: err}
		}

		results <- HandledResult{err: err}
	},
}
