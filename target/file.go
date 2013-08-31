package target

import (
	"io"
)

type Publisher interface {
	Publish(string, io.Reader) error
}
