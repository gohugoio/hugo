package transform

import (
	htmltran "code.google.com/p/go-html-transform/html/transform"
	"io"
)

type chain []*htmltran.Transform

func NewChain(trs ...*htmltran.Transform) chain {
	return trs
}

func (c *chain) Apply(w io.Writer, r io.Reader) (err error) {

	var tr *htmltran.Transformer

	if tr, err = htmltran.NewFromReader(r); err != nil {
		return
	}

	tr.ApplyAll(*c...)

	return tr.Render(w)
}
