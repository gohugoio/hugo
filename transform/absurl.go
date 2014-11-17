package transform

import (
	"bytes"
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
