package hugolib

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/tpl"
	"github.com/spf13/viper"
)

func pageFromString(in, filename string) (*Page, error) {
	return NewPageFrom(strings.NewReader(in), filename)
}

func CheckShortCodeMatch(t *testing.T, input, expected string, template tpl.Template) {

	p, _ := pageFromString(SIMPLE_PAGE, "simple.md")
	output, err := HandleShortcodes(input, p, template)

	if err != nil {
		t.Fatalf("Shortcode rendered error %s. Expected: %q, Got: %q", err, expected, output)
	}

	if output != expected {
		t.Fatalf("Shortcode render didn't match. got %q but exxpected %q", output, expected)
	}
}

func TestShortcodeGoFuzzReports(t *testing.T) {
	tem := tpl.New()

	tem.AddInternalShortcode("sc.html", `foo`)
	p, _ := pageFromString(SIMPLE_PAGE, "simple.md")

	for i, this := range []struct {
		data      string
		expectErr bool
	}{
		{"{{</*/", true},
	} {
		output, err := HandleShortcodes(this.data, p, tem)

		if this.expectErr && err == nil {
			t.Errorf("[%d] should have errored", i)
		}

		if !this.expectErr && err != nil {
			t.Errorf("[%d] should not have errored: %s", i, err)
		}

		if !this.expectErr && err == nil && len(output) == 0 {
			t.Errorf("[%d] empty result", i)
		}
	}

}

func TestNonSC(t *testing.T) {
	tem := tpl.New()
	// notice the syntax diff from 0.12, now comment delims must be added
	CheckShortCodeMatch(t, "{{%/* movie 47238zzb */%}}", "{{% movie 47238zzb %}}", tem)
}

// Issue #929
func TestHyphenatedSC(t *testing.T) {
	tem := tpl.New()
	tem.AddInternalShortcode("hyphenated-video.html", `Playing Video {{ .Get 0 }}`)

	CheckShortCodeMatch(t, "{{< hyphenated-video 47238zzb >}}", "Playing Video 47238zzb", tem)
}

func TestPositionalParamSC(t *testing.T) {
	tem := tpl.New()
	tem.AddInternalShortcode("video.html", `Playing Video {{ .Get 0 }}`)

	CheckShortCodeMatch(t, "{{< video 47238zzb >}}", "Playing Video 47238zzb", tem)
	CheckShortCodeMatch(t, "{{< video 47238zzb 132 >}}", "Playing Video 47238zzb", tem)
	CheckShortCodeMatch(t, "{{<video 47238zzb>}}", "Playing Video 47238zzb", tem)
	CheckShortCodeMatch(t, "{{<video 47238zzb    >}}", "Playing Video 47238zzb", tem)
	CheckShortCodeMatch(t, "{{<   video   47238zzb    >}}", "Playing Video 47238zzb", tem)
}

func TestNamedParamSC(t *testing.T) {
	tem := tpl.New()
	tem.AddInternalShortcode("img.html", `<img{{ with .Get "src" }} src="{{.}}"{{end}}{{with .Get "class"}} class="{{.}}"{{end}}>`)

	CheckShortCodeMatch(t, `{{< img src="one" >}}`, `<img src="one">`, tem)
	CheckShortCodeMatch(t, `{{< img class="aspen" >}}`, `<img class="aspen">`, tem)
	CheckShortCodeMatch(t, `{{< img src= "one" >}}`, `<img src="one">`, tem)
	CheckShortCodeMatch(t, `{{< img src ="one" >}}`, `<img src="one">`, tem)
	CheckShortCodeMatch(t, `{{< img src = "one" >}}`, `<img src="one">`, tem)
	CheckShortCodeMatch(t, `{{< img src = "one" class = "aspen grove" >}}`, `<img src="one" class="aspen grove">`, tem)
}

func TestInnerSC(t *testing.T) {
	tem := tpl.New()
	tem.AddInternalShortcode("inside.html", `<div{{with .Get "class"}} class="{{.}}"{{end}}>{{ .Inner }}</div>`)

	CheckShortCodeMatch(t, `{{< inside class="aspen" >}}`, `<div class="aspen"></div>`, tem)
	CheckShortCodeMatch(t, `{{< inside class="aspen" >}}More Here{{< /inside >}}`, "<div class=\"aspen\">More Here</div>", tem)
	CheckShortCodeMatch(t, `{{< inside >}}More Here{{< /inside >}}`, "<div>More Here</div>", tem)
}

