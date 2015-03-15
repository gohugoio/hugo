package transform

import (
	"bytes"
	bp "github.com/spf13/hugo/bufferpool"
	"net/url"
	"strings"
	"unicode/utf8"
)

type matchState int

const (
	matchStateNone matchState = iota
	matchStateWhitespace
	matchStatePartial
	matchStateFull
)

const (
	matchPrefixSrc int = iota
	matchPrefixHref
)

type contentlexer struct {
	content []byte

	pos   int // input position
	start int // item start position
	width int // width of last element

	matchers     []absurlMatcher
	state        stateFunc
	prefixLookup *prefixes

	b *bytes.Buffer
}

type stateFunc func(*contentlexer) stateFunc

type prefixRunes []rune

type prefixes struct {
	pr   []prefixRunes
	curr prefixRunes // current prefix lookup table
	i    int         // current index

	// first rune in potential match
	first rune

	// match-state:
	// none, whitespace, partial, full
	ms matchState
}

// match returns partial and full match for the prefix in play
// - it's a full match if all prefix runes has checked out in row
// - it's a partial match if it's on its way towards a full match
func (l *contentlexer) match(r rune) {
	p := l.prefixLookup
	if p.curr == nil {
		// assumes prefixes all start off on a different rune
		// works in this special case: href, src
		p.i = 0
		for _, pr := range p.pr {
			if pr[p.i] == r {
				fullMatch := len(p.pr) == 1
				p.first = r
				if !fullMatch {
					p.curr = pr
					l.prefixLookup.ms = matchStatePartial
				} else {
					l.prefixLookup.ms = matchStateFull
				}
				return
			}
		}
	} else {
		p.i++
		if p.curr[p.i] == r {
			fullMatch := len(p.curr) == p.i+1
			if fullMatch {
				p.curr = nil
				l.prefixLookup.ms = matchStateFull
			} else {
				l.prefixLookup.ms = matchStatePartial
			}
			return
		}

		p.curr = nil
	}

	l.prefixLookup.ms = matchStateNone
}

func (l *contentlexer) emit() {
	l.b.Write(l.content[l.start:l.pos])
	l.start = l.pos
}

var mainPrefixRunes = []prefixRunes{{'s', 'r', 'c', '='}, {'h', 'r', 'e', 'f', '='}}

type absurlMatcher struct {
	prefix      int
	match       []byte
	replacement []byte
}

func (a absurlMatcher) isSourceType() bool {
	return a.prefix == matchPrefixSrc
}

func checkCandidate(l *contentlexer) {
	isSource := l.prefixLookup.first == 's'
	for _, m := range l.matchers {

		if isSource && !m.isSourceType() || !isSource && m.isSourceType() {
			continue
		}

		if bytes.HasPrefix(l.content[l.pos:], m.match) {
			// check for schemaless urls
			posAfter := l.pos + len(m.match)
			if int(posAfter) >= len(l.content) {
				return
			}
			r, _ := utf8.DecodeRune(l.content[posAfter:])
			if r == '/' {
				// schemaless: skip
				return
			}
			if l.pos > l.start {
				l.emit()
			}
			l.pos += len(m.match)
			l.b.Write(m.replacement)
			l.start = l.pos
			return

		}
	}
}

func (l *contentlexer) replace() {
	contentLength := len(l.content)
	var r rune

	for {
		if int(l.pos) >= contentLength {
			l.width = 0
			break
		}

		var width int = 1
		r = rune(l.content[l.pos])
		if r >= utf8.RuneSelf {
			r, width = utf8.DecodeRune(l.content[l.pos:])
		}
		l.width = width
		l.pos += l.width

		if r == ' ' {
			l.prefixLookup.ms = matchStateWhitespace
		} else if l.prefixLookup.ms != matchStateNone {
			l.match(r)
			if l.prefixLookup.ms == matchStateFull {
				checkCandidate(l)
			}
		}

	}

	// Done!
	if l.pos > l.start {
		l.emit()
	}
}

func doReplace(content []byte, matchers []absurlMatcher) []byte {
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)

	lexer := &contentlexer{content: content,
		b:            b,
		prefixLookup: &prefixes{pr: mainPrefixRunes},
		matchers:     matchers}

	lexer.replace()

	return b.Bytes()
}

type absurlReplacer struct {
	htmlMatchers []absurlMatcher
	xmlMatchers  []absurlMatcher
}

func newAbsurlReplacer(baseUrl string) *absurlReplacer {
	u, _ := url.Parse(baseUrl)
	base := strings.TrimRight(u.String(), "/")

	// HTML
	dqHtmlMatch := []byte("\"/")
	sqHtmlMatch := []byte("'/")

	// XML
	dqXmlMatch := []byte("&#34;/")
	sqXmlMatch := []byte("&#39;/")

	dqHtml := []byte("\"" + base + "/")
	sqHtml := []byte("'" + base + "/")

	dqXml := []byte("&#34;" + base + "/")
	sqXml := []byte("&#39;" + base + "/")

	return &absurlReplacer{
		htmlMatchers: []absurlMatcher{
			{matchPrefixSrc, dqHtmlMatch, dqHtml},
			{matchPrefixSrc, sqHtmlMatch, sqHtml},
			{matchPrefixHref, dqHtmlMatch, dqHtml},
			{matchPrefixHref, sqHtmlMatch, sqHtml}},
		xmlMatchers: []absurlMatcher{
			{matchPrefixSrc, dqXmlMatch, dqXml},
			{matchPrefixSrc, sqXmlMatch, sqXml},
			{matchPrefixHref, dqXmlMatch, dqXml},
			{matchPrefixHref, sqXmlMatch, sqXml},
		}}

}

func (au *absurlReplacer) replaceInHtml(content []byte) []byte {
	return doReplace(content, au.htmlMatchers)
}

func (au *absurlReplacer) replaceInXml(content []byte) []byte {
	return doReplace(content, au.xmlMatchers)
}
