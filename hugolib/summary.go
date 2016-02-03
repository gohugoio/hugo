// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"regexp"

	"github.com/spf13/hugo/helpers"
)

// for enums define a type of the enum that one want
type Summerization int

// next we define the constants of type Summerization
// we use iota to provide int values to each of the constant types
const (
	PLAIN Summerization = 1 + iota
	HTML_FIRSTPARAGRAPH 
)

// Enums should be able to printout as strings
// so we declare the next best thing, a slice of strings
// for eg. the string value will be used in the println
var summerizations = [...]string {
 "PLAIN",
 "HTML_FIRSTPARAGRAPH",
}

// String() function will return the name
// that we want out constant Summerization be recognized as
func (summerization Summerization) String() string {
 return summerizations[summerization - 1]
}

func plainSummarizationStrategy(p *Page) (string, bool) {
	plain := p.PlainWords()
	return helpers.TruncateWordsToWholeSentence(plain, helpers.SummaryLength)
}

func htmlFirstParagraphSummerizationStrategy(p *Page) (string, bool) {
	content := string(p.renderBytes(p.rawContent))
	return regexp.MustCompile("<h[123456]").Split(content, 2)[0], len(content) != len(p.Summary)
}

func summaryStrategySwitch(p *Page) (string, bool) {
	if p.Site.SummaryStrategy == "html_firstparagraph" {
		return htmlFirstParagraphSummerizationStrategy(p)
	} else {
		return plainSummarizationStrategy(p)
	}
}

func Summarize(p *Page) {
	summary, truncated := summaryStrategySwitch(p)
	p.Summary = helpers.BytesToHTML([]byte(summary))
	p.Truncated = truncated
} 