func TestInnerSCWithMarkdown(t *testing.T) {
	tem := tpl.New()
	tem.AddInternalShortcode("inside.html", `<div{{with .Get "class"}} class="{{.}}"{{end}}>{{ .Inner }}</div>`)

	CheckShortCodeMatch(t, `{{% inside %}}
# More Here

[link](http://spf13.com) and text

{{% /inside %}}`, "<div><h1 id=\"more-here:bec3ed8ba720b9073ab75abcf3ba5d97\">More Here</h1>\n\n<p><a href=\"http://spf13.com\">link</a> and text</p>\n</div>", tem)
}

func TestInnerSCWithAndWithoutMarkdown(t *testing.T) {
	tem := tpl.New()
	tem.AddInternalShortcode("inside.html", `<div{{with .Get "class"}} class="{{.}}"{{end}}>{{ .Inner }}</div>`)

	CheckShortCodeMatch(t, `{{% inside %}}
# More Here

[link](http://spf13.com) and text

{{% /inside %}}

And then:

{{< inside >}}
# More Here

This is **plain** text.

{{< /inside >}}
`, "<div><h1 id=\"more-here:bec3ed8ba720b9073ab75abcf3ba5d97\">More Here</h1>\n\n<p><a href=\"http://spf13.com\">link</a> and text</p>\n</div>\n\nAnd then:\n\n<div>\n# More Here\n\nThis is **plain** text.\n\n</div>\n", tem)
}

func TestEmbeddedSC(t *testing.T) {
	tem := tpl.New()
	CheckShortCodeMatch(t, "{{% test %}}", "This is a simple Test", tem)
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" %}}`, "\n<figure class=\"bananas orange\">\n    \n        <img src=\"/found/here\" />\n    \n    \n</figure>\n", tem)
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" caption="This is a caption" %}}`, "\n<figure class=\"bananas orange\">\n    \n        <img src=\"/found/here\" alt=\"This is a caption\" />\n    \n    \n    <figcaption>\n        <p>\n        This is a caption\n        \n            \n        \n        </p> \n    </figcaption>\n    \n</figure>\n", tem)
}

func TestNestedSC(t *testing.T) {
	tem := tpl.New()
	tem.AddInternalShortcode("scn1.html", `<div>Outer, inner is {{ .Inner }}</div>`)
	tem.AddInternalShortcode("scn2.html", `<div>SC2</div>`)

	CheckShortCodeMatch(t, `{{% scn1 %}}{{% scn2 %}}{{% /scn1 %}}`, "<div>Outer, inner is <div>SC2</div>\n</div>", tem)

	CheckShortCodeMatch(t, `{{< scn1 >}}{{% scn2 %}}{{< /scn1 >}}`, "<div>Outer, inner is <div>SC2</div></div>", tem)
}

func TestNestedComplexSC(t *testing.T) {
	tem := tpl.New()
	tem.AddInternalShortcode("row.html", `-row-{{ .Inner}}-rowStop-`)
	tem.AddInternalShortcode("column.html", `-col-{{.Inner    }}-colStop-`)
	tem.AddInternalShortcode("aside.html", `-aside-{{    .Inner  }}-asideStop-`)

	CheckShortCodeMatch(t, `{{< row >}}1-s{{% column %}}2-**s**{{< aside >}}3-**s**{{< /aside >}}4-s{{% /column %}}5-s{{< /row >}}6-s`,
		"-row-1-s-col-2-<strong>s</strong>-aside-3-<strong>s</strong>-asideStop-4-s-colStop-5-s-rowStop-6-s", tem)

	// turn around the markup flag
	CheckShortCodeMatch(t, `{{% row %}}1-s{{< column >}}2-**s**{{% aside %}}3-**s**{{% /aside %}}4-s{{< /column >}}5-s{{% /row %}}6-s`,
		"-row-1-s-col-2-<strong>s</strong>-aside-3-<strong>s</strong>-asideStop-4-s-colStop-5-s-rowStop-6-s", tem)
}

