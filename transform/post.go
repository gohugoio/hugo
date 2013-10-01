package transform

import (
	"io"
)

type Transformer interface {
	Apply(io.Writer, io.Reader) error
}
