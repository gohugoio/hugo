package transform

import (
	"bytes"
	bp "github.com/spf13/hugo/bufferpool"
	"io"
)

type trans func(rw contentTransformer)

type link trans

type chain []link

func NewChain(trs ...link) chain {
	return trs
}

func NewEmptyTransforms() []link {
	return make([]link, 0, 20)
}

// contentTransformer is an interface that enables rotation  of pooled buffers
// in the transformer chain.
type contentTransformer interface {
	Content() []byte
	io.Writer
}

// Implements contentTransformer
// Content is read from the from-buffer and rewritten to to the to-buffer.
type fromToBuffer struct {
	from *bytes.Buffer
	to   *bytes.Buffer
}

func (ft fromToBuffer) Write(p []byte) (n int, err error) {
	return ft.to.Write(p)
}

func (ft fromToBuffer) Content() []byte {
	return ft.from.Bytes()
}

func (c *chain) Apply(w io.Writer, r io.Reader) error {

	b1 := bp.GetBuffer()
	defer bp.PutBuffer(b1)

	b1.ReadFrom(r)

	if len(*c) == 0 {
		b1.WriteTo(w)
		return nil
	}

	b2 := bp.GetBuffer()
	defer bp.PutBuffer(b2)

	fb := &fromToBuffer{from: b1, to: b2}

	for i, tr := range *c {
		if i > 0 {
			if fb.from == b1 {
				fb.from = b2
				fb.to = b1
				fb.to.Reset()
			} else {
				fb.from = b1
				fb.to = b2
				fb.to.Reset()
			}
		}

		tr(fb)
	}

	fb.to.WriteTo(w)
	return nil
}
