package transform

import (
	"io"

	bp "github.com/spf13/hugo/bufferpool"
)

type trans func([]byte) []byte

type link trans

type chain []link

func NewChain(trs ...link) chain {
	return trs
}

func NewEmptyTransforms() []link {
	return make([]link, 0, 20)
}

func (c *chain) Apply(w io.Writer, r io.Reader) (err error) {
	buffer := bp.GetBuffer()
	defer bp.PutBuffer(buffer)

	buffer.ReadFrom(r)
	b := buffer.Bytes()
	for _, tr := range *c {
		b = tr(b)
	}
	buffer.Reset()
	buffer.Write(b)
	buffer.WriteTo(w)
	return
}
