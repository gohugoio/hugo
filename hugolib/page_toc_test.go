package hugolib

import (
    "testing"
)

func TestTableOfContents(t *testing.T) {
    text := `
Blah blah blah blah blah.

## AA

Blah blah blah blah blah.

### AAA

Blah blah blah blah blah.

## BB

Blah blah blah blah blah.

### BBB

Blah blah blah blah blah.
`

    markdown := RemoveSummaryDivider([]byte(text))
    toc := string(tableOfContentsFromBytes(markdown))

    expected := `<nav>
<ul>
<li>
<ul>
<li><a href="#toc_0">AA</a>
<ul>
<li><a href="#toc_1">AAA</a></li>
</ul></li>
<li><a href="#toc_2">BB</a>
<ul>
<li><a href="#toc_3">BBB</a></li>
</ul></li>
</ul></li>
</ul>
</nav>
`

    if toc != expected {
        t.Errorf("Expected table of contents: %s, got: %s", expected, toc)
    }
}
