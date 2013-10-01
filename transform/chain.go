package transform

import (
	"io"
	"bytes"
)

type chain struct {
	transformers []Transformer
}

func NewChain(trs ...Transformer) Transformer {
	return &chain{transformers: trs}
}

func (c *chain) Apply(r io.Reader, w io.Writer) (err error) {
	in := r
	for _, tr := range c.transformers {
		out := new(bytes.Buffer)
		err = tr.Apply(in, out)
		if err != nil {
			return
		}
		in = bytes.NewBuffer(out.Bytes())
	}
	
	_, err = io.Copy(w, in)
	return
}
