package transform

import (
	htmltran "code.google.com/p/go-html-transform/html/transform"
	"io"
	"fmt"
)

type NavActive struct {
	Section string
}

func (n *NavActive) Apply(r io.Reader, w io.Writer) (err error) {
	var tr *htmltran.Transformer

	if n.Section == "" {
		_, err = io.Copy(w, r)
		return
	}

	if tr, err = htmltran.NewFromReader(r); err != nil {
		return
	}

	tr.Apply(htmltran.ModifyAttrib("class", "active"), fmt.Sprintf("li[data-nav=%s]", n.Section))

	return tr.Render(w)
}
