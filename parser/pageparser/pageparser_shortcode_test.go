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

package pageparser

import (
	"testing"
)

var (
	tstEOF            = nti(tEOF, "")
	tstLeftNoMD       = nti(tLeftDelimScNoMarkup, "{{<")
	tstRightNoMD      = nti(tRightDelimScNoMarkup, ">}}")
	tstLeftMD         = nti(tLeftDelimScWithMarkup, "{{%")
	tstRightMD        = nti(tRightDelimScWithMarkup, "%}}")
	tstSCClose        = nti(tScClose, "/")
	tstSC1            = nti(tScName, "sc1")
	tstSC1Inline      = nti(tScNameInline, "sc1.inline")
	tstSC2Inline      = nti(tScNameInline, "sc2.inline")
	tstSC2            = nti(tScName, "sc2")
	tstSC3            = nti(tScName, "sc3")
	tstSCSlash        = nti(tScName, "sc/sub")
	tstParam1         = nti(tScParam, "param1")
	tstParam2         = nti(tScParam, "param2")
	tstParamBoolTrue  = nti(tScParam, "true")
	tstParamBoolFalse = nti(tScParam, "false")
	tstParamInt       = nti(tScParam, "32")
	tstParamFloat     = nti(tScParam, "3.14")
	tstVal            = nti(tScParamVal, "Hello World")
	tstText           = nti(tText, "Hello World")
)

