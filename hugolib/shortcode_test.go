package hugolib

import (
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/viper"
	"strings"
	"testing"
)

func pageFromString(in, filename string) (*Page, error) {
	return NewPageFrom(strings.NewReader(in), filename)
}

func CheckShortCodeMatch(t *testing.T, input, expected string, template Template) {

	p, _ := pageFromString(SIMPLE_PAGE, "simple.md")
	output := ShortcodesHandle(input, p, template)

	if output != expected {
		t.Fatalf("Shortcode render didn't match. Expected: %q, Got: %q", expected, output)
	}
}

func TestNonSC(t *testing.T) {
	tem := NewTemplate()

	CheckShortCodeMatch(t, "{{% movie 47238zzb %}}", "{{% movie 47238zzb %}}", tem)
}

func TestPositionalParamSC(t *testing.T) {
	tem := NewTemplate()
	tem.AddInternalShortcode("video.html", `Playing Video {{ .Get 0 }}`)

	CheckShortCodeMatch(t, "{{% video 47238zzb %}}", "Playing Video 47238zzb", tem)
	CheckShortCodeMatch(t, "{{% video 47238zzb 132 %}}", "Playing Video 47238zzb", tem)
	CheckShortCodeMatch(t, "{{%video 47238zzb%}}", "Playing Video 47238zzb", tem)
	CheckShortCodeMatch(t, "{{%video 47238zzb    %}}", "Playing Video 47238zzb", tem)
	CheckShortCodeMatch(t, "{{%   video   47238zzb    %}}", "Playing Video 47238zzb", tem)
}

func TestNamedParamSC(t *testing.T) {
	tem := NewTemplate()
	tem.AddInternalShortcode("img.html", `<img{{ with .Get "src" }} src="{{.}}"{{end}}{{with .Get "class"}} class="{{.}}"{{end}}>`)

	CheckShortCodeMatch(t, `{{% img src="one" %}}`, `<img src="one">`, tem)
	CheckShortCodeMatch(t, `{{% img class="aspen" %}}`, `<img class="aspen">`, tem)
	CheckShortCodeMatch(t, `{{% img src= "one" %}}`, `<img src="one">`, tem)
	CheckShortCodeMatch(t, `{{% img src ="one" %}}`, `<img src="one">`, tem)
	CheckShortCodeMatch(t, `{{% img src = "one" %}}`, `<img src="one">`, tem)
	CheckShortCodeMatch(t, `{{% img src = "one" class = "aspen grove" %}}`, `<img src="one" class="aspen grove">`, tem)
}

func TestInnerSC(t *testing.T) {
	tem := NewTemplate()
	tem.AddInternalShortcode("inside.html", `<div{{with .Get "class"}} class="{{.}}"{{end}}>{{ .Inner }}</div>`)

	CheckShortCodeMatch(t, `{{% inside class="aspen" %}}`, `<div class="aspen"></div>`, tem)
	CheckShortCodeMatch(t, `{{% inside class="aspen" %}}More Here{{% /inside %}}`, "<div class=\"aspen\"><p>More Here</p>\n</div>", tem)
	CheckShortCodeMatch(t, `{{% inside %}}More Here{{% /inside %}}`, "<div><p>More Here</p>\n</div>", tem)
}

func TestInnerSCWithMarkdown(t *testing.T) {
	tem := NewTemplate()
	tem.AddInternalShortcode("inside.html", `<div{{with .Get "class"}} class="{{.}}"{{end}}>{{ .Inner }}</div>`)

	CheckShortCodeMatch(t, `{{% inside %}}
# More Here

[link](http://spf13.com) and text

{{% /inside %}}`, "<div><h1>More Here</h1>\n\n<p><a href=\"http://spf13.com\">link</a> and text</p>\n</div>", tem)
}

func TestEmbeddedSC(t *testing.T) {
	tem := NewTemplate()
	CheckShortCodeMatch(t, "{{% test %}}", "This is a simple Test", tem)
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" %}}`, "\n<figure class=\"bananas orange\">\n    \n        <img src=\"/found/here\"  />\n    \n    \n</figure>\n", tem)
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" caption="This is a caption" %}}`, "\n<figure class=\"bananas orange\">\n    \n        <img src=\"/found/here\" alt=\"This is a caption\" />\n    \n    \n    <figcaption>\n        <p>\n        This is a caption\n        \n            \n        \n        </p> \n    </figcaption>\n    \n</figure>\n", tem)
}

func TestFigureImgWidth(t *testing.T) {
	tem := NewTemplate()
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" width="100px" %}}`, "\n<figure class=\"bananas orange\">\n    \n        <img src=\"/found/here\" width=\"100px\" />\n    \n    \n</figure>\n", tem)
}

func TestUnbalancedQuotes(t *testing.T) {
	tem := NewTemplate()

	CheckShortCodeMatch(t, `{{% figure src="/uploads/2011/12/spf13-mongosv-speaking-copy-1024x749.jpg "Steve Francia speaking at OSCON 2012" alt="MongoSV 2011" %}}`, "\n<figure >\n    \n        <img src=\"/uploads/2011/12/spf13-mongosv-speaking-copy-1024x749.jpg%20%22Steve%20Francia%20speaking%20at%20OSCON%202012\" alt=\"MongoSV 2011\" />\n    \n    \n</figure>\n", tem)
}

func TestHighlight(t *testing.T) {
	if !helpers.HasPygments() {
		t.Skip("Skip test as Pygments is not installed")
	}
	defer viper.Set("PygmentsStyle", viper.Get("PygmentsStyle"))
	viper.Set("PygmentsStyle", "bw")

	tem := NewTemplate()

	code := `
{{% highlight java %}}
void do();
{{% /highlight %}}`
	CheckShortCodeMatch(t, code, "\n<div class=\"highlight\" style=\"background: #ffffff\"><pre style=\"line-height: 125%\"><span style=\"font-weight: bold\">void</span> do();\n</pre></div>\n", tem)
}
