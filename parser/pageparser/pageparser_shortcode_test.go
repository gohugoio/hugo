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

	qt "github.com/frankban/quicktest"
)

var (
	tstEOF       = nti(tEOF, "")
	tstLeftNoMD  = nti(tLeftDelimScNoMarkup, "{{<")
	tstRightNoMD = nti(tRightDelimScNoMarkup, ">}}")
	tstLeftMD    = nti(tLeftDelimScWithMarkup, "{{%")
	tstRightMD   = nti(tRightDelimScWithMarkup, "%}}")
	tstSCClose   = nti(tScClose, "/")
	tstSC1       = nti(tScName, "sc1")
	tstSC1Inline = nti(tScNameInline, "sc1.inline")
	tstSC2Inline = nti(tScNameInline, "sc2.inline")
	tstSC2       = nti(tScName, "sc2")
	tstSC3       = nti(tScName, "sc3")
	tstSCSlash   = nti(tScName, "sc/sub")
	tstParam1    = nti(tScParam, "param1")
	tstParam2    = nti(tScParam, "param2")
	tstVal       = nti(tScParamVal, "Hello World")
	tstText      = nti(tText, "Hello World")
)

var shortCodeLexerTests = []lexerTest{
	{"empty", "", []typeText{tstEOF}, nil},
	{"spaces", " \t\n", []typeText{nti(tText, " \t\n"), tstEOF}, nil},
	{"text", `to be or not`, []typeText{nti(tText, "to be or not"), tstEOF}, nil},
	{"no markup", `{{< sc1 >}}`, []typeText{tstLeftNoMD, tstSC1, tstRightNoMD, tstEOF}, nil},
	{"with EOL", "{{< sc1 \n >}}", []typeText{tstLeftNoMD, tstSC1, tstRightNoMD, tstEOF}, nil},

	{"forward slash inside name", `{{< sc/sub >}}`, []typeText{tstLeftNoMD, tstSCSlash, tstRightNoMD, tstEOF}, nil},

	{"simple with markup", `{{% sc1 %}}`, []typeText{tstLeftMD, tstSC1, tstRightMD, tstEOF}, nil},
	{"with spaces", `{{<     sc1     >}}`, []typeText{tstLeftNoMD, tstSC1, tstRightNoMD, tstEOF}, nil},
	{"indented on new line", "Hello\n    {{% sc1 %}}", []typeText{nti(tText, "Hello\n"), nti(tIndentation, "    "), tstLeftMD, tstSC1, tstRightMD, tstEOF}, nil},
	{"indented on new line tab", "Hello\n\t{{% sc1 %}}", []typeText{nti(tText, "Hello\n"), nti(tIndentation, "\t"), tstLeftMD, tstSC1, tstRightMD, tstEOF}, nil},
	{"indented on first line", "    {{% sc1 %}}", []typeText{nti(tIndentation, "    "), tstLeftMD, tstSC1, tstRightMD, tstEOF}, nil},
	{"mismatched rightDelim", `{{< sc1 %}}`, []typeText{
		tstLeftNoMD, tstSC1,
		nti(tError, "unrecognized character in shortcode action: U+0025 '%'. Note: Parameters with non-alphanumeric args must be quoted"),
	}, nil},
	{"inner, markup", `{{% sc1 %}} inner {{% /sc1 %}}`, []typeText{
		tstLeftMD,
		tstSC1,
		tstRightMD,
		nti(tText, " inner "),
		tstLeftMD,
		tstSCClose,
		tstSC1,
		tstRightMD,
		tstEOF,
	}, nil},
	{"close, but no open", `{{< /sc1 >}}`, []typeText{
		tstLeftNoMD, nti(tError, "got closing shortcode, but none is open"),
	}, nil},
	{"close wrong", `{{< sc1 >}}{{< /another >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose,
		nti(tError, "closing tag for shortcode 'another' does not match start tag"),
	}, nil},
	{"close, but no open, more", `{{< sc1 >}}{{< /sc1 >}}{{< /another >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose,
		nti(tError, "closing tag for shortcode 'another' does not match start tag"),
	}, nil},
	{"close with extra keyword", `{{< sc1 >}}{{< /sc1 keyword>}}`, []typeText{
		tstLeftNoMD, tstSC1, tstRightNoMD, tstLeftNoMD, tstSCClose, tstSC1,
		nti(tError, "unclosed shortcode"),
	}, nil},
	{"float param, positional", `{{< sc1 3.14 >}}`, []typeText{
		tstLeftNoMD, tstSC1, nti(tScParam, "3.14"), tstRightNoMD, tstEOF,
	}, nil},
	{"float param, named", `{{< sc1 param1=3.14 >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, nti(tScParamVal, "3.14"), tstRightNoMD, tstEOF,
	}, nil},
	{"named param, raw string", `{{< sc1 param1=` + "`" + "Hello World" + "`" + " >}}", []typeText{
		tstLeftNoMD, tstSC1, tstParam1, nti(tScParamVal, "Hello World"), tstRightNoMD, tstEOF,
	}, nil},
	{"float param, named, space before", `{{< sc1 param1= 3.14 >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, nti(tScParamVal, "3.14"), tstRightNoMD, tstEOF,
	}, nil},
	{"Youtube id", `{{< sc1 -ziL-Q_456igdO-4 >}}`, []typeText{
		tstLeftNoMD, tstSC1, nti(tScParam, "-ziL-Q_456igdO-4"), tstRightNoMD, tstEOF,
	}, nil},
	{"non-alphanumerics param quoted", `{{< sc1 "-ziL-.%QigdO-4" >}}`, []typeText{
		tstLeftNoMD, tstSC1, nti(tScParam, "-ziL-.%QigdO-4"), tstRightNoMD, tstEOF,
	}, nil},
	{"raw string", `{{< sc1` + "`" + "Hello World" + "`" + ` >}}`, []typeText{
		tstLeftNoMD, tstSC1, nti(tScParam, "Hello World"), tstRightNoMD, tstEOF,
	}, nil},
	{"raw string with newline", `{{< sc1` + "`" + `Hello 
	World` + "`" + ` >}}`, []typeText{
		tstLeftNoMD, tstSC1, nti(tScParam, `Hello 
	World`), tstRightNoMD, tstEOF,
	}, nil},
	{"raw string with escape character", `{{< sc1` + "`" + `Hello \b World` + "`" + ` >}}`, []typeText{
		tstLeftNoMD, tstSC1, nti(tScParam, `Hello \b World`), tstRightNoMD, tstEOF,
	}, nil},
	{"two params", `{{< sc1 param1   param2 >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, tstParam2, tstRightNoMD, tstEOF,
	}, nil},
	// issue #934
	{"self-closing", `{{< sc1 />}}`, []typeText{
		tstLeftNoMD, tstSC1, tstSCClose, tstRightNoMD, tstEOF,
	}, nil},
	// Issue 2498
	{"multiple self-closing", `{{< sc1 />}}{{< sc1 />}}`, []typeText{
		tstLeftNoMD, tstSC1, tstSCClose, tstRightNoMD,
		tstLeftNoMD, tstSC1, tstSCClose, tstRightNoMD, tstEOF,
	}, nil},
	{"self-closing with param", `{{< sc1 param1 />}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, tstSCClose, tstRightNoMD, tstEOF,
	}, nil},
	{"multiple self-closing with param", `{{< sc1 param1 />}}{{< sc1 param1 />}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, tstSCClose, tstRightNoMD,
		tstLeftNoMD, tstSC1, tstParam1, tstSCClose, tstRightNoMD, tstEOF,
	}, nil},
	{"multiple different self-closing with param", `{{< sc1 param1 />}}{{< sc2 param1 />}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, tstSCClose, tstRightNoMD,
		tstLeftNoMD, tstSC2, tstParam1, tstSCClose, tstRightNoMD, tstEOF,
	}, nil},
	{"nested simple", `{{< sc1 >}}{{< sc2 >}}{{< /sc1 >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstRightNoMD,
		tstLeftNoMD, tstSC2, tstRightNoMD,
		tstLeftNoMD, tstSCClose, tstSC1, tstRightNoMD, tstEOF,
	}, nil},
	{"nested complex", `{{< sc1 >}}ab{{% sc2 param1 %}}cd{{< sc3 >}}ef{{< /sc3 >}}gh{{% /sc2 %}}ij{{< /sc1 >}}kl`, []typeText{
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
	}, nil},

	{"two quoted params", `{{< sc1 "param nr. 1" "param nr. 2" >}}`, []typeText{
		tstLeftNoMD, tstSC1, nti(tScParam, "param nr. 1"), nti(tScParam, "param nr. 2"), tstRightNoMD, tstEOF,
	}, nil},
	{"two named params", `{{< sc1 param1="Hello World" param2="p2Val">}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, tstVal, tstParam2, nti(tScParamVal, "p2Val"), tstRightNoMD, tstEOF,
	}, nil},
	{"escaped quotes", `{{< sc1 param1=\"Hello World\"  >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, tstVal, tstRightNoMD, tstEOF,
	}, nil},
	{"escaped quotes, positional param", `{{< sc1 \"param1\"  >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, tstRightNoMD, tstEOF,
	}, nil},
	{"escaped quotes inside escaped quotes", `{{< sc1 param1=\"Hello \"escaped\" World\"  >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1,
		nti(tScParamVal, `Hello `), nti(tError, `got positional parameter 'escaped'. Cannot mix named and positional parameters`),
	}, nil},
	{
		"escaped quotes inside nonescaped quotes",
		`{{< sc1 param1="Hello \"escaped\" World"  >}}`,
		[]typeText{
			tstLeftNoMD, tstSC1, tstParam1, nti(tScParamVal, `Hello "escaped" World`), tstRightNoMD, tstEOF,
		},
		nil,
	},
	{
		"escaped quotes inside nonescaped quotes in positional param",
		`{{< sc1 "Hello \"escaped\" World"  >}}`,
		[]typeText{
			tstLeftNoMD, tstSC1, nti(tScParam, `Hello "escaped" World`), tstRightNoMD, tstEOF,
		},
		nil,
	},
	{"escaped raw string, named param", `{{< sc1 param1=` + `\` + "`" + "Hello World" + `\` + "`" + ` >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, nti(tError, "unrecognized escape character"),
	}, nil},
	{"escaped raw string, positional param", `{{< sc1 param1 ` + `\` + "`" + "Hello World" + `\` + "`" + ` >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, nti(tError, "unrecognized escape character"),
	}, nil},
	{"two raw string params", `{{< sc1` + "`" + "Hello World" + "`" + "`" + "Second Param" + "`" + ` >}}`, []typeText{
		tstLeftNoMD, tstSC1, nti(tScParam, "Hello World"), nti(tScParam, "Second Param"), tstRightNoMD, tstEOF,
	}, nil},
	{"unterminated quote", `{{< sc1 param2="Hello World>}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam2, nti(tError, "unterminated quoted string in shortcode parameter-argument: 'Hello World>}}'"),
	}, nil},
	{"unterminated raw string", `{{< sc1` + "`" + "Hello World" + ` >}}`, []typeText{
		tstLeftNoMD, tstSC1, nti(tError, "unterminated raw string in shortcode parameter-argument: 'Hello World >}}'"),
	}, nil},
	{"unterminated raw string in second argument", `{{< sc1` + "`" + "Hello World" + "`" + "`" + "Second Param" + ` >}}`, []typeText{
		tstLeftNoMD, tstSC1, nti(tScParam, "Hello World"), nti(tError, "unterminated raw string in shortcode parameter-argument: 'Second Param >}}'"),
	}, nil},
	{"one named param, one not", `{{< sc1 param1="Hello World" p2 >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, tstVal,
		nti(tError, "got positional parameter 'p2'. Cannot mix named and positional parameters"),
	}, nil},
	{"one named param, one quoted positional param, both raw strings", `{{< sc1 param1=` + "`" + "Hello World" + "`" + "`" + "Second Param" + "`" + ` >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, tstVal,
		nti(tError, "got quoted positional parameter. Cannot mix named and positional parameters"),
	}, nil},
	{"one named param, one quoted positional param", `{{< sc1 param1="Hello World" "And Universe" >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1, tstVal,
		nti(tError, "got quoted positional parameter. Cannot mix named and positional parameters"),
	}, nil},
	{"one quoted positional param, one named param", `{{< sc1 "param1" param2="And Universe" >}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1,
		nti(tError, "got named parameter 'param2'. Cannot mix named and positional parameters"),
	}, nil},
	{"ono positional param, one not", `{{< sc1 param1 param2="Hello World">}}`, []typeText{
		tstLeftNoMD, tstSC1, tstParam1,
		nti(tError, "got named parameter 'param2'. Cannot mix named and positional parameters"),
	}, nil},
	{"commented out", `{{</* sc1 */>}}`, []typeText{
		nti(tText, "{{<"), nti(tText, " sc1 "), nti(tText, ">}}"), tstEOF,
	}, nil},
	{"commented out, with asterisk inside", `{{</* sc1 "**/*.pdf" */>}}`, []typeText{
		nti(tText, "{{<"), nti(tText, " sc1 \"**/*.pdf\" "), nti(tText, ">}}"), tstEOF,
	}, nil},
	{"commented out, missing close", `{{</* sc1 >}}`, []typeText{
		nti(tError, "comment must be closed"),
	}, nil},
	{"commented out, misplaced close", `{{</* sc1 >}}*/`, []typeText{
		nti(tError, "comment must be closed"),
	}, nil},
	// Inline shortcodes
	{"basic inline", `{{< sc1.inline >}}Hello World{{< /sc1.inline >}}`, []typeText{tstLeftNoMD, tstSC1Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSCClose, tstSC1Inline, tstRightNoMD, tstEOF}, nil},
	{"basic inline with space", `{{< sc1.inline >}}Hello World{{< / sc1.inline >}}`, []typeText{tstLeftNoMD, tstSC1Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSCClose, tstSC1Inline, tstRightNoMD, tstEOF}, nil},
	{"inline self closing", `{{< sc1.inline >}}Hello World{{< /sc1.inline >}}Hello World{{< sc1.inline />}}`, []typeText{tstLeftNoMD, tstSC1Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSCClose, tstSC1Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSC1Inline, tstSCClose, tstRightNoMD, tstEOF}, nil},
	{"inline self closing, then a new inline", `{{< sc1.inline >}}Hello World{{< /sc1.inline >}}Hello World{{< sc1.inline />}}{{< sc2.inline >}}Hello World{{< /sc2.inline >}}`, []typeText{
		tstLeftNoMD, tstSC1Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSCClose, tstSC1Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSC1Inline, tstSCClose, tstRightNoMD,
		tstLeftNoMD, tstSC2Inline, tstRightNoMD, tstText, tstLeftNoMD, tstSCClose, tstSC2Inline, tstRightNoMD, tstEOF,
	}, nil},
	{"inline with template syntax", `{{< sc1.inline >}}{{ .Get 0 }}{{ .Get 1 }}{{< /sc1.inline >}}`, []typeText{tstLeftNoMD, tstSC1Inline, tstRightNoMD, nti(tText, "{{ .Get 0 }}"), nti(tText, "{{ .Get 1 }}"), tstLeftNoMD, tstSCClose, tstSC1Inline, tstRightNoMD, tstEOF}, nil},
	{"inline with nested shortcode (not supported)", `{{< sc1.inline >}}Hello World{{< sc1 >}}{{< /sc1.inline >}}`, []typeText{tstLeftNoMD, tstSC1Inline, tstRightNoMD, tstText, nti(tError, "inline shortcodes do not support nesting")}, nil},
	{"inline case mismatch", `{{< sc1.Inline >}}Hello World{{< /sc1.Inline >}}`, []typeText{tstLeftNoMD, nti(tError, "period in shortcode name only allowed for inline identifiers")}, nil},
}

func TestShortcodeLexer(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	for i, test := range shortCodeLexerTests {
		t.Run(test.name, func(t *testing.T) {
			items, err := collect([]byte(test.input), true, lexMainSection)
			c.Assert(err, qt.IsNil)
			if !equal(test.input, items, test.items) {
				got := itemsToString(items, []byte(test.input))
				expected := testItemsToString(test.items)
				c.Assert(got, qt.Equals, expected, qt.Commentf("Test %d: %s", i, test.name))
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
			_, err := collectWithConfig(input, true, lexMainSection, cfg)
			if err != nil {
				b.Fatal(err)
			}

		}
	}
}
