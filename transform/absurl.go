package transform

import (
	htmltran "code.google.com/p/go-html-transform/html/transform"
	"net/url"
)

func AbsURL(absURL string) (trs []*htmltran.Transform, err error) {
	var baseURL *url.URL

	if baseURL, err = url.Parse(absURL); err != nil {
		return
	}

	if trs, err = absUrlify(baseURL, elattr{"a", "href"}, elattr{"script", "src"}); err != nil {
		return
	}
	return
}

type elattr struct {
	tag, attr string
}

func absUrlify(baseURL *url.URL, selectors ...elattr) (trs []*htmltran.Transform, err error) {
	var inURL *url.URL

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
		mt := htmltran.MustTrans(htmltran.TransformAttrib(el.attr, replace), el.tag)
		trs = append(trs, mt)
	}

	return
}

func fragmentOnly(u *url.URL) bool {
	return u.Fragment != "" && u.Scheme == "" && u.Opaque == "" && u.User == nil && u.Host == "" && u.Path == "" && u.Path == "" && u.RawQuery == ""
}
