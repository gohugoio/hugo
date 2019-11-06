// TODO(bep) this is a very temporary fork based on https://github.com/yuin/goldmark-highlighting/pull/10
// MIT Licensed, Copyright Yusuke Inuzuka
package temphighlighting

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/yuin/goldmark/util"

	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
)

type preWrapper struct {
	language string
}

func (p preWrapper) Start(code bool, styleAttr string) string {
	w := &strings.Builder{}
	fmt.Fprintf(w, "<pre%s><code", styleAttr)
	if p.language != "" {
		fmt.Fprintf(w, " class=\"language-"+p.language)
	}
	fmt.Fprint(w, ">")
	return w.String()
}

func (p preWrapper) End(code bool) string {
	return "</code></pre>"
}

func TestHighlighting(t *testing.T) {
	var css bytes.Buffer
	markdown := goldmark.New(
		goldmark.WithExtensions(
			NewHighlighting(
				WithStyle("monokai"),
				WithCSSWriter(&css),
				WithFormatOptions(
					chromahtml.WithClasses(true),
					chromahtml.WithLineNumbers(false),
				),
				WithWrapperRenderer(func(w util.BufWriter, c CodeBlockContext, entering bool) {
					_, ok := c.Language()
					if entering {
						if !ok {
							w.WriteString("<pre><code>")
							return
						}
						w.WriteString(`<div class="highlight">`)
					} else {
						if !ok {
							w.WriteString("</pre></code>")
							return
						}
						w.WriteString(`</div>`)
					}
				}),
				WithCodeBlockOptions(func(c CodeBlockContext) []chromahtml.Option {
					if language, ok := c.Language(); ok {
						// Turn on line numbers for Go only.
						if string(language) == "go" {
							return []chromahtml.Option{
								chromahtml.WithLineNumbers(true),
							}
						}
					}
					return nil
				}),
			),
		),
	)
	var buffer bytes.Buffer
	if err := markdown.Convert([]byte(`
Title
=======
`+"``` go\n"+`func main() {
    fmt.Println("ok")
}
`+"```"+`
`), &buffer); err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(buffer.String()) != strings.TrimSpace(`
<h1>Title</h1>
<div class="highlight"><pre class="chroma"><span class="ln">1</span><span class="kd">func</span> <span class="nf">main</span><span class="p">(</span><span class="p">)</span> <span class="p">{</span>
<span class="ln">2</span>    <span class="nx">fmt</span><span class="p">.</span><span class="nf">Println</span><span class="p">(</span><span class="s">&#34;ok&#34;</span><span class="p">)</span>
<span class="ln">3</span><span class="p">}</span>
</pre></div>
`) {
		t.Error("failed to render HTML")
	}

	if strings.TrimSpace(css.String()) != strings.TrimSpace(`/* Background */ .chroma { color: #f8f8f2; background-color: #272822 }
/* Error */ .chroma .err { color: #960050; background-color: #1e0010 }
/* LineTableTD */ .chroma .lntd { vertical-align: top; padding: 0; margin: 0; border: 0; }
/* LineTable */ .chroma .lntable { border-spacing: 0; padding: 0; margin: 0; border: 0; width: auto; overflow: auto; display: block; }
/* LineHighlight */ .chroma .hl { display: block; width: 100%;background-color: #3c3d38 }
/* LineNumbersTable */ .chroma .lnt { margin-right: 0.4em; padding: 0 0.4em 0 0.4em;color: #7f7f7f }
/* LineNumbers */ .chroma .ln { margin-right: 0.4em; padding: 0 0.4em 0 0.4em;color: #7f7f7f }
/* Keyword */ .chroma .k { color: #66d9ef }
/* KeywordConstant */ .chroma .kc { color: #66d9ef }
/* KeywordDeclaration */ .chroma .kd { color: #66d9ef }
/* KeywordNamespace */ .chroma .kn { color: #f92672 }
/* KeywordPseudo */ .chroma .kp { color: #66d9ef }
/* KeywordReserved */ .chroma .kr { color: #66d9ef }
/* KeywordType */ .chroma .kt { color: #66d9ef }
/* NameAttribute */ .chroma .na { color: #a6e22e }
/* NameClass */ .chroma .nc { color: #a6e22e }
/* NameConstant */ .chroma .no { color: #66d9ef }
/* NameDecorator */ .chroma .nd { color: #a6e22e }
/* NameException */ .chroma .ne { color: #a6e22e }
/* NameFunction */ .chroma .nf { color: #a6e22e }
/* NameOther */ .chroma .nx { color: #a6e22e }
/* NameTag */ .chroma .nt { color: #f92672 }
/* Literal */ .chroma .l { color: #ae81ff }
/* LiteralDate */ .chroma .ld { color: #e6db74 }
/* LiteralString */ .chroma .s { color: #e6db74 }
/* LiteralStringAffix */ .chroma .sa { color: #e6db74 }
/* LiteralStringBacktick */ .chroma .sb { color: #e6db74 }
/* LiteralStringChar */ .chroma .sc { color: #e6db74 }
/* LiteralStringDelimiter */ .chroma .dl { color: #e6db74 }
/* LiteralStringDoc */ .chroma .sd { color: #e6db74 }
/* LiteralStringDouble */ .chroma .s2 { color: #e6db74 }
/* LiteralStringEscape */ .chroma .se { color: #ae81ff }
/* LiteralStringHeredoc */ .chroma .sh { color: #e6db74 }
/* LiteralStringInterpol */ .chroma .si { color: #e6db74 }
/* LiteralStringOther */ .chroma .sx { color: #e6db74 }
/* LiteralStringRegex */ .chroma .sr { color: #e6db74 }
/* LiteralStringSingle */ .chroma .s1 { color: #e6db74 }
/* LiteralStringSymbol */ .chroma .ss { color: #e6db74 }
/* LiteralNumber */ .chroma .m { color: #ae81ff }
/* LiteralNumberBin */ .chroma .mb { color: #ae81ff }
/* LiteralNumberFloat */ .chroma .mf { color: #ae81ff }
/* LiteralNumberHex */ .chroma .mh { color: #ae81ff }
/* LiteralNumberInteger */ .chroma .mi { color: #ae81ff }
/* LiteralNumberIntegerLong */ .chroma .il { color: #ae81ff }
/* LiteralNumberOct */ .chroma .mo { color: #ae81ff }
/* Operator */ .chroma .o { color: #f92672 }
/* OperatorWord */ .chroma .ow { color: #f92672 }
/* Comment */ .chroma .c { color: #75715e }
/* CommentHashbang */ .chroma .ch { color: #75715e }
/* CommentMultiline */ .chroma .cm { color: #75715e }
/* CommentSingle */ .chroma .c1 { color: #75715e }
/* CommentSpecial */ .chroma .cs { color: #75715e }
/* CommentPreproc */ .chroma .cp { color: #75715e }
/* CommentPreprocFile */ .chroma .cpf { color: #75715e }
/* GenericDeleted */ .chroma .gd { color: #f92672 }
/* GenericEmph */ .chroma .ge { font-style: italic }
/* GenericInserted */ .chroma .gi { color: #a6e22e }
/* GenericStrong */ .chroma .gs { font-weight: bold }
/* GenericSubheading */ .chroma .gu { color: #75715e }`) {
		t.Error("failed to render CSS")
	}

}

