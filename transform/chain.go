package transform

import (
	"bytes"
	"io"
)

type chain struct {
	transformers []Transformer
}

func NewChain(trs ...Transformer) Transformer {
	return &chain{transformers: trs}
}

func (c *chain) Apply(w io.Writer, r io.Reader) (err error) {
	in := r
	for _, tr := range c.transformers {
		out := new(bytes.Buffer)
		err = tr.Apply(out, in)
		if err != nil {
			return
		}
		in = bytes.NewBuffer(out.Bytes())
	}

	_, err = io.Copy(w, in)
	return
}
