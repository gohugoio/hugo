package transform

import (
	"io"
	"net/url"
	htmltran "code.google.com/p/go-html-transform/html/transform"
)

type Transformer struct {
	BaseURL string
}

func (t *Transformer) Apply(r io.Reader, w io.Writer) (err error) {
	var tr *htmltran.Transformer

	if tr, err = htmltran.NewFromReader(r); err != nil {
		return
	}

	if err = t.absUrlify(tr); err != nil {
		return
	}

	return tr.Render(w)
}

func (t *Transformer) absUrlify(tr *htmltran.Transformer) (err error) {
	var baseURL, inURL *url.URL

	if baseURL, err = url.Parse(t.BaseURL); err != nil {
		return
	}

	replace := func(in string) string {
		if inURL, err = url.Parse(in); err != nil {
			return in + "?"
		}
		return baseURL.ResolveReference(inURL).String()
	}

	if err = tr.Apply(htmltran.TransformAttrib("src", replace), "script"); err != nil {
		return
	}

	return tr.Apply(htmltran.TransformAttrib("href", replace), "a")
}
