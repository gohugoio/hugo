package hugolib

import (
	"strings"
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

func plainSummarizationStrategy(p *Page) ([]byte, bool) {
	plain := strings.TrimSpace(p.Plain())
	return []byte(helpers.TruncateWordsToWholeSentence(plain, helpers.SummaryLength)), len(p.Summary) != len(plain)
}

func htmlFirstParagraphSummerizationStrategy(p *Page) ([]byte, bool) {
	content := string(p.renderBytes(p.rawContent))
	return []byte(regexp.MustCompile("<h[123456]").Split(content, 2)[0]), len(content) != len(p.Summary)
}

func summaryStrategySwitch(p *Page) ([]byte, bool) {
	if p.Site.SummaryStrategy == "html_firstparagraph" {
		return htmlFirstParagraphSummerizationStrategy(p)
	} else {
		return plainSummarizationStrategy(p)
	}
}

func Summarize(p *Page) {
	summary, truncated := summaryStrategySwitch(p)
	p.Summary = helpers.BytesToHTML(summary)
	p.Truncated = truncated
} 