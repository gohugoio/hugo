package hugolib

import (
	"github.com/spf13/hugo/template/bundle"
	"strings"
	"testing"
)

func pageFromString(in, filename string) (*Page, error) {
	return ReadFrom(strings.NewReader(in), filename)
}

func CheckShortCodeMatch(t *testing.T, input, expected string, template bundle.Template) {

	p, _ := pageFromString(SIMPLE_PAGE, "simple.md")
	output := ShortcodesHandle(input, p, template)

	if output != expected {
		t.Fatalf("Shortcode render didn't match. Expected: %q, Got: %q", expected, output)
	}
}

func TestNonSC(t *testing.T) {
	tem := bundle.NewTemplate()

	CheckShortCodeMatch(t, "{{% movie 47238zzb %}}", "{{% movie 47238zzb %}}", tem)
}

func TestPositionalParamSC(t *testing.T) {
	tem := bundle.NewTemplate()
	tem.AddInternalShortcode("video.html", `Playing Video {{ .Get 0 }}`)

	CheckShortCodeMatch(t, "{{% video 47238zzb %}}", "Playing Video 47238zzb", tem)
	CheckShortCodeMatch(t, "{{% video 47238zzb 132 %}}", "Playing Video 47238zzb", tem)
	CheckShortCodeMatch(t, "{{%video 47238zzb%}}", "Playing Video 47238zzb", tem)
	CheckShortCodeMatch(t, "{{%video 47238zzb    %}}", "Playing Video 47238zzb", tem)
	CheckShortCodeMatch(t, "{{%   video   47238zzb    %}}", "Playing Video 47238zzb", tem)
}

func TestNamedParamSC(t *testing.T) {
	tem := bundle.NewTemplate()
	tem.AddInternalShortcode("img.html", `<img{{ with .Get "src" }} src="{{.}}"{{end}}{{with .Get "class"}} class="{{.}}"{{end}}>`)

	CheckShortCodeMatch(t, `{{% img src="one" %}}`, `<img src="one">`, tem)
	CheckShortCodeMatch(t, `{{% img class="aspen" %}}`, `<img class="aspen">`, tem)
	CheckShortCodeMatch(t, `{{% img src= "one" %}}`, `<img src="one">`, tem)
	CheckShortCodeMatch(t, `{{% img src ="one" %}}`, `<img src="one">`, tem)
	CheckShortCodeMatch(t, `{{% img src = "one" %}}`, `<img src="one">`, tem)
	CheckShortCodeMatch(t, `{{% img src = "one" class = "aspen grove" %}}`, `<img src="one" class="aspen grove">`, tem)
}

func TestInnerSC(t *testing.T) {
	tem := bundle.NewTemplate()
	tem.AddInternalShortcode("inside.html", `<div{{with .Get "class"}} class="{{.}}"{{end}}>{{ .Inner }}</div>`)

	CheckShortCodeMatch(t, `{{% inside class="aspen" %}}`, `<div class="aspen"></div>`, tem)
	CheckShortCodeMatch(t, `{{% inside class="aspen" %}}More Here{{% /inside %}}`, `<div class="aspen">More Here</div>`, tem)
	CheckShortCodeMatch(t, `{{% inside %}}More Here{{% /inside %}}`, `<div>More Here</div>`, tem)
}

func TestEmbeddedSC(t *testing.T) {
	tem := bundle.NewTemplate()
	CheckShortCodeMatch(t, "{{% test %}}", "This is a simple Test", tem)
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" %}}`, "\n<figure class=\"bananas orange\">\n    \n        <img src=\"/found/here\"  />\n    \n    \n</figure>\n", tem)
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" caption="This is a caption" %}}`, "\n<figure class=\"bananas orange\">\n    \n        <img src=\"/found/here\" alt=\"This is a caption\" />\n    \n    \n    <figcaption>\n        <p>\n        This is a caption\n        \n            \n        \n        </p> \n    </figcaption>\n    \n</figure>\n", tem)
}

func TestUnbalancedQuotes(t *testing.T) {
	tem := bundle.NewTemplate()

	CheckShortCodeMatch(t, `{{% figure src="/uploads/2011/12/spf13-mongosv-speaking-copy-1024x749.jpg "Steve Francia speaking at OSCON 2012" alt="MongoSV 2011" %}}`, "\n<figure >\n    \n        <img src=\"/uploads/2011/12/spf13-mongosv-speaking-copy-1024x749.jpg%20%22Steve%20Francia%20speaking%20at%20OSCON%202012\" alt=\"MongoSV 2011\" />\n    \n    \n</figure>\n", tem)
}