func TestHighlighting2(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			Highlighting,
		),
	)
	var buffer bytes.Buffer
	if err := markdown.Convert([]byte(`
Title
=======
`+"```"+`
func main() {
    fmt.Println("ok")
}
`+"```"+`
`), &buffer); err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(buffer.String()) != strings.TrimSpace(`
<h1>Title</h1>
<pre><code>func main() {
    fmt.Println(&quot;ok&quot;)
}
</code></pre>
`) {
		t.Error("failed to render HTML")
	}
}

func TestHighlighting3(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			Highlighting,
		),
	)
	var buffer bytes.Buffer
	if err := markdown.Convert([]byte(`
Title
=======

`+"```"+`cpp {hl_lines=[1,2]}
#include <iostream>
int main() {
    std::cout<< "hello" << std::endl;
}
`+"```"+`
`), &buffer); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(buffer.String()) != strings.TrimSpace(`
<h1>Title</h1>
<pre style="background-color:#fff"><span style="display:block;width:100%;background-color:#e5e5e5"><span style="color:#999;font-weight:bold;font-style:italic">#</span><span style="color:#999;font-weight:bold;font-style:italic">include</span> <span style="color:#999;font-weight:bold;font-style:italic">&lt;iostream&gt;</span><span style="color:#999;font-weight:bold;font-style:italic">
</span></span><span style="display:block;width:100%;background-color:#e5e5e5"><span style="color:#999;font-weight:bold;font-style:italic"></span><span style="color:#458;font-weight:bold">int</span> <span style="color:#900;font-weight:bold">main</span>() {
</span>    std<span style="color:#000;font-weight:bold">:</span><span style="color:#000;font-weight:bold">:</span>cout<span style="color:#000;font-weight:bold">&lt;</span><span style="color:#000;font-weight:bold">&lt;</span> <span style="color:#d14"></span><span style="color:#d14">&#34;</span><span style="color:#d14">hello</span><span style="color:#d14">&#34;</span> <span style="color:#000;font-weight:bold">&lt;</span><span style="color:#000;font-weight:bold">&lt;</span> std<span style="color:#000;font-weight:bold">:</span><span style="color:#000;font-weight:bold">:</span>endl;
}
</pre>
`) {
		t.Error("failed to render HTML")
	}
}

