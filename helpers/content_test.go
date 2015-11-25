package helpers

import (
	"html/template"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const tstHTMLContent = "<!DOCTYPE html><html><head><script src=\"http://two/foobar.js\"></script></head><body><nav><ul><li hugo-nav=\"section_0\"></li><li hugo-nav=\"section_1\"></li></ul></nav><article>content <a href=\"http://two/foobar\">foobar</a>. Follow up</article><p>This is some text.<br>And some more.</p></body></html>"

func TestStripHTML(t *testing.T) {
	type test struct {
		input, expected string
	}
	data := []test{
		{"<h1>strip h1 tag <h1>", "strip h1 tag "},
		{"<p> strip p tag </p>", " strip p tag \n"},
		{"</br> strip br<br>", " strip br\n"},
		{"</br> strip br2<br />", " strip br2\n"},
		{"This <strong>is</strong> a\nnewline", "This is a newline"},
	}
	for i, d := range data {
		output := StripHTML(d.input)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}
	}
}

func BenchmarkStripHTML(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StripHTML(tstHTMLContent)
	}
}

func TestStripEmptyNav(t *testing.T) {
	cleaned := StripEmptyNav([]byte("do<nav>\n</nav>\n\nbedobedo"))
	assert.Equal(t, []byte("dobedobedo"), cleaned)
}

func TestBytesToHTML(t *testing.T) {
	assert.Equal(t, template.HTML("dobedobedo"), BytesToHTML([]byte("dobedobedo")))
}

func TestTruncateWordsToWholeSentence(t *testing.T) {
	type test struct {
		input, expected string
		max             int
		truncated       bool
	}
	data := []test{
		{"a b c", "a b c", 12, false},
		{"a b c", "a b c", 3, false},
		{"a", "a", 1, false},
		{"This is a sentence.", "This is a sentence.", 5, false},
		{"This is also a sentence!", "This is also a sentence!", 1, false},
		{"To be. Or not to be. That's the question.", "To be.", 1, true},
		{" \nThis is not a sentence\n ", "This is not a", 4, true},
	}
	for i, d := range data {
		output, truncated := TruncateWordsToWholeSentence(strings.Fields(d.input), d.max)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}

		if d.truncated != truncated {
			t.Errorf("Test %d failed. Expected truncated=%t got %t", i, d.truncated, truncated)
		}
	}
}

func TestTruncateWordsByRune(t *testing.T) {
	type test struct {
		input, expected string
		max             int
		truncated       bool
	}
	data := []test{
		{"", "", 1, false},
		{"a b c", "a b c", 12, false},
		{"a b c", "a b c", 3, false},
		{"a", "a", 1, false},
		{"Hello 中国", "", 0, true},
		{"这是中文，全中文。", "这是中文，", 5, true},
		{"Hello 中国", "Hello 中", 2, true},
		{"Hello 中国", "Hello 中国", 3, false},
		{"Hello中国 Good 好的", "Hello中国 Good 好", 9, true},
		{"This is a sentence.", "This is", 2, true},
		{"This is also a sentence!", "This", 1, true},
		{"To be. Or not to be. That's the question.", "To be. Or not", 4, true},
		{" \nThis is    not a sentence\n ", "This is not", 3, true},
	}
	for i, d := range data {
		output, truncated := TruncateWordsByRune(strings.Fields(d.input), d.max)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}

		if d.truncated != truncated {
			t.Errorf("Test %d failed. Expected truncated=%t got %t", i, d.truncated, truncated)
		}
	}
}
