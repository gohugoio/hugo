package transform

import (
	"bytes"
)

func LiveReloadInject(ct contentTransformer) {
	match := []byte("</body>")
	replace := []byte(`<script data-no-instant>document.write('<script src="/livereload.js?mindelay=10"></' + 'script>')</script></body>`)

	newcontent := bytes.Replace(ct.Content(), match, replace, -1)
	if len(newcontent) == len(ct.Content()) {
		match := []byte("</BODY>")
		newcontent = bytes.Replace(ct.Content(), match, replace, -1)
	}

	ct.Write(newcontent)
}
