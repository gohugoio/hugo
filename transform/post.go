package transform

import (
	"io"
)

type Transformer interface {
	Apply(io.Reader, io.Writer) error
}