func TestFigureImgWidth(t *testing.T) {
	tem := tpl.New()
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" alt="apple" width="100px" %}}`, "\n<figure class=\"bananas orange\">\n    \n        <img src=\"/found/here\" alt=\"apple\" width=\"100px\" />\n    \n    \n</figure>\n", tem)
}

func TestHighlight(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	if !helpers.HasPygments() {
		t.Skip("Skip test as Pygments is not installed")
	}
	viper.Set("PygmentsStyle", "bw")

	tem := tpl.New()

	code := `
{{< highlight java >}}
void do();
{{< /highlight >}}`
	CheckShortCodeMatch(t, code, "\n<div class=\"highlight\" style=\"background: #ffffff\"><pre style=\"line-height: 125%\"><span style=\"font-weight: bold\">void</span> do();\n</pre></div>\n", tem)
}

const testScPlaceholderRegexp = "{@{@HUGOSHORTCODE-\\d+@}@}"

func TestExtractShortcodes(t *testing.T) {
	for i, this := range []struct {
		name             string
		input            string
		expectShortCodes string
		expect           interface{}
		expectErrorMsg   string
	}{
		{"text", "Some text.", "map[]", "Some text.", ""},
		{"invalid right delim", "{{< tag }}", "", false, "simple:4:.*unrecognized character.*}"},
		{"invalid close", "\n{{< /tag >}}", "", false, "simple:5:.*got closing shortcode, but none is open"},
		{"invalid close2", "\n\n{{< tag >}}{{< /anotherTag >}}", "", false, "simple:6: closing tag for shortcode 'anotherTag' does not match start tag"},
		{"unterminated quote 1", `{{< figure src="im caption="S" >}}`, "", false, "simple:4:.got pos.*"},
		{"unterminated quote 1", `{{< figure src="im" caption="S >}}`, "", false, "simple:4:.*unterm.*}"},
		{"one shortcode, no markup", "{{< tag >}}", "", testScPlaceholderRegexp, ""},
		{"one shortcode, markup", "{{% tag %}}", "", testScPlaceholderRegexp, ""},
		{"one pos param", "{{% tag param1 %}}", `tag([\"param1\"], true){[]}"]`, testScPlaceholderRegexp, ""},
		{"two pos params", "{{< tag param1 param2>}}", `tag([\"param1\" \"param2\"], false){[]}"]`, testScPlaceholderRegexp, ""},
		{"one named param", `{{% tag param1="value" %}}`, `tag([\"param1:value\"], true){[]}`, testScPlaceholderRegexp, ""},
		{"two named params", `{{< tag param1="value1" param2="value2" >}}`, `tag([\"param1:value1\" \"param2:value2\"], false){[]}"]`,
			testScPlaceholderRegexp, ""},
		{"inner", `Some text. {{< inner >}}Inner Content{{< / inner >}}. Some more text.`, `inner([], false){[Inner Content]}`,
			fmt.Sprintf("Some text. %s. Some more text.", testScPlaceholderRegexp), ""},
		// issue #934
		{"inner self-closing", `Some text. {{< inner />}}. Some more text.`, `inner([], false){[]}`,
			fmt.Sprintf("Some text. %s. Some more text.", testScPlaceholderRegexp), ""},
		{"close, but not inner", "{{< tag >}}foo{{< /tag >}}", "", false, "Shortcode 'tag' in page 'simple.md' has no .Inner.*"},
		{"nested inner", `Inner->{{< inner >}}Inner Content->{{% inner2 param1 %}}inner2txt{{% /inner2 %}}Inner close->{{< / inner >}}<-done`,
			`inner([], false){[Inner Content-> inner2([\"param1\"], true){[inner2txt]} Inner close->]}`,
			fmt.Sprintf("Inner->%s<-done", testScPlaceholderRegexp), ""},
		{"nested, nested inner", `Inner->{{< inner >}}inner2->{{% inner2 param1 %}}inner2txt->inner3{{< inner3>}}inner3txt{{</ inner3 >}}{{% /inner2 %}}final close->{{< / inner >}}<-done`,
			`inner([], false){[inner2-> inner2([\"param1\"], true){[inner2txt->inner3 inner3(%!q(<nil>), false){[inner3txt]}]} final close->`,
			fmt.Sprintf("Inner->%s<-done", testScPlaceholderRegexp), ""},
		{"two inner", `Some text. {{% inner %}}First **Inner** Content{{% / inner %}} {{< inner >}}Inner **Content**{{< / inner >}}. Some more text.`,
			`map["{@{@HUGOSHORTCODE-1@}@}:inner([], true){[First **Inner** Content]}" "{@{@HUGOSHORTCODE-2@}@}:inner([], false){[Inner **Content**]}"]`,
			fmt.Sprintf("Some text. %s %s. Some more text.", testScPlaceholderRegexp, testScPlaceholderRegexp), ""},
		{"closed without content", `Some text. {{< inner param1 >}}{{< / inner >}}. Some more text.`, `inner([\"param1\"], false){[]}`,
			fmt.Sprintf("Some text. %s. Some more text.", testScPlaceholderRegexp), ""},
		{"two shortcodes", "{{< sc1 >}}{{< sc2 >}}",
			`map["{@{@HUGOSHORTCODE-1@}@}:sc1([], false){[]}" "{@{@HUGOSHORTCODE-2@}@}:sc2([], false){[]}"]`,
			testScPlaceholderRegexp + testScPlaceholderRegexp, ""},
		{"mix of shortcodes", `Hello {{< sc1 >}}world{{% sc2 p2="2"%}}. And that's it.`,
			`map["{@{@HUGOSHORTCODE-1@}@}:sc1([], false){[]}" "{@{@HUGOSHORTCODE-2@}@}:sc2([\"p2:2\"]`,
			fmt.Sprintf("Hello %sworld%s. And that's it.", testScPlaceholderRegexp, testScPlaceholderRegexp), ""},
		{"mix with inner", `Hello {{< sc1 >}}world{{% inner p2="2"%}}Inner{{%/ inner %}}. And that's it.`,
			`map["{@{@HUGOSHORTCODE-1@}@}:sc1([], false){[]}" "{@{@HUGOSHORTCODE-2@}@}:inner([\"p2:2\"], true){[Inner]}"]`,
			fmt.Sprintf("Hello %sworld%s. And that's it.", testScPlaceholderRegexp, testScPlaceholderRegexp), ""},
	} {

		p, _ := pageFromString(SIMPLE_PAGE, "simple.md")
		tem := tpl.New()
		tem.AddInternalShortcode("tag.html", `tag`)
		tem.AddInternalShortcode("sc1.html", `sc1`)
		tem.AddInternalShortcode("sc2.html", `sc2`)
		tem.AddInternalShortcode("inner.html", `{{with .Inner }}{{ . }}{{ end }}`)
		tem.AddInternalShortcode("inner2.html", `{{.Inner}}`)
		tem.AddInternalShortcode("inner3.html", `{{.Inner}}`)

		content, shortCodes, err := extractShortcodes(this.input, p, tem)

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Fatalf("[%d] %s: ExtractShortcodes didn't return an expected error", i, this.name)
			} else {
				r, _ := regexp.Compile(this.expectErrorMsg)
				if !r.MatchString(err.Error()) {
					t.Fatalf("[%d] %s: ExtractShortcodes didn't return an expected error message, got %s but expected %s",
						i, this.name, err.Error(), this.expectErrorMsg)
				}
			}
			continue
		} else {
			if err != nil {
				t.Fatalf("[%d] %s: failed: %q", i, this.name, err)
			}
		}

		var expected string
		av := reflect.ValueOf(this.expect)
		switch av.Kind() {
		case reflect.String:
			expected = av.String()
		}

		r, err := regexp.Compile(expected)

		if err != nil {
			t.Fatalf("[%d] %s: Failed to compile regexp %q: %q", i, this.name, expected, err)
		}

		if strings.Count(content, shortcodePlaceholderPrefix) != len(shortCodes) {
			t.Fatalf("[%d] %s: Not enough placeholders, found %d", i, this.name, len(shortCodes))
		}

		if !r.MatchString(content) {
			t.Fatalf("[%d] %s: Shortcode extract didn't match. got %q but expected %q", i, this.name, content, expected)
		}

		for placeHolder, sc := range shortCodes {
			if !strings.Contains(content, placeHolder) {
				t.Fatalf("[%d] %s: Output does not contain placeholder %q", i, this.name, placeHolder)
			}

			if sc.params == nil {
				t.Fatalf("[%d] %s: Params is nil for shortcode '%s'", i, this.name, sc.name)
			}
		}

		if this.expectShortCodes != "" {
			shortCodesAsStr := fmt.Sprintf("map%q", collectAndSortShortcodes(shortCodes))
			if !strings.Contains(shortCodesAsStr, this.expectShortCodes) {
				t.Fatalf("[%d] %s: Short codes not as expected, got %s but expected %s", i, this.name, shortCodesAsStr, this.expectShortCodes)
			}
		}
	}
}

