package target

import (
	"bytes"
	"io"
)

type InMemoryTarget struct {
	Files map[string][]byte
}

func (t *InMemoryTarget) Publish(label string, reader io.Reader) (err error) {
	if t.Files == nil {
		t.Files = make(map[string][]byte)
	}
	bytes := new(bytes.Buffer)
	bytes.ReadFrom(reader)
	t.Files[label] = bytes.Bytes()
	return
}

func (t *InMemoryTarget) Translate(label string) (dest string, err error) {
	return label, nil
}
