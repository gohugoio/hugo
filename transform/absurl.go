package transform

import (
	"bytes"
	htmltran "code.google.com/p/go-html-transform/html/transform"
	"net/url"
	"strings"
)

func AbsURL(absURL string) (trs []link, err error) {
	var baseURL *url.URL

	if baseURL, err = url.Parse(absURL); err != nil {
		return
	}

	base := strings.TrimRight(baseURL.String(), "/")

	var (
		srcdq  = []byte(" src=\"" + base + "/")
		hrefdq = []byte(" href=\"" + base + "/")
		srcsq  = []byte(" src='" + base + "/")
		hrefsq = []byte(" href='" + base + "/")
	)
	trs = append(trs, func(content []byte) []byte {
		content = guardReplace(content, []byte(" src=\"//"), []byte(" src=\"/"), srcdq)
		content = guardReplace(content, []byte(" src='//"), []byte(" src='/"), srcsq)
		content = guardReplace(content, []byte(" href=\"//"), []byte(" href=\"/"), hrefdq)
		content = guardReplace(content, []byte(" href='//"), []byte(" href='/"), hrefsq)
		return content
	})
	return
}

func guardReplace(content, guard, match, replace []byte) []byte {
		if !bytes.Contains(content, guard) {
			content = bytes.Replace(content, match, replace, -1)
		}
		return content
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