func collectAndSortShortcodes(shortcodes map[string]shortcode) []string {
	var asArray []string

	for key, sc := range shortcodes {
		asArray = append(asArray, fmt.Sprintf("%s:%s", key, sc))
	}

	sort.Strings(asArray)
	return asArray

}

func BenchmarkReplaceShortcodeTokens(b *testing.B) {

	data := []struct {
		input        string
		replacements map[string]string
		expect       interface{}
	}{
		{"Hello {@{@HUGOSHORTCODE-1@}@}.", map[string]string{"{@{@HUGOSHORTCODE-1@}@}": "World"}, "Hello World."},
		{strings.Repeat("A", 100) + " {@{@HUGOSHORTCODE-1@}@}.", map[string]string{"{@{@HUGOSHORTCODE-1@}@}": "Hello World"}, strings.Repeat("A", 100) + " Hello World."},
		{strings.Repeat("A", 500) + " {@{@HUGOSHORTCODE-1@}@}.", map[string]string{"{@{@HUGOSHORTCODE-1@}@}": "Hello World"}, strings.Repeat("A", 500) + " Hello World."},
		{strings.Repeat("ABCD ", 500) + " {@{@HUGOSHORTCODE-1@}@}.", map[string]string{"{@{@HUGOSHORTCODE-1@}@}": "Hello World"}, strings.Repeat("ABCD ", 500) + " Hello World."},
		{strings.Repeat("A", 500) + " {@{@HUGOSHORTCODE-1@}@}." + strings.Repeat("BC", 500) + " {@{@HUGOSHORTCODE-1@}@}.", map[string]string{"{@{@HUGOSHORTCODE-1@}@}": "Hello World"}, strings.Repeat("A", 500) + " Hello World." + strings.Repeat("BC", 500) + " Hello World."},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i, this := range data {
			results, err := replaceShortcodeTokens([]byte(this.input), "HUGOSHORTCODE", this.replacements)

			if expectSuccess, ok := this.expect.(bool); ok && !expectSuccess {
				if err == nil {
					b.Fatalf("[%d] replaceShortcodeTokens didn't return an expected error", i)
				}
			} else {
				if err != nil {
					b.Fatalf("[%d] failed: %s", i, err)
					continue
				}
				if !reflect.DeepEqual(results, []byte(this.expect.(string))) {
					b.Fatalf("[%d] replaceShortcodeTokens, got \n%q but expected \n%q", i, results, this.expect)
				}
			}

		}

	}
}

