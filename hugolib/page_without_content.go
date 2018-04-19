// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"html/template"
)

// PageWithoutContent is sent to the shortcodes. They cannot access the content
// they're a part of. It would cause an infinite regress.
//
// Go doesn't support virtual methods, so this careful dance is currently (I think)
// the best we can do.
type PageWithoutContent struct {
	*Page
}

// Content returns an empty string.
func (p *PageWithoutContent) Content() (interface{}, error) {
	return "", nil
}

// Truncated always returns false.
func (p *PageWithoutContent) Truncated() bool {
	return false
}

// Summary returns an empty string.
func (p *PageWithoutContent) Summary() template.HTML {
	return ""
}

// WordCount always returns 0.
func (p *PageWithoutContent) WordCount() int {
	return 0
}

// ReadingTime always returns 0.
func (p *PageWithoutContent) ReadingTime() int {
	return 0
}

// FuzzyWordCount always returns 0.
func (p *PageWithoutContent) FuzzyWordCount() int {
	return 0
}

// Plain returns an empty string.
func (p *PageWithoutContent) Plain() string {
	return ""
}

// PlainWords returns an empty string slice.
func (p *PageWithoutContent) PlainWords() []string {
	return []string{}
}
