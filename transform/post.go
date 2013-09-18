package transform

import (
	htmltran "code.google.com/p/go-html-transform/html/transform"
	"io"
	"net/url"
)

type Transformer struct {
	BaseURL string
}

func (t *Transformer) Apply(r io.Reader, w io.Writer) (err error) {
	var tr *htmltran.Transformer

	if tr, err = htmltran.NewFromReader(r); err != nil {
		return
	}

	if err = t.absUrlify(tr, elattr{"a", "href"}, elattr{"script", "src"}); err != nil {
		return
	}

	return tr.Render(w)
}

type elattr struct {
	tag, attr string
}

func (t *Transformer) absUrlify(tr *htmltran.Transformer, selectors ...elattr) (err error) {
	var baseURL, inURL *url.URL

	if baseURL, err = url.Parse(t.BaseURL); err != nil {
		return
	}

	replace := func(in string) string {
		if inURL, err = url.Parse(in); err != nil {
			return in + "?"
		}
		if fragmentOnly(inURL) {
			return in
		}
		return baseURL.ResolveReference(inURL).String()
	}

	for _, el := range selectors {
		if err = tr.Apply(htmltran.TransformAttrib(el.attr, replace), el.tag); err != nil {
			return
		}
	}

	return
}

func fragmentOnly(u *url.URL) bool {
	return u.Fragment != "" && u.Scheme == "" && u.Opaque == "" && u.User == nil && u.Host == "" && u.Path == "" && u.Path == "" && u.RawQuery == ""
}