func TestReplaceShortcodeTokens(t *testing.T) {
	for i, this := range []struct {
		input        string
		prefix       string
		replacements map[string]string
		expect       interface{}
	}{
		{"Hello {@{@PREFIX-1@}@}.", "PREFIX", map[string]string{"{@{@PREFIX-1@}@}": "World"}, "Hello World."},
		{"A {@{@A-1@}@} asdf {@{@A-2@}@}.", "A", map[string]string{"{@{@A-1@}@}": "v1", "{@{@A-2@}@}": "v2"}, "A v1 asdf v2."},
		{"Hello {@{@PREFIX2-1@}@}. Go {@{@PREFIX2-2@}@}, Go, Go {@{@PREFIX2-3@}@} Go Go!.", "PREFIX2", map[string]string{"{@{@PREFIX2-1@}@}": "Europe", "{@{@PREFIX2-2@}@}": "Jonny", "{@{@PREFIX2-3@}@}": "Johnny"}, "Hello Europe. Go Jonny, Go, Go Johnny Go Go!."},
		{"A {@{@PREFIX-2@}@} {@{@PREFIX-1@}@}.", "PREFIX", map[string]string{"{@{@PREFIX-1@}@}": "A", "{@{@PREFIX-2@}@}": "B"}, "A B A."},
		{"A {@{@PREFIX-1@}@} {@{@PREFIX-2", "PREFIX", map[string]string{"{@{@PREFIX-1@}@}": "A"}, false},
		{"A {@{@PREFIX-1@}@} but not the second.", "PREFIX", map[string]string{"{@{@PREFIX-1@}@}": "A", "{@{@PREFIX-2@}@}": "B"}, "A A but not the second."},
		{"An {@{@PREFIX-1@}@}.", "PREFIX", map[string]string{"{@{@PREFIX-1@}@}": "A", "{@{@PREFIX-2@}@}": "B"}, "An A."},
		{"An {@{@PREFIX-1@}@} {@{@PREFIX-2@}@}.", "PREFIX", map[string]string{"{@{@PREFIX-1@}@}": "A", "{@{@PREFIX-2@}@}": "B"}, "An A B."},
		{"A {@{@PREFIX-1@}@} {@{@PREFIX-2@}@} {@{@PREFIX-3@}@} {@{@PREFIX-1@}@} {@{@PREFIX-3@}@}.", "PREFIX", map[string]string{"{@{@PREFIX-1@}@}": "A", "{@{@PREFIX-2@}@}": "B", "{@{@PREFIX-3@}@}": "C"}, "A A B C A C."},
		{"A {@{@PREFIX-1@}@} {@{@PREFIX-2@}@} {@{@PREFIX-3@}@} {@{@PREFIX-1@}@} {@{@PREFIX-3@}@}.", "PREFIX", map[string]string{"{@{@PREFIX-1@}@}": "A", "{@{@PREFIX-2@}@}": "B", "{@{@PREFIX-3@}@}": "C"}, "A A B C A C."},
		// Issue #1148 remove p-tags 10 =>
		{"Hello <p>{@{@PREFIX-1@}@}</p>. END.", "PREFIX", map[string]string{"{@{@PREFIX-1@}@}": "World"}, "Hello World. END."},
		{"Hello <p>{@{@PREFIX-1@}@}</p>. <p>{@{@PREFIX-2@}@}</p> END.", "PREFIX", map[string]string{"{@{@PREFIX-1@}@}": "World", "{@{@PREFIX-2@}@}": "THE"}, "Hello World. THE END."},
		{"Hello <p>{@{@PREFIX-1@}@}. END</p>.", "PREFIX", map[string]string{"{@{@PREFIX-1@}@}": "World"}, "Hello <p>World. END</p>."},
		{"<p>Hello {@{@PREFIX-1@}@}</p>. END.", "PREFIX", map[string]string{"{@{@PREFIX-1@}@}": "World"}, "<p>Hello World</p>. END."},
		{"Hello <p>{@{@PREFIX-1@}@}12", "PREFIX", map[string]string{"{@{@PREFIX-1@}@}": "World"}, "Hello <p>World12"},
		// Make sure the buffering expands as needed
		{"Hello {@{@P-1@}@}. {@{@P-1@}@}-{@{@P-1@}@} {@{@P-1@}@} {@{@P-1@}@} {@{@P-1@}@} END", "P", map[string]string{"{@{@P-1@}@}": strings.Repeat("BC", 100)},
			fmt.Sprintf("Hello %s. %s-%s %s %s %s END",
				strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100))},
	} {
		results, err := replaceShortcodeTokens([]byte(this.input), this.prefix, this.replacements)

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] replaceShortcodeTokens didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(results, []byte(this.expect.(string))) {
				t.Errorf("[%d] replaceShortcodeTokens, got \n%q but expected \n%q", i, results, this.expect)
			}
		}

	}

}
