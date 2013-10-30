package transform

import (
	"bytes"
	"strings"
	"testing"
)

const HTML_WITH_NAV = `<!DOCTYPE html>
<html>
<head></head>
<body>
<nav>
	<ul class="nav navbar-nav">
		<li hugo-nav="section_1"><a href="#">Section 1</a></li>
		<li hugo-nav="section_2"><a href="#">Section 2</a></li>
	</ul>
</nav>
</body>
</html>
`
const EXPECTED_HTML_WITH_NAV_1 = `<!DOCTYPE html><html><head></head>
<body>
<nav>
	<ul class="nav navbar-nav">
		<li hugo-nav="section_1"><a href="#">Section 1</a></li>
		<li hugo-nav="section_2" class="active"><a href="#">Section 2</a></li>
	</ul>
</nav>


</body></html>`

func TestDegenerateNoSectionSet(t *testing.T) {
	var (
		tr  = new(NavActive)
		out = new(bytes.Buffer)
	)

	if err := tr.Apply(out, strings.NewReader(HTML_WITH_NAV)); err != nil {
		t.Errorf("Unexpected error in NavActive.Apply: %s", err)
	}

	if out.String() != HTML_WITH_NAV {
		t.Errorf("NavActive.Apply should simply pass along the buffer unmodified.")
	}
}

func TestSetNav(t *testing.T) {
	tr := &NavActive{Section: "section_2"}
	out := new(bytes.Buffer)
	if err := tr.Apply(out, strings.NewReader(HTML_WITH_NAV)); err != nil {
		t.Errorf("Unexpected error in Apply() for NavActive: %s", err)
	}

	expected := EXPECTED_HTML_WITH_NAV_1
	if out.String() != expected {
		t.Errorf("NavActive.Apply output expected and got:\n%q\n%q", expected, out.String())
	}
}

func BenchmarkTransform(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tr := &NavActive{Section: "section_2"}
		out := new(bytes.Buffer)
		if err := tr.Apply(out, strings.NewReader(HTML_WITH_NAV)); err != nil {
			b.Errorf("Unexpected error in Apply() for NavActive: %s", err)
		}
	}
}