var shortCodeLexerTests = []lexerTest{
	{"empty", "", []Item{tstEOF}},
	{"spaces", " \t\n", []Item{nti(tText, " \t\n"), tstEOF}},
	{"text", `to be or not`, []Item{nti(tText, "to be or not"), tstEOF}},
	{"no markup", `{{< sc1 >}}`, []Item{tstLeftNoMD, tstSC1, tstRightNoMD, tstEOF}},
	{"with EOL", "{{< sc1 \n >}}", []Item{tstLeftNoMD, tstSC1, tstRightNoMD, tstEOF}},

	{"forward slash inside name", `{{< sc/sub >}}`, []Item{tstLeftNoMD, tstSCSlash, tstRightNoMD, tstEOF}},

	{"simple with markup", `{{% sc1 %}}`, []Item{tstLeftMD, tstSC1, tstRightMD, tstEOF}},
	{"with spaces", `{{<     sc1     >}}`, []Item{tstLeftNoMD, tstSC1, tstRightNoMD, tstEOF}},
	{"mismatched rightDelim", `{{< sc1 %}}`, []Item{tstLeftNoMD, tstSC1,
		nti(tError, "unrecognized character in shortcode action: U+0025 '%'. Note: Parameters with non-alphanumeric args must be quoted")}},
	{"inner, markup", `{{% sc1 %}} inner {{% /sc1 %}}`, []Item{
		tstLeftMD,
		tstSC1,
		tstRightMD,
		nti(tText, " inner "),
		tstLeftMD,
		tstSCClose,
		tstSC1,
		tstRightMD,
		tstEOF,
	}},
	{"close, but no open", `{{< /sc1 >}}`, []Item{
		tstLeftNoMD, nti(tError, "got closing shortcode, but none is open")}},
	{"close wrong", `{{< sc1 >}}{{< /another >}}`, []Item{
		tstLeftNoMD, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose,
		nti(tError, "closing tag for shortcode 'another' does not match start tag")}},
	{"close, but no open, more", `{{< sc1 >}}{{< /sc1 >}}{{< /another >}}`, []Item{
		tstLeftNoMD, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose,
		nti(tError, "closing tag for shortcode 'another' does not match start tag")}},
	{"close with extra keyword", `{{< sc1 >}}{{< /sc1 keyword>}}`, []Item{
		tstLeftNoMD, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose, tstSC1,
		nti(tError, "unclosed shortcode")}},
	{"float param, positional", `{{< sc1 3.14 >}}`, []Item{
		tstLeftNoMD, tstSC1, nti(tScParam, "3.14"), tstRightNoMD, tstEOF}},
	{"float param, named", `{{< sc1 param1=3.14 >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, nti(tScParamVal, "3.14"), tstRightNoMD, tstEOF}},
	{"named param, raw string", `{{< sc1 param1=` + "`" + "Hello World" + "`" + " >}}", []Item{
		tstLeftNoMD, tstSC1, tstParam1, nti(tScParamVal, "Hello World"), tstRightNoMD, tstEOF}},
	{"float param, named, space before", `{{< sc1 param1= 3.14 >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, nti(tScParamVal, "3.14"), tstRightNoMD, tstEOF}},
	{"Youtube id", `{{< sc1 -ziL-Q_456igdO-4 >}}`, []Item{
		tstLeftNoMD, tstSC1, nti(tScParam, "-ziL-Q_456igdO-4"), tstRightNoMD, tstEOF}},
	{"non-alphanumerics param quoted", `{{< sc1 "-ziL-.%QigdO-4" >}}`, []Item{
		tstLeftNoMD, tstSC1, nti(tScParam, "-ziL-.%QigdO-4"), tstRightNoMD, tstEOF}},
	{"raw string", `{{< sc1` + "`" + "Hello World" + "`" + ` >}}`, []Item{
		tstLeftNoMD, tstSC1, nti(tScParam, "Hello World"), tstRightNoMD, tstEOF}},
	{"raw string with newline", `{{< sc1` + "`" + `Hello 
	World` + "`" + ` >}}`, []Item{
		tstLeftNoMD, tstSC1, nti(tScParam, `Hello 
	World`), tstRightNoMD, tstEOF}},
	{"raw string with escape character", `{{< sc1` + "`" + `Hello \b World` + "`" + ` >}}`, []Item{
		tstLeftNoMD, tstSC1, nti(tScParam, `Hello \b World`), tstRightNoMD, tstEOF}},
	{"two params", `{{< sc1 param1   param2 >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstParam2, tstRightNoMD, tstEOF}},
	// issue #934
	{"self-closing", `{{< sc1 />}}`, []Item{
		tstLeftNoMD, tstSC1, tstSCClose, tstRightNoMD, tstEOF}},
	// Issue 2498
	{"multiple self-closing", `{{< sc1 />}}{{< sc1 />}}`, []Item{
		tstLeftNoMD, tstSC1, tstSCClose, tstRightNoMD,
		tstLeftNoMD, tstSC1, tstSCClose, tstRightNoMD, tstEOF}},
	{"self-closing with param", `{{< sc1 param1 />}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstSCClose, tstRightNoMD, tstEOF}},
	{"multiple self-closing with param", `{{< sc1 param1 />}}{{< sc1 param1 />}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstSCClose, tstRightNoMD,
		tstLeftNoMD, tstSC1, tstParam1, tstSCClose, tstRightNoMD, tstEOF}},
	{"multiple different self-closing with param", `{{< sc1 param1 />}}{{< sc2 param1 />}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstSCClose, tstRightNoMD,
		tstLeftNoMD, tstSC2, tstParam1, tstSCClose, tstRightNoMD, tstEOF}},
	{"nested simple", `{{< sc1 >}}{{< sc2 >}}{{< /sc1 >}}`, []Item{
		tstLeftNoMD, tstSC1, tstRightNoMD,
		tstLeftNoMD, tstSC2, tstRightNoMD,
		tstLeftNoMD, tstSCClose, tstSC1, tstRightNoMD, tstEOF}},
	{"nested complex", `{{< sc1 >}}ab{{% sc2 param1 %}}cd{{< sc3 >}}ef{{< /sc3 >}}gh{{% /sc2 %}}ij{{< /sc1 >}}kl`, []Item{
		tstLeftNoMD, tstSC1, tstRightNoMD,
		nti(tText, "ab"),
		tstLeftMD, tstSC2, tstParam1, tstRightMD,
		nti(tText, "cd"),
		tstLeftNoMD, tstSC3, tstRightNoMD,
		nti(tText, "ef"),
		tstLeftNoMD, tstSCClose, tstSC3, tstRightNoMD,
		nti(tText, "gh"),
		tstLeftMD, tstSCClose, tstSC2, tstRightMD,
		nti(tText, "ij"),
		tstLeftNoMD, tstSCClose, tstSC1, tstRightNoMD,
		nti(tText, "kl"), tstEOF,
	}},

	{"two quoted params", `{{< sc1 "param nr. 1" "param nr. 2" >}}`, []Item{
		tstLeftNoMD, tstSC1, nti(tScParam, "param nr. 1"), nti(tScParam, "param nr. 2"), tstRightNoMD, tstEOF}},
	{"two named params", `{{< sc1 param1="Hello World" param2="p2Val">}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstVal, tstParam2, nti(tScParamVal, "p2Val"), tstRightNoMD, tstEOF}},
	{"escaped quotes", `{{< sc1 param1=\"Hello World\"  >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstVal, tstRightNoMD, tstEOF}},
	{"escaped quotes, positional param", `{{< sc1 \"param1\"  >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstRightNoMD, tstEOF}},
	{"escaped quotes inside escaped quotes", `{{< sc1 param1=\"Hello \"escaped\" World\"  >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1,
		nti(tScParamVal, `Hello `), nti(tError, `got positional parameter 'escaped'. Cannot mix named and positional parameters`)}},
	{"escaped quotes inside nonescaped quotes",
		`{{< sc1 param1="Hello \"escaped\" World"  >}}`, []Item{
			tstLeftNoMD, tstSC1, tstParam1, nti(tScParamVal, `Hello "escaped" World`), tstRightNoMD, tstEOF}},
	{"escaped quotes inside nonescaped quotes in positional param",
		`{{< sc1 "Hello \"escaped\" World"  >}}`, []Item{
			tstLeftNoMD, tstSC1, nti(tScParam, `Hello "escaped" World`), tstRightNoMD, tstEOF}},
	{"escaped raw string, named param", `{{< sc1 param1=` + `\` + "`" + "Hello World" + `\` + "`" + ` >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, nti(tError, "unrecognized escape character")}},
	{"escaped raw string, positional param", `{{< sc1 param1 ` + `\` + "`" + "Hello World" + `\` + "`" + ` >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, nti(tError, "unrecognized escape character")}},
	{"two raw string params", `{{< sc1` + "`" + "Hello World" + "`" + "`" + "Second Param" + "`" + ` >}}`, []Item{
		tstLeftNoMD, tstSC1, nti(tScParam, "Hello World"), nti(tScParam, "Second Param"), tstRightNoMD, tstEOF}},
	{"unterminated quote", `{{< sc1 param2="Hello World>}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam2, nti(tError, "unterminated quoted string in shortcode parameter-argument: 'Hello World>}}'")}},
	{"unterminated raw string", `{{< sc1` + "`" + "Hello World" + ` >}}`, []Item{
		tstLeftNoMD, tstSC1, nti(tError, "unterminated raw string in shortcode parameter-argument: 'Hello World >}}'")}},
	{"unterminated raw string in second argument", `{{< sc1` + "`" + "Hello World" + "`" + "`" + "Second Param" + ` >}}`, []Item{
		tstLeftNoMD, tstSC1, nti(tScParam, "Hello World"), nti(tError, "unterminated raw string in shortcode parameter-argument: 'Second Param >}}'")}},
	{"one named param, one not", `{{< sc1 param1="Hello World" p2 >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstVal,
		nti(tError, "got positional parameter 'p2'. Cannot mix named and positional parameters")}},
	{"one named param, one quoted positional param, both raw strings", `{{< sc1 param1=` + "`" + "Hello World" + "`" + "`" + "Second Param" + "`" + ` >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstVal,
		nti(tError, "got quoted positional parameter. Cannot mix named and positional parameters")}},
	{"one named param, one quoted positional param", `{{< sc1 param1="Hello World" "And Universe" >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1, tstVal,
		nti(tError, "got quoted positional parameter. Cannot mix named and positional parameters")}},
	{"one quoted positional param, one named param", `{{< sc1 "param1" param2="And Universe" >}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1,
		nti(tError, "got named parameter 'param2'. Cannot mix named and positional parameters")}},
	{"ono positional param, one not", `{{< sc1 param1 param2="Hello World">}}`, []Item{
		tstLeftNoMD, tstSC1, tstParam1,
		nti(tError, "got named parameter 'param2'. Cannot mix named and positional parameters")}},
	{"commented out", `{{</* sc1 */>}}`, []Item{
		nti(tText, "{{<"), nti(tText, " sc1 "), nti(tText, ">}}"), tstEOF}},
	{"commented out, with asterisk inside", `{{</* sc1 "**/*.pdf" */>}}`, []Item{
		nti(tText, "{{<"), nti(tText, " sc1 \"**/*.pdf\" "), nti(tText, ">}}"), tstEOF}},
	{"commented out, missing close", `{{</* sc1 >}}`, []Item{
		nti(tError, "comment must be closed")}},
	{"commented out, misplaced close", `{{</* sc1 >}}*/`, []Item{
		nti(tError, "comment must be closed")}},
	// Inline shortcodes
	{"basic inline", `{{< sc1.inline >}}Hello World{{< /sc1.inline >}}`, []Item{tstLeftNoMD, tstSC1Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSCClose, tstSC1Inline, tstRightNoMD, tstEOF}},
	{"basic inline with space", `{{< sc1.inline >}}Hello World{{< / sc1.inline >}}`, []Item{tstLeftNoMD, tstSC1Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSCClose, tstSC1Inline, tstRightNoMD, tstEOF}},
	{"inline self closing", `{{< sc1.inline >}}Hello World{{< /sc1.inline >}}Hello World{{< sc1.inline />}}`, []Item{tstLeftNoMD, tstSC1Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSCClose, tstSC1Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSC1Inline, tstSCClose, tstRightNoMD, tstEOF}},
	{"inline self closing, then a new inline", `{{< sc1.inline >}}Hello World{{< /sc1.inline >}}Hello World{{< sc1.inline />}}{{< sc2.inline >}}Hello World{{< /sc2.inline >}}`, []Item{
		tstLeftNoMD, tstSC1Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSCClose, tstSC1Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSC1Inline, tstSCClose, tstRightNoMD,
		tstLeftNoMD, tstSC2Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSCClose, tstSC2Inline, tstRightNoMD, tstEOF}},
	{"inline with template syntax", `{{< sc1.inline >}}{{ .Get 0 }}{{ .Get 1 }}{{< /sc1.inline >}}`, []Item{tstLeftNoMD, tstSC1Inline, tstRightNoMD, nti(tText, "{{ .Get 0 }}"), nti(tText, "{{ .Get 1 }}"), tstLeftNoMD, tstSCClose, tstSC1Inline, tstRightNoMD, tstEOF}},
	{"inline with nested shortcode (not supported)", `{{< sc1.inline >}}Hello World{{< sc1 >}}{{< /sc1.inline >}}`, []Item{tstLeftNoMD, tstSC1Inline, tstRightNoMD, tstText, nti(tError, "inline shortcodes do not support nesting")}},
	{"inline case mismatch", `{{< sc1.Inline >}}Hello World{{< /sc1.Inline >}}`, []Item{tstLeftNoMD, nti(tError, "period in shortcode name only allowed for inline identifiers")}},
}

func TestShortcodeLexer(t *testing.T) {
	t.Parallel()
	for i, test := range shortCodeLexerTests {
		t.Run(test.name, func(t *testing.T) {
			items := collect([]byte(test.input), true, lexMainSection)
			if !equal(items, test.items) {
				t.Errorf("[%d] %s: got\n\t%v\nexpected\n\t%v", i, test.name, items, test.items)
			}
		})
	}
}

func BenchmarkShortcodeLexer(b *testing.B) {
	testInputs := make([][]byte, len(shortCodeLexerTests))
	for i, input := range shortCodeLexerTests {
		testInputs[i] = []byte(input.input)
	}
	var cfg Config
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range testInputs {
			items := collectWithConfig(input, true, lexMainSection, cfg)
			if len(items) == 0 {
			}

		}
	}
}