func TestHighlightingHlLines(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			NewHighlighting(
				WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
		),
	)

	for i, test := range []struct {
		attributes string
		expect     []int
	}{
		{`hl_lines=["2"]`, []int{2}},
		{`hl_lines=["2-3",5],linenostart=5`, []int{2, 3, 5}},
		{`hl_lines=["2-3"]`, []int{2, 3}},
	} {

		t.Run(fmt.Sprint(i), func(t *testing.T) {

			var buffer bytes.Buffer
			codeBlock := fmt.Sprintf(`bash {%s}
LINE1
LINE2
LINE3
LINE4
LINE5
LINE6
LINE7
LINE8
`, test.attributes)

			if err := markdown.Convert([]byte(`
`+"```"+codeBlock+"```"+`
`), &buffer); err != nil {
				t.Fatal(err)
			}

			for _, line := range test.expect {
				expectStr := fmt.Sprintf("<span class=\"hl\">LINE%d\n</span>", line)
				if !strings.Contains(buffer.String(), expectStr) {
					t.Fatal("got\n", buffer.String(), "\nexpected\n", expectStr)
				}
			}
		})
	}

}

func TestHighlightingLinenos(t *testing.T) {

	outputLineNumbersInTable := `<div class="chroma">
<table class="lntable"><tr><td class="lntd">
<span class="lnt">1
</span></td>
<td class="lntd">
LINE1
</td></tr></table>
</div>`

	for i, test := range []struct {
		attributes         string
		lineNumbers        bool
		lineNumbersInTable bool
		expect             string
	}{
		{`linenos=true`, false, false, `<span class="ln">1</span>LINE1`},
		{`linenos=false`, false, false, `LINE1`},
		{``, true, false, `<span class="ln">1</span>LINE1`},
		{``, true, true, outputLineNumbersInTable},
		{`linenos=inline`, true, true, `<span class="ln">1</span>LINE1`},
		{`linenos=foo`, false, false, `<span class="ln">1</span>LINE1`},
		{`linenos=table`, false, false, outputLineNumbersInTable},
	} {

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			markdown := goldmark.New(
				goldmark.WithExtensions(
					NewHighlighting(
						WithFormatOptions(
							chromahtml.WithLineNumbers(test.lineNumbers),
							chromahtml.LineNumbersInTable(test.lineNumbersInTable),
							chromahtml.PreventSurroundingPre(true),
							chromahtml.WithClasses(true),
						),
					),
				),
			)

			var buffer bytes.Buffer
			codeBlock := fmt.Sprintf(`bash {%s}
LINE1
`, test.attributes)

			content := "```" + codeBlock + "```"

			if err := markdown.Convert([]byte(content), &buffer); err != nil {
				t.Fatal(err)
			}

			s := strings.TrimSpace(buffer.String())

			if s != test.expect {
				t.Fatal("got\n", s, "\nexpected\n", test.expect)
			}

		})
	}

}
