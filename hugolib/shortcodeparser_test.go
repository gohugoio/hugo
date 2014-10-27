package hugolib

import (
	"testing"
)

type shortCodeLexerTest struct {
	name  string
	input string
	items []item
}

var (
	tstEOF       = item{tEOF, 0, ""}
	tstLeftNoMD  = item{tLeftDelimScNoMarkup, 0, "{{<"}
	tstRightNoMD = item{tRightDelimScNoMarkup, 0, ">}}"}
	tstLeftMD    = item{tLeftDelimScWithMarkup, 0, "{{%"}
	tstRightMD   = item{tRightDelimScWithMarkup, 0, "%}}"}
	tstSCClose   = item{tScClose, 0, "/"}
	tstSC1       = item{tScName, 0, "simple"}
	tstSC2       = item{tScName, 0, "shortcode2"}
	tstParam1    = item{tScParam, 0, "param1"}
	tstParam2    = item{tScParam, 0, "param2"}
	tstVal       = item{tScParamVal, 0, "Hello World"}
)

var shortCodeLexerTests = []shortCodeLexerTest{
	{"empty", "", []item{tstEOF}},
	{"spaces", " \t\n", []item{{tText, 0, " \t\n"}, tstEOF}},
	{"text", `to be or not`, []item{{tText, 0, "to be or not"}, tstEOF}},
	{"no markup", `{{< simple >}}`, []item{tstLeftNoMD, tstSC1, tstRightNoMD, tstEOF}},
	{"with EOL", "{{< simple \n >}}", []item{tstLeftNoMD, tstSC1, tstRightNoMD, tstEOF}},

	{"simple with markup", `{{% simple %}}`, []item{tstLeftMD, tstSC1, tstRightMD, tstEOF}},
	{"with spaces", `{{<     simple     >}}`, []item{tstLeftNoMD, tstSC1, tstRightNoMD, tstEOF}},
	{"mismatched rightDelim", `{{< simple %}}`, []item{tstLeftNoMD, tstSC1,
		{tError, 0, "unrecognized character in shortcode action: U+0025 '%'. Note: Parameters with non-alphanumeric args must be quoted"}}},
	{"inner, markup", `{{% simple %}} inner {{% /simple %}}`, []item{
		tstLeftMD,
		tstSC1,
		tstRightMD,
		{tText, 0, " inner "},
		tstLeftMD,
		tstSCClose,
		tstSC1,
		tstRightMD,
		tstEOF,
	}},
	{"close, but no open", `{{< /simple >}}`, []item{
		tstLeftNoMD, {tError, 0, "got closing shortcode, but none is open"}}},
	{"close wrong", `{{< simple >}}{{< /another >}}`, []item{
		tstLeftNoMD, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose,
		{tError, 0, "closing tag for shortcode 'another' does not match start tag"}}},
	{"close, but no open, more", `{{< simple >}}{{< /simple >}}{{< /another >}}`, []item{
		tstLeftNoMD, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose,
		{tError, 0, "closing tag for shortcode 'another' does not match start tag"}}},
	{"close with extra keyword", `{{< simple >}}{{< /simple keyword>}}`, []item{
		tstLeftNoMD, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose, tstSC1,
		{tError, 0, "unclosed shortcode"}}},
	{"Youtube id", `{{< simple -ziL-Q_456igdO-4 >}}`, []item{
		tstLeftNoMD, tstSC1, item{tScParam, 0, "-ziL-Q_456igdO-4"}, tstRightNoMD, tstEOF}},
	{"non-alphanumerics param quoted", `{{< simple "-ziL-.%QigdO-4" >}}`, []item{
		tstLeftNoMD, tstSC1, item{tScParam, 0, "-ziL-.%QigdO-4"}, tstRightNoMD, tstEOF}},
	{"two params", `{{< simple param1   param2 >}}`, []item{
		tstLeftNoMD, tstSC1, tstParam1, tstParam2, tstRightNoMD, tstEOF}},
	{"two quoted params", `{{< simple "param nr. 1" "param nr. 2" >}}`, []item{
		tstLeftNoMD, tstSC1, item{tScParam, 0, "param nr. 1"}, item{tScParam, 0, "param nr. 2"}, tstRightNoMD, tstEOF}},
	{"two named params", `{{< simple param1="Hello World" param2="p2Val">}}`, []item{
		tstLeftNoMD, tstSC1, tstParam1, tstVal, tstParam2, {tScParamVal, 0, "p2Val"}, tstRightNoMD, tstEOF}},
	{"escaped quotes", `{{< simple param1=\"Hello World\"  >}}`, []item{
		tstLeftNoMD, tstSC1, tstParam1, tstVal, tstRightNoMD, tstEOF}},
	{"escaped quotes inside escaped quotes", `{{< simple param1=\"Hello \"escaped\" World\"  >}}`, []item{
		tstLeftNoMD, tstSC1, tstParam1,
		item{tScParamVal, 0, `Hello `}, {tError, 0, `got positional parameter 'escaped'. Cannot mix named and positional parameters`}}},
	{"escaped quotes inside nonescaped quotes",
		`{{< simple param1="Hello \"escaped\" World"  >}}`, []item{
			tstLeftNoMD, tstSC1, tstParam1, item{tScParamVal, 0, `Hello "escaped" World`}, tstRightNoMD, tstEOF}},
	{"escaped quotes inside nonescaped quotes in positional param",
		`{{< simple "Hello \"escaped\" World"  >}}`, []item{
			tstLeftNoMD, tstSC1, item{tScParam, 0, `Hello "escaped" World`}, tstRightNoMD, tstEOF}},
	{"unterminated quote", `{{< simple param2="Hello World>}}`, []item{
		tstLeftNoMD, tstSC1, tstParam2, {tError, 0, "unterminated quoted string in shortcode parameter-argument: 'Hello World>}}'"}}},
	{"one named param, one not", `{{< simple param1="Hello World" p2 >}}`, []item{
		tstLeftNoMD, tstSC1, tstParam1, tstVal,
		{tError, 0, "got positional parameter 'p2'. Cannot mix named and positional parameters"}}},
	{"ono positional param, one not", `{{< simple param1 param2="Hello World">}}`, []item{
		tstLeftNoMD, tstSC1, tstParam1,
		{tError, 0, "got named parameter 'param2'. Cannot mix named and positional parameters"}}},
}

func TestPagelexer(t *testing.T) {
	for _, test := range shortCodeLexerTests {

		items := collect(&test)
		if !equal(items, test.items) {
			t.Errorf("%s: got\n\t%v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}

func collect(t *shortCodeLexerTest) (items []item) {
	l := newShortcodeLexer(t.name, t.input, 0)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == tEOF || item.typ == tError {
			break
		}
	}
	return
}

// no positional checking, for now ...
func equal(i1, i2 []item) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
	}
	return true
}
