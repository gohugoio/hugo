package source

import (
	"bytes"
	"fmt"
)

type ByteSource struct {
	Name    string
	Content []byte
}

func (b *ByteSource) String() string {
	return fmt.Sprintf("%s %s", b.Name, string(b.Content))
}

type InMemorySource struct {
	ByteSource []ByteSource
}

func (i *InMemorySource) Files() (files []*File) {
	files = make([]*File, len(i.ByteSource))
	for i, fake := range i.ByteSource {
		files[i] = NewFileWithContents(fake.Name, bytes.NewReader(fake.Content))
	}
	return
}
