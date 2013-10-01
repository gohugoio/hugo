package transform

import (
	htmltran "code.google.com/p/go-html-transform/html/transform"
	"fmt"
	"io"
)

type NavActive struct {
	Section  string
	AttrName string
}

func (n *NavActive) Apply(w io.Writer, r io.Reader) (err error) {
	var tr *htmltran.Transformer

	if n.Section == "" {
		_, err = io.Copy(w, r)
		return
	}

	if tr, err = htmltran.NewFromReader(r); err != nil {
		return
	}

	if n.AttrName == "" {
		n.AttrName = "hugo-nav"
	}

	err = tr.Apply(htmltran.ModifyAttrib("class", "active"), fmt.Sprintf("li[%s=%s]", n.AttrName, n.Section))
	if err != nil {
		return
	}

	return tr.Render(w)
}
